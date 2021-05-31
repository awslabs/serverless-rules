Amazon API Gateway Rules
========================

## Default Throttling

* __Level__: Error
* __cfn-lint__: ES2003
* __tflint (REST APIs)__: aws_apigateway_stage_throttling_rule
* __tflint (HTTP APIs)__: aws_apigatewayv2_stage_throttling_rule

Amazon API Gateway supports defining default limits for an API to prevent it from being overwhelmed by too many requests. This uses a [token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket), where a token counts for a single request.

__See:__

* [Throttle API requests for better throughput](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-request-throttling.html)
* [Throttling requests to your HTTP API](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-throttling.html)

## Logging

* __Level__: Error
* __cfn-lint__: ES2000
* __tflint (REST APIs)__: aws_apigateway_stage_logging_rule
* __tflint (HTTP APIs)__: aws_apigatewayv2_stage_logging_rule

Amazon API Gateway can send logs to Amazon CloudWatch Logs and Amazon Kinesis Data Firehose for centralization.

__See:__

* [Serverless Lens: Centralized and structured logging](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/centralized-and-structured-logging.html)
* [Monitoring REST APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/rest-api-monitor.html)
* [Monitoring your HTTP API](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-monitor.html)
* [Monitoring WebSocket APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/websocket-api-monitor.html)

## Structured Logging

* __Level__: Error
* __cfn-lint__: WS2001
* __tflint__: _Not implemented_

You can customize the log format that Amazon API Gateway uses to send logs. Structured logging makes it easier to derive queries to answer arbitrary questions about the health of your application.

__Why is this a warning?__

The rule in `serverless-rules` only check if the log structured is JSON-formatted.

While CloudWatch Logs Insights will automatically discover fields in JSON log entries, you can use the `parse` command to parse custom log entries to extract fields from custom format.

__See:__

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

__Why is this a warning?__

You might use [third party solutions](https://aws.amazon.com/lambda/partners/) for monitoring serverless applications. If this is the case, enabling tracing for API Gateway might be optional. Refer to the documentation of your monitoring solutions to see if you should enable AWS X-Ray tracing or not.

__See:__

* [Serverless Lens: Distributed Tracing](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/distributed-tracing.html)
* [Tracing user requests to REST APIs using X-Ray](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-xray.html)