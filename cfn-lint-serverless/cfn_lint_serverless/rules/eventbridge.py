"""
Rules for EventBridge resources
"""


from cfnlint.rules import CloudFormationLintRule, RuleMatch


class EventBridgeDLQRule(CloudFormationLintRule):
    """
    Ensure Event Bridge rules have a DLQ configured
    """

    id = "ES4000"  # noqa: VNE003
    shortdesc = "EventBridge DLQ"
    description = "Ensure Event Bridge rules have a DLQ configured"
    source_url = "https://awslabs.github.io/serverless-rules/rules/eventbridge/rule_without_dlq/"
    tags = ["stepfunctions"]

    _message = "EventBridge rule {} should have a DeadLetterConfig.Arn property for all its Targets."

    def match(self, cfn):
        """
        Match against EventBridge rules without a DeadLetterConfig for all targets
        """

        matches = []

        for key, value in cfn.get_resources(["AWS::Events::Rule"]).items():
            targets = value.get("Properties", {}).get("Targets", [])

            if not all(["Arn" in t.get("DeadLetterConfig", {}) for t in targets]):
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        return matches
