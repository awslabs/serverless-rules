package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLambdaEventSourceMappingFailureDestination checks if there is an on failure destination configured on event source mappings
type AwsLambdaEventSourceMappingFailureDestinationRule struct {
	resourceType  string
	block1Name    string
	block2Name    string
	attributeName string
}

// NewAwsLambdaEventSourceMappingFailureDestinationRule returns new rule with default attributes
func NewAwsLambdaEventSourceMappingFailureDestinationRule() *AwsLambdaEventSourceMappingFailureDestinationRule {
	return &AwsLambdaEventSourceMappingFailureDestinationRule{
		// TODO: Write resource type and attribute name here
		resourceType:  "aws_lambda_event_source_mapping",
		block1Name:    "destination_config",
		block2Name:    "on_failure",
		attributeName: "destination_arn",
	}
}

// Name returns the rule name
func (r *AwsLambdaEventSourceMappingFailureDestinationRule) Name() string {
	return "aws_lambda_event_source_mapping_failure_destination"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsLambdaEventSourceMappingFailureDestinationRule) Enabled() bool {
	// TODO: Determine whether the rule is enabled by default
	return true
}

// Severity returns the rule severity
func (r *AwsLambdaEventSourceMappingFailureDestinationRule) Severity() string {
	// TODO: Determine the rule's severiry
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsLambdaEventSourceMappingFailureDestinationRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/lambda/eventsourcemapping_failure_destination/"
}

// Check checks if aws_lambda_event_source_mapping as a destination on_failure configured
func (r *AwsLambdaEventSourceMappingFailureDestinationRule) Check(runner tflint.Runner) error {
	return runner.WalkResources(r.resourceType, func(resource *configs.Resource) error {

		// Block 1 - destination_config

		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{
				{
					Type: r.block1Name,
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		blocks := body.Blocks.OfType(r.block1Name)
		if len(blocks) != 1 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.block1Name),
				body.MissingItemRange,
			)

			return nil
		}

		// Block 2 - on_failure

		body, _, diags = blocks[0].Body.PartialContent(&hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{
				{
					Type: r.block2Name,
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		blocks = body.Blocks.OfType(r.block2Name)
		if len(blocks) != 1 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.block2Name),
				body.MissingItemRange,
			)

			return nil
		}

		// Attribute - destination_arn

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
