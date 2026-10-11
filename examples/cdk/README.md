
CDK Example
===========

Example on how to use cfn-lint-serverless with [AWS CDK](https://docs.aws.amazon.com/cdk/latest/guide/home.html).

Requirements
------------

- Python 3.10 or later
- AWS CDK CLI v2 (`npm install -g aws-cdk`), validated with 2.1145.0
- Python dependencies from `requirements.txt` (`aws-cdk-lib` 2.x and `constructs` 10.x)

Usage
-----

```bash
# Install the Python dependencies (for example in a virtual environment)
pip install -r requirements.txt

# Generate CloudFormation templates
cdk synth > template.yaml

# Run cfn-lint against the CloudFormation template
cfn-lint template.yaml -a cfn_lint_serverless.rules
```
