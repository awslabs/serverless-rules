"""
CLI-level regression tests against the actual cfn-lint executable
"""

import json
import os
import shutil
import subprocess
import sys
from pathlib import Path

import pytest

EXPECTED_RULE_IDS = {
    "ES1001",
    "ES1005",
    "ES1006",
    "ES1007",
    "ES2000",
    "ES2003",
    "ES4000",
    "ES6000",
    "ES7000",
    "WS1000",
    "WS1002",
    "WS1003",
    "WS1004",
    "WS2001",
    "WS2002",
    "WS3000",
    "WS5000",
}

CLI_FIXTURE = os.path.join(os.path.dirname(__file__), "cli", "sam_app.yaml")


def cfn_lint_executable() -> str:
    """
    Locate the cfn-lint executable of the active environment
    """

    candidate = Path(sys.executable).parent / "cfn-lint"

    if candidate.exists():
        return str(candidate)

    found = shutil.which("cfn-lint")

    if found is None:
        pytest.skip("cfn-lint executable not found in the active environment")

    return found


def run_cfn_lint(*args) -> subprocess.CompletedProcess:
    """
    Run cfn-lint; findings produce a nonzero exit code by design
    """

    env = dict(os.environ)
    repo_root = Path(__file__).resolve().parent.parent
    env["PYTHONPATH"] = f"{repo_root}{os.pathsep}" + env.get("PYTHONPATH", "")

    return subprocess.run(
        [cfn_lint_executable(), *args],
        capture_output=True,
        text=True,
        env=env,
        check=False,
    )


def serverless_ids(output: str):
    ids = set()

    for line in output.splitlines():
        token = line.split()[0].rstrip(":") if line.split() else ""

        if len(token) > 1 and token[1] == "S" and token[0] in "WEI":
            ids.add(token)

    return ids


def test_rule_discovery_lists_all_serverless_rules():
    result = run_cfn_lint("--list-rules", "--append-rules", "cfn_lint_serverless.rules")

    assert result.returncode == 0, result.stderr

    found = serverless_ids(result.stdout)

    assert EXPECTED_RULE_IDS <= found


def test_direct_sam_cli_regression():
    result = run_cfn_lint(CLI_FIXTURE, "-a", "cfn_lint_serverless.rules", "-f", "json")

    # Findings cause a nonzero exit code; that is expected here
    assert result.returncode != 0
    findings = json.loads(result.stdout)
    by_rule = {}
    for finding in findings:
        rule = finding["Rule"]["Id"]
        by_rule.setdefault(rule, []).append(finding)

    # Expected direct-SAM findings
    assert "WS1000" in by_rule
    assert "WS1004" in by_rule
    assert "WS2002" in by_rule

    # WS1000 is only reported for the untraced function; the suppressed and
    # healthy functions do not appear
    ws1000_locations = {f["Filename"] + str(f["Location"]["Start"]) for f in by_rule["WS1000"]}
    assert len(by_rule["WS1000"]) == 1
    assert all("SuppressedFunction" not in f.get("Message", "") for f in by_rule["WS1000"])
    assert len(ws1000_locations) == 1

    # Inherited memory/timeout and declared API settings produce no findings
    for absent in ["ES1005", "ES1006", "ES2000", "ES2003", "E0002"]:
        assert absent not in by_rule

    # No rule processing crashed
    assert "E0002" not in result.stdout


def test_command_line_ignore_applies_to_sam_findings():
    result = run_cfn_lint(
        CLI_FIXTURE, "-a", "cfn_lint_serverless.rules", "-f", "json", "-i", "WS1000", "WS1004", "WS2002"
    )

    findings = json.loads(result.stdout) if result.stdout.strip() else []
    serverless_findings = [f for f in findings if f["Rule"]["Id"][1] == "S"]

    assert serverless_findings == []


def test_template_metadata_suppression():
    result = run_cfn_lint(CLI_FIXTURE, "-a", "cfn_lint_serverless.rules", "-f", "json")
    findings = json.loads(result.stdout)

    assert not [f for f in findings if "SuppressedFunction" in f.get("Message", "") and f["Rule"]["Id"] == "WS1000"]
