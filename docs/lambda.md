AWS Lambda Rules
================

## EventSourceMapping Failure Destination

An AWS Lambda event source mapping is used to read from streams and poll-based event sources. You can configure the event source mapping to send invocation records to another services such as Amazon SNS or Amazon SQS when it discards an event batch.

__See:__

* [AWS Lambda event source mappings](https://docs.aws.amazon.com/lambda/latest/dg/invocation-eventsourcemapping.html)

## Tracing

AWS Lambda can emit traces to AWS X-Ray, which enable visualizing service maps for faster troubleshooting.

__See:__

* [Serverless Lens: Distributed Tracing](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/distributed-tracing.html)
* [Using AWS Lambda with X-Ray](https://docs.aws.amazon.com/lambda/latest/dg/services-xray.html)