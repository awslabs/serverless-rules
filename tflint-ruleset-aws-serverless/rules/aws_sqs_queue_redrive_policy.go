package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsSqsQueueRedrivePolicy checks if an SQS Queue has a redrive policy configured
type AwsSqsQueueRedrivePolicyRule struct {
	resourceType  string
	attributeName string
}

// NewAwsSqsQueueRedrivePolicyRule returns new rule with default attributes
func NewAwsSqsQueueRedrivePolicyRule() *AwsSqsQueueRedrivePolicyRule {
	return &AwsSqsQueueRedrivePolicyRule{
		resourceType:  "aws_sqs_queue",
		attributeName: "redrive_policy",
	}
}

// Name returns the rule name
func (r *AwsSqsQueueRedrivePolicyRule) Name() string {
	return "aws_sqs_queue_redrive_policy"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSqsQueueRedrivePolicyRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsSqsQueueRedrivePolicyRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsSqsQueueRedrivePolicyRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/sqs/redrive_policy/"
}

// Check checks if an SQS Queue has a redrive policy configured
func (r *AwsSqsQueueRedrivePolicyRule) Check(runner tflint.Runner) error {
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
