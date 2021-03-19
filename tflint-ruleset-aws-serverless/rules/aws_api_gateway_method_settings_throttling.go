package rules

import (
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsAPIGatewayMethodSettingsThrottlingRule checks whether there is a default "aws_api_gateway_method_settings" resource with throttling values
type AwsAPIGatewayMethodSettingsThrottlingRule struct{}

func NewAwsAPIGatewayMethodSettingsThrottlingRule() *AwsAPIGatewayMethodSettingsThrottlingRule {
	return &AwsAPIGatewayMethodSettingsThrottlingRule{}
}

// Name returns the rule name
func (r *AwsAPIGatewayMethodSettingsThrottlingRule) Name() string {
	return "aws_api_gateay_method_settings_throttling_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsAPIGatewayMethodSettingsThrottlingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsAPIGatewayMethodSettingsThrottlingRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsAPIGatewayMethodSettingsThrottlingRule) Link() string {
	return ""
}

// Check checks whether default "aws_api_gateway_method_settings" have throttling values
func (r *AwsAPIGatewayMethodSettingsThrottlingRule) Check(runner tflint.Runner) error {
	return runner.WalkResources("aws_api_gateway_method_settings", func(resource *configs.Resource) error {
		// Load resource body
		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name:     "method_path",
					Required: true,
				},
			},
			Blocks: []hcl.BlockHeaderSchema{
				{
					Type: "settings",
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		// Only looking at the default method settings
		attribute := body.Attributes["method_path"]
		var methodPath string
		err := runner.EvaluateExpr(attribute.Expr, &methodPath, nil)

		if err != nil {
			return err
		}

		if methodPath != "*/*" {
			return nil
		}

		// Load 'settings' block
		blocks := body.Blocks.OfType("settings")
		if len(blocks) != 1 {
			runner.EmitIssue(
				r,
				"\"settings\" is not present.",
				body.MissingItemRange,
			)

			return nil
		}

		blockBody, _, diags := blocks[0].Body.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name: "throttling_burst_limit",
				},
				{
					Name: "throttling_rate_limit",
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		// Check throttling limits
		var throttlingBurstLimit int
		throttlingBurstLimitAttribute, burstOk := blockBody.Attributes["throttling_burst_limit"]
		if !burstOk {
			runner.EmitIssue(
				r,
				"\"throttling_burst_limit\" is not present.",
				blockBody.MissingItemRange,
			)
		} else {
			err = runner.EvaluateExpr(throttlingBurstLimitAttribute.Expr, &throttlingBurstLimit, nil)
			if err != nil {
				return err
			}
		}

		var throttlingRateLimit int
		throttlingRateLimitAttribute, rateOk := blockBody.Attributes["throttling_rate_limit"]
		if !rateOk {
			runner.EmitIssue(
				r,
				"\"throttling_rate_limit\" is not present.",
				blockBody.MissingItemRange,
			)
		} else {
			err = runner.EvaluateExpr(throttlingRateLimitAttribute.Expr, &throttlingRateLimit, nil)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
