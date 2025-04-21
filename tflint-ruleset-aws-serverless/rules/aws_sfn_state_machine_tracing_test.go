package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSfnStateMachineTracing(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "missing tracing_configuration",
			Content: `
resource "aws_sfn_state_machine" "this" {
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSfnStateMachineTracingRule(),
					Message: "\"tracing_configuration\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 40},
					},
				},
			},
		},
		{
			Name: "missing enabled",
			Content: `
resource "aws_sfn_state_machine" "this" {
  tracing_configuration {}
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSfnStateMachineTracingRule(),
					Message: "\"enabled\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 24},
					},
				},
			},
		},
		{
			Name: "enabled false",
			Content: `
resource "aws_sfn_state_machine" "this" {
  tracing_configuration {
	enabled = false
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSfnStateMachineTracingRule(),
					Message: "\"enabled\" should be set to true.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 12},
						End:      hcl.Pos{Line: 4, Column: 17},
					},
				},
			},
		},
		{
			Name: "valid",
			Content: `
resource "aws_sfn_state_machine" "this" {
  tracing_configuration {
	enabled = true
  }
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsSfnStateMachineTracingRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
