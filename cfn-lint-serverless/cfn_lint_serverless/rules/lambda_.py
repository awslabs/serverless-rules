"""
Rules for Lambda resources
"""


import re
from collections import defaultdict
from typing import List

from cfnlint.rules import CloudFormationLintRule, RuleMatch


class LambdaTracingRule(CloudFormationLintRule):
    """
    Ensure Lambda functions have tracing enabled
    """

    id = "WS1000"  # noqa: VNE003
    shortdesc = "Lambda Tracing"
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


class LambdaESMDestinationRule(CloudFormationLintRule):
    """
    Ensure Lambda event source mappings have a destination configured
    """

    id = "ES1001"  # noqa: VNE003
    shortdesc = "Lambda Event Source Mapping Destination"
    description = "Ensure Lambda event source mappings have a destination configured"
    tags = ["lambda"]

    _message = "Lambda event source mapping {} should have a DestinationConfig.OnFailure.Destination property."

    def match(self, cfn):
        """
        Match against Event Source Mappings without a destination configured
        """

        matches = []

        for key, value in cfn.get_resources(["AWS::Lambda::EventSourceMapping"]).items():
            destination = (
                value.get("Properties", {}).get("DestinationConfig", {}).get("OnFailure", {}).get("Destination", False)
            )

            if not destination:
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        return matches


class LambdaPermissionPrincipalsRule(CloudFormationLintRule):
    """
    Ensure that Lambda functions do not have Lambda permissions with different principals
    """

    id = "WS1002"  # noqa: VNE003
    shortdesc = "Lambda Permission Principals"
    description = "Ensure that Lambda functions do not have Lambda permissions with different principals"
    tags = ["lambda"]
    _message = "Lambda function {} has Lambda permissions with different principals"

    def _get_permissions(self, cfn):
        """
        Parse all AWS::Lambda::Permissions in the template
        """

        permissions = defaultdict(list)

        for _, value in cfn.get_resources(["AWS::Lambda::Permission"]).items():
            function_name = value.get("Properties", {}).get("FunctionName", "")
            principal = value.get("Properties", {}).get("Principal", "")

            if isinstance(function_name, dict):
                if "Ref" in function_name:
                    function_name = function_name["Ref"]
                elif "Fn::Ref" in function_name:
                    function_name = function_name["Fn::Ref"]
                elif "GetAtt" in function_name:
                    function_name = function_name["GetAtt"][0]
                elif "Fn::GetAtt" in function_name:
                    function_name = function_name["Fn::GetAtt"][0]

            permissions[str(function_name)].append(principal)

        return permissions

    def match(self, cfn):
        """
        Match against Lambda functions with multiple Principal for Permissions
        """

        matches = []

        for key, value in self._get_permissions(cfn).items():
            if len(set(value)) > 1:
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        return matches


class LambdaStarPermissionRule(CloudFormationLintRule):
    """
    Ensure that Lambda functions don't have stars in IAM policy actions
    """

    id = "WS1003"  # noqa: VNE003
    shortdesc = "Lambda Star Permission"
    description = "Ensure that Lambda functions don't have stars in IAM policy actions"
    tags = ["lambda"]

    _message = "IAM Role {} with Lambda as principal has policy actions with stars"

    def _get_principals(self, properties) -> List[str]:
        """
        Retrieve principals from assume role policy documents
        """

        principals = []

        for statement in properties.get("AssumeRolePolicyDocument", {}).get("Statement", []):
            if "Service" not in statement.get("Principal", {}):
                continue

            services = statement.get("Principal", {}).get("Service")

            if isinstance(services, str):
                principals.append(services)
            elif isinstance(services, list):
                principals.extend(services)
            # Ignored for now if it's not a list of str

        return principals

    def _get_actions(self, properties) -> List[str]:
        """
        Retrieve actions from policy documents
        """

        actions = []

        for policy in properties.get("Policies", []):
            for statement in policy.get("PolicyDocument", {}).get("Statement", []):

                action = statement.get("Action")

                if isinstance(action, str):
                    actions.append(action)
                elif isinstance(action, list):
                    actions.extend(action)
                # Ignored for now if it's not a list of str

        return actions

    def match(self, cfn):
        """
        Match against IAM roles with Lambda as principal and stars in actions
        """

        matches = []

        for key, value in cfn.get_resources(["AWS::IAM::Role"]).items():
            principals = self._get_principals(value.get("Properties", {}))
            actions = self._get_actions(value.get("Properties", {}))

            if "lambda.amazonaws.com" in principals and any(["*" in a for a in actions]):
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        return matches


class LambdaLogRetentionRule(CloudFormationLintRule):
    """
    Ensure that Lambda functions have a corresponding Log Group with retention
    """

    id = "WS1004"  # noqa: VNE003
    shortdesc = "Lambda Log Retention"
    description = "Ensure that Lambda functions have a corresponding Log Group with retention"
    tags = ["lambda"]

    _message = "Lambda function {} does not have a corresponding log group with a Retention property"

    def _get_valid_functions(self, log_groups):
        """
        Return function names with valid LogGroups
        """

        known = []

        # Scan log groups for resource names
        for resource in log_groups.values():
            # This use an autogenerated log group name
            if "LogGroupName" not in resource.get("Properties", {}):
                continue

            log_group_name = resource.get("Properties", {}).get("LogGroupName")
            retention = log_group_name = resource.get("Properties", {}).get("RetentionInDays", None)

            # No retention
            if retention is None:
                continue

            # Look for Substitution function
            if not isinstance(log_group_name, dict):
                continue

            if "Fn::Sub" in log_group_name:
                log_group_name = log_group_name["Fn::Sub"]
            elif "Sub" in log_group_name:
                log_group_name = log_group_name["Sub"]
            else:
                continue

            match = re.search(r"/aws/lambda/\${(?P<func>[^}]+)}", log_group_name["Fn::Sub"])
            if match is not None:
                known.append(match["func"])

        return known

    def match(self, cfn):
        """
        Match against Lambda functions without a LogGroup with a Retention property
        """

        matches = []

        functions = cfn.get_resources(["AWS::Lambda::Function"])
        log_groups = cfn.get_resources(["AWS::Logs::LogGroup"])

        known = self._get_valid_functions(log_groups)

        # Scan functions against log groups
        for function in functions.keys():
            if function not in known:
                matches.append(RuleMatch(["Resources", function], self._message.format(function)))

        return matches
