
CDK Example
===========

Example on how to use cfn-lint-serverless with [AWS CDK](https://docs.aws.amazon.com/cdk/latest/guide/home.html).

Usage
-----

```bash
# Generate CloudFormation templates
cdk synth > template.yaml

# Run cfn-lint against the CloudFormation template
cfn-lint template.yaml -a cfn_lint_serverless.rules
```