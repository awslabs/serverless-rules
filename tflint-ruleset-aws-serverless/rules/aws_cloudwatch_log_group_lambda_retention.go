package rules

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type awsLambdaLogGroup struct {
	functionName string
	resourceName string
	found        bool
	blockRange   hcl.Range
}

// AwsCloudwatchLogGroupLambdaRetention checks if Lambda functions have a corresponding log group with retention configured
type AwsCloudwatchLogGroupLambdaRetentionRule struct {
	tflint.DefaultRule
	functionResourceType string
	logGroupResourceType string
	functionNameAttrName string
	nameAttrName         string
	retentionAttrName    string
}

// NewAwsCloudwatchLogGroupLambdaRetentionRule returns new rule with default attributes
func NewAwsCloudwatchLogGroupLambdaRetentionRule() *AwsCloudwatchLogGroupLambdaRetentionRule {
	return &AwsCloudwatchLogGroupLambdaRetentionRule{
		functionResourceType: "aws_lambda_function",
		logGroupResourceType: "aws_cloudwatch_log_group",
		functionNameAttrName: "function_name",
		nameAttrName:         "name",
		retentionAttrName:    "retention_in_days",
	}
}

// Name returns the rule name
func (r *AwsCloudwatchLogGroupLambdaRetentionRule) Name() string {
	return "aws_cloudwatch_log_group_lambda_retention"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsCloudwatchLogGroupLambdaRetentionRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsCloudwatchLogGroupLambdaRetentionRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsCloudwatchLogGroupLambdaRetentionRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/lambda/log_retention.html"
}

// Check checks if Lambda functions have a corresponding log group with retention configured
func (r *AwsCloudwatchLogGroupLambdaRetentionRule) Check(runner tflint.Runner) error {

	// Gather all Lambda functions
	var functions []awsLambdaLogGroup
	resources, err := runner.GetResourceContent(r.functionResourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.functionNameAttrName},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		function := awsLambdaLogGroup{
			resourceName: resource.Labels[0],
			functionName: "",
			found:        false,
			blockRange:   resource.DefRange,
		}

		// Function name attribute
		attribute, ok := resource.Body.Attributes[r.functionNameAttrName]
		if ok {
			var value string
			err := runner.EvaluateExpr(attribute.Expr, &value, nil)
			if err != nil {
				return err
			}
			function.functionName = value
		}

		functions = append(functions, function)
	}

	// Lookup log groups
	logGroups, err := runner.GetResourceContent(r.logGroupResourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.nameAttrName},
			{Name: r.retentionAttrName},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range logGroups.Blocks {
		// Get log group attributes
		nameAttr, ok := resource.Body.Attributes[r.nameAttrName]
		// No need to check further if there are no name, early return
		if !ok {
			continue
		}

		var nameValue string
		err := runner.EvaluateExpr(nameAttr.Expr, &nameValue, nil)
		if err != nil {
			return err
		}

		retentionAttr, ok := resource.Body.Attributes[r.retentionAttrName]
		// No need to check further if there are no retention, early return
		if !ok {
			continue
		}

		var retentionValue string
		err = runner.EvaluateExpr(retentionAttr.Expr, &retentionValue, nil)
		if err != nil {
			return err
		}

		// Parse log group names
		re := regexp.MustCompile(`^/aws/lambda/(.*)$`)
		m := re.FindAllStringSubmatch(nameValue, -1)
		// Log group name doesn't match pattern, early return
		if len(m) > 1 {
			continue
		}
		functionName := m[0][1]

		for i := range functions {
			if functions[i].functionName == functionName {
				functions[i].found = true
			}
			if fmt.Sprintf("${aws_lambda_function.%s.function_name}", functions[i].resourceName) == functionName {
				functions[i].found = true
			}
		}
	}

	for _, function := range functions {
		if !function.found {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is missing a log group with retention_in_days.", r.functionResourceType),
				function.blockRange,
			)
		}
	}

	return nil
}
