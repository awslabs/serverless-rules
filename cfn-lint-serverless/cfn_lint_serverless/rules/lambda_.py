"""
Rules for Lambda resources
"""

from collections import defaultdict
from typing import Dict, List, Union

from cfnlint.rules import CloudFormationLintRule, RuleMatch

from .. import sam
from ..utils import Value, ensure_list


class LambdaTracingRule(CloudFormationLintRule):
    """
    Ensure Lambda functions have tracing enabled
    """

    id = "WS1000"  # noqa: N815
    shortdesc = "Lambda Tracing"
    description = "Ensure that Lambda functions have tracing enabled"
    source_url = "https://awslabs.github.io/serverless-rules/rules/lambda/tracing/"
    tags = ["lambda"]

    _message = "Lambda function {} should have Tracing property set to 'Active'."

    def match(self, cfn):
        """
        Match against Lambda functions without tracing enabled
        """

        matches = []

        for key, value in cfn.get_resources(["AWS::Lambda::Function"]).items():
            tracing_mode = value.get("Properties", {}).get("TracingConfig", {}).get("Mode", None)

            if tracing_mode != "Active":
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        for key, value in sam.get_functions(cfn).items():
            tracing_mode = value.get("Properties", {}).get("Tracing", None)

            if tracing_mode != "Active":
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        return matches


class LambdaESMDestinationRule(CloudFormationLintRule):
    """
    Ensure Lambda event source mappings have a destination configured
    """

    id = "ES1001"  # noqa: N815
    shortdesc = "Lambda Event Source Mapping Destination"
    description = "Ensure Lambda event source mappings have a destination configured"
    source_url = "https://awslabs.github.io/serverless-rules/rules/lambda/eventsourcemapping_failure_destination/"
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

        # SAM Kinesis and DynamoDB events back stream-based event source
        # mappings. SQS polling events do not support this destination.
        for key, value in sam.get_functions(cfn).items():
            for name, event_type, properties, path in sam.iter_events(key, value):
                if event_type not in sam.STREAM_EVENT_TYPES:
                    continue

                destination = (
                    properties.get("DestinationConfig", {}).get("OnFailure", {}).get("Destination", False)
                )

                if not destination:
                    matches.append(RuleMatch(path, self._message.format(f"{key}/{name}")))

        return matches


class LambdaPermissionPrincipalsRule(CloudFormationLintRule):
    """
    Ensure that Lambda functions do not have Lambda permissions with different principals
    """

    id = "WS1002"  # noqa: N815
    shortdesc = "Lambda Permission Principals"
    description = "Ensure that Lambda functions do not have Lambda permissions with different principals"
    source_url = "https://awslabs.github.io/serverless-rules/rules/lambda/permission_multiple_principals/"
    tags = ["lambda"]
    _message = "Lambda function {} has Lambda permissions with different principals"

    def _get_permissions(self, cfn):
        """
        Parse all AWS::Lambda::Permissions in the template
        """

        permissions = defaultdict(list)

        for _, value in cfn.get_resources(["AWS::Lambda::Permission"]).items():
            principal = Value(value.get("Properties", {}).get("Principal", ""))
            function_name = Value(value.get("Properties", {}).get("FunctionName", ""))

            for reference in function_name.references:
                permissions[reference].append(principal.id)

        return permissions

    def match(self, cfn):
        """
        Match against Lambda functions with multiple Principal for Permissions
        """

        matches = []

        permissions = self._get_permissions(cfn)

        # Supported SAM function events also imply resource-based invocation
        # permissions. Polling event sources are not included in the mapping.
        for key, value in sam.get_functions(cfn).items():
            for _, event_type, _, _ in sam.iter_events(key, value):
                principal = sam.EVENT_PRINCIPALS.get(event_type, None)

                if principal is not None:
                    permissions[key].append(principal)

        for key, value in permissions.items():
            if len(set(value)) > 1:
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        return matches


