package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLambdaEventInvokeConfigAsyncOnFailure(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "missing aws_lambda_function_event_invoke_config",
			Content: `
resource "aws_lambda_permission" "this" {
	action        = "lambda:InvokeFunction"
	function_name = "my-lambda-function"
	principal     = "sns.amazonaws.com"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaEventInvokeConfigAsyncOnFailureRule(),
					Message: "missing \"aws_lambda_function_event_invoke_config\" resource for function_name.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 18},
						End:      hcl.Pos{Line: 4, Column: 38},
					},
				},
			},
		},
		{
			Name: "missing destination_config",
			Content: `
resource "aws_lambda_permission" "this" {
	action        = "lambda:InvokeFunction"
	function_name = "my-lambda-function"
	principal     = "sns.amazonaws.com"
}

resource "aws_lambda_function_event_invoke_config" "example" {
	function_name = "my-lambda-function"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaEventInvokeConfigAsyncOnFailureRule(),
					Message: "\"destination_config\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 8, Column: 62},
						End:      hcl.Pos{Line: 8, Column: 62},
					},
				},
			},
		},
		{
			Name: "missing on_failure",
			Content: `
resource "aws_lambda_permission" "this" {
	action        = "lambda:InvokeFunction"
	function_name = "my-lambda-function"
	principal     = "sns.amazonaws.com"
}

resource "aws_lambda_function_event_invoke_config" "example" {
	function_name = "my-lambda-function"

	destination_config {}
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaEventInvokeConfigAsyncOnFailureRule(),
					Message: "\"on_failure\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 11, Column: 21},
						End:      hcl.Pos{Line: 11, Column: 21},
					},
				},
			},
		},
		{
			Name: "missing destination",
			Content: `
resource "aws_lambda_permission" "this" {
	action        = "lambda:InvokeFunction"
	function_name = "my-lambda-function"
	principal     = "sns.amazonaws.com"
}

resource "aws_lambda_function_event_invoke_config" "example" {
	function_name = "my-lambda-function"

	destination_config {
		on_failure {}
	}
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaEventInvokeConfigAsyncOnFailureRule(),
					Message: "\"destination\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 12, Column: 14},
						End:      hcl.Pos{Line: 12, Column: 14},
					},
				},
			},
		},
		{
			Name: "valid",
			Content: `
resource "aws_lambda_permission" "this" {
	action        = "lambda:InvokeFunction"
	function_name = "my-lambda-function"
	principal     = "sns.amazonaws.com"
}

resource "aws_lambda_function_event_invoke_config" "example" {
	function_name = "my-lambda-function"

	destination_config {
		on_failure {
			destination = "arn:aws:sqs:us-east-1:111122223333:my-dlq"
		}
	}
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "valid no async no event invoke config",
			Content: `
resource "aws_lambda_permission" "this" {
	action        = "lambda:InvokeFunction"
	function_name = "my-lambda-function"
	principal     = "apigateway.amazonaws.com"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "valid no async",
			Content: `
resource "aws_lambda_permission" "this" {
	action        = "lambda:InvokeFunction"
	function_name = "my-lambda-function"
	principal     = "apigateway.amazonaws.com"
}

resource "aws_lambda_function_event_invoke_config" "example" {
	function_name = "my-lambda-function"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLambdaEventInvokeConfigAsyncOnFailureRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
