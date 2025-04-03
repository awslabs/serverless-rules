package rules

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsApigatewayV2StageStructuredLogging checks if API Gateway logging format is in JSON
type AwsApigatewayV2StageStructuredLoggingRule struct {
	tflint.DefaultRule
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
		blocks := resource.Body.Blocks.OfType(r.blockName)
		if len(blocks) == 0 {
			continue
		}

		attribute, ok := blocks[0].Body.Attributes[r.attributeName]
		if !ok {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				blocks[0].DefRange,
			)
			continue
		}

		var attrValue string
		err := runner.EvaluateExpr(attribute.Expr, &attrValue, nil)
		if err != nil {
			return err
		}

		attrValue = re.ReplaceAllLiteralString(attrValue, "4")

		var js map[string]interface{}
		if json.Unmarshal([]byte(attrValue), &js) != nil {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not valid JSON.", r.attributeName),
				attribute.Expr.Range(),
			)
		}
	}

	return nil
}
