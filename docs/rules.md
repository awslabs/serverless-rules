Rules
=====

## Rule levels

A rule can have one of the follow three rule levels: Error, Warning, or Info.

An __Error__ level for a rule means this is a recommended practice for the vast majority of circumstances.

A __Warning__ level means that this is a recommended practice, but you can achieve similar results through a different implementation. For example, you can create alarms through [third party offering](https://aws.amazon.com/lambda/partners/), rather than using AWS CloudWatch.

An __Info__ level means that this does not necessarily align with recommended practices, but can point out potential issues or misconfiguration. For example, an Amazon EventBridge event bus without any rules associated to it, as you might create those rules throught a different template.

## AWS Lambda

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
| Warning | [Lambda Tracing](lambda.md#tracing)                                 | WS1000   | Y      |
| Error   | [EventSourceMapping Failure Destination](lambda.md#eventsourcemapping-failure-destination) | ES1001   |        |
| Warning | [Lambda Permission Multiple Principals](lambda.md#permission-multiple-principals) | WS1002   |        |
| Warning | [Lambda Star Permissions](lambda.md#star-permissions) | WS1003   |        |
| Warning | [Lambda Log Retention](lambda.md#log-retention) | WS1004   |        |
| Error   | Lambda Deprecated Runtime                                           |          |        |
| Error   | Lambda No Error Alarm                                               |          |        |
| Error   | Async Lambda No Failure Destination                                 |          |        |
| Warning | Sync Lambda No Duration Alarm                                       |          |        |
| Warning | Sync Lambda With Destination                                        |          |        |
| Error   | SQS Lambda ReservedConcurrency < 5                                  |          |        |

## Amazon API Gateway REST APIs

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
| Error   | [API Gateway Logging](api_gateway.md#logging)                       | ES2000   | Y      |
| Warning | [API Gateway Structured Logging](api_gateway.md#structured-logging) | WS2001   |        |
| Warning | [API Gateway Tracing](api_gateway.md#tracing)                       | WS2002   | Y      |
| Warning | [API Gateway Default Throttling](api_gateway.md#default-throttling) | ES2003   | Y      |

## Amazon API Gateway HTTP APIs

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
| Error   | [API Gateway Logging](api_gateway.md#logging)                       | ES2000   | Y      |
| Warning | [API Gateway Structured Logging](api_gateway.md#structured-logging) | WS2001   |        |
| Warning | [API Gateway Default Throttling](api_gateway.md#default-throttling) | ES2003   | Y      |

## AWS AppSync

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
| Error   | [AppSync Tracing](appsync.md#tracing) | WS3000   | Y      |

## Amazon EventBridge

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
| Error   | [EventBridge Rule No DLQ](eventbridge.md#rule-without-dlq) | ES4000   |        |
| Info    | EventBridge Bus No Rule                                             |          |        |

## Amazon Step Functions

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
| Warning | [Step Functions Tracing](step_functions.md#tracing) | WS5000   |        |

## Amazon SQS

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
| Error   | SQS Queue No DLQ                                                    |          |        |

## Amazon SNS

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
