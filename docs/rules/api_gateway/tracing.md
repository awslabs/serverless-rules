# API Gateway Tracing

__Level__: Warning
{: class="badge badge-yellow" }

__Initial version__: 0.1.3
{: class="badge badge-blue" }

__cfn-lint__: WS2002
{: class="badge" }

__tflint (REST)__: aws_apigateway_stage_tracing_rule
{: class="badge" }

__tflint (HTTP)__: _Not supported_
{: class="badge" }

Amazon API Gateway can emit traces to AWS X-Ray, which enable visualizing service maps for faster troubleshooting.

## Implementations for REST APIs

=== "CDK"

    ```typescript
    import { RestApi } from '@aws-cdk/aws-apigateway';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const myApi = new RestApi(
          scope, 'MyApi',
          {
            // Enable tracing on API Gateway
            deployOptions: {
              tracingEnabled: true,
            },
          }
        );
      }
    }
    ```

=== "CloudFormation (JSON)"

    ```json
    {
      "Resources": {
        "Api": {
          "Type": "AWS::Serverless::Api",
          "Properties": {
            "DefinitionUri": "openapi.yaml",
            "StageName": "prod",

            // Enable tracing on API Gateway
            "TracingEnabled": true
          }
        }
      }
    }
    ```

=== "CloudFormation (YAML)"

    ```yaml
    Resources:
      Api:
        Type: AWS::Serverless::Api
        Properties:
          DefinitionUri: openapi.yaml
          StageName: prod

          # Enable tracing on API Gateway
          TracingEnabled: true
    ```

=== "Serverless Framework"

    ```yaml
    provider:
      name: aws
      # Enable tracing on API Gateway
      tracing:
        apiGateway: true
    ```

=== "Terraform"

    ```tf
    resource "aws_api_gateway_rest_api" "this" {
      body = file("openapi.yaml")
    }

    resource "aws_api_gateway_deployment" "this" {
      rest_api_id = aws_api_gateway_rest_api.this.id

      triggers = {
        redeployment = sha1(jsonencode(aws_api_gateway_rest_api.this.body))
      }

      lifecycle {
        create_before_destroy = true
      }
    }

    resource "aws_api_gateway_stage" "this" {
      deployment_id = aws_api_gateway_deployment.this.id
      rest_api_id   = aws_api_gateway_rest_api.this.id
      stage_name    = "prod"

      # Enable tracing on API Gateway
      xray_tracing_enabled = true
    }
    ```

## Implementations for HTTP APIs

__Remark__: HTTP APIs do not support tracing at the moment.

## Why is this a warning?

You might use [third party solutions](https://aws.amazon.com/lambda/partners/) for monitoring serverless applications. If this is the case, enabling tracing for API Gateway might be optional. Refer to the documentation of your monitoring solutions to see if you should enable AWS X-Ray tracing or not.

## See also

* [Serverless Lens: Distributed Tracing](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/distributed-tracing.html)
* [Tracing user requests to REST APIs using X-Ray](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-xray.html)