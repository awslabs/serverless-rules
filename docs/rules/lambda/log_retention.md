# Lambda Log Retention

__Level__: Warning
{: class="badge badge-yellow" }

__Initial version__: 0.1.3
{: class="badge badge-blue" }

__cfn-lint__: WS1004
{: class="badge" }

__tflint__: aws_cloudwatch_log_group_lambda_retention
{: class="badge" }

By default, CloudWatch log groups created by Lambda functions have an unlimited retention time. For cost optimization purposes, you should set a retention duration on all log groups. For log archival, export and set cost-effective storage classes that best suit your needs.

??? warning "Referencing the function name in the log group"

    This rule works by matching a Lambda function name in the CloudWatch log group name. For CloudFormation, it supports `Fn::Join`, `Fn::Sub`, and hard-coding the Lambda function name into the log group name.

    Here are some examples of valid implementations:

    === "With Fn::Join"
    
        ```yaml
        Resources:
          Function:
            Type: AWS::Serverless::Function
            Properties:
              # Omitting other properties

          LogGroup:
            Type: AWS::Logs::LogGroup
            Properties:
              LogGroupName:
                Fn::Join:
                - ""
                - - "/aws/lambda/"
                  - !Ref Function
              RetentionInDays: 7
        ```

    === "With Fn::Sub"
    
        ```yaml
        Resources:
          Function:
            Type: AWS::Serverless::Function
            Properties:
              # Omitting other properties

          LogGroup:
            Type: AWS::Logs::LogGroup
            Properties:
              LogGroupName: !Sub "/aws/lambda/${Function}"
              RetentionInDays: 7
        ```

    === "With function name"
    
        ```yaml
        Resources:
          Function:
            Type: AWS::Serverless::Function
            Properties:
              # Omitting other properties
              FunctionName: my_function_name

          LogGroup:
            Type: AWS::Logs::LogGroup
            Properties:
              LogGroupName: "/aws/lambda/my_function_name
              RetentionInDays: 7
        ```

## Why is this a warning?

Since `serverless-rules` evaluate infrastructure-as-code template, it cannot check if you use a solution that will automatically change the configuration of log groups after the fact.

## Implementations

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
      name              = "/aws/lambda/${aws_lambda_function.this.function_name}
      # Explicit retention time
      retention_in_days = 7
    }
    ```

## See also

* [Serverless Lens: Logging Ingestion and Storage](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/logging-ingestion-and-storage.html)
* [__CloudFormation__: AWS::Logs::LogGroup](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-logs-loggroup.html)
* [__Terraform__: aws_cloudwatch_log_group](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group)