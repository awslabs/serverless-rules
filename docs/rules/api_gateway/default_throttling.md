# API Gateway Default Throttling

__Level__: Error
{: class="badge badge-red" }

__Initial version__: 0.1.3
{: class="badge badge-blue" }

__cfn-lint__: ES2003
{: class="badge" }

__tflint (REST)__: aws_apigateway_stage_throttling_rule
{: class="badge" }

__tflint (HTTP)__: aws_apigatewayv2_stage_throttling_rule
{: class="badge" }

Amazon API Gateway supports defining default limits for an API to prevent it from being overwhelmed by too many requests. This uses a [token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket), where a token counts for a single request.

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
            deployOptions: {
              // Throttling for default methods
              methodOptions: {
                '*/*': {
                  throttlingBurstLimit: 1000,
                  throttlingRateLimit: 10,
                }
              }
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
        "MyApi": {
          "Type": "AWS::Serverless::Api",
          "Properties": {
            "DefinitionUri": "openapi.yaml",
            "StageName": "prod",
            // Throttling for default methods by setting HttpMethod  to '*' and
            // ResourcePath to '/*'
            "MethodSettings": [{
              "HttpMethod": "*",
              "ResourcePath": "/*",
              "ThrottlingRateLimit": 10,
              "ThrottlingBurstLimit": 1000
            }]
          }
        }
      }
    }
    ```

=== "CloudFormation (YAML)"

    ```yaml
    Resources:
      MyApi:
        Type: AWS::Serverless::Api
        Properties:
          DefinitionUri: openapi.yaml
          StageName: prod

          # Throttling for default methods by setting HttpMethod  to '*' and
          # ResourcePath to '/*'
          MethodSettings:
            - HttpMethod: "*"
              ResourcePath: "/*"
              ThrottlingRateLimit: 10
              ThrottlingBurstLimit: 1000
    ```

=== "Serverless Framework"

    ```yaml
    resources:
      Resources:
        MyApi:
          Type: AWS::Serverless::Api
          Properties:
            DefinitionUri: openapi.yaml
            StageName: prod

            # Throttling for default methods by setting HttpMethod  to '*' and
            # ResourcePath to '/*'
            MethodSettings:
              - HttpMethod: "*"
                ResourcePath: "/*"
                ThrottlingRateLimit: 10
                ThrottlingBurstLimit: 1000
    ```

=== "Terraform"

    ```tf
    resource "aws_api_gateway_stage" "this" {
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
    }

    # Throttling for default methods by setting method_path to '*/*'
    resource "aws_api_gateway_method_settings" "this" {
      rest_api_id = aws_api_gateway_rest_api.this.id
      stage_name  = aws_api_gateway_stage.this.stage_name
      method_path = "*/*"

      settings {
        throttling_burst_limit = 1000
        throttling_rate_limit  = 10
      }
    }
    ```

## Implementations for HTTP APIs

=== "CDK"

    __Remark__: this is currently not supported in AWS CDK as [an L2 construct](https://docs.aws.amazon.com/cdk/latest/guide/constructs.html) at the moment.

    ```typescript
    import { CfnStage, HttpApi } from '@aws-cdk/aws-apigatewayv2';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const myApi = new HttpApi(
          scope, 'MyApi'
        );

        // Throttling for default methods by setting method_path to '*/* using escape hatch.

        // See https://docs.aws.amazon.com/cdk/latest/guide/cfn_layer.html#cfn_layer_resource
        // for more information.
        const defaultStage = myApi.defaultStage.node.defaultChild as CfnStage;
        defaultStage.defaultRouteSettings = {
          throttlingBurstLimit = 1000,
          throttlingRateLimit = 10,
        };
      }
    }
    ```

=== "CloudFormation (JSON)"

    ```json
    {
      "Resources": {
        "MyApi": {
          "Type": "AWS::Serverless::HttpApi",
          "Properties": {
            "DefinitionUri": "openapi.yaml",
            "StageName": "prod",
            "DefaultRouteSettings": {
              "ThrottlingBurstLimit": 1000,
              "ThrottlingRateLimit": 10
            }
          }
        }
      }
    }
    ```

=== "CloudFormation (YAML)"

    ```yaml
    Resources:
      MyApi:
        Type: AWS::Serverless::HttpApi
        Properties:
          DefinitionUri: "openapi.yaml"
          StageName: prod
          DefaultRouteSettings:
            ThrottlingBurstLimit: 1000
            ThrottlingRateLimit: 10
    ```

=== "Serverless Framework"

    ```yaml
    resources:
      Resources:
        MyApi:
          Type: AWS::Serverless::HttpApi
          Properties:
            DefinitionUri: "openapi.yaml"
            StageName: prod
            DefaultRouteSettings:
              ThrottlingBurstLimit: 1000
              ThrottlingRateLimit: 10
    ```

=== "Terraform"

    ```tf
    resource "aws_apigatewayv2_api" "this" {
      name          = "my-api"
      protocol_type = "HTTP"
      body          = file("openapi.yaml") 
    }

    resource "aws_apigatewayv2_stage" "this" {
      api_id = aws_apigatewayv2_api.this.id
      name   = "prod"

      # Default throttling settings
      default_route_settings {
        throttling_burst_limit = 1000
        throttling_rate_limit  = 10
      }
    }
    ```

## See also

* [Throttle API requests for better throughput](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-request-throttling.html)
* [Throttling requests to your HTTP API](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-throttling.html)
