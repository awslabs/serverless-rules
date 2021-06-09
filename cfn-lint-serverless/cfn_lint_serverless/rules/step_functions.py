"""
Rules for Step Functions resources
"""


from cfnlint.rules import CloudFormationLintRule, RuleMatch


class StepFunctionsTracingRule(CloudFormationLintRule):
    """
    Ensure Step Functions state machines have tracing enabled
    """

    id = "WS5000"  # noqa: VNE003
    shortdesc = "Step Functions Tracing"
    description = "Ensure that Step Functions state machines have tracing enabled"
    source_url = "https://awslabs.github.io/serverless-rules/rules/step_functions/tracing/"
    tags = ["stepfunctions"]

    _message = "Step Functions state machine {} should have TracingConfiguration.Enabled set to true."

    def match(self, cfn):
        """
        Match against Step Functions state machine without tracing enabled
        """

        matches = []

        for key, value in cfn.get_resources(["AWS::StepFunctions::StateMachine"]).items():
            tracing_enabled = value.get("Properties", {}).get("TracingConfiguration", {}).get("Enabled", False)

            if not tracing_enabled:
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        return matches
