AWS Step Functions Rules
========================

## Tracing

* __Level__: Warning
* __cfn-lint__: WS5000
* __tflint__: _Not implemented_

AWS Step Functions can emit traces to AWS X-Ray, which enable visualizing service maps for faster troubleshooting.

### Why is this a warning?

You might use [third party solutions](https://aws.amazon.com/lambda/partners/) for monitoring serverless applications. If this is the case, enabling tracing for Step Functions might be optional. Refer to the documentation of your monitoring solutions to see if you should enable AWS X-Ray tracing or not.

### Implementations

<details>
<summary>CDK</summary>

```typescript
import { Function } from '@aws-cdk/aws-lambda';
import { StateMachine } from '@aws-cdk/aws-stepfunctions';
import * as tasks from '@aws-cdk/aws-stepfunctions-

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

    const myJob = tasks.LambdaInvoke(
      scope, 'MyJob',
      {
        lambdaFunction: myFunction,
      },
    );

    const myStateMachine = new StateMachine(
      scope, 'MyStateMachine',
      {
        definition: myJob,
        // Enable tracing on Step Functions
        tracingEnabled: true,
      }
    );
  }
}
```
</details>

<details>
<summary>CloudFormation/SAM</summary>

__JSON__

```json
{
  "Resources": {
    "StateMachine": {
      "Type": "AWS::StepFunctions::StateMachine",
      "Properties": {
        "DefinitionString": "{ \"StartAt\": \"HelloWorld\", \"States\": { \"HelloWorld\": { \"Type\": \"Task\", \"Resource\": \"arn:aws:lambda:us-east-1:111122223333:function:HelloFunction\", \"End\": \"true\" }}}",
        "RoleArn": "arn:aws:iam::111122223333:role/service-role/StatesExecutionRole",

        // Enable active tracing for Step Functions
        "TracingConfiguration": {
          "Enabled": true
        }
      }
    }
  }
}
```

__YAML__

```yaml
Resources:
  StateMachine:
    Type: AWS::StepFunctions::StateMachine
    Properties:
      DefinitionString: |
        {
          "StartAt": "HelloWorld",
          "States": {
            "HelloWorld": {
              "Type": "Task",
              "Resource": "arn:aws:lambda:us-east-1:111122223333:function:HelloFunction",
              "End": "true"
            }
          }
        }
      RoleArn: arn:aws:iam::111122223333:role/service-role/StatesExecutionRole

      # Enable active tracing for Step Functions
      TracingConfiguration:
        Enabled: true
```

</details>

<details>
<summary>Serverless Framework</summary>

```yaml
resources:
  Resources:
    StateMachine:
      Type: AWS::StepFunctions::StateMachine
      Properties:
        DefinitionString: |
          {
            "StartAt": "HelloWorld",
            "States": {
              "HelloWorld": {
                "Type": "Task",
                "Resource": "arn:aws:lambda:us-east-1:111122223333:function:HelloFunction",
                "End": "true"
              }
            }
          }
        RoleArn: arn:aws:iam::111122223333:role/service-role/StatesExecutionRole

        # Enable active tracing for Step Functions
        TracingConfiguration:
          Enabled: true
```
</details>

<details>
<summary>Terraform</summary>

```hcl
resource "aws_sfn_state_machine" "this" {
  name     = "my-state-machine"
  role_arn = "arn:aws:iam::111122223333:role/my-state-machine-role"

  definition = <<EOF
{
  "StartAt": "HelloWorld",
  "States": {
    "HelloWorld": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:us-east-1:111122223333:function:HelloFunction",
      "End": "true"
    }
  }
}
EOF

  # Enable active tracing for Step Functions
  tracing_configuration {
    enabled = true
  }
}
```
</details>

### See also

* [Serverless Lens: Distributed Tracing](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/distributed-tracing.html)
* [AWS X-Ray and Step Functions](https://docs.aws.amazon.com/step-functions/latest/dg/concepts-xray-tracing.html)