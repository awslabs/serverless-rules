Rules
=====

## AWS Lambda

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
| Warning | [Lambda Tracing](lambda.md#tracing)                                 | WS1000   | Y      |
| Error   | Lambda Deprecated Runtime                                           |          |        |
| Error   | Lambda No Error Alarm                                               |          |        |
| Warning | Lambda Star Permissions                                             |          |        |
| Warning | Lambda Multiple Triggers                                            |          |        |
| Warning | Lambda LogGroup No Retention                                        |          |        |
| Error   | Async Lambda No Failure Destination                                 |          |        |
| Warning | Sync Lambda No Duration Alarm                                       |          |        |
| Warning | Sync Lambda With Destination                                        |          |        |
| Error   | EventSourceMapping No Failure Destination                           |          |        |
| Error   | SQS Lambda ReservedConcurrency < 5                                  |          |        |

## Amazon API Gateway REST APIs

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
| Error   | [API Gateway Logging](api_gateway.md#logging)                       | ES2000   | Y      |
| Warning | [API Gateway Structured Logging](api_gateway.md#structured-logging) | WS2001   |        |
| Warning | [API Gateway Tracing](api_gateway.md#tracing)                       | WS2002   | Y      |
| Warning | [API Gateway Default Throttling](api_gateway.md#default-throttling) | ES2003   |        |

## Amazon API Gateway HTTP APIs

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
| Error   | [API Gateway Logging](api_gateway.md#logging)                       | ES2000   | Y      |
| Warning | [API Gateway Structured Logging](api_gateway.md#structured-logging) | WS2001   |        |
| Warning | [API Gateway Default Throttling](api_gateway.md#default-throttling) | ES2003   |        |

## AWS AppSync

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
| Error   | AppSync Tracing                                                     |          |        |

## Amazon EventBridge

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
| Warning | EventBridge Rule No DLQ                                             |          |        |
| Info    | EventBridge Bus No Rule                                             |          |        |

## Amazon Step Functions

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
| Error   | Step Functions Tracing                                              |          |        |

## Amazon SQS

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|
| Warning | SQS Queue No DLQ                                                    |          |        |

## Amazon SNS

| Level   | Name                                                                | cfn-lint | tflint |
|---------|---------------------------------------------------------------------|----------|--------|