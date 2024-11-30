# Lambda End-of-life Runtime

__Level__: Error
{: class="badge badge-red" }

__Initial version__: 0.1.7
{: class="badge badge-blue" }

__cfn-lint__: E2531
{: class="badge" }

__tflint__: aws_lambda_function_eol_runtime
{: class="badge" }

Managed Lambda runtimes for .zip file archives are built around a combination of operating system, programming language, and software libraries that are subject to maintenance and security updates. When security updates are no longer available for a component of a runtime, Lambda deprecates the runtime.

!!! info

    This rule is implemented natively in `cfn-lint` as rule number __E2531__.

## Implementations

=== "CDK"

    ```typescript
    import { Code, Function, Runtime } from '@aws-cdk/aws-lambda';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const myFunction = new Function(
          scope, 'MyFunction',
          {
            code: Code.fromAsset('src/hello/'),
            handler: 'main.handler',
            // Select a runtime that is not deprecated
            runtime: Runtime.PYTHON_3_8,
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
            "CodeUri": ".",
            // Select a runtime that is not deprecated
            "Runtime": "python3.8",
            "Handler": "main.handler"
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
          # Select a runtime that is not deprecated
          Runtime: python3.12
          Handler: main.handler
    ```

=== "Serverless Framework"

    ```yaml
    provider:
      name: aws
      # Select a runtime that is not deprecated
      runtime: nodejs14.x

    functions:
      hello:
        handler: handler.hello
    ```

=== "Terraform"

    ```tf
    resource "aws_lambda_function" "this" {
      function_name = "my-function"
      # Select a runtime that is not deprecated
      runtime       = "python3.8"
      handler       = "main.handler"
      filename      = "function.zip"
    }
    ```

## See also

* [Runtime support policy](https://docs.aws.amazon.com/lambda/latest/dg/runtime-support-policy.html)