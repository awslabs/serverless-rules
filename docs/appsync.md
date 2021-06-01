AWS AppSync Rules
=================

## Tracing

* __Level__: Warning
* __cfn-lint__: WS3000
* __tflint__: aws_appsync_graphql_api_tracing_rule

AWS AppSync can emit traces to AWS X-Ray, which enable visualizing service maps for faster troubleshooting.

### Why is this a warning?

You might use [third party solutions](https://aws.amazon.com/lambda/partners/) for monitoring serverless applications. If this is the case, enabling tracing for AppSync APIs might be optional. Refer to the documentation of your monitoring solutions to see if you should enable AWS X-Ray tracing or not.

### Implementations

<details>
<summary>CDK</summary>

```typescript
import { GraphqlApi } from '@aws-cdk/aws-appsync';

export class MyStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    myApi = new GraphqlApi(
      scope, 'MyApi',
      {
        name: 'my-api',
        // Enable active tracing
        xrayEnabled: true,
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
    "GraphQLApi": {
    "Type": "AWS::AppSync::GraphQLApi",
    "Properties": {
      "Name": "api",
      "AuthenticationType": "AWS_IAM",

      // Enable active tracing
      "XrayEnabled": true
    }
    }
  }
}
```

__YAML__

```yaml
Resources:
  GraphQLApi:
    Type: AWS::AppSync::GraphQLApi
    Properties:
      Name: api
      AuthenticationType: AWS_IAM

      # Enable active tracing
      XrayEnabled: true
```
</details>

<details>
<summary>Serverless Framework</summary>

```yaml
resources:
  Resources:
    GraphQLApi:
      Type: AWS::AppSync::GraphQLApi
      Properties:
        Name: api
        AuthenticationType: AWS_IAM

        # Enable active tracing
        XrayEnabled: true
```
</details>

<details>
<summary>Terraform</summary>

```hcl
resource "aws_appsync_graphql_api" "this" {
  name                = "api"
  authentication_type = "AWS_IAM"

  # Enable active tracing
  xray_enabled = true
}
```
</details>

### See also

* [Serverless Lens: Distributed Tracing](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/distributed-tracing.html)
* [Tracing with AWS X-Ray](https://docs.aws.amazon.com/appsync/latest/devguide/x-ray-tracing.html)