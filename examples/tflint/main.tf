terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.42.0"
    }

    random = {
      source = "hashicorp/random"
      version = "3.1.0"
    }
  }
}

resource "random_pet" "this" {
  length = 2
}

resource "aws_lambda_function" "this" {
  function_name = "${random_pet.this.id}-function"

  runtime = "python3.8"
  handler = "main.handler"

  filename = "function.zip"

  tracing_config {
    mode = "Active"
  }
}