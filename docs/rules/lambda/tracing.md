# Lambda Tracing

__Level__: Warning
{: class="badge badge-yellow" }

__Initial version__: 0.1.3
{: class="badge badge-blue" }

__cfn-lint__: WS1000
{: class="badge" }

__tflint__: aws_lambda_function_tracing_rule
{: class="badge" }

AWS Lambda can emit traces to AWS X-Ray, which enables visualizing service maps for faster troubleshooting.

## Why is this a warning?

You might use [third party solutions](https://aws.amazon.com/lambda/partners/) for monitoring serverless applications. If this is the case, enabling tracing for your AWS Lambda functions might be optional. Refer to the documentation of your monitoring solutions to see if you should enable AWS X-Ray tracing or not.

## Implementations

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
          Runtime: python3.12
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

## See also

* [Serverless Lens: Distributed Tracing](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/distributed-tracing.html)
* [Using AWS Lambda with X-Ray](https://docs.aws.amazon.com/lambda/latest/dg/services-xray.html)
