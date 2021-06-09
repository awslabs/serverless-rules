Amazon SNS Rules
================

## No Redrive Policy

* __Level__: Error
* __cfn-lint__: ES7000
* __tflint__: aws_sns_topic_subscription_redrive_policy

You can configure the redrive policy on an Amazon SNS subscription. If SNS cannot deliver the message after the number of attempts set in its delivery policy, SNS will send it to the dead-letter queue specified in the redrive policy.

### Implementations

=== "CDK"

    ```typescript
    import { Queue } from '@aws-cdk/aws-sqs';
    import { Topic } from '@aws-cdk/aws-sns';
    import { UrlSubscription } from '@aws-cdk-aws-sns-subscriptions';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // Dead letter queue
        const myDLQ = new Queue(
          scope, 'MyDLQ',
        );

        // SNS Topic
        const myTopic = new Topic(scope, 'MyTopic');

        // Adding an URL subscription
        myTopic.addSubscription(new UrlSubscription(
          'https://example.com',
          {
            // Configure the redrive policy for the subscription
            deadLetterQueue: myDLQ,
          }
        ));
      }
    }
    ```

=== "CloudFormation (JSON)"

    ```json
    {
      "Resources": {
        "MySubscription": {
        "Type": "AWS::SNS::Subscription",
        "Properties": {
          "Protocol": "https",
          "Endpoint": "https://example.com/"
          "TopicArn": "my-topic-arn",

          // Configure the redrive policy for the subscription
          "RedrivePolicy": "{ \"deadLetterTargetArn\": \"arn:aws:sqs:us-east-2:123456789012:MyDeadLetterQueue\"}"
          }
        }
      }
    }
    ```

=== "CloudFormation (YAML)"

    ```yaml
    Resources:
      MySubscription:
        Type: AWS::SNS::Subscription
        Properties:
          Protocol: https
          Endpoint: https://example.com/
          TopicArn: "my-topic-arn"

          # Configure the redrive policy for the subscription
          RedrivePolicy: |
            {
              "deadLetterTargetArn": "arn:aws:sqs:us-east-2:123456789012:MyDeadLetterQueue"
            }
    ```

=== "Serverless Framework"

    ```yaml
    # For subscriptions to Lambda function endpoints
    functions:
      MyFunction:
        handler: hello.handler
        events:
          - sns:
              topicName: my-topic
              # Configure the redrive policy for the subscription to the Lambda function
              redrivePolicy:
                deadLetterTargetArn: arn:aws:sqs:us-east-2:123456789012:MyDeadLetterQueue

    # For subscriptions to other types of endpoint
    resources:
      Resources:
        MySubscription:
          Type: AWS::SNS::Subscription
          Properties:
            Protocol: https
            Endpoint: https://example.com/
            TopicArn: "my-topic-arn"

            # Configure the redrive policy for the subscription to another type of resource
            RedrivePolicy: |
              {
                "deadLetterTargetArn": "arn:aws:sqs:us-east-2:123456789012:MyDeadLetterQueue"
              }
    ```

=== "Terraform"

    ```tf
    resource "aws_sns_topic_subscription" "this" {
      endpoint = "https://example.com/"
      protocol = "https"
      topic_arn = "my-topic-arn"

      # Configure the redrive policy for the subscription
      redrive_policy = <<EOF
    {
      "deadLetterTargetArn": "arn:aws:sqs:us-east-2:123456789012:MyDeadLetterQueue"
    }
    EOF
    }
    ```

### See also

* [Amazon SNS message delivery retries](https://docs.aws.amazon.com/sns/latest/dg/sns-message-delivery-retries.html)
* [Amazon SNS dead-letter queues (DLQs)](https://docs.aws.amazon.com/sns/latest/dg/sns-dead-letter-queues.html)
* [__CloudFormation__: AWS::SNS::Subscription](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-sns-subscription.html)
* [__Terraform__: aws_sns_topic_subscription](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sns_topic_subscription)