# EventBridge Rule without DLQ

__Level__: Error
{: class="badge badge-red" }

__Initial version__: 0.1.3
{: class="badge badge-blue" }

__cfn-lint__: ES4000
{: class="badge" }

__tflint__: aws_cloudwatch_event_target_no_dlq
{: class="badge" }

Sometimes, an event isn't successfully delivered to the target(s) specified in a rule. By default, EventBridge will retry for 24 hours and up to 185 times, but you can customize the retry policy.

If EventBridge cannot deliver an event after all its retries, it can send it to a dead-letter queue. You can then inspect the event and remediate the underlying issue.

## Implementations

=== "CDK"

    ```typescript
    import { Function } from '@aws-cdk/aws-lambda';
    import { Rule } from '@aws-cdk/aws-events';
    import * as targets from '@aws-cdk/aws-events-targets';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const myFunction = new Function(
          scope, 'MyFunction',
          {
            code: Code.fromAsset('src/hello/'),
            handler: 'main.handler',
            runtime: Runtime.PYTHON_3_8,
          }
        );

        const myRule = new Rule(
          scope, 'MyRule',
          {
            eventPattern: {
              source: ['my-source'],
            }
          }
        );

        myRule.addTarget(new targets.LambdaFunction(
          myfunction,
          // Add a DLQ to the 'myFunction' target
          {
            deadLetterQueue: myQueue,
          }
        ));
      }
    }
    ```

=== "CloudFormation (JSON)"

    ```json
    {
      "Resources": {
        "MyRule": {
          "Type": "AWS::Events::Rule",
          "Properties": {
            "EventBusName": "default",
            "EventPattern": "{\"source\": [\"my-source\"]}",
            "Targets": [{
              "Id": "MyFunction",
              "Arn": "arn:aws:lambda:us-east-1:111122223333:function:MyFunction",
              // Add a DLQ to the 'MyFunction' target
              "DeadLetterConfig": {
                "Arn": "arn:aws:sqs:us-east-1:111122223333:dlq"
              }
            }]
          }
        }
      }
    }
    ```

=== "CloudFormation (YAML)"

    ```yaml
    Resources:
      MyRule:
        Type: AWS::Events::Rule
        Properties:
          EventBusName: default
          EventPattern: |
            {
              "source": ["my-source"]
            }
          Targets:
            - Id: MyFunction
              Arn: arn:aws:lambda:us-east-1:111122223333:function:MyFunction
              # Add a DLQ to the 'MyFunction' target
              DeadLetterConfig:
                Arn: arn:aws:sqs:us-east-1:111122223333:dlq
    ```

=== "Serverless Framework"

    ```yaml
    resources:
      Resources:
        MyRule:
          Type: AWS::Events::Rule
          Properties:
            EventBusName: default
            EventPattern: |
              {
                "source": ["my-source"]
              }
            Targets:
              - Id: MyFunction
                Arn: arn:aws:lambda:us-east-1:111122223333:function:MyFunction
                # Add a DLQ to the 'MyFunction' target
                DeadLetterConfig:
                  Arn: arn:aws:sqs:us-east-1:111122223333:dlq
    ```

=== "Terraform"

    ```tf
    resource "aws_cloudwatch_event_rule" "this" {
      event_pattern = <<EOF
    {
      "source": ["my-source"]
    }
    EOF
    }

    resource "aws_cloudwatch_event_target" "this" {
      rule      = aws_cloudwatch_event_rule.this.name
      target_id = "MyFunction"
      arn       = "arn:aws:lambda:us-east-1:111122223333:function:MyFunction"

      # Add a DLQ to the 'MyFunction' target
      dead_letter_config {
        arn = "arn:aws:sqs:us-east-1:111122223333:dlq"
      }
    }
    ```

## See also

* [Event retry policy and using dead-letter queues](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-rule-dlq.html)