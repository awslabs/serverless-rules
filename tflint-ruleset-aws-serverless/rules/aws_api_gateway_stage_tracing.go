package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsAPIGatewayStageTracingRule checks whether "aws_api_gateway_stage" has tracing enabled.
type AwsAPIGatewayStageTracingRule struct {
	resourceType  string
	attributeName string
}

// NewAwsAPIGatewayStageTracingRule returns new rule
func NewAwsAPIGatewayStageTracingRule() *AwsAPIGatewayStageTracingRule {
	return &AwsAPIGatewayStageTracingRule{
		resourceType:  "aws_api_gateway_stage",
		attributeName: "xray_tracing_enabled",
	}
}

// Name returns the rule name
func (r *AwsAPIGatewayStageTracingRule) Name() string {
	return "aws_apigateway_stage_tracing_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsAPIGatewayStageTracingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsAPIGatewayStageTracingRule) Severity() string {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsAPIGatewayStageTracingRule) Link() string {
	return ""
}

// Check checks whether "aws_api_gateway_stage" has tracing enabled
func (r *AwsAPIGatewayStageTracingRule) Check(runner tflint.Runner) error {
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