class LambdaStarPermissionRule(CloudFormationLintRule):
    """
    Ensure that Lambda functions don't have stars in IAM policy actions
    """

    id = "WS1003"  # noqa: N815
    shortdesc = "Lambda Star Permission"
    description = "Ensure that Lambda functions don't have stars in IAM policy actions"
    source_url = "https://awslabs.github.io/serverless-rules/rules/lambda/star_permissions/"
    tags = ["lambda"]

    _message = "IAM Role {} with Lambda as principal has policy actions with stars"
    _message_sam = "Lambda function {} has an inline policy with stars in actions"

    def _get_principals(self, properties) -> List[str]:
        """
        Retrieve principals from assume role policy documents
        """

        principals = []

        for statement in ensure_list(properties.get("AssumeRolePolicyDocument", {}).get("Statement")):
            if not isinstance(statement, dict):
                continue

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

        for policy in ensure_list(properties.get("Policies")):
            if not isinstance(policy, dict):
                continue

            for statement in ensure_list(policy.get("PolicyDocument", {}).get("Statement")):
                if not isinstance(statement, dict):
                    continue

                action = statement.get("Action")

                if isinstance(action, str):
                    actions.append(action)
                elif isinstance(action, list):
                    actions.extend(action)
                # Ignored for now if it's not a list of str

        return actions

    def _get_sam_actions(self, properties) -> List[str]:
        """
        Retrieve actions from inline policy documents in SAM function Policies

        Named managed policies and SAM policy templates are not expanded, and
        an explicit Role means SAM does not deploy the function Policies.
        """

        actions = []

        if properties.get("Role", None) is not None:
            return actions

        for policy in ensure_list(properties.get("Policies")):
            # Managed policy name or ARN
            if not isinstance(policy, dict):
                continue

            # SAM policy template, e.g. DynamoDBWritePolicy
            if "Statement" not in policy:
                continue

            for statement in ensure_list(policy.get("Statement")):
                if not isinstance(statement, dict):
                    continue

                action = statement.get("Action")

                if isinstance(action, str):
                    actions.append(action)
                elif isinstance(action, list):
                    actions.extend(action)

        return actions

    @staticmethod
    def _has_star(actions: List[str]) -> bool:
        """
        Return True if any action is a literal or service-level wildcard
        """

        return any([a == "*" or ":*" in a for a in actions])

    def match(self, cfn):
        """
        Match against IAM roles with Lambda as principal and stars in actions,
        and against inline SAM function policies with stars in actions
        """

        matches = []

        for key, value in cfn.get_resources(["AWS::IAM::Role"]).items():
            principals = self._get_principals(value.get("Properties", {}))
            actions = self._get_actions(value.get("Properties", {}))

            if "lambda.amazonaws.com" in principals and self._has_star(actions):
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        for key, value in sam.get_functions(cfn).items():
            actions = self._get_sam_actions(value.get("Properties", {}))

            if self._has_star(actions):
                matches.append(RuleMatch(["Resources", key], self._message_sam.format(key)))

        return matches


class LambdaLogRetentionRule(CloudFormationLintRule):
    """
    Ensure that Lambda functions have a corresponding Log Group with retention
    """

    id = "WS1004"  # noqa: N815
    shortdesc = "Lambda Log Retention"
    description = "Ensure that Lambda functions have a corresponding Log Group with retention"
    source_url = "https://awslabs.github.io/serverless-rules/rules/lambda/log_retention/"
    tags = ["lambda"]

    _message = "Lambda function {} does not have a corresponding log group with a Retention property"

    def _get_valid_functions(self, log_groups):
        """
        Return function names with valid LogGroups
        """

        known_refs = []
        known_names = []

        # Scan log groups for resource names
        for resource in log_groups.values():
            log_group_name = Value(resource.get("Properties", {}).get("LogGroupName", None))
            retention = resource.get("Properties", {}).get("RetentionInDays", None)

            # No retention or log group name, break early
            if log_group_name is None or retention is None:
                continue

            # Searching for references in log group name
            if len(log_group_name.references) > 0:
                known_refs.extend(log_group_name.references)
            # Otherwise, check if this is a hardcoded name
            elif log_group_name.id.find("/aws/lambda/") == 0:
                known_names.append(log_group_name.id[12:])

        return {"ref": known_refs, "name": known_names}

    def match(self, cfn):
        """
        Match against Lambda functions without a LogGroup with a Retention property
        """

        matches = []

        functions = cfn.get_resources(["AWS::Lambda::Function", sam.SAM_FUNCTION_TYPE])
        log_groups = cfn.get_resources(["AWS::Logs::LogGroup"])

        known = self._get_valid_functions(log_groups)

        # Scan functions against log groups
        for function_ref, function in functions.items():
            log_group_found = False

            if function_ref in known["ref"]:
                log_group_found = True

            function_name = function.get("Properties", {}).get("FunctionName", None)
            if function_name is not None and function_name in known["name"]:
                log_group_found = True

            if not log_group_found:
                matches.append(RuleMatch(["Resources", function_ref], self._message.format(function_ref)))

        return matches


class LambdaDefaultMemorySizeRule(CloudFormationLintRule):
    """
    Ensure that Lambda functions have an explicit memory value
    """

    id = "ES1005"  # noqa: N815
    shortdesc = "Lambda Default Memory Size"
    description = "Ensure that Lambda functions have an explicit memory value"
    source_url = "https://awslabs.github.io/serverless-rules/rules/lambda/default_memory_size/"
    tags = ["lambda"]

    _message = "Lambda function {} does not have a MemorySize property"

    def match(self, cfn):
        """
        Match against Lambda functions without an explicity MemorySize
        """

        matches = []

        for key, value in cfn.get_resources(["AWS::Lambda::Function"]).items():
            memory_size = value.get("Properties", {}).get("MemorySize", None)

            if memory_size is None:
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        for key, value in sam.get_functions(cfn).items():
            memory_size = value.get("Properties", {}).get("MemorySize", None)

            if memory_size is None:
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        return matches


