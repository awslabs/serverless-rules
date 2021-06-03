AWS Lambda Rules
================

## EventSourceMapping Failure Destination

* __Level__: Error
* __cfn-lint__: ES1001
* __tflint__: _Not implemented_

An AWS Lambda event source mapping is used to read from streams and poll-based event sources. You can configure the event source mapping to send invocation records to another services such as Amazon SNS or Amazon SQS when it discards an event batch.

### Implementations

=== "CDK"

    ```typescript
    import { EventSourceMapping, SqsDlq, StartingPosition } from '@aws-cdk/aws-lambda';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        new EventSourceMapping(scope, "MyEventSourceMapping", {
          target: myFunction,
          eventSourceArn: 'arn:aws:dynamodb:us-east-1:111122223333:table/my-table/stream/my-stream',
          startingPosition: StartingPosition.LATEST,
          onFailure: SqsDlq(mySqsQueue),
        });
      }
    }
    ```

=== "CloudFormation (JSON)"

    ```json
    {
      "Resource": {
        "MyEventSourceMapping": {
          "Type": "AWS::Lambda::EventSourceMapping",
          "Properties": {
            // Required properties
            "FunctionName": "my-function",
            "EventSourceArn": "arn:aws:dynamodb:us-east-1:111122223333:table/my-table/stream/my-stream",
            "StartingPosition": "LATEST",

            // Add an OnFailure destination on the event source mapping
            "DestinationConfig": {
              "OnFailure": {
                "Destination": "arn:aws:sqs:us-east-1:111122223333:my-dlq"
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
      MyEventSourceMapping:
        Type: AWS::Lambda::EventSourceMapping
        Properties:
          # Required properties
          FunctionName: my-function
          EventSourceArn: arn:aws:dynamodb:us-east-1:111122223333:table/my-table/stream/my-stream
          StartingPosition: LATEST

          # Add an OnFailure destination on the event source mapping
          DestinationConfig:
            OnFailure:
              Destination: arn:aws:sqs:us-east-1:111122223333:my-dlq 
    ```

=== "Serverless Framework"

    ```yaml
    functions:
      MyFunction:
        handler: hello.handler

    resources:
      Resources:
        MyEventSourceMapping:
          Type: AWS::Lambda::EventSourceMapping
          Properties:
            # Required properties
            FunctionName:
              Fn::Ref: MyFunction
            EventSourceArn: arn:aws:dynamodb:us-east-1:111122223333:table/my-table/stream/my-stream
            StartingPosition: LATEST

            # Add an OnFailure destination on the event source mapping
            DestinationConfig:
              OnFailure:
                Destination: arn:aws:sqs:us-east-1:111122223333:my-dlq 
    ```

=== "Terraform"

    ```tf
    resource "aws_lambda_event_source_mapping" "this" {
      # Required fields
      event_source_arn  = "arn:aws:dynamodb:us-east-1:111122223333:table/my-table/stream/my-stream"
      function_name     = "my-function"
      starting_position = "LATEST"

      # Add an OnFailure destination on the event source mapping
      destination_config {
        on_failure {
          destination_arn = "arn:aws:sqs:us-east-1:111122223333:my-dlq"
        }
      }
    }
    ```


### See also

