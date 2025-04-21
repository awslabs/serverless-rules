package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsCloudwatchEventTargetNoDlq(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "missing dead_letter_config",
			Content: `
resource "aws_cloudwatch_event_target" "this" {
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsCloudwatchEventTargetNoDlqRule(),
					Message: "\"dead_letter_config\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 46},
					},
				},
			},
		},
		{
			Name: "missing arn",
			Content: `
resource "aws_cloudwatch_event_target" "this" {
  dead_letter_config {}
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsCloudwatchEventTargetNoDlqRule(),
					Message: "\"arn\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 21},
					},
				},
			},
		},
		{
			Name: "correct",
			Content: `
resource "aws_cloudwatch_event_target" "this" {
  dead_letter_config {
	arn = "my-sqs-queue-arn"
  }
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsCloudwatchEventTargetNoDlqRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
