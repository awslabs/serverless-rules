Creating rules
==============

If you are thinking of creating or proposing a new rule, please follow the process outline below. The first step before adding a new rule is to [submit an issue](#create-an-issue) to collect feedback from other members of the community.

## Create an issue

Before starting the implementation of a new rule, please [create an issue using the _New rule_ template](https://github.com/awslabs/serverless-rules/issues/new?assignees=&labels=feature-request%2C+triage&template=rule.md&title=). This will allow members of the community to provide feedback on its implementation, if it meets the needs of most serverless users, if it's the right level, etc.

## Template for `cfn-lint` rules

```python
# TODO: set the rule name
class __Rule(CloudFormationLintRule):
    # TODO: set docstring
    """
    Ensure that ...
    """

    # TODO: update these values
    id = "..." # noqa: N815
    shortdesc = "..."
    description = "Ensure that ..."
    source_url = "..."
    tags = ["lambda"]

    _message = "... {} ..."

    def match(self, cfn):
        # TODO: update docstring
        """
        Match against ...
        """

        matches = []

        # TODO: set resource type
        for key, value in cfn.get_resources(["..."]).items():
            # TODO: set property name
            prop = value.get("Properties", {}).get("...", None)

            if prop is None:
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        return matches
```

## Template for documentation

Please use the following template when writing documentation for a rule. Each rule goes into a separate markdown file into the relevant service folder. For example, a rule for AWS Lambda would go into the `docs/rules/lambda/` folder.

~~~markdown
# _Service Name Rule name_

__Level__: _Rule level_
{: class="badge badge-red" }

__Initial version__: _release version_
{: class="badge badge-blue" }

__cfn-lint__: _cfn-lint rule ID_
{: class="badge" }

__tflint__: _tflint rule name_
{: class="badge" }

_Short explanation on the rule_

## Implementations

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

## See also

* _List of links to the relevant documentation, from sources such as AWS Well-Architected, service documentation, etc._
~~~