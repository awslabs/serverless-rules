Amazon SQS Rules
================

## No Redrive Policy

* __Level__: Error
* __cfn-lint__: ES6000
* __tflint__: aws_sqs_queue_redrive_policy

You can configure the redrive policy on an Amazon SQS queue. With a redrive policy, you can define how many times the SQS will make the messages available for consumers. After that, SQS will send it to the dead-letter queue specified in the policy.

### Implementations

=== "CDK"

    ```typescript
    import { Queue } from '@aws-cdk/aws-sqs';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // Dead letter queue
        const myDLQ = new Queue(
          scope, "MyDLQ",
        );

        const myQueue = new Queue(
          scope, "MyQueue",
          {
            // Configure the redrive policy for MyQueue
            deadLetterQueue: {
              maxReceiveCount: 4,
              queue: myDLQ,
            },
          }
        );
      }
    }
    ```

=== "CloudFormation (JSON)"

    ```json
    {
      "Resources": {
        "MyQueue": {
          "Type": "AWS::SQS::Queue",
          "Properties": {
            // Configure the redrive policy for MyQueue
            "RedrivePolicy": "{ \"deadLetterTargetArn\": \"arn:aws:sqs:us-east-2:111122224444:my-dlq\", \"maxReceiveCount\": 4 }"
          }
        }
      }
    }
    ```

=== "CloudFormation (YAML)"

    ```yaml
    Resources:
      MyQueue:
        Type: AWS::SQS::Queue
        Properties:
          # Configure the redrive policy for MyQueue
          RedrivePolicy: |
            {
              "deadLetterTargetArn": "arn:aws:sqs:us-east-2:111122224444:my-dlq",
              "maxReceiveCount": 4
            }
    ```

=== "Serverless Framework"

    ```yaml
    resources:
      Resources:
        MyQueue:
          Type: AWS::SQS::Queue
          Properties:
            # Configure the redrive policy for MyQueue
            RedrivePolicy: |
              {
                "deadLetterTargetArn": "arn:aws:sqs:us-east-2:111122224444:my-dlq",
                "maxReceiveCount": 4
              }
    ```

=== "Terraform"

    ```tf
    resource "aws_sqs_queue" "terraform_queue" {
      redrive_policy = jsonencode({
        deadLetterTargetArn = "arn:aws:sqs:us-east-2:111122224444:my-dlq"
        maxReceiveCount     = 4
      })
    }
    ```

### See also

* [Serverless Lens: Failure Management](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/failure-management.html)
* [Amazon SQS dead-letter-queues](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-dead-letter-queues.html)
* [__CloudFormation__: AWS::SQS::Queue](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sqs-queues.html)
* [__Terraform__: aws_sqs_queue](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue)