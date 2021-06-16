package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsIamRoleLambdaNoStar(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "inline has star",
			Content: `
resource "aws_iam_role" "this" {
	name = "my-function-role"
	assume_role_policy = <<EOF
{
	"Version": "2012-10-17",
	"Statement": [{
		"Action": "sts:AssumeRole",
		"Effect": "Allow",
		"Sid": "",
		"Principal": {
			"Service": "lambda.amazonaws.com"
		}
	}]
}
EOF
	
	inline_policy {
		name = "FunctionPolicy"
		policy = <<EOF
{
	"Version": "2012-10-17",
	"Statement": [{
		"Effect": "Allow",
		"Action": ["dynamodb:*"],
		"Resource": "*"
	}]
}
EOF
	}
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsIamRoleLambdaNoStarRule(),
					Message: "Inline policy for role with Lambda as principal has policy actions with stars.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 20, Column: 12},
						End:      hcl.Pos{Line: 29, Column: 4},
					},
				},
			},
		},
		{
			Name: "inline valid",
			Content: `
resource "aws_iam_role" "this" {
	name = "my-function-role"
	assume_role_policy = <<EOF
{
	"Version": "2012-10-17",
	"Statement": [{
		"Action": "sts:AssumeRole",
		"Effect": "Allow",
		"Sid": "",
		"Principal": {
			"Service": "lambda.amazonaws.com"
		}
	}]
}
EOF
	
	inline_policy {
		name = "FunctionPolicy"
		policy = <<EOF
{
	"Version": "2012-10-17",
	"Statement": [{
		"Effect": "Allow",
		"Action": ["dynamodb:Query"],
		"Resource": "*"
	}]
}
EOF
	}
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsIamRoleLambdaNoStarRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
