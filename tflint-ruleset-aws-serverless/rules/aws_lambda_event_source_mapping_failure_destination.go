package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLambdaEventSourceMappingFailureDestination checks if there is an on failure destination configured on event source mappings
type AwsLambdaEventSourceMappingFailureDestinationRule struct {
	tflint.DefaultRule
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
func (r *AwsLambdaEventSourceMappingFailureDestinationRule) Severity() tflint.Severity {
	return tflint.ERROR // TODO: Determine the rule's severity
}

// Link returns the rule reference link
func (r *AwsLambdaEventSourceMappingFailureDestinationRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/lambda/eventsourcemapping_failure_destination/"
}

// Check checks if aws_lambda_event_source_mapping as a destination on_failure configured
func (r *AwsLambdaEventSourceMappingFailureDestinationRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: r.block1Name,
				Body: &hclext.BodySchema{
					Blocks: []hclext.BlockSchema{
						{
							Type: r.block2Name,
							Body: &hclext.BodySchema{
								Attributes: []hclext.AttributeSchema{
									{Name: r.attributeName},
								},
							},
						},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		// Check destination_config block
		destConfigBlocks := resource.Body.Blocks.OfType(r.block1Name)
		if len(destConfigBlocks) == 0 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.block1Name),
				resource.DefRange,
			)
			continue
		}

		// Check on_failure block
		onFailureBlocks := destConfigBlocks[0].Body.Blocks.OfType(r.block2Name)
		if len(onFailureBlocks) == 0 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.block2Name),
				destConfigBlocks[0].DefRange,
			)
			continue
		}

		// Check destination_arn attribute
		_, exists := onFailureBlocks[0].Body.Attributes[r.attributeName]
		if !exists {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				onFailureBlocks[0].DefRange,
			)
		}
	}

	return nil
}
