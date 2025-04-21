package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type AwsAppsyncGraphqlAPITracingRule struct {
	tflint.DefaultRule
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
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.attributeName},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				resource.DefRange,
			)
			continue
		}

		var xrayTracingEnabled string
		err := runner.EvaluateExpr(attribute.Expr, &xrayTracingEnabled, nil)
		if err != nil {
			return err
		}

		if xrayTracingEnabled != "true" {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" should be set to true.", r.attributeName),
				attribute.Expr.Range(),
			)
		}
	}

	return nil
}
