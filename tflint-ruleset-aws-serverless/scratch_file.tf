provider "aws" {
  region = "us-west-2"
}

resource "aws_sqs_queue" "test_queue" {
  name = "test-queue"

  # Intentionally omitting redrive_policy to test the rule
  visibility_timeout_seconds = 30
  message_retention_seconds = 345600  # 4 days
}

# Example with redrive_policy (for testing both cases)
resource "aws_sqs_queue" "dlq" {
  name = "test-dlq"
}

resource "aws_sqs_queue" "queue_with_dlq" {
  name = "test-queue-with-dlq"
  
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq.arn
    maxReceiveCount     = 4
  })
}
