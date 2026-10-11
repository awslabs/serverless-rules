"""
Rules for API Gateway resources
"""

import json
import re
from typing import List

from cfnlint.rules import CloudFormationLintRule, RuleMatch

from .. import sam


class ApiGatewayLoggingRule(CloudFormationLintRule):
    """
    Ensure API Gateway REST and HTTP APIs have logging enabled
    """

    id = "ES2000"  # noqa: N815
    shortdesc = "API Gateway Logging"
    description = "Ensure that API Gateway REST and HTTP APIs have logging enabled"
    source_url = "https://awslabs.github.io/serverless-rules/rules/api_gateway/logging/"
    tags = ["apigateway"]

    _message_access_log_settings = "API Gateway stage {} is missing the AccessLogSetting property."
    _message_destination_arn = "API Gateway stage {} is missing the AccessLogSetting.DestinationArn property."

    def _match_properties(self, label: str, path: list, properties: dict, kind: str) -> List[RuleMatch]:
        """
        Match effective API properties without log settings
        """

        if kind == "rest":
            log_settings = properties.get("AccessLogSetting", None)
        else:
            log_settings = properties.get("AccessLogSettings", None)

        if log_settings is None:
            return [RuleMatch(path, self._message_access_log_settings.format(label))]

        if not isinstance(log_settings, dict) or log_settings.get("DestinationArn", None) is None:
            return [RuleMatch(path, self._message_destination_arn.format(label))]

        return []

    def match(self, cfn) -> List[RuleMatch]:
        """
        Match against API Gateway stages and SAM APIs without log settings
        """

        matches = []

        for key, value in cfn.get_resources(["AWS::ApiGateway::Stage", "AWS::ApiGatewayV2::Stage"]).items():
            kind = "rest" if value["Type"] == "AWS::ApiGateway::Stage" else "http"

            matches.extend(self._match_properties(key, ["Resources", key], value.get("Properties", {}), kind))

        # Explicit SAM APIs
        for key, value in sam.get_apis(cfn).items():
            matches.extend(
                self._match_properties(key, ["Resources", key], value.get("Properties", {}), "rest")
            )

        for key, value in sam.get_http_apis(cfn).items():
            matches.extend(
                self._match_properties(
                    key, ["Resources", key], value.get("Properties", {}), "http"
                )
            )

        # Implicit APIs created by function events, reported per owning function
        for view in sam.get_implicit_apis(cfn):
            matches.extend(self._match_properties(view["owner"], view["path"], view["properties"], view["kind"]))

        return matches


class ApiGatewayStructuredLoggingRule(CloudFormationLintRule):
    """
    Ensure API Gateway REST and HTTP APIs use structured logging
    """

    id = "WS2001"  # noqa: N815
    shortdesc = "API Gateway Structured Logging"
    description = "Ensure that API Gateway REST and HTTP APIs are using structured logs"
    source_url = "https://awslabs.github.io/serverless-rules/rules/api_gateway/structured_logging/"
    tags = ["apigateway"]

    _message = "API Gateway stage {} is not using structured logging for the AccessLogSetting.Format property."
    _log_format_pattern = re.compile(r"\$context\.[a-zA-Z\.]+")

    def _check_log_format(self, log_format: str) -> bool:
        """
        Validate if a log Format is valid JSON

        As users can use non-string keys like '$context.integration.latency', we
        need to transform the log format into something that can be decoded into
        JSON first.

        If the log format is not a string (e.g., it's an intrinsic function like
        !Ref or !Sub), we skip validation and return True.
        """

        # Skip validation when Format is an intrinsic function (Ref/Sub/etc.)
        if not isinstance(log_format, str):
            return True

        # Strip any leading/trailing quotes and whitespace that might be in the string
        log_format = log_format.strip()
        if (log_format.startswith("'") and log_format.endswith("'")) or (
            log_format.startswith('"') and log_format.endswith('"')
        ):
            log_format = log_format[1:-1]

        # Replace context variables with a simple value
        log_format = self._log_format_pattern.sub("0", log_format)

        try:
            json.loads(log_format)
            return True
        except json.decoder.JSONDecodeError:
            return False

    def _match_properties(self, label: str, path: list, properties: dict, kind: str) -> List[RuleMatch]:
        """
        Match effective API properties not using a JSON log format
        """

        settings_key = "AccessLogSetting" if kind == "rest" else "AccessLogSettings"
        log_settings = properties.get(settings_key, {})

        if not isinstance(log_settings, dict):
            return []

        log_format = log_settings.get("Format", None)

        # Ignore if it's not set. Another rule should catch it.
        if log_format is None:
            return []

        # Unresolved intrinsic expressions cannot be asserted as valid JSON
        if not isinstance(log_format, str):
            return []

        if not self._check_log_format(log_format):
            return [RuleMatch(path, self._message.format(label))]

        return []

    def match(self, cfn):
        """
        Match against API Gateway stages and SAM APIs not using structured logging
        """

        matches = []

        for key, value in cfn.get_resources(["AWS::ApiGateway::Stage", "AWS::ApiGatewayV2::Stage"]).items():
            kind = "rest" if value["Type"] == "AWS::ApiGateway::Stage" else "http"

            matches.extend(self._match_properties(key, ["Resources", key], value.get("Properties", {}), kind))

        # Explicit SAM APIs
        for key, value in sam.get_apis(cfn).items():
            matches.extend(
                self._match_properties(key, ["Resources", key], value.get("Properties", {}), "rest")
            )

        for key, value in sam.get_http_apis(cfn).items():
            matches.extend(
                self._match_properties(
                    key, ["Resources", key], value.get("Properties", {}), "http"
                )
            )

        # Implicit APIs created by function events, reported per owning function
        for view in sam.get_implicit_apis(cfn):
            matches.extend(self._match_properties(view["owner"], view["path"], view["properties"], view["kind"]))

        return matches


