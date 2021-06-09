# API Gateway Logging

__Level__: Error
{: class="badge badge-red" }

__Initial version__: 0.1.3
{: class="badge badge-blue" }

__cfn-lint__: ES2000
{: class="badge" }

__tflint (REST)__: aws_apigateway_stage_logging_rule
{: class="badge" }

__tflint (HTTP)__: aws_apigatewayv2_stage_logging_rule
{: class="badge" }

Amazon API Gateway can send logs to Amazon CloudWatch Logs and Amazon Kinesis Data Firehose for centralization.

## Implementations for REST APIs

=== "CDK" 

    ```typescript
    import { LogGroup } from '@aws-cdk/aws-logs';
    import { LogGroupLogDestination, RestApi } from '@aws-cdk/aws-apigateway';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const myLogGroup = new LogGroup(
          scope, 'MyLogGroup'
        );

        const myApi = new RestApi(
          scope, 'MyApi',
          {
            deployOptions: {
              // Setup logging for API Gateway
              accessLogDestination: new LogGroupLogDestination(myLogGroup),
              accessLogFormat: JSON.stringify({
                "stage" : "$context.stage",
                "request_id" : "$context.requestId",
                "api_id" : "$context.apiId",
                "resource_path" : "$context.resourcePath",
                "resource_id" : "$context.resourceId",
                "http_method" : "$context.httpMethod",
                "source_ip" : "$context.identity.sourceIp",
                "user-agent" : "$context.identity.userAgent",
                "account_id" : "$context.identity.accountId",
                "api_key" : "$context.identity.apiKey",
                "caller" : "$context.identity.caller",
                "user" : "$context.identity.user",
                "user_arn" : "$context.identity.userArn",
                "integration_latency": $context.integration.latency
              }),
            }
          }
        );
      }
    }
    ```

=== "CloudFormation (JSON)"

    ```json
    {
      "Resource": {
        "Type": "AWS::Serverless::Api",
        "Properties": {
          "DefinitionUri": "openapi.yaml",
          "StageName": "prod",

          // Setup logging for API Gateway
          "AccessLogSetting":{
            "DestinationArn": "arn:aws:logs:eu-west-1:123456789012:log-group:my-log-group",
            "Format": "{ \"stage\" : \"$context.stage\", \"request_id\" : \"$context.requestId\", \"api_id\" : \"$context.apiId\", \"resource_path\" : \"$context.resourcePath\", \"resource_id\" : \"$context.resourceId\", \"http_method\" : \"$context.httpMethod\", \"source_ip\" : \"$context.identity.sourceIp\", \"user-agent\" : \"$context.identity.userAgent\", \"account_id\" : \"$context.identity.accountId\", \"api_key\" : \"$context.identity.apiKey\", \"caller\" : \"$context.identity.caller\", \"user\" : \"$context.identity.user\", \"user_arn\" : \"$context.identity.userArn\", \"integration_latency\": $context.integration.latency }"
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

          # Setup logging for API Gateway
          AccessLogSetting:
            DestinationArn: "arn:aws:logs:eu-west-1:123456789012:log-group:my-log-group"
            Format: |
              {
                "stage" : "$context.stage",
                "request_id" : "$context.requestId",
                "api_id" : "$context.apiId",
                "resource_path" : "$context.resourcePath",
                "resource_id" : "$context.resourceId",
                "http_method" : "$context.httpMethod",
                "source_ip" : "$context.identity.sourceIp",
                "user-agent" : "$context.identity.userAgent",
                "account_id" : "$context.identity.accountId",
                "api_key" : "$context.identity.apiKey",
                "caller" : "$context.identity.caller",
                "user" : "$context.identity.user",
                "user_arn" : "$context.identity.userArn",
                "integration_latency": $context.integration.latency
              }
    ```

=== "Serverless Framework"

    ```yaml
    provider:
      name: aws
      logs:
        # Setup logging for API Gateway
        restApi:
          accessLogging: true
          format: |
            {
              "stage" : "$context.stage",
              "request_id" : "$context.requestId",
              "api_id" : "$context.apiId",
              "resource_path" : "$context.resourcePath",
              "resource_id" : "$context.resourceId",
              "http_method" : "$context.httpMethod",
              "source_ip" : "$context.identity.sourceIp",
              "user-agent" : "$context.identity.userAgent",
              "account_id" : "$context.identity.accountId",
              "api_key" : "$context.identity.apiKey",
              "caller" : "$context.identity.caller",
              "user" : "$context.identity.user",
              "user_arn" : "$context.identity.userArn",
              "integration_latency": $context.integration.latency
            }
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

      # Setup logging for API Gateway
      access_log_settings {
        destination_arn = "arn:aws:logs:eu-west-1:123456789012:log-group:my-log-group"
        format = <<EOF
    {
      "stage" : "$context.stage",
      "request_id" : "$context.requestId",
      "api_id" : "$context.apiId",
      "resource_path" : "$context.resourcePath",
      "resource_id" : "$context.resourceId",
      "http_method" : "$context.httpMethod",
      "source_ip" : "$context.identity.sourceIp",
      "user-agent" : "$context.identity.userAgent",
      "account_id" : "$context.identity.accountId",
      "api_key" : "$context.identity.apiKey",
      "caller" : "$context.identity.caller",
      "user" : "$context.identity.user",
      "user_arn" : "$context.identity.userArn",
      "integration_latency": $context.integration.latency
    }
    EOF
      }
    }
    ```

## Implementations for HTTP APIs

