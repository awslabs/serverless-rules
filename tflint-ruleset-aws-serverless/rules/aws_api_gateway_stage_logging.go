package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsAPIGatewayStageLoggingRule checks whether "aws_api_gateway_stage" has Logging enabled.
type AwsAPIGatewayStageLoggingRule struct {
	resourceType string
	blockName    string
	tflint.DefaultRule
}

// NewAwsAPIGatewayStageLoggingRule returns new rule
func NewAwsAPIGatewayStageLoggingRule() *AwsAPIGatewayStageLoggingRule {
	return &AwsAPIGatewayStageLoggingRule{
		resourceType: "aws_api_gateway_stage",
		blockName:    "access_log_settings",
	}
}

// Name returns the rule name
func (r *AwsAPIGatewayStageLoggingRule) Name() string {
	return "aws_apigateway_stage_logging_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsAPIGatewayStageLoggingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsAPIGatewayStageLoggingRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsAPIGatewayStageLoggingRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/api_gateway/logging/"
}

// Metadata returns the rule metadata
func (r *AwsAPIGatewayStageLoggingRule) Metadata() interface{} {
	return struct {
		Name     string
		Severity tflint.Severity
		Link     string
	}{
		Name:     r.Name(),
		Severity: r.Severity(),
		Link:     r.Link(),
	}
}

// Check checks whether "aws_api_gateway_stage" has logging enabled
func (r *AwsAPIGatewayStageLoggingRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: r.blockName,
				Body: &hclext.BodySchema{},
			},
		},
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to get resource content: %w", err)
	}

	for _, resource := range resources.Blocks {
		blocks := resource.Body.Blocks.OfType(r.blockName)
		if len(blocks) == 0 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.blockName),
				(*resource).DefRange,
			)
		}
	}

	return nil
}
