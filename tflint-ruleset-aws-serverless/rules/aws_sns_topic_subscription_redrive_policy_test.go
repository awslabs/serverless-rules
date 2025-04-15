package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSnsTopicSubscriptionRedrivePolicy(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "missing redrive_policy",
			Content: `
resource "aws_sns_topic_subscription" "this" {
  endpoint = "https://example.com/"
  protocol = "https"
  topic_arn = "my-topic-arn"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSnsTopicSubscriptionRedrivePolicyRule(),
					Message: "\"redrive_policy\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 45},
					},
				},
			},
		},
		{
			Name: "valid",
			Content: `
resource "aws_sns_topic_subscription" "this" {
  endpoint = "https://example.com/"
  protocol = "https"
  topic_arn = "my-topic-arn"

  redrive_policy = <<EOF
{
"deadLetterTargetArn": "arn:aws:sqs:us-east-2:123456789012:MyDeadLetterQueue"
}
EOF
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsSnsTopicSubscriptionRedrivePolicyRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
