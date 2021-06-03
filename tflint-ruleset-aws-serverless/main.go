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
			Version: "0.1.5",
			Rules: []tflint.Rule{
				rules.NewAwsAPIGatewayMethodSettingsThrottlingRule(),
				rules.NewAwsAPIGatewayStageLoggingRule(),
				rules.NewAwsAPIGatewayStageTracingRule(),
				rules.NewAwsAPIGatewayStageV2LoggingRule(),
				rules.NewAwsApigatewayV2StageThrottlingRule(),
				rules.NewAwsAppsyncGraphqlAPITracingRule(),
				rules.NewAwsLambdaFunctionTracingRule(),
			},
		},
	})
}
