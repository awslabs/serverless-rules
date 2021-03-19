package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLambdaFunctionTracingRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "block not present",
			Content: `
resource "aws_lambda_function" "this" {}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaFunctionTracingRule(),
					Message: "\"tracing_config\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 39},
						End:      hcl.Pos{Line: 2, Column: 39},
					},
				},
			},
		},
		{
			Name: "attribute not present",
			Content: `
		resource "aws_lambda_function" "this" {
			tracing_config {}
		}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaFunctionTracingRule(),
					Message: "\"tracing_config.mode\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 19},
						End:      hcl.Pos{Line: 3, Column: 19},
					},
				},
			},
		},
		{
			Name: "invalid value",
			Content: `
resource "aws_lambda_function" "this" {
	tracing_config {
		mode = "PassThrough"
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaFunctionTracingRule(),
					Message: "\"tracing_config.mode\" should be set to Active.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 10},
						End:      hcl.Pos{Line: 4, Column: 23},
					},
				},
			},
		},
		{
			Name: "valid value",
			Content: `
resource "aws_lambda_function" "this" {
	tracing_config {
		mode = "Active"
	}
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLambdaFunctionTracingRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
