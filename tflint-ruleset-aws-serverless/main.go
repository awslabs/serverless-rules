package main

import (
	"github.com/awslabs/serverless-rules/tflint-ruleset-aws-serverless/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "aws-serverless",
			Version: "0.3.5",
			Rules:   rules.Rules,
		},
	})
}
