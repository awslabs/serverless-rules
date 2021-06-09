# AppSync Tracing

__Level__: Warning
{: class="badge badge-yellow" }

__Initial version__: 0.1.3
{: class="badge badge-blue" }

__cfn-lint__: WS3000
{: class="badge" }

__tflint__: aws_appsync_graphql_api_tracing_rule
{: class="badge" }

AWS AppSync can emit traces to AWS X-Ray, which enables visualizing service maps for faster troubleshooting.

## Why is this a warning?

You might use [third party solutions](https://aws.amazon.com/lambda/partners/) for monitoring serverless applications. If this is the case, enabling tracing for AppSync APIs might be optional. Refer to the documentation of your monitoring solutions to see if you should enable AWS X-Ray tracing or not.

## Implementations

=== "CDK"

    ```typescript
    import { GraphqlApi } from '@aws-cdk/aws-appsync';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const myApi = new GraphqlApi(
          scope, 'MyApi',
          {
            name: 'my-api',
            // Enable active tracing
            xrayEnabled: true,
          }
        );
      }
    }
    ```

=== "CloudFormation (JSON)"

    ```json
    {
      "Resources": {
        "GraphQLApi": {
        "Type": "AWS::AppSync::GraphQLApi",
        "Properties": {
          "Name": "api",
          "AuthenticationType": "AWS_IAM",

          // Enable active tracing
          "XrayEnabled": true
        }
        }
      }
    }
    ```

=== "CloudFormation (YAML)"

    ```yaml
    Resources:
      GraphQLApi:
        Type: AWS::AppSync::GraphQLApi
        Properties:
          Name: api
          AuthenticationType: AWS_IAM

          # Enable active tracing
          XrayEnabled: true
    ```

=== "Serverless Framework"

    ```yaml
    resources:
      Resources:
        GraphQLApi:
          Type: AWS::AppSync::GraphQLApi
          Properties:
            Name: api
            AuthenticationType: AWS_IAM

            # Enable active tracing
            XrayEnabled: true
    ```

=== "Terraform"

    ```tf
    resource "aws_appsync_graphql_api" "this" {
      name                = "api"
      authentication_type = "AWS_IAM"

      # Enable active tracing
      xray_enabled = true
    }
    ```

## See also

* [Serverless Lens: Distributed Tracing](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/distributed-tracing.html)
* [Tracing with AWS X-Ray](https://docs.aws.amazon.com/appsync/latest/devguide/x-ray-tracing.html)