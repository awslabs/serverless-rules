Serverless Framework Example
============================

Example on how to use cfn-lint-serverless with the [Serverless Framework](https://www.serverless.com/).

Requirements
------------

- Node.js 20 or later
- Serverless Framework v4 (`npm install -g serverless`), validated with 4.43.0
- Serverless Framework v4 requires authentication for `sls package`: run `serverless login`, or set `SERVERLESS_ACCESS_KEY` or `SERVERLESS_LICENSE_KEY` for unattended runs

Usage
-----

```bash
# Generate CloudFormation templates
sls package

# Run cfn-lint against the CloudFormation template
cfn-lint .serverless/cloudformation-template-update-stack.json -a cfn_lint_serverless.rules
```
