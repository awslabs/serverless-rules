plugin "aws" {
  enabled = true
}

plugin "aws-serverless" {
  enabled = true
  # Uncomment those lines if you are using tflint 0.29 or later
  # source = "github.com/awslabs/serverless-rules"
  # version = "0.1.8"
}