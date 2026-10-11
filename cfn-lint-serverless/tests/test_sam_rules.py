"""
Integration tests for direct SAM rule behavior
"""

import yaml
import cfnlint.core
import cfnlint.decode.cfn_yaml
import pytest


@pytest.fixture(scope="session")
def rules():
    return cfnlint.core.get_rules(["cfn_lint_serverless.rules"], [], [])


def run_template(rules, tmp_path, template: dict):
    """
    Run all rules against an in-memory template and return the matches
    """

    filename = str(tmp_path / "template.yaml")
    with open(filename, "w") as fp:
        yaml.safe_dump(template, fp)

    template_dict = cfnlint.decode.cfn_yaml.load(filename)

    return cfnlint.core.run_checks(filename, template_dict, rules, ["eu-west-1"])


def implicit_api_template(events_function_a, functions_extra=None, globals_=None):
    template = {
        "Transform": "AWS::Serverless-2016-10-31",
        "Globals": globals_ or {"Function": {"Runtime": "python3.12"}},
        "Resources": {
            "FunctionA": {
                "Type": "AWS::Serverless::Function",
                "Properties": {
                    "CodeUri": ".",
                    "Handler": "main.handler",
                    "Events": events_function_a,
                },
            },
        },
    }

    if functions_extra:
        template["Resources"].update(functions_extra)

    return template


def test_shared_implicit_api_reports_per_owner(rules, tmp_path):
    """
    Two functions sharing the implicit REST API each receive the finding
    """

    template = implicit_api_template(
        {"Root": {"Type": "Api", "Properties": {"Method": "GET", "Path": "/a"}}},
        functions_extra={
            "FunctionB": {
                "Type": "AWS::Serverless::Function",
                "Properties": {
                    "CodeUri": ".",
                    "Handler": "main.handler",
                    "Events": {"Root": {"Type": "Api", "Properties": {"Method": "GET", "Path": "/b"}}},
                },
            }
        },
    )

    matches = [m for m in run_template(rules, tmp_path, template) if m.rule.id == "WS2002"]

    assert len(matches) == 2
    assert {m.path[1] for m in matches} == {"FunctionA", "FunctionB"}


def test_multiple_routes_produce_one_finding_per_owner(rules, tmp_path):
    """
    Several Api events on one function produce a single finding per rule
    """

    template = implicit_api_template(
        {
            "Root": {"Type": "Api", "Properties": {"Method": "GET", "Path": "/"}},
            "Other": {"Type": "Api", "Properties": {"Method": "POST", "Path": "/other"}},
        }
    )

    matches = run_template(rules, tmp_path, template)

    for rule_id in ["ES2000", "ES2003", "WS2002"]:
        rule_matches = [m for m in matches if m.rule.id == rule_id]
        assert len(rule_matches) == 1
        assert rule_matches[0].path[:2] == ["Resources", "FunctionA"]


def test_explicit_api_reference_creates_no_implicit_findings(rules, tmp_path):
    template = implicit_api_template(
        {
            "Root": {
                "Type": "Api",
                "Properties": {"Method": "GET", "Path": "/", "RestApiId": {"Ref": "Api"}},
            }
        }
    )
    template["Resources"]["Api"] = {
        "Type": "AWS::Serverless::Api",
        "Properties": {
            "StageName": "prod",
            "AccessLogSetting": {
                "DestinationArn": "arn:aws:logs:eu-west-1:123456789012:log-group:my-log-group",
                "Format": '{"request_id": "$context.requestId"}',
            },
            "MethodSettings": [
                {
                    "HttpMethod": "*",
                    "ResourcePath": "/*",
                    "ThrottlingBurstLimit": 1000,
                    "ThrottlingRateLimit": 10,
                }
            ],
            "TracingEnabled": True,
        },
    }

    matches = run_template(rules, tmp_path, template)

    assert not [m for m in matches if m.rule.id in ["ES2000", "ES2003", "WS2001", "WS2002"]]


def test_intrinsic_log_format_is_not_parsed(rules, tmp_path):
    """
    An unresolved intrinsic access log format does not raise and is not
    treated as a passing JSON format
    """

    template = {
        "Transform": "AWS::Serverless-2016-10-31",
        "Resources": {
            "Api": {
                "Type": "AWS::Serverless::Api",
                "Properties": {
                    "StageName": "prod",
                    "AccessLogSetting": {
                        "DestinationArn": "arn:aws:logs:eu-west-1:123456789012:log-group:my-log-group",
                        "Format": {"Fn::Sub": '{"request_id": "${AWS::RequestId}"}'},
                    },
                    "MethodSettings": [
                        {
                            "HttpMethod": "*",
                            "ResourcePath": "/*",
                            "ThrottlingBurstLimit": 1000,
                            "ThrottlingRateLimit": 10,
                        }
                    ],
                    "TracingEnabled": True,
                },
            }
        },
    }

    matches = run_template(rules, tmp_path, template)
    match_ids = [m.rule.id for m in matches]

    assert "E0002" not in match_ids
    assert "WS2001" not in match_ids


def test_async_event_finding_points_at_event_declaration(rules, tmp_path):
    template = {
        "Transform": "AWS::Serverless-2016-10-31",
        "Resources": {
            "Function": {
                "Type": "AWS::Serverless::Function",
                "Properties": {
                    "CodeUri": ".",
                    "Handler": "main.handler",
                    "Runtime": "python3.12",
                    "Events": {
                        "Topic": {
                            "Type": "SNS",
                            "Properties": {"Topic": "arn:aws:sns:us-east-1:111122223333:my-topic"},
                        }
                    },
                },
            }
        },
    }

    matches = [m for m in run_template(rules, tmp_path, template) if m.rule.id == "ES1007"]

    assert len(matches) == 1
    assert matches[0].path == ["Resources", "Function", "Properties", "Events", "Topic"]
