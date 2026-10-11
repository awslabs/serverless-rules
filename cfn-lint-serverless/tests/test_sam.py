"""
Focused tests for the read-only SAM helper views
"""

import copy

import pytest
from cfnlint.template import Template

from cfn_lint_serverless import sam


def make_cfn(template: dict) -> Template:
    """
    Build a transformed cfn-lint Template, as rules would receive it
    """

    cfn = Template("test.yaml", template, ["eu-west-1"])
    cfn.transform()

    return cfn


def test_get_functions_returns_sam_functions_with_effective_properties():
    """
    Globals are already merged by the host when rules inspect resources
    """

    cfn = make_cfn(
        {
            "Transform": "AWS::Serverless-2016-10-31",
            "Globals": {"Function": {"Runtime": "python3.12", "Timeout": 30, "Tags": {"env": "prod"}}},
            "Resources": {
                "Fn": {
                    "Type": "AWS::Serverless::Function",
                    "Properties": {"CodeUri": ".", "Handler": "main.handler", "Timeout": 60},
                },
                "Native": {"Type": "AWS::Lambda::Function", "Properties": {}},
            },
        }
    )

    functions = sam.get_functions(cfn)

    assert list(functions.keys()) == ["Fn"]
    properties = functions["Fn"]["Properties"]
    assert properties["Runtime"] == "python3.12"
    # Local value overrides the inherited one
    assert properties["Timeout"] == 60
    # Maps merge recursively
    assert properties["Tags"] == {"env": "prod"}


def test_get_apis_and_http_apis_split_by_type():
    cfn = make_cfn(
        {
            "Transform": "AWS::Serverless-2016-10-31",
            "Resources": {
                "RestApi": {"Type": "AWS::Serverless::Api", "Properties": {"StageName": "prod"}},
                "HttpApi": {"Type": "AWS::Serverless::HttpApi", "Properties": {}},
                "Fn": {"Type": "AWS::Serverless::Function", "Properties": {"CodeUri": ".", "Handler": "h"}},
            },
        }
    )

    assert list(sam.get_apis(cfn).keys()) == ["RestApi"]
    assert list(sam.get_http_apis(cfn).keys()) == ["HttpApi"]


def test_iter_events_yields_original_paths():
    cfn = make_cfn(
        {
            "Transform": "AWS::Serverless-2016-10-31",
            "Resources": {
                "Fn": {
                    "Type": "AWS::Serverless::Function",
                    "Properties": {
                        "CodeUri": ".",
                        "Handler": "main.handler",
                        "Events": {
                            "Topic": {"Type": "SNS", "Properties": {"Topic": "arn:aws:sns:us-east-1:111122223333:t"}},
                            "Broken": {"Properties": {}},
                        },
                    },
                }
            },
        }
    )

    events = list(sam.iter_events("Fn", sam.get_functions(cfn)["Fn"]))

    assert [(name, event_type) for name, event_type, _, _ in events] == [("Topic", "SNS"), ("Broken", None)]
    assert events[0][2] == {"Topic": "arn:aws:sns:us-east-1:111122223333:t"}
    assert events[0][3] == ["Resources", "Fn", "Properties", "Events", "Topic"]


def test_iter_events_handles_missing_or_unsupported_shapes():
    function = {"Type": "AWS::Serverless::Function", "Properties": {"CodeUri": "."}}

    assert list(sam.iter_events("Fn", function)) == []
    assert list(sam.iter_events("Fn", {"Properties": {"Events": []}})) == []
    assert list(sam.iter_events("Fn", {"Properties": {"Events": {"Odd": "not-a-mapping"}}})) == []


