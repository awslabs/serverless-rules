AWS AppSync Rules
=================

## Tracing

* __Level__: Warning
* __cfn-lint__: WS3000
* __tflint__: aws_appsync_graphql_api_tracing_rule

AWS AppSync can emit traces to AWS X-Ray, which enable visualizing service maps for faster troubleshooting.

__Why is this a warning?__

You might use [third party solutions](https://aws.amazon.com/lambda/partners/) for monitoring serverless applications. If this is the case, enabling tracing for AppSync APIs might be optional. Refer to the documentation of your monitoring solutions to see if you should enable AWS X-Ray tracing or not.

__See:__

* [Serverless Lens: Distributed Tracing](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/distributed-tracing.html)
* [Tracing with AWS X-Ray](https://docs.aws.amazon.com/appsync/latest/devguide/x-ray-tracing.html)