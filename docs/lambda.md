AWS Lambda Rules
================

## EventSourceMapping Failure Destination

* __Level__: Error
* __cfn-lint__: ES1001
* __tflint__: _Not implemented_

An AWS Lambda event source mapping is used to read from streams and poll-based event sources. You can configure the event source mapping to send invocation records to another services such as Amazon SNS or Amazon SQS when it discards an event batch.

__See:__

* [AWS Lambda event source mappings](https://docs.aws.amazon.com/lambda/latest/dg/invocation-eventsourcemapping.html)

## Log Retention

* __Level__: Warning
* __cfn-lint__: WS1004
* __tflint__: _Not implemented_

By default, CloudWatch log groups created by Lambda functions have an unlimited retention time. For cost optimization purposes, you should set a retention duration on all log groups. For log archival, export and set cost-effective storage classes that best suit your needs.

__Why is this a warning?__

Since `serverless-rules` evaluate infrastructure as code template, it cannot check if you use a solution that will automatically change the configuration of log groups after the fact.

__See:__

* [Serverless Lens: Logging Ingestion and Storage](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/logging-ingestion-and-storage.html)

## Permission Multiple Principals

* __Level__: Warning
* __cfn-lint__: WS1002
* __tflint__: _Not implemented_

You can use resource-based policies to grant permission to other AWS services to invoke your Lambda functions. Different AWS services usually send different payloads to Lambda functions. If multiple services can invoke the same function, this function needs to handle the different types of payload properly, or this could cause unexpected behavior.

In general, it's better to deploy multiple Lambda functions with different function handlers for each invocation source.

__Why is this a warning?__

You might have a valid reason for invoking a Lambda function from different event sources or AWS services. If this is the case and you know what you are doing, you might ignore this rule.

__See:__
* [Using resource-based policies for AWS Lambda](https://docs.aws.amazon.com/lambda/latest/dg/access-control-resource-based.html)

## Star Permissions

* __Level__: Warning
* __cfn-lint__: WS1003
* __tflint__: _Not implemented_

 With Lambda functions, it’s recommended that you follow least-privileged access and only allow the access needed to perform a given operation. Attaching a role with more permissions than necessary can open up your systems for abuse.

 __Why is this a warning?__

If your Lambda function need a broad range of permissions, you do not know ahead of time which permissions you will need, and you have evaluated the risks of using broad permissions for this function, you might ignore this rule.

__See:__
* [Serverless Lens: Identity and Access Management](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/identity-and-access-management.html)
* [AWS Lambda execution role](https://docs.aws.amazon.com/lambda/latest/dg/lambda-intro-execution-role.html)

## Tracing

* __Level__: Warning
* __cfn-lint__: WS1000
* __tflint__: aws_lambda_function_tracing_rule

AWS Lambda can emit traces to AWS X-Ray, which enable visualizing service maps for faster troubleshooting.

__Why is this a warning?__

You might use [third party solutions](https://aws.amazon.com/lambda/partners/) for monitoring serverless applications. If this is the case, enabling tracing for your AWS Lambda functions might be optional. Refer to the documentation of your monitoring solutions to see if you should enable AWS X-Ray tracing or not.

__See:__

* [Serverless Lens: Distributed Tracing](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/distributed-tracing.html)
* [Using AWS Lambda with X-Ray](https://docs.aws.amazon.com/lambda/latest/dg/services-xray.html)