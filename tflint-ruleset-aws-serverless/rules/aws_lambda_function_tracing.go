package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLambdaFunctionTracingRule checks whether "aws_lambda_function" has tracing enabled.
type AwsLambdaFunctionTracingRule struct {
	resourceType  string
	blockName     string
	attributeName string
	tflint.DefaultRule
}

// NewAwsLambdaFunctionTracingRule returns new rule
func NewAwsLambdaFunctionTracingRule() *AwsLambdaFunctionTracingRule {
	return &AwsLambdaFunctionTracingRule{
		resourceType:  "aws_lambda_function",
		blockName:     "tracing_config",
		attributeName: "mode",
	}
}

// Name returns the rule name
func (r *AwsLambdaFunctionTracingRule) Name() string {
	return "aws_lambda_function_tracing_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsLambdaFunctionTracingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsLambdaFunctionTracingRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsLambdaFunctionTracingRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/lambda/tracing/"
}

// Metadata returns the rule metadata
func (r *AwsLambdaFunctionTracingRule) Metadata() interface{} {
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

// Check checks whether "aws_lambda_function" has tracing enabled
func (r *AwsLambdaFunctionTracingRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: r.blockName,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{
							Name:     r.attributeName,
							Required: false,
						},
					},
				},
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
			continue
		}

		block := blocks[0]
		attribute, ok := block.Body.Attributes[r.attributeName]
		if !ok {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s.%s\" is not present.", r.blockName, r.attributeName),
				block.DefRange,
			)
			continue
		}

		var xrayTracingEnabled string
		err := runner.EvaluateExpr(attribute.Expr, &xrayTracingEnabled, nil)
		if err != nil {
			return err
		}

		if xrayTracingEnabled != "Active" {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s.%s\" should be set to Active.", r.blockName, r.attributeName),
				attribute.Expr.Range(),
			)
		}
	}

	return nil
}
