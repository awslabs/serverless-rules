package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsAPIGatewayRestStageTracingRule checks whether "aws_api_gateway_stage" has tracing enabled.
type AwsAPIGatewayRestStageTracingRule struct {
	resourceType  string
	attributeName string
}

// NewAwsAPIGatewayRestStageTracingRule returns new rule
func NewAwsAPIGatewayRestStageTracingRule() *AwsAPIGatewayRestStageTracingRule {
	return &AwsAPIGatewayRestStageTracingRule{
		resourceType:  "aws_api_gateway_stage",
		attributeName: "xray_tracing_enabled",
	}
}

// Name returns the rule name
func (r *AwsAPIGatewayRestStageTracingRule) Name() string {
	return "aws_api_gateway_rest_api_tracing_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsAPIGatewayRestStageTracingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsAPIGatewayRestStageTracingRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsAPIGatewayRestStageTracingRule) Link() string {
	return ""
}

// Check checks whether "aws_api_gateway_stage" has tracing enabled
func (r *AwsAPIGatewayRestStageTracingRule) Check(runner tflint.Runner) error {
	return runner.WalkResources(r.resourceType, func(resource *configs.Resource) error {
		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name:     r.attributeName,
					Required: true,
				},
			},
		})

		if diags.HasErrors() {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				body.MissingItemRange,
			)
			return nil
		}

		attribute := body.Attributes[r.attributeName]
		var xrayTracingEnabled string
		err := runner.EvaluateExpr(attribute.Expr, &xrayTracingEnabled, nil)

		return runner.EnsureNoError(err, func() error {
			if xrayTracingEnabled != "true" {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("\"%s\" should be set to true.", r.attributeName),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
