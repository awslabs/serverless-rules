package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLambdaFunctionDefaultTimeout checks if there is an explicit timeout
type AwsLambdaFunctionDefaultTimeoutRule struct {
	resourceType  string
	attributeName string
	tflint.DefaultRule
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
func (r *AwsLambdaFunctionDefaultTimeoutRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsLambdaFunctionDefaultTimeoutRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/lambda/default_timeout/"
}

// Metadata returns the rule metadata
func (r *AwsLambdaFunctionDefaultTimeoutRule) Metadata() interface{} {
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

// Check checks if there is an explicit timeout
func (r *AwsLambdaFunctionDefaultTimeoutRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{
				Name: r.attributeName,
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, ok := resource.Body.Attributes[r.attributeName]
		if !ok {
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