* [AWS Lambda event source mappings](https://docs.aws.amazon.com/lambda/latest/dg/invocation-eventsourcemapping.html)
* [__CloudFormation__: AWS::Lambda::EventSourceMapping](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-eventsourcemapping.html#cfn-lambda-eventsourcemapping-destinationconfig)
* [__Terraform__: aws_lambda_event_source_mapping](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_event_source_mapping)

## Log Retention

* __Level__: Warning
* __cfn-lint__: WS1004
* __tflint__: _Not implemented_

By default, CloudWatch log groups created by Lambda functions have an unlimited retention time. For cost optimization purposes, you should set a retention duration on all log groups. For log archival, export and set cost-effective storage classes that best suit your needs.

### Why is this a warning?

Since `serverless-rules` evaluate infrastructure as code template, it cannot check if you use a solution that will automatically change the configuration of log groups after the fact.

### Implementations

=== "CDK"

    ```typescript
    import { Code, Function, Runtime } from '@aws-cdk/aws-lambda';
    import { LogGroup, RetentionDays } from '@aws-cdk/aws-logs';

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

        // Explicit log group that refers to the Lambda function
        const myLogGroup = new LogGroup(
          scope, 'MyLogGroup',
          {
            logGroupName: `/aws/lambda/${myFunction.functionName}`,
            retention: RetentionDays.ONE_WEEK,
          }
        );
      }
    }
    ```

=== "CloudFormation (JSON)"

    ```json
    {
      "Resources": {
        // Lambda function
        "Function": {
          "Type": "AWS::Serverless::Function",
          "Properties": {
            "CodeUri": ".",
            "Runtime": "python3.8",
            "Handler": "main.handler",
            "Tracing": "Active"
          }
        },

        // Explicit log group that refers to the Lambda function
        "LogGroup": {
          "Type": "AWS::Logs::LogGroup",
          "Properties": {
            "LogGroupName": {
              "Fn::Sub": "/aws/lambda/${Function}"
            },
            // Explicit retention time
            "RetentionInDays": 7
          }
        }
      }
    }
    ```

=== "CloudFormation (YAML)"

    ```yaml
    Resources:
      Function:
        Type: AWS::Serverless::Function
        Properties:
          CodeUri: .
          Runtime: python3.8
          Handler: main.handler
          Tracing: Active

      # Explicit log group that refers to the Lambda function
      LogGroup:
        Type: AWS::Logs::LogGroup
        Properties:
          LogGroupName: !Sub "/aws/lambda/${Function}"
          # Explicit retention time
          RetentionInDays: 7
    ```

=== "Serverless Framework"

    ```yaml
    provider:
      name: aws
      runtime: python3.8
      lambdaHashingVersion: '20201221'
      # This will automatically create the log group with retention
      logRetentionInDays: 14
        
    functions:
      hello:
        handler: handler.hello
    ```

=== "Terraform"

    ```tf
    resource "aws_lambda_function" "this" {
      function_name = "my-function"
      handler       = "main.handler"
      runtime       = "python3.8"
      filename      = "function.zip"
      role          = "arn:aws:iam::111122223333:role/my-function-role"
    }

    # Explicit log group
    resource "aws_cloudwatch_log_group" "this" {
      name              = "/aws/lambda/{aws_lambda_function.this.function_name}
      # Explicit retention time
      retention_in_days = 7
    }
    ```

### See also

* [Serverless Lens: Logging Ingestion and Storage](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/logging-ingestion-and-storage.html)
* [__CloudFormation__: AWS::Logs::LogGroup](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-logs-loggroup.html)
* [__Terraform__: aws_cloudwatch_log_group](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group)

## Permission Multiple Principals

* __Level__: Warning
* __cfn-lint__: WS1002
* __tflint__: _Not implemented_

You can use resource-based policies to grant permission to other AWS services to invoke your Lambda functions. Different AWS services usually send different payloads to Lambda functions. If multiple services can invoke the same function, this function needs to handle the different types of payload properly, or this could cause unexpected behavior.

In general, it's better to deploy multiple Lambda functions with different function handlers for each invocation source.

### Implementations

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
          Runtime: python3.8
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

### Why is this a warning?

You might have a valid reason for invoking a Lambda function from different event sources or AWS services. If this is the case and you know what you are doing, you might ignore this rule.

### See also
* [Using resource-based policies for AWS Lambda](https://docs.aws.amazon.com/lambda/latest/dg/access-control-resource-based.html)

## Star Permissions

* __Level__: Warning
* __cfn-lint__: WS1003
* __tflint__: _Not implemented_

 With Lambda functions, it’s recommended that you follow least-privileged access and only allow the access needed to perform a given operation. Attaching a role with more permissions than necessary can open up your systems for abuse.

 ### Why is this a warning?

If your Lambda function need a broad range of permissions, you do not know ahead of time which permissions you will need, and you have evaluated the risks of using broad permissions for this function, you might ignore this rule.


### Implementations

=== "CDK"

    ```typescript
    import { AttributeType, Table } from '@aws-cdk/aws-dynamodb';
    import { Code, Function, Runtime } from '@aws-cdk/aws-lambda';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const myTable = new Table(
          scope, 'MyTable',
          {
            partitionKey: {
              name: 'id',
              type: AttributeType.STRING,
            }
          },
        );

        const myFunction = new Function(
          scope, 'MyFunction',
          {
            code: Code.fromAsset('src/hello/'),
            handler: 'main.handler',
            runtime: Runtime.PYTHON_3_8,
          }
        );

        // Grant read access to the DynamoDB table
        table.grantReadData(myFunction);
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

            "Policies": [{
              "Version": "2012-10-17",
              "Statement": [{
                "Effect": "Allow",
                // Tightly scoped permissions to just 's3:GetObject'
                // instead of 's3:*' or '*'
                "Action": "s3:GetObject",
                "Resource": "arn:aws:s3:::my-bucket/*"
              }]
            }]
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
          Runtime: python3.8
          Handler: main.handler

          Policies:
            - Version: "2012-10-17"
              Statement:
                - Effect: Allow
                  # Tightly scoped permissions to just 's3:GetObject'
                  # instead of 's3:*' or '*'
                  Action: s3:GetObject
                  Resource: "arn:aws:s3:::my-bucket/*"
    ```

=== "Serverless Framework"

    ```yaml
    provider:
      name: aws
      iam:
        role:
          name: my-function-role
          statements:
            - Effect: Allow
              # Tightly scoped permissions to just 's3:GetObject'
              # instead of 's3:*' or '*'
              Action: s3:GetObject
              Resource: "arn:aws:s3:::my-bucket/*"
        
    functions:
      hello:
        handler: handler.hello
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
    ```

### See also
* [Serverless Lens: Identity and Access Management](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/identity-and-access-management.html)
* [AWS Lambda execution role](https://docs.aws.amazon.com/lambda/latest/dg/lambda-intro-execution-role.html)

## Tracing

* __Level__: Warning
* __cfn-lint__: WS1000
* __tflint__: aws_lambda_function_tracing_rule

AWS Lambda can emit traces to AWS X-Ray, which enable visualizing service maps for faster troubleshooting.

### Why is this a warning?

You might use [third party solutions](https://aws.amazon.com/lambda/partners/) for monitoring serverless applications. If this is the case, enabling tracing for your AWS Lambda functions might be optional. Refer to the documentation of your monitoring solutions to see if you should enable AWS X-Ray tracing or not.

### Implementations

=== "CDK"

    ```typescript
    import { Code, Function, Runtime, Tracing } from '@aws-cdk/aws-lambda';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const myFunction = new Function(
          scope, 'MyFunction',
          {
            code: Code.fromAsset('src/hello/'),
            handler: 'main.handler',
            runtime: Runtime.PYTHON_3_8,
            // Enable active tracing
            tracing: Tracing.ACTIVE,
          }
        );
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
            // Required properties
            "CodeUri": ".",
            "Runtime": "python3.8",
            "Handler": "main.handler",

            // Enable active tracing
            "Tracing": "Active"
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
          # Required properties
          CodeUri: .
          Runtime: python3.8
          Handler: main.handler

          # Enable active tracing
          Tracing: Active
    ```

=== "Serverless Framework"

    ```yaml
    provider:
      tracing:
        # Enable active tracing for Lambda functions
        lambda: true

    functions:
      hello:
        handler: handler.hello
    ```

=== "Terraform"

    ```tf
    resource "aws_lambda_function" "this" {
      function_name = "my-function"
      runtime       = "python3.8"
      handler       = "main.handler"
      filename      = "function.zip"

      # Enable active tracing
      tracing_config {
        mode = "Active"
      }
    }
    ```

### See also

* [Serverless Lens: Distributed Tracing](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/distributed-tracing.html)
* [Using AWS Lambda with X-Ray](https://docs.aws.amazon.com/lambda/latest/dg/services-xray.html)