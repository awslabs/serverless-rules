package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLambdaPermissionMultiplePrincipals(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "multiple principals",
			Content: `
resource "aws_lambda_permission" "a" {
	function_name = "my-function"
	principal = "events.amazonaws.com"
}

resource "aws_lambda_permission" "b" {
	function_name = "my-function"
	principal = "sns.amazonaws.com"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaPermissionMultiplePrincipalsRule(),
					Message: `different "principal" values for the same function_name.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 14},
						End:      hcl.Pos{Line: 4, Column: 36},
					},
				},
				{
					Rule:    NewAwsLambdaPermissionMultiplePrincipalsRule(),
					Message: `different "principal" values for the same function_name.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 9, Column: 14},
						End:      hcl.Pos{Line: 9, Column: 33},
					},
				},
			},
		},
		{
			Name: "multiple principals",
			Content: `
resource "aws_lambda_permission" "a" {
	function_name = "my-function"
	principal = "events.amazonaws.com"
}

resource "aws_lambda_permission" "b" {
	function_name = "my-function"
	principal = "sns.amazonaws.com"
}

resource "aws_lambda_permission" "c" {
	function_name = "my-function"
	principal = "apigateway.amazonaws.com"
}

resource "aws_lambda_permission" "d" {
	function_name = "my-function"
	principal = "sqs.amazonaws.com"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaPermissionMultiplePrincipalsRule(),
					Message: `different "principal" values for the same function_name.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 14},
						End:      hcl.Pos{Line: 4, Column: 36},
					},
				},
				{
					Rule:    NewAwsLambdaPermissionMultiplePrincipalsRule(),
					Message: `different "principal" values for the same function_name.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 9, Column: 14},
						End:      hcl.Pos{Line: 9, Column: 33},
					},
				},
				{
					Rule:    NewAwsLambdaPermissionMultiplePrincipalsRule(),
					Message: `different "principal" values for the same function_name.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 14, Column: 14},
						End:      hcl.Pos{Line: 14, Column: 40},
					},
				},
				{
					Rule:    NewAwsLambdaPermissionMultiplePrincipalsRule(),
					Message: `different "principal" values for the same function_name.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 19, Column: 14},
						End:      hcl.Pos{Line: 19, Column: 33},
					},
				},
			},
		},
		{
			Name: "multiple principals",
			Content: `
resource "aws_lambda_permission" "a" {
	function_name = "my-function"
	principal = "events.amazonaws.com"
}

resource "aws_lambda_permission" "b" {
	function_name = "my-function"
	principal = "sns.amazonaws.com"
}

resource "aws_lambda_permission" "c" {
	function_name = "my-function-c"
	principal = "sns.amazonaws.com"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaPermissionMultiplePrincipalsRule(),
					Message: `different "principal" values for the same function_name.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 14},
						End:      hcl.Pos{Line: 4, Column: 36},
					},
				},
				{
					Rule:    NewAwsLambdaPermissionMultiplePrincipalsRule(),
					Message: `different "principal" values for the same function_name.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 9, Column: 14},
						End:      hcl.Pos{Line: 9, Column: 33},
					},
				},
			},
		},
		{
			Name: "multiple principals",
			Content: `
resource "aws_lambda_permission" "a" {
	function_name = "my-function"
	principal = "events.amazonaws.com"
}

resource "aws_lambda_permission" "b" {
	function_name = "my-function"
	principal = "events.amazonaws.com"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "multiple principals",
			Content: `
resource "aws_lambda_permission" "a" {
	function_name = "my-function-a"
	principal = "events.amazonaws.com"
}

resource "aws_lambda_permission" "b" {
	function_name = "my-function-b"
	principal = "sns.amazonaws.com"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLambdaPermissionMultiplePrincipalsRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
