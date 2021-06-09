"""
Rules for SNS resources
"""


from cfnlint.rules import CloudFormationLintRule, RuleMatch


class SnsNoRedrivePolicyRule(CloudFormationLintRule):
    """
    Ensure SNS subscriptions have a redrive policy configured
    """

    id = "ES7000"  # noqa: VNE003
    shortdesc = "SNS No Redrive Policy"
    description = "Ensure SNS subscriptions have a redrive policy configured"
    source_url = "https://awslabs.github.io/serverless-rules/rules/sns/redrive_policy/"
    tags = ["sns"]

    _message = "SNS Subscription {} should have a RedrivePolicy property configured."

    def match(self, cfn):
        """
        Match against SNS subscriptions without RedrivePolicy
        """

        matches = []

        for key, value in cfn.get_resources(["AWS::SNS::Subscription"]).items():
            redrive_policy = value.get("Properties", {}).get("RedrivePolicy", None)

            if redrive_policy is None:
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        return matches