class ApiGatewayDefaultThrottlingRule(CloudFormationLintRule):
    """
    Ensure API Gateway REST APIs have throttling enabled
    """

    id = "ES2003"  # noqa: N815
    shortdesc = "API Gateway Throttling"
    description = "Ensure that API Gateway REST APIs have default throttling limits set."
    source_url = "https://awslabs.github.io/serverless-rules/rules/api_gateway/default_throttling/"
    tags = ["apigateway"]

    _message_method_settings = "API Gateway stage {} does not have a default MethodSettings property with ThrottlingBurstLimit and ThrottlingRateLimit."  # noqa: E501
    _message_default_route_settings = "API Gateway stage {} does not have a default DefaultRouteSettings property with ThrottlingBurstLimit and ThrottlingRateLimit."  # noqa: E501

    def _match_rest_api(self, label: str, path: list, properties: dict) -> List[RuleMatch]:
        """
        Match effective REST API properties
        """

        method_settings = [
            ms
            for ms in properties.get("MethodSettings", [])
            if isinstance(ms, dict) and ms.get("HttpMethod") == "*" and ms.get("ResourcePath") == "/*"
        ] or [{}]
        default_method_setting = method_settings[0]

        if "ThrottlingBurstLimit" not in default_method_setting or "ThrottlingRateLimit" not in default_method_setting:
            return [RuleMatch(path, self._message_method_settings.format(label))]

        return []

    def _match_http_api(self, label: str, path: list, properties: dict) -> List[RuleMatch]:
        """
        Match effective HTTP API properties
        """

        route_settings = properties.get("DefaultRouteSettings", {})

        if not isinstance(route_settings, dict):
            route_settings = {}

        if "ThrottlingBurstLimit" not in route_settings or "ThrottlingRateLimit" not in route_settings:
            return [RuleMatch(path, self._message_default_route_settings.format(label))]

        return []

    def match(self, cfn) -> List[RuleMatch]:
        """
        Match against API Gateway stages and SAM APIs without default throttling
        """

        matches = []

        # API Gateway REST APIs
        for key, value in cfn.get_resources(["AWS::ApiGateway::Stage"]).items():
            matches.extend(self._match_rest_api(key, ["Resources", key], value.get("Properties", {})))

        # API Gateway HTTP APIs
        for key, value in cfn.get_resources(["AWS::ApiGatewayV2::Stage"]).items():
            matches.extend(self._match_http_api(key, ["Resources", key], value.get("Properties", {})))

        # Explicit SAM APIs
        for key, value in sam.get_apis(cfn).items():
            matches.extend(
                self._match_rest_api(key, ["Resources", key], value.get("Properties", {}))
            )

        for key, value in sam.get_http_apis(cfn).items():
            matches.extend(
                self._match_http_api(key, ["Resources", key], value.get("Properties", {}))
            )

        # Implicit APIs created by function events, reported per owning function
        for view in sam.get_implicit_apis(cfn):
            if view["kind"] == "rest":
                matches.extend(self._match_rest_api(view["owner"], view["path"], view["properties"]))
            else:
                matches.extend(self._match_http_api(view["owner"], view["path"], view["properties"]))

        return matches


class ApiGatewayTracingRule(CloudFormationLintRule):
    """
    Ensure API Gateway REST APIs have tracing enabled
    """

    id = "WS2002"  # noqa: N815
    shortdesc = "API Gateway Tracing"
    description = "Ensure that API Gateway REST APIs have tracing enabled"
    source_url = "https://awslabs.github.io/serverless-rules/rules/api_gateway/tracing/"
    tags = ["apigateway"]

    _message = "API Gateway stage {} does not have the TracingEnabled property set to true."

    def _match_properties(self, label: str, path: list, properties: dict) -> List[RuleMatch]:
        """
        Match effective REST API properties without tracing enabled
        """

        if not properties.get("TracingEnabled", False):
            return [RuleMatch(path, self._message.format(label))]

        return []

    def match(self, cfn):
        """
        Match against API Gateway stages and SAM REST APIs without tracing enabled
        """

        matches = []

        # HTTP APIs don't support X-Ray
        for key, value in cfn.get_resources(["AWS::ApiGateway::Stage"]).items():
            matches.extend(self._match_properties(key, ["Resources", key], value.get("Properties", {})))

        # Explicit SAM REST APIs
        for key, value in sam.get_apis(cfn).items():
            matches.extend(self._match_properties(key, ["Resources", key], value.get("Properties", {})))

        # Implicit REST APIs created by function events, reported per owning function
        for view in sam.get_implicit_apis(cfn):
            if view["kind"] == "rest":
                matches.extend(self._match_properties(view["owner"], view["path"], view["properties"]))

        return matches
