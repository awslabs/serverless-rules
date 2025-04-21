package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsApigatewayStageStructuredLogging(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "missing format",
			Content: `
resource "aws_api_gateway_stage" "valid" {
	access_log_settings {
		destination_arn = "ARN"
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsApigatewayStageStructuredLoggingRule(),
					Message: "\"format\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 2},
						End:      hcl.Pos{Line: 3, Column: 21},
					},
				},
			},
		},
		{
			Name: "non-json",
			Content: `
resource "aws_api_gateway_stage" "valid" {
	access_log_settings {
		destination_arn = "ARN"
		format = "FORMAT"
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsApigatewayStageStructuredLoggingRule(),
					Message: "\"format\" is not valid JSON.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 12},
						End:      hcl.Pos{Line: 5, Column: 20},
					},
				},
			},
		},
		{
			Name: "malformed json",
			Content: `
resource "aws_api_gateway_stage" "valid" {
	access_log_settings {
		destination_arn = "ARN"
		format = <<EOF
{
	"stage" : "$context.stage",
	"request_id" : "$context.requestId",
	"api_id" : "$context.apiId",
	"resource_path" : "$context.resourcePath",
	"resource_id" : "$context.resourceId",
	"http_method" : "$context.httpMethod",
	"source_ip" : "$context.identity.sourceIp",
	"user-agent" : "$context.identity.userAgent",
	"account_id" : "$context.identity.accountId",
	"api_key" : "$context.identity.apiKey",
	"caller" : "$context.identity.caller",
	"user" : "$context.identity.user",
	"user_arn" : "$context.identity.userArn",
	"integration_latency": $context.integration.latency
EOF
	}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsApigatewayStageStructuredLoggingRule(),
					Message: "\"format\" is not valid JSON.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 12},
						End:      hcl.Pos{Line: 21, Column: 4},
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
		format = <<EOF
{
	"stage" : "$context.stage",
	"request_id" : "$context.requestId",
	"api_id" : "$context.apiId",
	"resource_path" : "$context.resourcePath",
	"resource_id" : "$context.resourceId",
	"http_method" : "$context.httpMethod",
	"source_ip" : "$context.identity.sourceIp",
	"user-agent" : "$context.identity.userAgent",
	"account_id" : "$context.identity.accountId",
	"api_key" : "$context.identity.apiKey",
	"caller" : "$context.identity.caller",
	"user" : "$context.identity.user",
	"user_arn" : "$context.identity.userArn",
	"integration_latency": $context.integration.latency
}
EOF
	}
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsApigatewayStageStructuredLoggingRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
