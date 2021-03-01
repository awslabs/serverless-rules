package main

import (
	"github.com/aws-samples/serverless-rules/tflint-ruleset-aws-serverless/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "aws-serverless",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewAwsAPIGatewayStageLoggingRule(),
				rules.NewAwsAPIGatewayStageV2LoggingRule(),
				rules.NewAwsAPIGatewayStageTracingRule(),
				rules.NewAwsLambdaFunctionTracingRule(),
			},
		},
	})
}
