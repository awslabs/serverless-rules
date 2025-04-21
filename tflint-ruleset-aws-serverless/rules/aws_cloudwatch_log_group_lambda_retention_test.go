package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsCloudwatchLogGroupLambdaRetention(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "missing aws_cloudwatch_log_group",
			Content: `
resource "aws_lambda_function" "this" {
	function_name = "my-function-name"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsCloudwatchLogGroupLambdaRetentionRule(),
					Message: `"aws_lambda_function" is missing a log group with retention_in_days.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 38},
					},
				},
			},
		},
		{
			Name: "missing retention_in_days",
			Content: `
resource "aws_lambda_function" "this" {
	function_name = "my-function-name"
}

resource "aws_cloudwatch_log_group" "this" {
	name = "/aws/lambda/my-function-name"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsCloudwatchLogGroupLambdaRetentionRule(),
					Message: `"aws_lambda_function" is missing a log group with retention_in_days.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 38},
					},
				},
			},
		},
		{
			Name: "valid",
			Content: `
resource "aws_lambda_function" "this" {
	function_name = "my-function-name"
}

resource "aws_cloudwatch_log_group" "this" {
	name = "/aws/lambda/my-function-name"
	retention_in_days = 7
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "non-lambda",
			Content: `
resource "aws_cloudwatch_log_group" "this" {
	name = "not-lambda"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsCloudwatchLogGroupLambdaRetentionRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
