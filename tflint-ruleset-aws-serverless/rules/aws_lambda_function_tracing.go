package rules

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLambdaFunctionTracingRule checks whether "aws_lambda_function" has tracing enabled.
type AwsLambdaFunctionTracingRule struct {
	resourceType  string
	blockName     string
	attributeName string
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

// Check checks whether "aws_lambda_function" has tracing enabled
func (r *AwsLambdaFunctionTracingRule) Check(runner tflint.Runner) error {
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

		// Check if the block exists
		blocks := body.Blocks.OfType(r.blockName)
		if len(blocks) != 1 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.blockName),
				body.MissingItemRange,
			)
			return nil
		}

		// Retrieve the block body
		blockBody, _, diags := blocks[0].Body.PartialContent(&hcl.BodySchema{
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
				fmt.Sprintf("\"%s.%s\" is not present.", r.blockName, r.attributeName),
				blockBody.MissingItemRange,
			)
			return nil
		}

		attribute := blockBody.Attributes[r.attributeName]

		var xrayTracingEnabled string
		err := runner.EvaluateExpr(attribute.Expr, &xrayTracingEnabled, nil)

		return runner.EnsureNoError(err, func() error {
			if xrayTracingEnabled != "Active" {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("\"%s.%s\" should be set to Active.", r.blockName, r.attributeName),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
