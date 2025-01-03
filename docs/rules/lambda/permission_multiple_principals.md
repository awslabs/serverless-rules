# Lambda Permission Multiple Principals

__Level__: Warning
{: class="badge badge-yellow" }

__Initial version__: 0.1.3
{: class="badge badge-blue" }

__cfn-lint__: WS1002
{: class="badge" }

__tflint__: aws_lambda_permission_multiple_principals
{: class="badge" }

You can use resource-based policies to grant permission to other AWS services to invoke your Lambda functions. Different AWS services usually send different payloads to Lambda functions. If multiple services can invoke the same function, this function needs to handle the different types of payload properly, or this could cause unexpected behavior.

In general, it's better to create multiple Lambda functions with different function handlers for each invocation source.

## Implementations

=== "CDK"

    ```typescript
    import { ServicePrincipal } from '@aws-cdk/aws-iam';
    import { Function } from '@aws-cdk/aws-lambda';
    import { SnsEventSource } from '@aws-cdk/aws-lambda-event-sources';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const myFunction = new Function(
          scope, 'MyFunction',
          {
            code: Code.fromAsset('src/hello/'),
            handler: 'main.handler',
            runtime: Runtime.PYTHON_3_8,
          }
        );

        // This will implicitely grant those SNS topics the permission to invoke
        // the Lambda function. As both come from SNS, this is a single principal
        // ('sns.amazonaws.com') and thus will not trigger the rule.
        myFunction.addEventSource(new SnsEventSource(myTopic1));
        myFunction.addEventSource(new SnsEventSource(myTopic2));
      }
    }
    ```

=== "CloudFormation (JSON)"

    ```json
    {
      "Resources": {
        "MyFunction": {
          "Type": "AWS::Serverless::Function",
          "Properties": {
            "CodeUri": ".",
            "Runtime": "python3.8",
            "Handler": "main.handler",
            // This will implicitely grant those SNS topics the permission to invoke
            // the Lambda function. As both come from SNS, this is a single principal
            // ('sns.amazonaws.com') and thus will not trigger the rule.
            "Events": {
              "Topic1": {
                "Type": "SNS",
                "Properties": {
                  "Topic": arn:aws:sns:us-east-1:111122223333:topic1
                }
              }
              "Topic2": {
                "Type": "SNS",
                "Properties": {
                  "Topic": arn:aws:sns:us-east-1:111122223333:topic2
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
    Resources:
      MyFunction:
        Type: AWS::Serverless::Function
        Properties:
          CodeUri: .
          Runtime: python3.12
          Handler: main.handler
          Tracing: Active
          # This will implicitely grant those SNS topics the permission to invoke
          # the Lambda function. As both come from SNS, this is a single principal
          # ('sns.amazonaws.com') and thus will not trigger the rule.
          Events:
            Topic1:
              Type: SNS
              Properties:
                Topic: arn:aws:sns:us-east-1:111122223333:topic1
            Topic2:
              Type: SNS
              Properties:
                Topic: arn:aws:sns:us-east-1:111122223333:topic2
    ```

=== "Serverless Framework"

    ```yaml
    functions:
      hello:
        handler: handler.hello
        # This will implicitely grant those SNS topics the permission to invoke
        # the Lambda function. As both come from SNS, this is a single principal
        # ('sns.amazonaws.com') and thus will not trigger the rule.
        events:
          - sns: topic1
          - sns: topic2
    ```

=== "Terraform"

    ```tf
    resource "aws_iam_role" "this" {
      name = "my-function-role"
      assume_role_policy = data.aws_iam_policy_document.assume.json

      inline_policy {
        name = "FunctionPolicy"
        policy = data.aws_iam_policy_document.this.json
      }
    }

    data "aws_iam_policy_document" "assume" {
      statement {
        actions = ["sts:AssumeRole"]
        principals {
          type       = "Service"
          identifiers = ["lambda.amazonaws.com"]
        }
      }
    }

    data "aws_iam_policy_document" "this" {
      statement {
        # Tightly scoped permissions to just 'dynamodb:Query'
        # instead of 'dynamodb:*' or '*'
        actions = ["dynamodb:Query"]
        resources = ["arn:aws:dynamodb:eu-west-1:111122223333:table/my-table"]
      }
    }

    resource "aws_lambda_function" "this" {
      function_name = "my-function"
      handler       = "main.handler"
      runtime       = "python3.8"
      filename      = "function.zip"
      role          = aws_iam_role.this.arn
    }

    # Add a Lambda permission for Amazon EventBridge
    resource "aws_lambda_permission" "this" {
      statement_id  = "AllowExecutionFromEventBridge"
      action        = "lambda:InvokeFunction"
      function_name = aws_lambda_function.this.function_name
      principal     = "events.amazonaws.com"
    }
    ```

## Why is this a warning?

You might have a valid reason for invoking a Lambda function from different event sources or AWS services. If this is the case and you know what you are doing, you might ignore this rule.

## See also
* [Using resource-based policies for AWS Lambda](https://docs.aws.amazon.com/lambda/latest/dg/access-control-resource-based.html)