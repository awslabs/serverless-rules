package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLambdaFunctionEolRuntime(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "EOL dotnetcore2.1",
			Content: `
resource "aws_lambda_function" "this" {
  runtime = "dotnetcore2.1"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaFunctionEolRuntimeRule(),
					Message: "\"dotnetcore2.1\" is an end-of-life runtime.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 13},
						End:      hcl.Pos{Line: 3, Column: 28},
					},
				},
			},
		},
		{
			Name: "EOL python2.7",
			Content: `
resource "aws_lambda_function" "this" {
  runtime = "python2.7"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaFunctionEolRuntimeRule(),
					Message: "\"python2.7\" is an end-of-life runtime.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 13},
						End:      hcl.Pos{Line: 3, Column: 24},
					},
				},
			},
		},
		{
			Name: "EOL ruby2.5",
			Content: `
resource "aws_lambda_function" "this" {
  runtime = "ruby2.5"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaFunctionEolRuntimeRule(),
					Message: "\"ruby2.5\" is an end-of-life runtime.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 13},
						End:      hcl.Pos{Line: 3, Column: 22},
					},
				},
			},
		},
		{
			Name: "EOL nodejs10.x",
			Content: `
resource "aws_lambda_function" "this" {
  runtime = "nodejs10.x"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaFunctionEolRuntimeRule(),
					Message: "\"nodejs10.x\" is an end-of-life runtime.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 13},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
			},
		},
		{
			Name: "not EOL",
			Content: `
resource "aws_lambda_function" "this" {
  runtime = "python3.8"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLambdaFunctionEolRuntimeRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
