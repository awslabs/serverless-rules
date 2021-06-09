# Lambda Default Timeout

__Level__: Error
{: class="badge badge-red" }

__Initial version__: 0.1.8
{: class="badge badge-blue" }

__cfn-lint__: ES1006
{: class="badge" }

__tflint__: aws_lambda_function_default_timeout
{: class="badge" }

You can define the timeout value, which restricts the maximum duration of a single invocation of your Lambda functions.

If your timeout value is too short, Lambda might terminate invocations prematurely. On the other side, setting the timeout much higher than the average execution may cause functions to execute for longer upon code malfunction, resulting in higher costs and possibly reaching concurrency limits depending on how such functions are invoked.

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
            runtime: Runtime.PYTHON_3_8,
            // Change the function timeout
            timeout: 10,
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

            // Change the function timeout
            "Timeout": 10
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

          # Change the function timeout
          Timeout: 10

    ```

=== "Serverless Framework"

    ```yaml
    provider:
      name: aws
      # Change the timeout across all functions
      timeout: 10

    functions:
      hello:
        handler: handler.hello
        # Change the timeout for one function
        timeout: 15
    ```

=== "Terraform"

    ```tf
    resource "aws_lambda_function" "this" {
      function_name = "my-function"
      runtime       = "python3.8"
      handler       = "main.handler"
      filename      = "function.zip"

      # Change the function timeout
      timeout = 10
    }
    ```

## See also

* [AWS Lambda execution environment](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-context.html)
* [Serverless Lens: Optimize](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/optimize.html)
* [__CloudFormation__: AWS::Lambda::Function](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-function.html)
* [__Terraform__: aws_lambda_function](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function)