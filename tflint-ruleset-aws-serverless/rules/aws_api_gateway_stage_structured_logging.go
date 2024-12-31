package rules

import (
	"encoding/json"
	"fmt"
	"regexp"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsApigatewayStageStructuredLogging checks if API Gateway logging format is in JSON
type AwsApigatewayStageStructuredLoggingRule struct {
	resourceType  string
	blockName     string
	attributeName string
}

// NewAwsApigatewayStageStructuredLoggingRule returns new rule with default attributes
func NewAwsApigatewayStageStructuredLoggingRule() *AwsApigatewayStageStructuredLoggingRule {
	return &AwsApigatewayStageStructuredLoggingRule{
		resourceType:  "aws_api_gateway_stage",
		blockName:     "access_log_settings",
		attributeName: "format",
	}
}

// Name returns the rule name
func (r *AwsApigatewayStageStructuredLoggingRule) Name() string {
	return "aws_api_gateway_stage_structured_logging"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsApigatewayStageStructuredLoggingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsApigatewayStageStructuredLoggingRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsApigatewayStageStructuredLoggingRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/api_gateway/structured_logging/"
}

// Check checks if API Gateway logging format is in JSON
func (r *AwsApigatewayStageStructuredLoggingRule) Check(runner tflint.Runner) error {
	// Regexp to substitute all $context. variables
	re := regexp.MustCompile(`\$context\.[a-zA-Z\.]+`)

	return runner.WalkResourceBlocks(r.resourceType, r.blockName, func(block *hcl.Block) error {
		body, _, diags := block.Body.PartialContent(&hcl.BodySchema{
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

			attrValue = re.ReplaceAllLiteralString(attrValue, "4")

			// TODO: test if JSON
			var js map[string]interface{}
			if json.Unmarshal([]byte(attrValue), &js) != nil {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("\"%s\" is not valid JSON.", r.attributeName),
					attribute.Expr,
				)
			}
		}

		return nil
	})
}
