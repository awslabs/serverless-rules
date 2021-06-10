Rules
=====

## Rule levels

A rule can have one of the following three rule levels: Error, Warning, or Info.

An __Error__{: class="badge badge-red" } level for a rule means this is a recommended practice for the vast majority of circumstances.

A __Warning__{: class="badge badge-yellow" } level means that this is a recommended practice, but you can achieve similar results through a different implementation. For example, you can create alarms through [third party offering](https://aws.amazon.com/lambda/partners/), rather than using AWS CloudWatch.

An __Info__{: class="badge badge-blue" } level means that this does not necessarily align with recommended practices but can point out potential issues or misconfiguration. For example, an Amazon EventBridge event bus without any rules associated with it, as you might create those rules through a different template.

## AWS Lambda

| Level                                      | Name                                                                | cfn-lint | tflint |
|:------------------------------------------:|---------------------------------------------------------------------|:--------:|:------:|
| __Warning__{: class="badge badge-yellow" } | [Lambda Tracing](lambda/tracing.md)                                 | WS1000   | aws_lambda_function_tracing_rule |
| __Error__{: class="badge badge-red" }      | [EventSourceMapping Failure Destination](lambda/eventsourcemapping_failure_destination.md) | ES1001   | aws_lambda_event_source_mapping_failure_destination |
| __Warning__{: class="badge badge-yellow" } | [Lambda Permission Multiple Principals](lambda/permission_multiple_principals.md) | WS1002   |_Not implemented_|
| __Warning__{: class="badge badge-yellow" } | [Lambda Star Permissions](lambda/star_permissions.md)               | WS1003   |_Not implemented_|
| __Warning__{: class="badge badge-yellow" } | [Lambda Log Retention](lambda/log_retention.md)                     | WS1004   |_Not implemented_|
| __Error__{: class="badge badge-red" }      | [Lambda Default Memory Size](lambda/default_memory_size.md)         | ES1005   | aws_lambda_function_default_memory |
| __Error__{: class="badge badge-red" }      | [Lambda Default Timeout](lambda/default_timeout.md)                 | ES1006   | aws_lambda_function_default_timeout |
| __Error__{: class="badge badge-red" }      | [Async Lambda Failure Destination](lambda/async_failure_destination.md) | ES1007 |_Not implemented_|
| __Error__{: class="badge badge-red" }      | [Lambda EOL Runtime](lambda/end_of_life_runtime.md)                 | _E2531_  | aws_lambda_function_eol_runtime |
| __Error__{: class="badge badge-red" }      | Lambda No Error Alarm                                               |_Not implemented_|_Not implemented_|
| __Warning__{: class="badge badge-yellow" } | Sync Lambda No Duration Alarm                                       |_Not implemented_|_Not implemented_|
| __Warning__{: class="badge badge-yellow" } | Sync Lambda With Destination                                        |_Not implemented_|_Not implemented_|
| __Error__{: class="badge badge-red" }      | SQS Lambda ReservedConcurrency < 5                                  |_Not implemented_|_Not implemented_|

## Amazon API Gateway REST APIs

| Level                                      | Name                                                                | cfn-lint | tflint |
|:------------------------------------------:|---------------------------------------------------------------------|:--------:|:------:|
| __Error__{: class="badge badge-red" }      | [API Gateway Logging](api_gateway/logging.md)                       | ES2000   | aws_apigateway_stage_logging_rule |
| __Warning__{: class="badge badge-yellow" } | [API Gateway Structured Logging](api_gateway/structured_logging.md) | WS2001   | aws_api_gateway_stage_structured_logging |
| __Warning__{: class="badge badge-yellow" } | [API Gateway Tracing](api_gateway/tracing.md)                       | WS2002   | aws_apigateway_stage_tracing_rule |
| __Warning__{: class="badge badge-yellow" } | [API Gateway Default Throttling](api_gateway/default_throttling.md) | ES2003   | aws_apigateway_stage_throttling_rule |

## Amazon API Gateway HTTP APIs

| Level                                      | Name                                                                | cfn-lint | tflint |
|:------------------------------------------:|---------------------------------------------------------------------|:--------:|:------:|
| __Error__{: class="badge badge-red" }      | [API Gateway Logging](api_gateway/logging.md)                       | ES2000   | aws_apigatewayv2_stage_logging_rule |
| __Warning__{: class="badge badge-yellow" } | [API Gateway Structured Logging](api_gateway/structured_logging.md) | WS2001   | aws_apigatewayv2_stage_structured_logging |
| __Warning__{: class="badge badge-yellow" } | [API Gateway Default Throttling](api_gateway/default_throttling.md) | ES2003   | aws_apigatewayv2_stage_throttling_rule |

## AWS AppSync

| Level                                      | Name                                                                | cfn-lint | tflint |
|:------------------------------------------:|---------------------------------------------------------------------|:--------:|:------:|
| __Error__{: class="badge badge-red" }      | [AppSync Tracing](appsync/tracing.md)                               | WS3000   | aws_appsync_graphql_api_tracing_rule |

## Amazon EventBridge

| Level                                      | Name                                                                | cfn-lint | tflint |
|:------------------------------------------:|---------------------------------------------------------------------|:--------:|:------:|
| __Error__{: class="badge badge-red" }      | [EventBridge Rule Without DLQ](eventbridge/rule_without_dlq.md)     | ES4000   | aws_cloudwatch_event_target_no_dlq |
| __Info__{: class="badge badge-blue" }      | EventBridge Bus No Rule                                             |_Not implemented_|_Not implemented_|

## Amazon SNS

| Level                                      | Name                                                                | cfn-lint | tflint |
|:------------------------------------------:|---------------------------------------------------------------------|:--------:|:------:|
| __Error__{: class="badge badge-red" }      | [SNS Redrive Policy](sns/redrive_policy.md)                         | ES7000 | aws_sns_topic_subscription_redrive_policy |

## Amazon SQS

| Level                                      | Name                                                                | cfn-lint | tflint |
|:------------------------------------------:|---------------------------------------------------------------------|:--------:|:------:|
| __Error__{: class="badge badge-red" }      | [SQS Redrive Policy](sqs/redrive_policy.md)                         | ES6000   | aws_sqs_queue_redrive_policy |

## Amazon Step Functions

| Level                                      | Name                                                                | cfn-lint | tflint |
|:------------------------------------------:|---------------------------------------------------------------------|:--------:|:------:|
| __Warning__{: class="badge badge-yellow" } | [Step Functions Tracing](step_functions/tracing.md)                 | WS5000   | aws_sfn_state_machine_tracing |