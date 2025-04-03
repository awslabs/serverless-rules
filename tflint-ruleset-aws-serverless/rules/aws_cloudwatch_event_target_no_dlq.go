package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsCloudwatchEventTargetNoDlq checks if there is a DLQ configured on EventBridge targets
type AwsCloudwatchEventTargetNoDlqRule struct {
	tflint.DefaultRule
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
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: r.blockName,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: r.attributeName},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		// Check for block
		blocks := resource.Body.Blocks.OfType(r.blockName)
		if len(blocks) == 0 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.blockName),
				resource.DefRange,
			)
			continue
		}

		// Check for attribute
		_, exists := blocks[0].Body.Attributes[r.attributeName]
		if !exists {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				blocks[0].DefRange,
			)
		}
	}

	return nil
}
