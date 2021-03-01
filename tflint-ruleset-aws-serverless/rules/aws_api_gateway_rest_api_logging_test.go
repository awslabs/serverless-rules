package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsAPIGatewayRestStageLoggingRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "missing access_log_settings is invalid",
			Content: `
resource "aws_api_gateway_stage" "missing" {
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAPIGatewayRestStageLoggingRule(),
					Message: "\"access_log_settings\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 44},
						End:      hcl.Pos{Line: 2, Column: 44},
					},
				},
			},
		},
		{
			Name: "valid",
			Content: `
resource "aws_api_gateway_stage" "valid" {
	access_log_settings {
		destination_arn = "ARN"
		format = "FORMAT"
	}
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsAPIGatewayRestStageLoggingRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