class LambdaDefaultTimeoutRule(CloudFormationLintRule):
    """
    Ensure that Lambda functions have an explicit timeout value
    """

    id = "ES1006"  # noqa: N815
    shortdesc = "Lambda Default Timeout"
    description = "Ensure that Lambda functions have an explicit timeout value"
    source_url = "https://awslabs.github.io/serverless-rules/rules/lambda/default_timeout/"
    tags = ["lambda"]

    _message = "Lambda function {} does not have a Timeout property"

    def match(self, cfn):
        """
        Match against Lambda functions without an explicity Timeout
        """

        matches = []

        for key, value in cfn.get_resources(["AWS::Lambda::Function"]).items():
            timeout = value.get("Properties", {}).get("Timeout", None)

            if timeout is None:
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        for key, value in sam.get_functions(cfn).items():
            timeout = value.get("Properties", {}).get("Timeout", None)

            if timeout is None:
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        return matches


class LambdaAsyncNoDestinationRule(CloudFormationLintRule):
    """
    Ensure that Lambda functions invoked asynchronously have a destination configured
    """

    id = "ES1007"  # noqa: N815
    shortdesc = "Lambda Async Destination"
    description = "Ensure that Lambda functions invoked asynchronously have a destination configured"
    source_url = "https://awslabs.github.io/serverless-rules/rules/lambda/async_failure_destination/"
    tags = ["lambda"]

    _message = "Lambda permission {} has an asynchronous permission but doesn't have an EventInvokeConfig resource related to it"  # noqa: E501
    _message_sam = "Lambda function event {} is invoked asynchronously but doesn't have an on-failure destination configured"  # noqa: E501

    _async_principals = [
        "events.amazonaws.com",
        "events.amazonaws.com.cn",
        "iot.amazonaws.com",
        "iot.amazonaws.com.cn",
        "s3.amazonaws.com",
        "s3.amazonaws.com.cn",
        "sns.amazonaws.com",
        "sns.amazonaws.com.cn",
    ]

    def _get_async_functions(self, permissions: Dict[str, dict]) -> Dict[str, Union[dict, str]]:
        """
        Return a list of FunctionName properties for permissions with principals that invoke Lambda
        functions asynchronously
        """

        async_functions = {}

        for key, value in permissions.items():
            function_name = value.get("Properties", {}).get("FunctionName", "")
            principal = value.get("Properties", {}).get("Principal", "")

            # No FunctionName, we cannot evaluate this rule
            if not function_name:
                continue

            if principal in self._async_principals:
                async_functions[key] = function_name

        return async_functions

    def match(self, cfn):
        """
        Match against Lambda functions without an explicity Timeout
        """

        matches = []

        permissions = cfn.get_resources(["AWS::Lambda::Permission"])
        event_invoke_configs = cfn.get_resources(["AWS::Lambda::EventInvokeConfig"])

        async_functions = self._get_async_functions(permissions)

        found = []

        for value in event_invoke_configs.values():
            function_name = value.get("Properties", {}).get("FunctionName", "")
            on_failure_destination = (
                value.get("Properties", {}).get("DestinationConfig", {}).get("OnFailure", {}).get("Destination", None)
            )

            if function_name in async_functions.values() and on_failure_destination is not None:
                found.append(function_name)

        for key, value in async_functions.items():
            if value not in found:
                matches.append(RuleMatch(["Resources", key], self._message.format(key)))

        # SAM functions with supported asynchronous event types
        for key, value in sam.get_functions(cfn).items():
            matches.extend(self._match_sam_function(key, value))

        return matches

    @staticmethod
    def _has_sam_failure_destination(properties: dict) -> bool:
        """
        Return True for an explicit destination or a SAM-managed destination declaration
        """

        event_invoke_config = properties.get("EventInvokeConfig", {})
        if not isinstance(event_invoke_config, dict):
            event_invoke_config = {}

        on_failure = event_invoke_config.get("DestinationConfig", {}).get("OnFailure", {})
        if not isinstance(on_failure, dict):
            on_failure = {}

        return (
            on_failure.get("Destination", None) is not None
            or on_failure.get("Type", None) in sam.SAM_MANAGED_DESTINATION_TYPES
        )

    def _match_sam_function(self, key: str, value: dict):
        """
        Match against SAM function events invoked asynchronously without a destination
        """

        properties = value.get("Properties", {})

        if self._has_sam_failure_destination(properties):
            return []

        return [
            RuleMatch(path, self._message_sam.format(f"{key}/{name}"))
            for name, event_type, _, path in sam.iter_events(key, value)
            if event_type in sam.ASYNC_EVENT_TYPES
        ]
