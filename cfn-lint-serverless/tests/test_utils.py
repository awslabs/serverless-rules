"""
Testing utility functions
"""


import pytest

from cfn_lint_serverless import utils

value_test_cases = [
    # str
    {"input": "MyString", "id": "MyString", "references": []},
    # Ref
    {"input": {"Ref": "MyResource"}, "id": "MyResource", "references": ["MyResource"]},
    # Fn::GetAtt
    {"input": {"Fn::GetAtt": ["MyResource", "Arn"]}, "id": "MyResource.Arn", "references": ["MyResource"]},
    # Fn::Join
    {"input": {"Fn::Join": ["/", ["ABC", "DEF"]]}, "id": "ABC/DEF", "references": []},
    # Fn::Join with references
    {
        "input": {"Fn::Join": ["/", ["ABC", {"Ref": "MyResource"}]]},
        "id": "ABC/MyResource",
        "references": ["MyResource"],
    },
    # Fn::Sub
    {"input": {"Fn::Sub": "abc-${MyResource}"}, "id": "abc-${MyResource}", "references": ["MyResource"]},
    # Fn::Sub with variables
    {"input": {"Fn::Sub": ["abc-${MyVar}", {"MyVar": "MyResource"}]}, "id": "abc-${MyVar}", "references": []},
    # Fn::Sub with variables and references
    {
        "input": {"Fn::Sub": ["abc-${MyVar}", {"MyVar": {"Ref": "MyResource"}}]},
        "id": "abc-${MyVar}",
        "references": ["MyResource"],
    },
]


@pytest.mark.parametrize("case", value_test_cases)
def test_value(case):
    """
    Test Value()
    """

    print(case)

    output = utils.Value(case["input"])

    assert case["id"] == output.id
    assert case["references"] == output.references
