package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsAPIGatewayMethodSettingsThrottlingRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "missing settings block is invalid",
			Content: `
resource "aws_api_gateway_method_settings" "missing" {
	method_path = "*/*"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAPIGatewayMethodSettingsThrottlingRule(),
					Message: "\"settings\" block is required for default method settings",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 53},
					},
				},
			},
		},
		{
			Name: "empty settings block is invalid",
			Content: `
resource "aws_api_gateway_method_settings" "empty" {
	method_path = "*/*"
	settings {}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAPIGatewayMethodSettingsThrottlingRule(),
					Message: "\"throttling_burst_limit\" is required for default method settings",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 2},
						End:      hcl.Pos{Line: 4, Column: 10},
					},
				},
				{
					Rule:    NewAwsAPIGatewayMethodSettingsThrottlingRule(),
					Message: "\"throttling_rate_limit\" is required for default method settings",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 2},
						End:      hcl.Pos{Line: 4, Column: 10},
					},
				},
			},
		},
		{
			Name: "missing throttling rate limit is invalid",
			Content: `
resource "aws_api_gateway_method_settings" "missingrate" {
	method_path = "*/*"
	settings {
		throttling_burst_limit = 1000
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAPIGatewayMethodSettingsThrottlingRule(),
					Message: "\"throttling_rate_limit\" is required for default method settings",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 2},
						End:      hcl.Pos{Line: 4, Column: 10},
					},
				},
			},
		},
		{
			Name: "valid",
			Content: `
resource "aws_api_gateway_method_settings" "valid" {
	method_path = "*/*"
	settings {
		throttling_burst_limit = 1000
		throttling_rate_limit = 100
	}
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsAPIGatewayMethodSettingsThrottlingRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
