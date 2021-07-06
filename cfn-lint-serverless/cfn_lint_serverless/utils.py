"""
Utilities 
"""


import re
from typing import List, Tuple, Union

SUB_PATTERN = re.compile(r"\${(?P<ref>[^}]+)}")


class Value:

    id = ""  # noqa: VNE003
    references = None

    def __init__(self, value: Union[dict, str]):
        """
        Parse a CloudFormation value

        This handles intrinsic functions, such as 'Fn::Sub' and 'Fn::Join' and
        returns an object that contains both a uniquely identifiable string and
        references to other resources.
        """

        self.references = []

        # String
        if isinstance(value, str):
            self.id = value

        # Not a dict - return an error here
        elif not isinstance(value, dict):
            raise ValueError("'value' should be of type str or dict, got '%s'" % type(value))

        # 'Ref' intrinsic function
        elif "Ref" in value:
            self.id, self.references = self._get_from_ref(value["Ref"])

        # 'Fn::GetAtt' intrinsic function
        elif "Fn::GetAtt" in value:
            self.id, self.references = self._get_from_getatt(value["Fn::GetAtt"])

        # 'Fn::Join' intrisic function
        elif "Fn::Join" in value:
            self.id, self.references = self._get_from_join(value["Fn::Join"])

        # 'Fn::Sub' intrisic function
        elif "Fn::Sub" in value:
            self.id, self.references = self._get_from_sub(value["Fn::Sub"])

    def _get_from_ref(self, value: str) -> Tuple[str, List[str]]:
        """
        Return the name and references from a 'Ref' intrinsic function
        """

        return [value, [value]]

    def _get_from_getatt(self, value: list) -> Tuple[str, List[str]]:
        """
        Return the name and references from a 'Fn::GetAtt' intrinsic function
        """

        id_ = ".".join(value)
        references = [value[0]]

        return (id_, references)

    def _get_from_join(self, value: list) -> Tuple[str, List[str]]:
        """
        Return the name and references from a 'Fn::Join' intrinsic function
        """

        delimiter = value[0]
        # Using Value() here to get nested references
        sub_values = [Value(v) for v in value[1]]

        id_ = delimiter.join([v.id for v in sub_values])
        references = []
        for sub_value in sub_values:
            references.extend(sub_value.references)

        return (id_, references)

    def _get_from_sub(self, value: Union[str, list]) -> Tuple[str, List[str]]:
        """
        Return the name and references from a 'Fn::Sub' intrinsic function
        """

        pattern = value
        variables = {}

        if isinstance(value, list):
            pattern = value[0]
            # Using Value() here to get nested references
            variables = {k: Value(v) for k, v in value[1].items()}

        references = []

        for match in SUB_PATTERN.findall(pattern):
            if match in variables:
                references.extend(variables[match].references)
            else:
                references.append(match)

        return (pattern, references)
