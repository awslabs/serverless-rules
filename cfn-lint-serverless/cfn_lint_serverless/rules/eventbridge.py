"""
Rules for EventBridge resources
"""

from cfnlint.rules import CloudFormationLintRule, RuleMatch

from .. import sam


class EventBridgeDLQRule(CloudFormationLintRule):
    """
    Ensure Event Bridge rules have a DLQ configured
    """

    id = "ES4000"  # noqa: N815
    shortdesc = "EventBridge DLQ"
    description = "Ensure Event Bridge rules have a DLQ configured"
    source_url = "https://awslabs.github.io/serverless-rules/rules/eventbridge/rule_without_dlq/"
    tags = ["stepfunctions"]

    _message = "EventBridge rule {} should have a DeadLetterConfig.Arn property for all its Targets."
    _message_sam = "Event for Lambda function {} creates an EventBridge rule target without a dead-letter queue."  # noqa: E501

    @staticmethod
    def _has_dlq(properties: dict) -> bool:
        """
        Return True for an explicit queue ARN or a valid SAM-managed DLQ declaration
        """

        dead_letter_config = properties.get("DeadLetterConfig", {})

        if not isinstance(dead_letter_config, dict):
            return False

        return "Arn" in dead_letter_config or dead_letter_config.get("Type", None) == "SQS"

    def _match_sam_events(self, key: str, value: dict):
        """
        Match against SAM events whose generated EventBridge target supports a
        DLQ. ScheduleV2 uses a different mechanism and is not covered.
        """

        matches = []

        for name, event_type, properties, path in sam.iter_events(key, value):
            if event_type not in sam.DLQ_EVENT_TYPES:
                continue

            if not self._has_dlq(properties):
                matches.append(RuleMatch(path, self._message_sam.format(f"{key}/{name}")))

        return matches

    def match(self, cfn):
        """
        Match against EventBridge rules without a DeadLetterConfig for all targets
        """

        matches = []

        for key, value in cfn.get_resources(["AWS::Events::Rule"]).items():
            targets = value.get("Properties", {}).get("Targets", [])

            if not all(["Arn" in t.get("DeadLetterConfig", {}) for t in targets]):
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        for key, value in sam.get_functions(cfn).items():
            matches.extend(self._match_sam_events(key, value))

        return matches
