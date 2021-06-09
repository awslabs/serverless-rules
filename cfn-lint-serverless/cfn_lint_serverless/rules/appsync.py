"""
Rules for AppSync resources
"""


from cfnlint.rules import CloudFormationLintRule, RuleMatch


class AppSyncTracingRule(CloudFormationLintRule):
    """
    Ensure AppSync GraphQL APIs have tracing enabled
    """

    id = "WS3000"  # noqa: VNE003
    shortdesc = "AppSync Tracing"
    description = "Ensure AppSync GraphQL APIs have tracing enabled"
    source_url = "https://awslabs.github.io/serverless-rules/rules/appsync/tracing/"
    tags = ["appsync"]

    _message = "AppSync GraphQL API {} should have XrayEnabled set to true."

    def match(self, cfn):
        """
        Match against AppSync GraphQL Apis without tracing enabled
        """

        matches = []

        for key, value in cfn.get_resources(["AWS::AppSync::GraphQLApi"]).items():
            xray_enabled = value.get("Properties", {}).get("XrayEnabled", False)

            if not xray_enabled:
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        return matches
