package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsAppsyncGraphqlAPITracingRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "not present",
			Content: `
resource "aws_appsync_graphql_api" "false" {}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAppsyncGraphqlAPITracingRule(),
					Message: "\"xray_enabled\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 44},
						End:      hcl.Pos{Line: 2, Column: 44},
					},
				},
			},
		},
		{
			Name: "false is invalid",
			Content: `
		resource "aws_appsync_graphql_api" "false" {
			xray_enabled = false
		}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAppsyncGraphqlAPITracingRule(),
					Message: "\"xray_enabled\" should be set to true.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 19},
						End:      hcl.Pos{Line: 3, Column: 24},
					},
				},
			},
		},
		{
			Name: "true is valid",
			Content: `
			resource "aws_appsync_graphql_api" "true" {
				xray_enabled = true
			}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsAppsyncGraphqlAPITracingRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
