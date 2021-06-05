plugin "aws" {
  enabled = true
}

plugin "aws-serverless" {
  enabled = true
  source = "github.com/awslabs/serverless-rules"
  version = "0.1.6"
}