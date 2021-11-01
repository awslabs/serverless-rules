package rules

import (
	"fmt"
	"regexp"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type awsLambdaLogGroup struct {
	functionName  string
	resourceName  string
	found         bool
	resourceRange hcl.Range
}

// AwsCloudwatchLogGroupLambdaRetention checks if Lambda functions have a corresponding log group with retention configured
type AwsCloudwatchLogGroupLambdaRetentionRule struct {
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
func (r *AwsCloudwatchLogGroupLambdaRetentionRule) Severity() string {
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
	err := runner.WalkResources(r.functionResourceType, func(resource *configs.Resource) error {
		function := awsLambdaLogGroup{
			resourceName:  resource.Name,
			functionName:  "",
			found:         false,
			resourceRange: resource.Config.MissingItemRange(),
		}

		// Function name attribute
		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name: r.functionNameAttrName,
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		attribute, ok := body.Attributes[r.functionNameAttrName]

		if ok {
			var value string
			err := runner.EvaluateExpr(attribute.Expr, &value, nil)
			if err != nil {
				return err
			}
			function.functionName = value
		}

		functions = append(functions, function)

		return nil
	})

	if err != nil {
		return err
	}

	// Lookup log groups
	err = runner.WalkResources(r.logGroupResourceType, func(resource *configs.Resource) error {
		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name: r.nameAttrName,
				},
				{
					Name: r.retentionAttrName,
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		// Get log group attributes

		nameAttr, ok := body.Attributes[r.nameAttrName]
		// No need to check further if there are no name, early return
		if !ok {
			return nil
		}

		var nameValue string
		err := runner.EvaluateExpr(nameAttr.Expr, &nameValue, nil)
		if err != nil {
			return err
		}

		retentionAttr, ok := body.Attributes[r.retentionAttrName]
		// No need to check further if there are no retention, early return
		if !ok {
			return nil
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
			return nil
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

		return nil
	})

	if err != nil {
		return err
	}

	for _, function := range functions {
		if !function.found {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is missing a log group with retention_in_days.", r.functionResourceType),
				function.resourceRange,
			)
		}
	}

	return nil
}
