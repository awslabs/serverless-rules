package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsSnsTopicSubscriptionRedrivePolicyRule checks that an SNS subscription has a redrive policy configured
type AwsSnsTopicSubscriptionRedrivePolicyRule struct {
	resourceType  string
	attributeName string
	tflint.DefaultRule
}

// NewAwsSnsTopicSubscriptionRedrivePolicyRule returns new rule with default attributes
func NewAwsSnsTopicSubscriptionRedrivePolicyRule() *AwsSnsTopicSubscriptionRedrivePolicyRule {
	return &AwsSnsTopicSubscriptionRedrivePolicyRule{
		resourceType:  "aws_sns_topic_subscription",
		attributeName: "redrive_policy",
	}
}

// Name returns the rule name
func (r *AwsSnsTopicSubscriptionRedrivePolicyRule) Name() string {
	return "aws_sns_topic_subscription_redrive_policy"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSnsTopicSubscriptionRedrivePolicyRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsSnsTopicSubscriptionRedrivePolicyRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsSnsTopicSubscriptionRedrivePolicyRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/sns/redrive_policy/"
}

// Check checks that an SNS subscription has a redrive policy configured
func (r *AwsSnsTopicSubscriptionRedrivePolicyRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.attributeName},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attr, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				resource.DefRange,
			)
			continue
		}

		var attrValue string
		if err := runner.EvaluateExpr(attr.Expr, &attrValue, nil); err != nil {
			return err
		}
	}
	return nil
}
