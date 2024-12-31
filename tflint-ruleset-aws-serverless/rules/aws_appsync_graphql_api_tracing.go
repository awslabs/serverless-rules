package rules

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type AwsAppsyncGraphqlAPITracingRule struct {
	resourceType  string
	attributeName string
}

func NewAwsAppsyncGraphqlAPITracingRule() *AwsAppsyncGraphqlAPITracingRule {
	return &AwsAppsyncGraphqlAPITracingRule{
		resourceType:  "aws_appsync_graphql_api",
		attributeName: "xray_enabled",
	}
}

// Name returns the rule name
func (r *AwsAppsyncGraphqlAPITracingRule) Name() string {
	return "aws_appsync_graphql_api_tracing_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsAppsyncGraphqlAPITracingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsAppsyncGraphqlAPITracingRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsAppsyncGraphqlAPITracingRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/appsync/tracing/"
}

// Check checks whether "aws_appsync_graphql_api" has tracing enabled
func (r *AwsAppsyncGraphqlAPITracingRule) Check(runner tflint.Runner) error {
	return runner.WalkResources(r.resourceType, func(resource *configs.Resource) error {
		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name:     r.attributeName,
					Required: true,
				},
			},
		})

		if diags.HasErrors() {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				body.MissingItemRange,
			)
			return nil
		}

		attribute := body.Attributes[r.attributeName]
		var xrayTracingEnabled string
		err := runner.EvaluateExpr(attribute.Expr, &xrayTracingEnabled, nil)

		return runner.EnsureNoError(err, func() error {
			if xrayTracingEnabled != "true" {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("\"%s\" should be set to true.", r.attributeName),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
