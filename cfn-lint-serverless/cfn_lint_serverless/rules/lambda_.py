"""
Rules for Lambda resources
"""


from cfnlint.rules import CloudFormationLintRule, RuleMatch


class LambdaTracingRule(CloudFormationLintRule):
    """
    Ensure Lambda functions have tracing enabled
    """

    id = "WS1000"  # noqa: VNE003
    shortdesc = "Lambda Logging"
    description = "Ensure that Lambda functions have tracing enabled"
    tags = ["lambda"]

    _message = "Lambda function {} should have TracingConfig.Mode set to 'Active'."

    def match(self, cfn):
        """
        Match against Lambda functions without tracing enabled
        """

        matches = []

        for key, value in cfn.get_resources(["AWS::Lambda::Function"]).items():
            tracing_mode = value.get("Properties", {}).get("TracingConfig", {}).get("Mode", None)

            if tracing_mode != "Active":
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        return matches
