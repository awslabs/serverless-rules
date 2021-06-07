package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSqsQueueRedrivePolicy(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "missing redrive_policy",
			Content: `
resource "aws_sqs_queue" "this" {
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSqsQueueRedrivePolicyRule(),
					Message: "\"redrive_policy\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 33},
						End:      hcl.Pos{Line: 2, Column: 33},
					},
				},
			},
		},
		{
			Name: "valid",
			Content: `
resource "aws_sqs_queue" "this" {
  redrive_policy = <<EOF
{
  "deadLetterTargetArn": "my-sqs-arn",
  "maxReceiveCount": 4
}
EOF
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsSqsQueueRedrivePolicyRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
