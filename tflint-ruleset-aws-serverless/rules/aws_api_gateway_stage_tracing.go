package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsAPIGatewayStageTracingRule checks whether "aws_api_gateway_stage" has tracing enabled.
type AwsAPIGatewayStageTracingRule struct {
	tflint.DefaultRule
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
func (r *AwsAPIGatewayStageTracingRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsAPIGatewayStageTracingRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/api_gateway/tracing/"
}

// Check checks whether "aws_api_gateway_stage" has tracing enabled
func (r *AwsAPIGatewayStageTracingRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.attributeName},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				resource.DefRange,
			)
			continue
		}

		var xrayTracingEnabled string
		err := runner.EvaluateExpr(attribute.Expr, &xrayTracingEnabled, nil)
		if err != nil {
			return err
		}

		if xrayTracingEnabled != "true" {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" should be set to true.", r.attributeName),
				attribute.Expr.Range(),
			)
		}
	}

	return nil
}
