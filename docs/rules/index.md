Rules
=====

## Rule levels

A rule can have one of the following three rule levels: Error, Warning, or Info.

An __Error__ level for a rule means this is a recommended practice for the vast majority of circumstances.

A __Warning__ level means that this is a recommended practice, but you can achieve similar results through a different implementation. For example, you can create alarms through [third party offering](https://aws.amazon.com/lambda/partners/), rather than using AWS CloudWatch.

An __Info__ level means that this does not necessarily align with recommended practices but can point out potential issues or misconfiguration. For example, an Amazon EventBridge event bus without any rules associated with it, as you might create those rules through a different template.

## AWS Lambda

| Level       | Name                                                                | cfn-lint | tflint |
|:-----------:|---------------------------------------------------------------------|:--------:|:------:|
| __Warning__ | [Lambda Tracing](lambda.md#tracing)                                 | WS1000   | aws_lambda_function_tracing_rule |
| __Error__   | [EventSourceMapping Failure Destination](lambda.md#eventsourcemapping-failure-destination) | ES1001   | aws_lambda_event_source_mapping_failure_destination |
| __Warning__ | [Lambda Permission Multiple Principals](lambda.md#permission-multiple-principals) | WS1002   |_Not implemented_|
| __Warning__ | [Lambda Star Permissions](lambda.md#star-permissions) | WS1003   |_Not implemented_|
| __Warning__ | [Lambda Log Retention](lambda.md#log-retention) | WS1004   |_Not implemented_|
| __Error__   | Lambda Deprecated Runtime                                           |_Not implemented_|_Not implemented_|
| __Error__   | Lambda No Error Alarm                                               |_Not implemented_|_Not implemented_|
| __Error__   | Async Lambda No Failure Destination                                 |_Not implemented_|_Not implemented_|
| __Warning__ | Sync Lambda No Duration Alarm                                       |_Not implemented_|_Not implemented_|
| __Warning__ | Sync Lambda With Destination                                        |_Not implemented_|_Not implemented_|
| __Error__   | SQS Lambda ReservedConcurrency < 5                                  |_Not implemented_|_Not implemented_|

## Amazon API Gateway REST APIs

| Level       | Name                                                                | cfn-lint | tflint |
|:-----------:|---------------------------------------------------------------------|:--------:|:------:|
| __Error__   | [API Gateway Logging](api_gateway.md#logging)                       | ES2000   | aws_apigateway_stage_logging_rule |
| __Warning__ | [API Gateway Structured Logging](api_gateway.md#structured-logging) | WS2001   |_Not implemented_|
| __Warning__ | [API Gateway Tracing](api_gateway.md#tracing)                       | WS2002   | aws_apigateway_stage_tracing_rule |
| __Warning__ | [API Gateway Default Throttling](api_gateway.md#default-throttling) | ES2003   | aws_apigateway_stage_throttling_rule |

## Amazon API Gateway HTTP APIs

| Level       | Name                                                                | cfn-lint | tflint |
|:-----------:|---------------------------------------------------------------------|:--------:|:------:|
| __Error__   | [API Gateway Logging](api_gateway.md#logging)                       | ES2000   | aws_apigatewayv2_stage_logging_rule |
| __Warning__ | [API Gateway Structured Logging](api_gateway.md#structured-logging) | WS2001   |_Not implemented_|
| __Warning__ | [API Gateway Default Throttling](api_gateway.md#default-throttling) | ES2003   | aws_apigatewayv2_stage_throttling_rule |

## AWS AppSync

| Level       | Name                                                                | cfn-lint | tflint |
|:-----------:|---------------------------------------------------------------------|:--------:|:------:|
| __Error__   | [AppSync Tracing](appsync.md#tracing)                               | WS3000   | aws_appsync_graphql_api_tracing_rule |

## Amazon EventBridge

| Level       | Name                                                                | cfn-lint | tflint |
|:-----------:|---------------------------------------------------------------------|:--------:|:------:|
| __Error__   | [EventBridge Rule No DLQ](eventbridge.md#rule-without-dlq)          | ES4000   | aws_cloudwatch_event_target_no_dlq |
| __Info__    | EventBridge Bus No Rule                                             |_Not implemented_|_Not implemented_|

## Amazon Step Functions

| Level       | Name                                                                | cfn-lint | tflint |
|:-----------:|---------------------------------------------------------------------|:--------:|:------:|
| __Warning__ | [Step Functions Tracing](step_functions.md#tracing) | WS5000   | aws_sfn_state_machine_tracing |

## Amazon SQS

| Level       | Name                                                                | cfn-lint | tflint |
|:-----------:|---------------------------------------------------------------------|:--------:|:------:|
| __Error__   | [SQS Queue No Redrive Policy](sqs.md#no-redrive-policy)             | ES6000   | aws_sqs_queue_redrive_policy |

## Amazon SNS

| Level       | Name                                                                | cfn-lint | tflint |
|:-----------:|---------------------------------------------------------------------|:--------:|:------:|
| __Error__   | SNS Topic No DLQ                                                    |_Not implemented_|_Not implemented_|