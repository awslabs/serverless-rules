package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsApigatewayV2StageThrottlingRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "missing default_route_settings is invalid",
			Content: `
resource "aws_apigatewayv2_stage" "missing" {
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsApigatewayV2StageThrottlingRule(),
					Message: "\"default_route_settings\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 44},
					},
				},
			},
		},
		{
			Name: "empty default_route_settings is invalid",
			Content: `
resource "aws_apigatewayv2_stage" "empty" {
	default_route_settings {}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsApigatewayV2StageThrottlingRule(),
					Message: "\"throttling_burst_limit\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 2},
						End:      hcl.Pos{Line: 3, Column: 24},
					},
				},
				{
					Rule:    NewAwsApigatewayV2StageThrottlingRule(),
					Message: "\"throttling_rate_limit\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 2},
						End:      hcl.Pos{Line: 3, Column: 24},
					},
				},
			},
		},
		{
			Name: "missing throttling_rate_limit is invalid",
			Content: `
resource "aws_apigatewayv2_stage" "missingrate" {
	default_route_settings {
		throttling_burst_limit = 1000
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsApigatewayV2StageThrottlingRule(),
					Message: "\"throttling_rate_limit\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 2},
						End:      hcl.Pos{Line: 3, Column: 24},
					},
				},
			},
		},
		{
			Name: "valid",
			Content: `
resource "aws_apigatewayv2_stage" "valid" {
	default_route_settings {
		throttling_burst_limit = 1000
		throttling_rate_limit = 100
	}
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsApigatewayV2StageThrottlingRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
