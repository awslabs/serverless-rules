package rules

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsApigatewayStageStructuredLogging checks if API Gateway logging format is in JSON
type AwsApigatewayStageStructuredLoggingRule struct {
	tflint.DefaultRule
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
