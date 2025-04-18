terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.40.0"
    }

    random = {
      source = "hashicorp/random"
      version = "3.7.1"
    }
  }
}

resource "random_pet" "this" {
  length = 2
}

resource "aws_iam_role" "lambda_role" {
  name = "${random_pet.this.id}-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_lambda_function" "this" {
  function_name = "${random_pet.this.id}-function"

  runtime = "python3.9"
  handler = "main.handler"
  role    = aws_iam_role.lambda_role.arn

  filename = "function.zip"

  tracing_config {
    mode = "Active"
  }
}