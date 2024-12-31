package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsSfnStateMachineTracing checks if tracing is enabled for Step functions
type AwsSfnStateMachineTracingRule struct {
	resourceType  string
	blockName     string
	attributeName string
	tflint.DefaultRule
}

// NewAwsSfnStateMachineTracingRule returns new rule with default attributes
func NewAwsSfnStateMachineTracingRule() *AwsSfnStateMachineTracingRule {
	return &AwsSfnStateMachineTracingRule{
		resourceType:  "aws_sfn_state_machine",
		blockName:     "tracing_configuration",
		attributeName: "enabled",
	}
}

// Name returns the rule name
func (r *AwsSfnStateMachineTracingRule) Name() string {
	return "aws_sfn_state_machine_tracing"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSfnStateMachineTracingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsSfnStateMachineTracingRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsSfnStateMachineTracingRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/step_functions/tracing/"
}

// TODO: Write the details of the inspection
// Check checks if tracing is enabled for Step functions
func (r *AwsSfnStateMachineTracingRule) Check(runner tflint.Runner) error {
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
		blocks := resource.Body.Blocks
		if len(blocks) == 0 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.blockName),
				resource.DefRange,
			)
			continue
		}

		for _, block := range blocks {
			var attrValue string
			attr, exists := block.Body.Attributes[r.attributeName]
			if !exists {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is not present.", r.attributeName),
					block.DefRange,
				)
				continue
			}

			if err := runner.EvaluateExpr(attr.Expr, &attrValue, nil); err != nil {
				return err
			}

			if attrValue != "true" {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" should be set to true.", r.attributeName),
					attr.Expr.Range(),
				)
			}
		}
	}
	return nil
}