=== "CDK"

    __Remark__: this is currently not supported in AWS CDK as [an L2 construct](https://docs.aws.amazon.com/cdk/latest/guide/constructs.html) at the moment. See [this GitHub issue](https://github.com/aws/aws-cdk/issues/11100) for more details.

    ```typescript
    import { CfnStage, HttpApi } from '@aws-cdk/aws-apigatewayv2';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const myApi = new HttpApi(
          scope, 'MyApi'
        );

        // Setup logging for API Gateway using escape hatch.

        // See https://github.com/aws/aws-cdk/issues/11100 and
        // https://docs.aws.amazon.com/cdk/latest/guide/cfn_layer.html#cfn_layer_resource
        // for more information.
        const defaultStage = myApi.defaultStage.node.defaultChild as CfnStage;
        defaultStage.accessLogSettings = {
          destinationArn: '',
          format: JSON.stringify({
            "stage" : "$context.stage",
            "request_id" : "$context.requestId",
            "api_id" : "$context.apiId",
            "resource_path" : "$context.resourcePath",
            "resource_id" : "$context.resourceId",
            "http_method" : "$context.httpMethod",
            "source_ip" : "$context.identity.sourceIp",
            "user-agent" : "$context.identity.userAgent",
            "account_id" : "$context.identity.accountId",
            "api_key" : "$context.identity.apiKey",
            "caller" : "$context.identity.caller",
            "user" : "$context.identity.user",
            "user_arn" : "$context.identity.userArn",
            "integration_latency": $context.integration.latency
          }),
        };
      }
    }
    ```


=== "CloudFormation (JSON)"

    ```json
    {
      "Resource": {
        "Type": "AWS::Serverless::HttpApi",
        "Properties": {
          "DefinitionUri": "openapi.yaml",
          "StageName": "prod",

          // Setup logging for API Gateway
          "AccessLogSettings":{
            "DestinationArn": "arn:aws:logs:eu-west-1:123456789012:log-group:my-log-group",
            "Format": "{ \"stage\" : \"$context.stage\", \"request_id\" : \"$context.requestId\", \"api_id\" : \"$context.apiId\", \"resource_path\" : \"$context.resourcePath\", \"resource_id\" : \"$context.resourceId\", \"http_method\" : \"$context.httpMethod\", \"source_ip\" : \"$context.identity.sourceIp\", \"user-agent\" : \"$context.identity.userAgent\", \"account_id\" : \"$context.identity.accountId\", \"api_key\" : \"$context.identity.apiKey\", \"caller\" : \"$context.identity.caller\", \"user\" : \"$context.identity.user\", \"user_arn\" : \"$context.identity.userArn\", \"integration_latency\": $context.integration.latency }"
          }
        }
      }
    }
    ```

=== "CloudFormation (YAML)"

    ```yaml
    Resources:
      Api:
        Type: AWS::Serverless::HttpApi
        Properties:
          DefinitionUri: openapi.yaml
          StageName: prod

          # Setup logging for API Gateway
          AccessLogSettings:
            DestinationArn: "arn:aws:logs:eu-west-1:123456789012:log-group:my-log-group"
            Format: |
              {
                "stage" : "$context.stage",
                "request_id" : "$context.requestId",
                "api_id" : "$context.apiId",
                "resource_path" : "$context.resourcePath",
                "resource_id" : "$context.resourceId",
                "http_method" : "$context.httpMethod",
                "source_ip" : "$context.identity.sourceIp",
                "user-agent" : "$context.identity.userAgent",
                "account_id" : "$context.identity.accountId",
                "api_key" : "$context.identity.apiKey",
                "caller" : "$context.identity.caller",
                "user" : "$context.identity.user",
                "user_arn" : "$context.identity.userArn",
                "integration_latency": $context.integration.latency
              }
    ```

=== "Serverless Framework"

    ```yaml
    provider:
      name: aws
      logs:
        httpApi:
          format: |
            {
              "stage" : "$context.stage",
              "request_id" : "$context.requestId",
              "api_id" : "$context.apiId",
              "resource_path" : "$context.resourcePath",
              "resource_id" : "$context.resourceId",
              "http_method" : "$context.httpMethod",
              "source_ip" : "$context.identity.sourceIp",
              "user-agent" : "$context.identity.userAgent",
              "account_id" : "$context.identity.accountId",
              "api_key" : "$context.identity.apiKey",
              "caller" : "$context.identity.caller",
              "user" : "$context.identity.user",
              "user_arn" : "$context.identity.userArn",
              "integration_latency": $context.integration.latency
            }
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

      # Setup logging for API Gateway
      access_log_settings {
        destination_arn = "arn:aws:logs:eu-west-1:123456789012:log-group:my-log-group"
        format = <<EOF
    {
      "stage" : "$context.stage",
      "request_id" : "$context.requestId",
      "api_id" : "$context.apiId",
      "resource_path" : "$context.resourcePath",
      "resource_id" : "$context.resourceId",
      "http_method" : "$context.httpMethod",
      "source_ip" : "$context.identity.sourceIp",
      "user-agent" : "$context.identity.userAgent",
      "account_id" : "$context.identity.accountId",
      "api_key" : "$context.identity.apiKey",
      "caller" : "$context.identity.caller",
      "user" : "$context.identity.user",
      "user_arn" : "$context.identity.userArn",
      "integration_latency": $context.integration.latency
    }
    EOF
      }
    }
    ```

## See also

* [Serverless Lens: Centralized and structured logging](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/centralized-and-structured-logging.html)
* [Monitoring REST APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/rest-api-monitor.html)
* [Monitoring your HTTP API](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-monitor.html)
* [Monitoring WebSocket APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/websocket-api-monitor.html)