# Lambda EventSourceMapping Failure Destination

__Level__: Error
{: class="badge badge-red" }

__Initial version__: 0.1.6
{: class="badge badge-blue" }

__cfn-lint__: ES1001
{: class="badge" }

__tflint__: aws_lambda_event_source_mapping_failure_destination
{: class="badge" }

An AWS Lambda event source mapping reads from streams and poll-based event sources to invoke your functions. You can configure the event source mapping to send invocation records to another service such as Amazon SNS or Amazon SQS when it discards an event batch.

## Implementations

=== "CDK"

    ```typescript
    import { EventSourceMapping, SqsDlq, StartingPosition } from '@aws-cdk/aws-lambda';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        new EventSourceMapping(scope, "MyEventSourceMapping", {
          target: myFunction,
          eventSourceArn: 'arn:aws:dynamodb:us-east-1:111122223333:table/my-table/stream/my-stream',
          startingPosition: StartingPosition.LATEST,
          onFailure: SqsDlq(mySqsQueue),
        });
      }
    }
    ```

=== "CloudFormation (JSON)"

    ```json
    {
      "Resource": {
        "MyEventSourceMapping": {
          "Type": "AWS::Lambda::EventSourceMapping",
          "Properties": {
            // Required properties
            "FunctionName": "my-function",
            "EventSourceArn": "arn:aws:dynamodb:us-east-1:111122223333:table/my-table/stream/my-stream",
            "StartingPosition": "LATEST",

            // Add an OnFailure destination on the event source mapping
            "DestinationConfig": {
              "OnFailure": {
                "Destination": "arn:aws:sqs:us-east-1:111122223333:my-dlq"
              }
            }
          }
        }
      }
    }
    ```

=== "CloudFormation (YAML)"

    ```yaml
    Resources:
      MyEventSourceMapping:
        Type: AWS::Lambda::EventSourceMapping
        Properties:
          # Required properties
          FunctionName: my-function
          EventSourceArn: arn:aws:dynamodb:us-east-1:111122223333:table/my-table/stream/my-stream
          StartingPosition: LATEST

          # Add an OnFailure destination on the event source mapping
          DestinationConfig:
            OnFailure:
              Destination: arn:aws:sqs:us-east-1:111122223333:my-dlq 
    ```

=== "Serverless Framework"

    ```yaml
    functions:
      MyFunction:
        handler: hello.handler

    resources:
      Resources:
        MyEventSourceMapping:
          Type: AWS::Lambda::EventSourceMapping
          Properties:
            # Required properties
            FunctionName:
              Fn::Ref: MyFunction
            EventSourceArn: arn:aws:dynamodb:us-east-1:111122223333:table/my-table/stream/my-stream
            StartingPosition: LATEST

            # Add an OnFailure destination on the event source mapping
            DestinationConfig:
              OnFailure:
                Destination: arn:aws:sqs:us-east-1:111122223333:my-dlq 
    ```

=== "Terraform"

    ```tf
    resource "aws_lambda_event_source_mapping" "this" {
      # Required fields
      event_source_arn  = "arn:aws:dynamodb:us-east-1:111122223333:table/my-table/stream/my-stream"
      function_name     = "my-function"
      starting_position = "LATEST"

      # Add an OnFailure destination on the event source mapping
      destination_config {
        on_failure {
          destination_arn = "arn:aws:sqs:us-east-1:111122223333:my-dlq"
        }
      }
    }
    ```


## See also

* [AWS Lambda event source mappings](https://docs.aws.amazon.com/lambda/latest/dg/invocation-eventsourcemapping.html)
* [__CloudFormation__: AWS::Lambda::EventSourceMapping](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-eventsourcemapping.html#cfn-lambda-eventsourcemapping-destinationconfig)
* [__Terraform__: aws_lambda_event_source_mapping](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_event_source_mapping)