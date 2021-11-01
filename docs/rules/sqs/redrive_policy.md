# SQS Redrive Policy

__Level__: Error
{: class="badge badge-red" }

__Initial version__: 0.1.7
{: class="badge badge-blue" }

__cfn-lint__: ES6000
{: class="badge" }

__tflint__: aws_sqs_queue_redrive_policy
{: class="badge" }

You can configure the redrive policy on an Amazon SQS queue. With a redrive policy, you can define how many times SQS will make the messages available for consumers. After that, SQS will send it to the dead-letter queue specified in the policy.

??? bug "Disabled for Terraform"

    This rule is disabled for Terraform, as the current linter only support static values in expressions. See [this issue](https://github.com/awslabs/serverless-rules/issues/107) for more information.

## Implementations

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
    resource "aws_sqs_queue" "this" {
      # Configure the redrive policy for the queue
      redrive_policy = jsonencode({
        deadLetterTargetArn = "arn:aws:sqs:us-east-2:111122224444:my-dlq"
        maxReceiveCount     = 4
      })
    }
    ```

## See also

* [Serverless Lens: Failure Management](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/failure-management.html)
* [Amazon SQS dead-letter-queues](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-dead-letter-queues.html)
* [__CloudFormation__: AWS::SQS::Queue](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sqs-queues.html)
* [__Terraform__: aws_sqs_queue](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue)