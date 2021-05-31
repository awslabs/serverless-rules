Serverless Framework Example
============================

Example on how to use cfn-lint-serverless with the [Serverless Framework](https://www.serverless.com/).

Usage
-----

```bash
# Generate CloudFormation templates
sls package

# Run cfn-lint against the CloudFormation template
cfn-lint .serverless/cloudformation-template-update-stack.json -a cfn_lint_serverless.rules
```