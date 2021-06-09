# API Gateway Structured Logging

__Level__: Warning
{: class="badge badge-yellow" }

__Initial version__: 0.1.3
{: class="badge badge-blue" }

__cfn-lint__: WS2001
{: class="badge" }

__tflint (REST)__: aws_api_gateway_stage_structured_logging
{: class="badge" }

__tflint (HTTP)__: aws_apigatewayv2_stage_structured_logging
{: class="badge" }

You can customize the log format that Amazon API Gateway uses to send logs. Structured logging makes it easier to derive queries to answer arbitrary questions about the health of your application.

## Why is this a warning?

The rule in `serverless-rules` only checks if the structured log is JSON-formatted.

While CloudWatch Logs Insights will automatically discover fields in JSON log entries, you can use the `parse` command to parse custom log entries to extract fields from custom format.

## Implementations

See the implementations for [Logging on API Gateway](logging.md).

## See also

* [Serverless Lens: Centralized and structured logging](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/centralized-and-structured-logging.html)
* [Monitoring REST APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/rest-api-monitor.html)
* [Monitoring your HTTP API](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-monitor.html)
* [Monitoring WebSocket APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/websocket-api-monitor.html)
* [Amazon CloudWatch Logs: Supported Logs and Discovered Fields](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CWL_AnalyzeLogData-discoverable-fields.html)
* [Amazon CloudWatch Logs: Logs Insights Query Syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CWL_QuerySyntax.html)