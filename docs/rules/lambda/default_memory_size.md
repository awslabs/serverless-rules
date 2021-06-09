# Lambda Default Memory Size

__Level__: Error
{: class="badge badge-red" }

__Initial version__: 0.1.8
{: class="badge badge-blue" }

__cfn-lint__: ES1005
{: class="badge" }

__tflint__: aws_lambda_function_default_memory
{: class="badge" }

Lambda allocates CPU power in proportion to the amount of memory configured. By default, your functions have 128 MB of memory allocated. You can increase that value up to 10 GB. With more CPU resources, your Lambda function's duration might decrease.

You can use tools such as [AWS Lambda Power Tuning](https://github.com/alexcasalboni/aws-lambda-power-tuning) to test your function at different memory settings to find the one that matches your cost and performance requirements the best.

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
            // Change the function memory size
            memorySize: 2048,
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

            // Change the function memory size
            "MemorySize": 2048
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

          # Change the function memory size
          MemorySize: 2048

    ```

=== "Serverless Framework"

    ```yaml
    provider:
      name: aws
      # Change the memory size across all functions
      memorySize: 2048

    functions:
      hello:
        handler: handler.hello
        # Change the memory size for one function
        memorySize: 512
    ```

=== "Terraform"

    ```tf
    resource "aws_lambda_function" "this" {
      function_name = "my-function"
      runtime       = "python3.8"
      handler       = "main.handler"
      filename      = "function.zip"

      # Change the default memory size value
      memory_size = 2048
    }
    ```

## See also

* [Configuring Lambda function memory](https://docs.aws.amazon.com/lambda/latest/dg/configuration-memory.html)
* [Serverless Lens: Optimize](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/optimize.html)
* [AWS Lambda Power Tuning](https://github.com/alexcasalboni/aws-lambda-power-tuning)
* [__CloudFormation__: AWS::Lambda::Function](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-function.html)
* [__Terraform__: aws_lambda_function](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function)