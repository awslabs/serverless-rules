package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsSqsQueueRedrivePolicyRule checks if an SQS Queue has a redrive policy configured
type AwsSqsQueueRedrivePolicyRule struct {
	resourceType  string
	attributeName string
	tflint.DefaultRule
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
	return false
}

// Severity returns the rule severity
func (r *AwsSqsQueueRedrivePolicyRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsSqsQueueRedrivePolicyRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/sqs/redrive_policy/"
}

// Check checks if an SQS Queue has a redrive policy configured
func (r *AwsSqsQueueRedrivePolicyRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.attributeName},
		},
	}, nil)
	if err != nil {
		return fmt.Errorf("error getting resource content: %w", err)
	}

	logger.Debug(fmt.Sprintf("Found %d aws_sqs_queue resources", len(resources.Blocks)))

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

		var redrivePolicy string
		err := runner.EvaluateExpr(attr.Expr, &redrivePolicy, nil)
		if err != nil {
			return fmt.Errorf("failed to evaluate redrive policy: %w", err)
		}

		// Validate that the redrive policy contains required fields
		if redrivePolicy == "" {
			runner.EmitIssue(
				r,
				"redrive_policy is empty",
				attr.Expr.Range(),
			)
		}
	}
	return nil
}

//func (r *AwsSqsQueueRedrivePolicyRule) Check(runner tflint.Runner) error {
//	return runner.WalkResources(r.resourceType, func(resource *configs.Resource) error {
//		// Attribute
//		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
//			Attributes: []hcl.AttributeSchema{
//				{
//					Name: r.attributeName,
//				},
//			},
//		})
//
//		if diags.HasErrors() {
//			return diags
//		}
//
//		var attrValue string
//		attribute, ok := body.Attributes[r.attributeName]
//		if !ok {
//			runner.EmitIssue(
//				r,
//				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
//				body.MissingItemRange,
//			)
//		} else {
//			err := runner.EvaluateExpr(attribute.Expr, &attrValue, nil)
//			if err != nil {
//				return err
//			}
//		}
//
//		return nil
//	})
//}
