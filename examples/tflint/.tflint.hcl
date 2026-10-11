plugin "aws" {
  enabled = true
  source  = "github.com/terraform-linters/tflint-ruleset-aws"
  version = "0.49.0"
}

plugin "aws-serverless" {
  enabled = true
  source  = "github.com/awslabs/serverless-rules"
  version = "0.3.6"
}
