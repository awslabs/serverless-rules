package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsApigatewayV2StageThrottlingRule checks whether "aws_apigatewayv2_stage" has default throttling values.
type AwsApigatewayV2StageThrottlingRule struct {
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
func (r *AwsApigatewayV2StageThrottlingRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsApigatewayV2StageThrottlingRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/api_gateway/default_throttling/"
}

// Check checks whether "aws_apigatewayv2_stage" has has default throttling values
func (r *AwsApigatewayV2StageThrottlingRule) Check(runner tflint.Runner) error {
	return runner.WalkResources(r.resourceType, func(resource *configs.Resource) error {
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

		blockBody, _, diags := blocks[0].Body.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name: r.burstAttributeName,
				},
				{
					Name: r.rateAttributeName,
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		// Check throttling limits
		var burstLimit int
		burstLimitAttribute, burstOk := blockBody.Attributes[r.burstAttributeName]
		if !burstOk {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.burstAttributeName),
				blockBody.MissingItemRange,
			)
		} else {
			err := runner.EvaluateExpr(burstLimitAttribute.Expr, &burstLimit, nil)
			if err != nil {
				return err
			}
		}

		var rateLimit int
		rateLimitAttribute, rateOk := blockBody.Attributes[r.rateAttributeName]
		if !rateOk {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.rateAttributeName),
				blockBody.MissingItemRange,
			)
		} else {
			err := runner.EvaluateExpr(rateLimitAttribute.Expr, &rateLimit, nil)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
