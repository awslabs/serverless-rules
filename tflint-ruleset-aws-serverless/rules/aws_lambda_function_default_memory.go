package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLambdaFunctionDefaultMemory checks if there is an explicit memory size
type AwsLambdaFunctionDefaultMemoryRule struct {
	tflint.DefaultRule
	resourceType  string
	attributeName string
}

// NewAwsLambdaFunctionDefaultMemoryRule returns new rule with default attributes
func NewAwsLambdaFunctionDefaultMemoryRule() *AwsLambdaFunctionDefaultMemoryRule {
	return &AwsLambdaFunctionDefaultMemoryRule{
		resourceType:  "aws_lambda_function",
		attributeName: "memory_size",
	}
}

// Name returns the rule name
func (r *AwsLambdaFunctionDefaultMemoryRule) Name() string {
	return "aws_lambda_function_default_memory"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsLambdaFunctionDefaultMemoryRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsLambdaFunctionDefaultMemoryRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsLambdaFunctionDefaultMemoryRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/lambda/default_memory_size/"
}

// Metadata returns the metadata of the rule
func (r *AwsLambdaFunctionDefaultMemoryRule) Metadata() interface{} {
	return nil
}

// Check checks if there is an explicit memory size
func (r *AwsLambdaFunctionDefaultMemoryRule) Check(runner tflint.Runner) error {
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

		var attrValue string
		err := runner.EvaluateExpr(attribute.Expr, &attrValue, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
