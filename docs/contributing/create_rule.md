Creating rules
==============

If you are thinking of creating or proposing a new rule, please follow the process outline below. The first step before any rule creating is to [submit an issue](#create-an-issue) to collect feedback from other members of the community.

## Create an issue

Before starting the implementation of a new rule, please [create an issue using the _New rule_ template](https://github.com/awslabs/serverless-rules/issues/new?assignees=&labels=feature-request%2C+triage&template=rule.md&title=). This will allow members of the community to provide feedback on its implementation, if it meets the needs of most serverless users, if it's the right level, etc.

## Template for documentation

Please use the following template when writing documentation for a rule. Each block goes into the proper Markdown file for that service. For example, a new rule for AWS Lambda goes into the [lambda.md](../rules/lambda.md) file.

~~~markdown
## _Rule name_

* __Level__: _Rule level_
* __cfn-lint__: _cfn-lint rule ID_
* __tflint__: _tflint rule name_

_Short explanation on the rule_

### Implementations

=== "CDK"

    ```typescript
    // Imports here

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // Implementation
      }
    }
    ```

=== "CloudFormation (JSON)"

    ```json
    {
      "Resources": {
        // Add resources here
      }
    }
    ```

=== "CloudFormation (YAML)"

    ```yaml
    Resources:
      # Add resources here
    ```

=== "Serverless Framework"

    ```yaml
    provider:
      name: aws
      # Add provider-specific configuration here

    resources:
      # Add resources here
    ```

=== "Terraform"

    ```tf
    # Add Terraform resources here
    ```

### See also

* _List of links to the relevant documentation, from sources such as AWS Well-Architected, service documentation, etc._
~~~