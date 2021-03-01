package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsAPIGatewayStageV2LoggingRule checks whether "aws_api_gateway_stage" has Logging enabled.
type AwsAPIGatewayStageV2LoggingRule struct {
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
	return "aws_api_gateway_rest_api_Logging_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsAPIGatewayStageV2LoggingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsAPIGatewayStageV2LoggingRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsAPIGatewayStageV2LoggingRule) Link() string {
	return ""
}

// Check checks whether "aws_api_gateway_stage" has logging enabled
func (r *AwsAPIGatewayStageV2LoggingRule) Check(runner tflint.Runner) error {
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
