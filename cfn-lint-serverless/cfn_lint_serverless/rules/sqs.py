"""
Rules for SQS resources
"""


from cfnlint.rules import CloudFormationLintRule, RuleMatch

from ..utils import Value


class SqsNoRedrivePolicyRule(CloudFormationLintRule):
    """
    Ensure SQS queues have a redrive policy configured
    """

    id = "ES6000"  # noqa: VNE003
    shortdesc = "SQS No Redrive Policy"
    description = "Ensure SQS queues have a redrive policy configured"
    source_url = "https://awslabs.github.io/serverless-rules/rules/sqs/redrive_policy/"
    tags = ["sqs"]

    _message = "SQS Queue {} should have a RedrivePolicy property configured."

    def match(self, cfn):
        """
        Match against SQS queues without RedrivePolicy
        """

        matches = {}
        dlqs = []

        for key, value in cfn.get_resources(["AWS::SQS::Queue"]).items():
            redrive_policy = value.get("Properties", {}).get("RedrivePolicy", None)

            if redrive_policy is None:
                matches[key] = RuleMatch(["Resources", key], self._message.format(key))

            else:
                redrive_policy = Value(redrive_policy)
                # If a queue is used as a DLQ, it doesn't need a redrive policy
                # See https://github.com/awslabs/serverless-rules/issues/79
                dlqs.extend(redrive_policy.references)

        return [v for k, v in matches.items() if k not in dlqs]
