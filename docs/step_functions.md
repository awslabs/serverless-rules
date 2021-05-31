AWS Step Functions Rules
========================

## Tracing

* __Level__: Warning
* __cfn-lint__: WS5000
* __tflint__: _Not implemented_

AWS Step Functions can emit traces to AWS X-Ray, which enable visualizing service maps for faster troubleshooting.

__Why is this a warning?__

You might use [third party solutions](https://aws.amazon.com/lambda/partners/) for monitoring serverless applications. If this is the case, enabling tracing for Step Functions might be optional. Refer to the documentation of your monitoring solutions to see if you should enable AWS X-Ray tracing or not.

__See:__

* [Serverless Lens: Distributed Tracing](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/distributed-tracing.html)
* [AWS X-Ray and Step Functions](https://docs.aws.amazon.com/step-functions/latest/dg/concepts-xray-tracing.html)