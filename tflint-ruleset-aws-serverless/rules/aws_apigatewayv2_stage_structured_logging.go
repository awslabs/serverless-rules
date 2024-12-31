package rules

import (
	"encoding/json"
	"fmt"
	"regexp"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsApigatewayV2StageStructuredLogging checks if API Gateway logging format is in JSON
type AwsApigatewayV2StageStructuredLoggingRule struct {
	resourceType  string
	blockName     string
	attributeName string
}

// NewAwsApigatewayV2StageStructuredLoggingRule returns new rule with default attributes
func NewAwsApigatewayV2StageStructuredLoggingRule() *AwsApigatewayV2StageStructuredLoggingRule {
	return &AwsApigatewayV2StageStructuredLoggingRule{
		resourceType:  "aws_apigatewayv2_stage",
		blockName:     "access_log_settings",
		attributeName: "format",
	}
}

// Name returns the rule name
func (r *AwsApigatewayV2StageStructuredLoggingRule) Name() string {
	return "aws_apigatewayv2_stage_structured_logging"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsApigatewayV2StageStructuredLoggingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsApigatewayV2StageStructuredLoggingRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsApigatewayV2StageStructuredLoggingRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/api_gateway/structured_logging/"
}

// Check checks if API Gateway logging format is in JSON
func (r *AwsApigatewayV2StageStructuredLoggingRule) Check(runner tflint.Runner) error {
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
