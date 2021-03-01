package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsAPIGatewayRestStageLoggingRule checks whether "aws_api_gateway_stage" has Logging enabled.
type AwsAPIGatewayRestStageLoggingRule struct {
	resourceType string
	blockName    string
}

// NewAwsAPIGatewayRestStageLoggingRule returns new rule
func NewAwsAPIGatewayRestStageLoggingRule() *AwsAPIGatewayRestStageLoggingRule {
	return &AwsAPIGatewayRestStageLoggingRule{
		resourceType: "aws_api_gateway_stage",
		blockName:    "access_log_settings",
	}
}

// Name returns the rule name
func (r *AwsAPIGatewayRestStageLoggingRule) Name() string {
	return "aws_api_gateway_rest_api_Logging_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsAPIGatewayRestStageLoggingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsAPIGatewayRestStageLoggingRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsAPIGatewayRestStageLoggingRule) Link() string {
	return ""
}

// Check checks whether "aws_api_gateway_stage" has logging enabled
func (r *AwsAPIGatewayRestStageLoggingRule) Check(runner tflint.Runner) error {
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
