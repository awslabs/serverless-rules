Amazon API Gateway Rules
========================

## Default Throttling

* __Level__: Error
* __cfn-lint__: ES2003
* __tflint (REST APIs)__: aws_apigateway_stage_throttling_rule
* __tflint (HTTP APIs)__: aws_apigatewayv2_stage_throttling_rule

Amazon API Gateway supports defining default limits for an API to prevent it from being overwhelmed by too many requests. This uses a [token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket), where a token counts for a single request.

### Implementations for REST APIs

<details>
<summary>CDK</summary>

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
</details>

<details>
<summary>CloudFormation/SAM</summary>

__JSON__

```json
{
  "Resources": {
    "MyApi": {
      "Type": "AWS::Serverless::Api",
      "Properties": {
        "DefinitionUri": "openapi.yaml",
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

__YAML__

```yaml
Resources:
  MyApi:
    Type: AWS::Serverless::Api
    Properties:
      DefinitionUri: openapi.yaml

      # Throttling for default methods by setting HttpMethod  to '*' and
      # ResourcePath to '/*'
      MethodSettings:
        - HttpMethod: "*"
          ResourcePath: "/*"
          ThrottlingRateLimit: 10
          ThrottlingBurstLimit: 1000
```
</details>

<details>
<summary>Serverless Framework</summary>

```yaml
resources:
  Resources:
    MyApi:
      Type: AWS::Serverless::Api
      Properties:
        DefinitionUri: openapi.yaml

        # Throttling for default methods by setting HttpMethod  to '*' and
        # ResourcePath to '/*'
        MethodSettings:
          - HttpMethod: "*"
            ResourcePath: "/*"
            ThrottlingRateLimit: 10
            ThrottlingBurstLimit: 1000
```
</details>

<details>
<summary>Terraform</summary>

```hcl
resource "aws_api_gateway_stage" "this" {
  body = filename("openapi.yaml") 
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
</details>

### Implementations for HTTP APIs

<details>
<summary>CDK</summary>

_TODO_

</details>

<details>
<summary>CloudFormation/SAM</summary>

__JSON__

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

__YAML__

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
</details>

<details>
<summary>Serverless Framework</summary>

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
</details>

<details>
<summary>Terraform</summary>

```hcl
resource "aws_apigatewayv2_api" "this" {
  name          = "my-api"
  protocol_type = "HTTP"
  body          = filename("openapi.yaml") 
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
</details>

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

<details>
<summary>CDK</summary>

```typescript
export class MyStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    
  }
}
```
</details>

<details>
<summary>CloudFormation/SAM</summary>

__JSON__

```json
```

__YAML__

```yaml
```
</details>

<details>
<summary>Serverless Framework</summary>

```yaml
```
</details>

<details>
<summary>Terraform</summary>

```hcl
```
</details>

### Implementations for HTTP APIs

<details>
<summary>CDK</summary>

```typescript
export class MyStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    
  }
}
```
</details>

<details>
<summary>CloudFormation/SAM</summary>

__JSON__

```json
```

__YAML__

```yaml
```
</details>

<details>
<summary>Serverless Framework</summary>

```yaml
```
</details>

<details>
<summary>Terraform</summary>

```hcl
```
</details>

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

The rule in `serverless-rules` only check if the log structured is JSON-formatted.

While CloudWatch Logs Insights will automatically discover fields in JSON log entries, you can use the `parse` command to parse custom log entries to extract fields from custom format.

### Implementations for REST APIs

<details>
<summary>CDK</summary>

```typescript
export class MyStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    
  }
}
```
</details>

<details>
<summary>CloudFormation/SAM</summary>

__JSON__

```json
```

__YAML__

```yaml
```
</details>

<details>
<summary>Serverless Framework</summary>

```yaml
```
</details>

<details>
<summary>Terraform</summary>

```hcl
```
</details>

### Implementations for HTTP APIs

<details>
<summary>CDK</summary>

```typescript
export class MyStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    
  }
}
```
</details>

<details>
<summary>CloudFormation/SAM</summary>

__JSON__

```json
```

__YAML__

```yaml
```
</details>

<details>
<summary>Serverless Framework</summary>

```yaml
```
</details>

<details>
<summary>Terraform</summary>

```hcl
```
</details>

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

<details>
<summary>CDK</summary>

```typescript
export class MyStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    
  }
}
```
</details>

<details>
<summary>CloudFormation/SAM</summary>

__JSON__

```json
```

__YAML__

```yaml
```
</details>

<details>
<summary>Serverless Framework</summary>

```yaml
```
</details>

<details>
<summary>Terraform</summary>

```hcl
```
</details>

### Implementations for HTTP APIs

<details>
<summary>CDK</summary>

```typescript
export class MyStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    
  }
}
```
</details>

<details>
<summary>CloudFormation/SAM</summary>

__JSON__

```json
```

__YAML__

```yaml
```
</details>

<details>
<summary>Serverless Framework</summary>

```yaml
```
</details>

<details>
<summary>Terraform</summary>

```hcl
```
</details>

### Why is this a warning?

You might use [third party solutions](https://aws.amazon.com/lambda/partners/) for monitoring serverless applications. If this is the case, enabling tracing for API Gateway might be optional. Refer to the documentation of your monitoring solutions to see if you should enable AWS X-Ray tracing or not.

### See also

* [Serverless Lens: Distributed Tracing](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/distributed-tracing.html)
* [Tracing user requests to REST APIs using X-Ray](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-xray.html)