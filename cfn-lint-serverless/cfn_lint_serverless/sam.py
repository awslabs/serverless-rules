"""
Read-only views over AWS SAM resources.

These helpers expose SAM functions, APIs, and function events to the rules
without modifying the host template and without relying on resources that
SAM would generate at deployment time. cfn-lint has already merged effective
`Globals` settings into the resource properties returned here, so helpers
must not apply inheritance a second time.
"""

from typing import Dict, Iterator, List, Tuple

SAM_FUNCTION_TYPE = "AWS::Serverless::Function"
SAM_API_TYPE = "AWS::Serverless::Api"
SAM_HTTP_API_TYPE = "AWS::Serverless::HttpApi"

# SAM function event types that invoke Lambda through a resource-based
# policy, mapped to the service principal they imply. Polling event sources
# such as SQS, Kinesis, and DynamoDB are intentionally absent.
EVENT_PRINCIPALS = {
    "Api": "apigateway.amazonaws.com",
    "HttpApi": "apigateway.amazonaws.com",
    "CloudWatchEvent": "events.amazonaws.com",
    "EventBridgeRule": "events.amazonaws.com",
    "Schedule": "events.amazonaws.com",
    "IoTRule": "iot.amazonaws.com",
    "S3": "s3.amazonaws.com",
    "SNS": "sns.amazonaws.com",
}

# Event types that invoke a function asynchronously.
ASYNC_EVENT_TYPES = frozenset(["CloudWatchEvent", "EventBridgeRule", "IoTRule", "S3", "Schedule", "SNS"])

# Event types backed by stream-based event source mappings, which support an
# on-failure destination. SQS polling events do not support that property.
STREAM_EVENT_TYPES = frozenset(["Kinesis", "DynamoDB"])

# Event types whose generated EventBridge target supports a dead-letter queue.
# ScheduleV2 targets use a different mechanism and are intentionally excluded.
DLQ_EVENT_TYPES = frozenset(["CloudWatchEvent", "EventBridgeRule", "Schedule"])

# SAM-managed destination types that do not need an explicit destination ARN.
SAM_MANAGED_DESTINATION_TYPES = frozenset(["SQS", "SNS"])


def get_functions(cfn) -> Dict[str, dict]:
    """
    Return all SAM functions in the template, keyed by logical ID.
    """

    return cfn.get_resources([SAM_FUNCTION_TYPE])


def get_apis(cfn) -> Dict[str, dict]:
    """
    Return all explicit SAM REST APIs in the template, keyed by logical ID.
    """

    return cfn.get_resources([SAM_API_TYPE])


def get_http_apis(cfn) -> Dict[str, dict]:
    """
    Return all explicit SAM HTTP APIs in the template, keyed by logical ID.
    """

    return cfn.get_resources([SAM_HTTP_API_TYPE])


def iter_events(logical_id: str, function: dict) -> Iterator[Tuple[str, str, dict, list]]:
    """
    Yield (name, type, properties, path) for each event of a SAM function.

    The path points to the original event declaration in the template so
    findings and resource-level suppressions stay attached to the function.
    """

    events = function.get("Properties", {}).get("Events", {})

    if not isinstance(events, dict):
        return

    for name, event in events.items():
        if not isinstance(event, dict):
            continue

        event_type = event.get("Type", None)
        properties = event.get("Properties", {})

        if not isinstance(properties, dict):
            properties = {}

        yield name, event_type, properties, ["Resources", logical_id, "Properties", "Events", name]


# Event type -> (API kind, Globals section, explicit reference property)
_IMPLICIT_API_EVENTS = {
    "Api": ("rest", "Api", "RestApiId"),
    "HttpApi": ("http", "HttpApi", "ApiId"),
}


def _implicit_api_kind(event_type: str, properties: dict):
    """
    Classify an event as creating an implicit API, if applicable
    """

    config = _IMPLICIT_API_EVENTS.get(event_type, None)

    if config is None:
        return None

    kind, globals_key, reference_property = config

    # An explicit or unresolved API reference does not create an implicit API
    if reference_property in properties:
        return None

    return kind, globals_key


def get_implicit_apis(cfn) -> List[dict]:
    """
    Return read-only views of the implicit APIs created by function events.

    An `Api` event without `RestApiId` and an `HttpApi` event without `ApiId`
    make SAM create an implicit API, whose stage settings come from the
    applicable `Globals.Api` or `Globals.HttpApi` section. Multiple routes on
    one function yield a single view so each rule reports at most one finding
    per owning function.
    """

    globals_ = cfn.template.get("Globals", {})

    if not isinstance(globals_, dict):
        globals_ = {}

    views = {}

    for logical_id, function in get_functions(cfn).items():
        for name, event_type, properties, path in iter_events(logical_id, function):
            classification = _implicit_api_kind(event_type, properties)

            if classification is None:
                continue

            kind, globals_key = classification

            if (logical_id, kind) in views:
                continue

            global_properties = globals_.get(globals_key, {})
            if not isinstance(global_properties, dict):
                global_properties = {}

            views[(logical_id, kind)] = {
                "owner": logical_id,
                "event": name,
                "kind": kind,
                "properties": global_properties,
                "path": path,
            }

    return list(views.values())
