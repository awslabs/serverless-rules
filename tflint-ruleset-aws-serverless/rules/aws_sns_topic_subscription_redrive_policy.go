package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TODO: Write the rule's description here
// AwsSnsTopicSubscriptionRedrivePolicy checks that an SNS subscription has a redrive policy configured
type AwsSnsTopicSubscriptionRedrivePolicyRule struct {
	resourceType  string
	attributeName string
}

// NewAwsSnsTopicSubscriptionRedrivePolicyRule returns new rule with default attributes
func NewAwsSnsTopicSubscriptionRedrivePolicyRule() *AwsSnsTopicSubscriptionRedrivePolicyRule {
	return &AwsSnsTopicSubscriptionRedrivePolicyRule{
		// TODO: Write resource type and attribute name here
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
func (r *AwsSnsTopicSubscriptionRedrivePolicyRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsSnsTopicSubscriptionRedrivePolicyRule) Link() string {
	return ""
}

// Check checks that an SNS subscription has a redrive policy configured
func (r *AwsSnsTopicSubscriptionRedrivePolicyRule) Check(runner tflint.Runner) error {
	return runner.WalkResources(r.resourceType, func(resource *configs.Resource) error {
		// Attribute
		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name: r.attributeName,
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		var attrValue string
		attribute, ok := body.Attributes[r.attributeName]
		if !ok {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				body.MissingItemRange,
			)
		} else {
			err := runner.EvaluateExpr(attribute.Expr, &attrValue, nil)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
