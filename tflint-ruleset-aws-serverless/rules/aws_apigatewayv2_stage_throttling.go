package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsApigatewayV2StageThrottlingRule checks whether "aws_apigatewayv2_stage" has default throttling values.
type AwsApigatewayV2StageThrottlingRule struct {
	tflint.DefaultRule
	resourceType       string
	blockName          string
	burstAttributeName string
	rateAttributeName  string
}

// NewAwsApigatewayV2StageThrottlingRule returns new rule
func NewAwsApigatewayV2StageThrottlingRule() *AwsApigatewayV2StageThrottlingRule {
	return &AwsApigatewayV2StageThrottlingRule{
		resourceType:       "aws_apigatewayv2_stage",
		blockName:          "default_route_settings",
		burstAttributeName: "throttling_burst_limit",
		rateAttributeName:  "throttling_rate_limit",
	}
}

// Name returns the rule name
func (r *AwsApigatewayV2StageThrottlingRule) Name() string {
	return "aws_apigatewayv2_stage_throttling_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsApigatewayV2StageThrottlingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsApigatewayV2StageThrottlingRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsApigatewayV2StageThrottlingRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/api_gateway/default_throttling/"
}

// Check checks whether "aws_apigatewayv2_stage" has has default throttling values
func (r *AwsApigatewayV2StageThrottlingRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: r.blockName,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: r.burstAttributeName},
						{Name: r.rateAttributeName},
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

		// Check throttling limits
		_, burstOk := blocks[0].Body.Attributes[r.burstAttributeName]
		if !burstOk {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.burstAttributeName),
				blocks[0].DefRange,
			)
		}

		_, rateOk := blocks[0].Body.Attributes[r.rateAttributeName]
		if !rateOk {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.rateAttributeName),
				blocks[0].DefRange,
			)
		}
	}

	return nil
}
