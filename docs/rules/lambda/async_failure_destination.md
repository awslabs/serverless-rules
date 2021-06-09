# Lambda Async Failure Destination

__Level__: Error
{: class="badge badge-red" }

__Initial version__: 0.1.8
{: class="badge badge-blue" }

__cfn-lint__: ES1007
{: class="badge" }

__tflint__: _Not implemented_
{: class="badge" }

Several AWS services, such as Amazon S3, Amazon SNS, or Amazon EventBridge, invoke Lambda functions asynchronously to process events. When you invoke a function asynchronously, you don't wait for a response from the function code. You hand off the event to Lambda and Lambda handles the rest.

When an asynchronous calls fail, they should be captured and retried whenever possible. For this purpose, you can set a destination where Lambda will send events for successful or failed invocations.

??? warning "Matching function name between resources"

    This rule works by comparing _Lambda Permission_ resources with _Lambda Event Invoke Config_ resources. For this rule to work correctly, you must set the function name on both resources in the exact same way.

    For example, in CloudFormation, if you use the `Fn::Ref` intrinsic function to refer to your Lambda function on both resources, this rule will work normally. If you use `Fn::Ref` on one, and `Fn::Join` on another, this rule will not work.

    Here are some examples of valid implementation in CloudFormation:

    === "With Fn::Ref"

        ```yaml
        Resources:
          Permission:
            Type: AWS::Lambda::Permission
            Properties:
              # Other properties omitted
              FunctionName: !Ref MyFunction

          EventInvokeConfig:
            Type: AWS::Lambda::EventInvokeConfig
            Properties:
              # Other properties omitted
              FunctionName: !Ref MyFunction
        ```

    === "With Fn::Sub"

        ```yaml
        Resources:
          Permission:
            Type: AWS::Lambda::Permission
            Properties:
              # Other properties omitted
              FunctionName: !Sub "arn:${AWS::Partition}:lambda:${AWS::Region}:${AWS::AccountId}:function:${MyFunction}"

          EventInvokeConfig:
            Type: AWS::Lambda::EventInvokeConfig
            Properties:
              # Other properties omitted
              FunctionName: !Sub "arn:${AWS::Partition}:lambda:${AWS::Region}:${AWS::AccountId}:function:${MyFunction}"
        ```

    === "With a static value"

        ```yaml
        Resources:
          Permission:
            Type: AWS::Lambda::Permission
            Properties:
              # Other properties omitted
              FunctionName: my-lambda-function

          EventInvokeConfig:
            Type: AWS::Lambda::EventInvokeConfig
            Properties:
              # Other properties omitted
              FunctionName: my-lambda-function
        ```

    By comparison, this implementation will return an error:

    === "With mixed references"

        ```yaml
        Resources:
          Permission:
            Type: AWS::Lambda::Permission
            Properties:
              # Other properties omitted
              FunctionName: !Ref MyFunction

          EventInvokeConfig:
            Type: AWS::Lambda::EventInvokeConfig
            Properties:
              # Other properties omitted
              FunctionName: my-lambda-function
        ```

## Implementations

=== "CDK"

    ```typescript
    import { Code, Function, Runtime } from '@aws-cdk/aws-lambda';
    import { SnsEventSource } from '@aws-cdk/aws-lambda-event-sources';
    import { SqsDestination } from '@aws-cdk/aws-lambda-destinations';
    import { Topic } from '@aws-cdk/aws-sns';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const myTopic = new Topic(
          scope, 'MyTopic',
        );

        const myDLQ = new Queue(
          scope, 'MyDLQ',
        );

        const myFunction = new Function(
          scope, 'MyFunction',
          {
            code: Code.fromAsset('src/hello/'),
            handler: 'main.handler',
            runtime: Runtime.PYTHON_3_8,

            onFailure: new SqsDestination(myDLQ),
          }
        );

        // SNS will trigger the function asynchronously
        myFunction.addEventSource(new SnsEventSource(myTopic));


      }
    }
    ```

=== "CloudFormation (JSON)"

    ```json
    {
      "Resources": {
        "SNSFunction": {
          "Type": "AWS::Serverless::Function",
          "Properties": {
            "CodeUri": ".",
            // SNS will trigger the function asynchronously
            "Events": {
              "SNS": {
                "Type": "SNS",
                "Properties": {
                  "Topic": "my-sns-topic"
                }
              }
            },
            // Configure a failure destination for the function
            "EventInvokeConfig": {
              "DestinationConfig": {
                "OnFailure": {
                  "Type": "SQS",
                  "Destination": "arn:aws:sqs:us-east-1:111122223333:my-dlq"
                }
              }
            }
          }
        }
      }
    }
    ```

=== "CloudFormation (YAML)"

    ```yaml
    SNSFunction:
      Type: AWS::Serverless::Function
      Properties:
        CodeUri: .
        # SNS will trigger the function asynchronously
        Events:
          SNS:
            Type: SNS
            Properties:
              Topic: my-sns-topic
        # Configure a failure destination for the function
        EventInvokeConfig:
          DestinationConfig:
            OnFailure:
              Type: SQS
              Destination: arn:aws:sqs:us-east-1:111122223333:my-dlq
    ```

=== "Serverless Framework"

    ```yaml
    functions:
      hello:
        handler: main.handler
        # SNS will trigger the function asynchronously
        events:
          - sns:
              topicName: my-sns-topic
        # Configure a failure destination for the function
        destinations:
          onFailure: arn:aws:sqs:us-east-1:111122223333:my-dlq
    ```

=== "Terraform"

    ```tf
    resource "aws_lambda_function" "this" {
      function_name = "my-function"
      runtime       = "python3.8"
      handler       = "main.handler"
      filename      = "function.zip"
    }

    resource "aws_lambda_permission" "this" {
      action        = "lambda:InvokeFunction"
      function_name = aws_lambda_function.this.function_name
      # Grants the permission to SNS to invoke this function
      # SNS will trigger the function asynchronously
      principal     = "sns.amazonaws.com"
    }

    resource "aws_lambda_function_event_invoke_config" "example" {
      function_name = aws_lambda_alias.example.function_name

      # Configure a failure destination for the function
      destination_config {
        on_failure {
          destination = "arn:aws:sqs:us-east-1:111122223333:my-dlq"
        }
      }
    }
    ```

## See also

* [Asynchronous invocation](https://docs.aws.amazon.com/lambda/latest/dg/invocation-async.html)
* [Serverless Lens: Failure Management](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/failure-management.html)
* [__CloudFormation__: AWS::Lambda::EventInvokeConfig](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-eventinvokeconfig.html)
* [__Terraform__: aws_lambda_function_event_invoke_config](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function_event_invoke_config)