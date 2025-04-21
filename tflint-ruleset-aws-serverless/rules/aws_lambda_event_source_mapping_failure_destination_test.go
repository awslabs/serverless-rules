package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLambdaEventSourceMappingFailureDestination(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "missing destination_config",
			Content: `
resource "aws_lambda_event_source_mapping" "this" {
  event_source_arn = "my-event-source-arn"
  function_name = "my-lambda-function"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaEventSourceMappingFailureDestinationRule(),
					Message: "\"destination_config\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 50},
					},
				},
			},
		},
		{
			Name: "missing on_failure",
			Content: `
resource "aws_lambda_event_source_mapping" "this" {
  event_source_arn = "my-event-source-arn"
  function_name = "my-lambda-function"

  destination_config {}
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaEventSourceMappingFailureDestinationRule(),
					Message: "\"on_failure\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 21},
					},
				},
			},
		},
		{
			Name: "missing destination_arn",
			Content: `
resource "aws_lambda_event_source_mapping" "this" {
  event_source_arn = "my-event-source-arn"
  function_name = "my-lambda-function"

  destination_config {
	on_failure {}
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaEventSourceMappingFailureDestinationRule(),
					Message: "\"destination_arn\" is not present.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 2},
						End:      hcl.Pos{Line: 7, Column: 12},
					},
				},
			},
		},
		{
			Name: "valid",
			Content: `
resource "aws_lambda_event_source_mapping" "this" {
  event_source_arn = "my-event-source-arn"
  function_name = "my-lambda-function"

  destination_config {
	on_failure {
	  destination_arn = "my-destination-arn"
	}
  }
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLambdaEventSourceMappingFailureDestinationRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
