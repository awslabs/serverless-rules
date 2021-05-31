Amazon EventBridge Rules
========================

## Rule without DLQ

* __Level__: Error
* __cfn-lint__: ES4000
* __tflint__: _Not implemented_

Sometimes, an event isn't successfully delivered to the target(s) specified in a rule. By default, EventBridge will retry for 24 hours and up to 185 times, but you can customize the retry policy.

If EventBridge was not able to deliver the event after all retries, it can send that event to a dead-letter queue to prevent the loss of that event, and allow you to inspect and remediate the underlying issue.

__See:__

* [Event retry policy and using dead-letter queues](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-rule-dlq.html)