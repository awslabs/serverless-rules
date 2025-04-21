package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLambdaFunctionDefaultMemory(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "missing",
			Content: `
resource "aws_lambda_function" "this" {
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaFunctionDefaultMemoryRule(),
					Message: "\"memory_size\" is not present.",
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
	memory_size = 2048
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLambdaFunctionDefaultMemoryRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
