package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsAPIGatewayStageV2LoggingRule checks whether "aws_api_gateway_stage" has Logging enabled.
type AwsAPIGatewayStageV2LoggingRule struct {
	tflint.DefaultRule
	resourceType string
	blockName    string
}

// NewAwsAPIGatewayStageV2LoggingRule returns new rule
func NewAwsAPIGatewayStageV2LoggingRule() *AwsAPIGatewayStageV2LoggingRule {
	return &AwsAPIGatewayStageV2LoggingRule{
		resourceType: "aws_api_gatewayv2_stage",
		blockName:    "access_log_settings",
	}
}

// Name returns the rule name
func (r *AwsAPIGatewayStageV2LoggingRule) Name() string {
	return "aws_apigatewayv2_stage_logging_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsAPIGatewayStageV2LoggingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsAPIGatewayStageV2LoggingRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsAPIGatewayStageV2LoggingRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/api_gateway/logging/"
}

// Check checks whether "aws_api_gateway_stage" has logging enabled
func (r *AwsAPIGatewayStageV2LoggingRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: r.blockName,
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		blocks := resource.Body.Blocks.OfType(r.blockName)
		if len(blocks) == 0 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.blockName),
				resource.DefRange,
			)
		}
	}

	return nil
}
