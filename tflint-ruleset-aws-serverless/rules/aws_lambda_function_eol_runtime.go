package rules

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TODO: Write the rule's description here
// AwsLambdaFunctionEolRuntime checks if the runtime is marked as end-of-life
type AwsLambdaFunctionEolRuntimeRule struct {
	resourceType  string
	attributeName string
	enum          []string
}

// NewAwsLambdaFunctionEolRuntimeRule returns new rule with default attributes
func NewAwsLambdaFunctionEolRuntimeRule() *AwsLambdaFunctionEolRuntimeRule {
	return &AwsLambdaFunctionEolRuntimeRule{
		// TODO: Write resource type and attribute name here
		resourceType:  "aws_lambda_function",
		attributeName: "runtime",
		enum: []string{
			"dotnetcore2.1",
			"python2.7",
			"ruby2.5",
			"nodejs10.x",
			"nodejs8.10",
			"nodejs6.10",
			"nodejs4.3-edge",
			"nodejs4.3",
			"nodejs",
			"dotnetcore2.0",
			"dotnetcore1.0",
		},
	}
}

// Name returns the rule name
func (r *AwsLambdaFunctionEolRuntimeRule) Name() string {
	return "aws_lambda_function_eol_runtime"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsLambdaFunctionEolRuntimeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsLambdaFunctionEolRuntimeRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsLambdaFunctionEolRuntimeRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/lambda/end_of_life_runtime/"
}

// Check checks if the runtime is marked as end-of-life
func (r *AwsLambdaFunctionEolRuntimeRule) Check(runner tflint.Runner) error {
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val, nil)

		return runner.EnsureNoError(err, func() error {
			found := false
			for _, item := range r.enum {
				if item == val {
					found = true
				}
			}
			if found {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf(`"%s" is an end-of-life runtime.`, val),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
