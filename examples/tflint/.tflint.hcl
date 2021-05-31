plugin "aws" {
  enabled = true
}

plugin "aws-serverless" {
  enabled = true
  source = "github.com/aws-samples/serverless-rules/tflint-ruleset-aws-serverless"
  version = "0.1.0"
}