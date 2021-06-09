package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsSfnStateMachineTracing checks if tracing is enabled for Step functions
type AwsSfnStateMachineTracingRule struct {
	resourceType  string
	blockName     string
	attributeName string
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
func (r *AwsSfnStateMachineTracingRule) Severity() string {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsSfnStateMachineTracingRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/step_functions/tracing/"
}

// TODO: Write the details of the inspection
// Check checks if tracing is enabled for Step functions
func (r *AwsSfnStateMachineTracingRule) Check(runner tflint.Runner) error {
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

			if attrValue != "true" {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("\"%s\" should be set to true.", r.attributeName),
					attribute.Expr,
				)
			}
		}

		return nil
	})
}
