package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsAPIGatewayMethodSettingsThrottlingRule checks whether there is a default "aws_api_gateway_method_settings" resource with throttling values
type AwsAPIGatewayMethodSettingsThrottlingRule struct {
	tflint.DefaultRule
}

func NewAwsAPIGatewayMethodSettingsThrottlingRule() *AwsAPIGatewayMethodSettingsThrottlingRule {
	return &AwsAPIGatewayMethodSettingsThrottlingRule{}
}

// Name returns the rule name
func (r *AwsAPIGatewayMethodSettingsThrottlingRule) Name() string {
	return "aws_api_gateway_method_settings_throttling_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsAPIGatewayMethodSettingsThrottlingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsAPIGatewayMethodSettingsThrottlingRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsAPIGatewayMethodSettingsThrottlingRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/api_gateway/default_throttling/"
}

// Check checks whether default "aws_api_gateway_method_settings" have throttling values
func (r *AwsAPIGatewayMethodSettingsThrottlingRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("aws_api_gateway_method_settings", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "method_path"},
		},
		Blocks: []hclext.BlockSchema{
			{
				Type: "settings",
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "throttling_burst_limit"},
						{Name: "throttling_rate_limit"},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to get resource content: %w", err)
	}

	for _, resource := range resources.Blocks {
		// Only looking at the default method settings
		methodPath, exists := resource.Body.Attributes["method_path"]
		if !exists {
			continue
		}

		var path string
		if err := runner.EvaluateExpr(methodPath.Expr, &path, nil); err != nil {
			return fmt.Errorf("failed to evaluate method_path: %w", err)
		}

		if path != "*/*" {
			continue
		}

		// Load 'settings' block
		settingsBlocks := resource.Body.Blocks
		if len(settingsBlocks) == 0 {
			runner.EmitIssue(
				r,
				"\"settings\" block is required for default method settings",
				resource.DefRange,
			)
			continue
		}

		settings := settingsBlocks[0].Body
		// Check throttling limits
		if _, exists := settings.Attributes["throttling_burst_limit"]; !exists {
			runner.EmitIssue(
				r,
				"\"throttling_burst_limit\" is required for default method settings",
				settingsBlocks[0].DefRange,
			)
		}

		if _, exists := settings.Attributes["throttling_rate_limit"]; !exists {
			runner.EmitIssue(
				r,
				"\"throttling_rate_limit\" is required for default method settings",
				settingsBlocks[0].DefRange,
			)
		}
	}

	return nil
}
