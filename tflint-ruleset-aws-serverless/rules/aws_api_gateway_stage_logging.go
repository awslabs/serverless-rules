package rules

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsAPIGatewayStageLoggingRule checks whether "aws_api_gateway_stage" has Logging enabled.
type AwsAPIGatewayStageLoggingRule struct {
	resourceType string
	blockName    string
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

// Check checks whether "aws_api_gateway_stage" has logging enabled
func (r *AwsAPIGatewayStageLoggingRule) Check(runner tflint.Runner) error {
	return runner.WalkResources(r.resourceType, func(resource *configs.Resource) error {
		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{
				{
					Type: r.blockName,
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		blocks := body.Blocks.OfType(r.blockName)
		if len(blocks) != 1 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.blockName),
				body.MissingItemRange,
			)
		}

		return nil
	})
}
