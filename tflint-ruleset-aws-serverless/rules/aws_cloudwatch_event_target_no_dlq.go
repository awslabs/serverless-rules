package rules

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsCloudwatchEventTargetNoDlq checks if there is a DLQ configured on EventBridge targets
type AwsCloudwatchEventTargetNoDlqRule struct {
	resourceType  string
	blockName     string
	attributeName string
}

// NewAwsCloudwatchEventTargetNoDlqRule returns new rule with default attributes
func NewAwsCloudwatchEventTargetNoDlqRule() *AwsCloudwatchEventTargetNoDlqRule {
	return &AwsCloudwatchEventTargetNoDlqRule{
		resourceType:  "aws_cloudwatch_event_target",
		blockName:     "dead_letter_config",
		attributeName: "arn",
	}
}

// Name returns the rule name
func (r *AwsCloudwatchEventTargetNoDlqRule) Name() string {
	return "aws_cloudwatch_event_target_no_dlq"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsCloudwatchEventTargetNoDlqRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsCloudwatchEventTargetNoDlqRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsCloudwatchEventTargetNoDlqRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/eventbridge/rule_without_dlq/"
}

// Check checks if there is a DLQ configured on EventBridge targets
func (r *AwsCloudwatchEventTargetNoDlqRule) Check(runner tflint.Runner) error {
	return runner.WalkResources(r.resourceType, func(resource *configs.Resource) error {
		// Block

		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{
				{
					Type: r.blockName,
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		blocks := body.Blocks.OfType(r.blockName)
		if len(blocks) != 1 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.blockName),
				body.MissingItemRange,
			)

			return nil
		}

		// Attribute
		body, _, diags = blocks[0].Body.PartialContent(&hcl.BodySchema{
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
