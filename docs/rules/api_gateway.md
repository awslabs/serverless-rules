Amazon API Gateway Rules
========================

## Default Throttling

* __Level__: Error
* __cfn-lint__: ES2003
* __tflint (REST APIs)__: aws_apigateway_stage_throttling_rule
* __tflint (HTTP APIs)__: aws_apigatewayv2_stage_throttling_rule

Amazon API Gateway supports defining default limits for an API to prevent it from being overwhelmed by too many requests. This uses a [token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket), where a token counts for a single request.

### Implementations for REST APIs

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
                  throttlingRateLimite: 10,
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

### Implementations for HTTP APIs

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

### See also

* [Throttle API requests for better throughput](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-request-throttling.html)
* [Throttling requests to your HTTP API](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-throttling.html)

## Logging

* __Level__: Error
* __cfn-lint__: ES2000
* __tflint (REST APIs)__: aws_apigateway_stage_logging_rule
* __tflint (HTTP APIs)__: aws_apigatewayv2_stage_logging_rule

Amazon API Gateway can send logs to Amazon CloudWatch Logs and Amazon Kinesis Data Firehose for centralization.

### Implementations for REST APIs

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

### Implementations for HTTP APIs

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

### See also

* [Serverless Lens: Centralized and structured logging](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/centralized-and-structured-logging.html)
* [Monitoring REST APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/rest-api-monitor.html)
* [Monitoring your HTTP API](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-monitor.html)
* [Monitoring WebSocket APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/websocket-api-monitor.html)

## Structured Logging

* __Level__: Error
* __cfn-lint__: WS2001
* __tflint__: _Not implemented_

You can customize the log format that Amazon API Gateway uses to send logs. Structured logging makes it easier to derive queries to answer arbitrary questions about the health of your application.

### Why is this a warning?

The rule in `serverless-rules` only checks if the structured log is JSON-formatted.

While CloudWatch Logs Insights will automatically discover fields in JSON log entries, you can use the `parse` command to parse custom log entries to extract fields from custom format.

### Implementations

See the implementations for [Logging on API Gateway](#logging).

### See also

* [Serverless Lens: Centralized and structured logging](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/centralized-and-structured-logging.html)
* [Monitoring REST APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/rest-api-monitor.html)
* [Monitoring your HTTP API](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-monitor.html)
* [Monitoring WebSocket APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/websocket-api-monitor.html)
* [Amazon CloudWatch Logs: Supported Logs and Discovered Fields](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CWL_AnalyzeLogData-discoverable-fields.html)
* [Amazon CloudWatch Logs: Logs Insights Query Syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CWL_QuerySyntax.html)

## Tracing

* __Level__: Warning
* __cfn-lint__: WS2002
* __tflint (REST APIs)__: aws_apigateway_stage_tracing_rule
* __tflint (HTTP APIs)__: _Not supported_

Amazon API Gateway can emit traces to AWS X-Ray, which enable visualizing service maps for faster troubleshooting.

### Implementations for REST APIs

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

### Implementations for HTTP APIs

__Remark__: HTTP APIs do not support tracing at the moment.

### Why is this a warning?

You might use [third party solutions](https://aws.amazon.com/lambda/partners/) for monitoring serverless applications. If this is the case, enabling tracing for API Gateway might be optional. Refer to the documentation of your monitoring solutions to see if you should enable AWS X-Ray tracing or not.

### See also

* [Serverless Lens: Distributed Tracing](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/distributed-tracing.html)
* [Tracing user requests to REST APIs using X-Ray](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-xray.html)