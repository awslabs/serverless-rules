package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsAPIGatewayRestStageTracingRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "not present",
			Content: `
resource "aws_api_gateway_stage" "false" {}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAPIGatewayRestStageTracingRule(),
					Message: "\"xray_tracing_enabled\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 42},
						End:      hcl.Pos{Line: 2, Column: 42},
					},
				},
			},
		},
		{
			Name: "false is invalid",
			Content: `
		resource "aws_api_gateway_stage" "false" {
			xray_tracing_enabled = false
		}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAPIGatewayRestStageTracingRule(),
					Message: "\"xray_tracing_enabled\" should be set to true.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 27},
						End:      hcl.Pos{Line: 3, Column: 32},
					},
				},
			},
		},
		{
			Name: "true is valid",
			Content: `
			resource "aws_api_gateway_stage" "true" {
				xray_tracing_enabled = true
			}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsAPIGatewayRestStageTracingRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
