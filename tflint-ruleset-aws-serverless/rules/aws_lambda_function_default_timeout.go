package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLambdaFunctionDefaultTimeout checks if there is an explicit timeout
type AwsLambdaFunctionDefaultTimeoutRule struct {
	resourceType  string
	attributeName string
}

// NewAwsLambdaFunctionDefaultTimeoutRule returns new rule with default attributes
func NewAwsLambdaFunctionDefaultTimeoutRule() *AwsLambdaFunctionDefaultTimeoutRule {
	return &AwsLambdaFunctionDefaultTimeoutRule{
		resourceType:  "aws_lambda_function",
		attributeName: "timeout",
	}
}

// Name returns the rule name
func (r *AwsLambdaFunctionDefaultTimeoutRule) Name() string {
	return "aws_lambda_function_default_timeout"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsLambdaFunctionDefaultTimeoutRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsLambdaFunctionDefaultTimeoutRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsLambdaFunctionDefaultTimeoutRule) Link() string {
	return ""
}

// Check checks if there is an explicit timeout
func (r *AwsLambdaFunctionDefaultTimeoutRule) Check(runner tflint.Runner) error {
	return runner.WalkResources(r.resourceType, func(resource *configs.Resource) error {
		// Attribute
		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name: r.attributeName,
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		var attrValue string
		attribute, ok := body.Attributes[r.attributeName]
		if !ok {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				body.MissingItemRange,
			)
		} else {
			err := runner.EvaluateExpr(attribute.Expr, &attrValue, nil)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