def test_get_implicit_apis_from_events_and_globals():
    cfn = make_cfn(
        {
            "Transform": "AWS::Serverless-2016-10-31",
            "Globals": {"Api": {"TracingEnabled": True}},
            "Resources": {
                "Fn": {
                    "Type": "AWS::Serverless::Function",
                    "Properties": {
                        "CodeUri": ".",
                        "Handler": "main.handler",
                        "Events": {
                            "Root": {"Type": "Api", "Properties": {"Method": "GET", "Path": "/"}},
                            "Other": {"Type": "Api", "Properties": {"Method": "POST", "Path": "/other"}},
                            "Hook": {"Type": "HttpApi", "Properties": {"Method": "POST", "Path": "/hook"}},
                        },
                    },
                },
                "ExplicitFn": {
                    "Type": "AWS::Serverless::Function",
                    "Properties": {
                        "CodeUri": ".",
                        "Handler": "main.handler",
                        "Events": {
                            "Root": {
                                "Type": "Api",
                                "Properties": {"Method": "GET", "Path": "/", "RestApiId": {"Ref": "Api"}},
                            }
                        },
                    },
                },
                "Api": {"Type": "AWS::Serverless::Api", "Properties": {"StageName": "prod"}},
            },
        }
    )

    views = sam.get_implicit_apis(cfn)

    # Repeated routes on one function produce a single REST view; the explicit
    # RestApiId reference does not create an implicit API
    assert {(v["owner"], v["kind"]) for v in views} == {("Fn", "rest"), ("Fn", "http")}

    rest = next(v for v in views if v["kind"] == "rest")
    assert rest["properties"] == {"TracingEnabled": True}
    assert rest["event"] == "Root"
    assert rest["path"] == ["Resources", "Fn", "Properties", "Events", "Root"]

    http = next(v for v in views if v["kind"] == "http")
    assert http["properties"] == {}


def test_ignore_globals_is_applied_by_the_host():
    """
    A resource-level IgnoreGlobals attribute excludes inherited values, and
    already merged lists are not duplicated
    """

    cfn = make_cfn(
        {
            "Transform": "AWS::Serverless-2016-10-31",
            "Globals": {
                "Function": {
                    "Tracing": "Active",
                    "Timeout": 30,
                    "Policies": ["AWSLambdaBasicExecutionRole"],
                }
            },
            "Resources": {
                "Fn": {
                    "Type": "AWS::Serverless::Function",
                    "IgnoreGlobals": ["Tracing"],
                    "Properties": {
                        "CodeUri": ".",
                        "Handler": "main.handler",
                        "Policies": [
                            {
                                "Version": "2012-10-17",
                                "Statement": {"Effect": "Allow", "Action": "s3:GetObject", "Resource": "*"},
                            }
                        ],
                    },
                }
            },
        }
    )

    properties = sam.get_functions(cfn)["Fn"]["Properties"]

    # Inherited value excluded by IgnoreGlobals
    assert "Tracing" not in properties
    # Non-ignored globals are inherited
    assert properties["Timeout"] == 30
    # Already merged lists are not duplicated
    assert properties["Policies"] == [
        "AWSLambdaBasicExecutionRole",
        {"Version": "2012-10-17", "Statement": {"Effect": "Allow", "Action": "s3:GetObject", "Resource": "*"}},
    ]


def test_helpers_do_not_change_the_host_template():
    template = {
        "Transform": "AWS::Serverless-2016-10-31",
        "Globals": {"Function": {"Timeout": 30}, "Api": {"TracingEnabled": True}},
        "Resources": {
            "Fn": {
                "Type": "AWS::Serverless::Function",
                "Properties": {
                    "CodeUri": ".",
                    "Handler": "main.handler",
                    "Events": {"Root": {"Type": "Api", "Properties": {"Method": "GET", "Path": "/"}}},
                },
            }
        },
    }

    cfn = make_cfn(template)
    before = copy.deepcopy(cfn.template)

    sam.get_functions(cfn)
    sam.get_apis(cfn)
    sam.get_http_apis(cfn)
    sam.get_implicit_apis(cfn)
    for logical_id, function in sam.get_functions(cfn).items():
        list(sam.iter_events(logical_id, function))

    assert cfn.template == before
