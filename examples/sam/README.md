SAM Example
===========

Example on how to use cfn-lint-serverless with the [AWS Serverless Application Model (SAM)](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/what-is-sam.html).

Requirements
------------

- Python 3.10 or later with `cfn-lint` 1.57.0 or later and `cfn-lint-serverless` installed
- Optionally, [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html) 1.x, validated with 1.167.0 or later, for `sam validate`

Usage
-----

```bash
# Optionally validate the template with the SAM CLI
sam validate -t template.yaml

# Run cfn-lint against the SAM template
cfn-lint template.yaml -a cfn_lint_serverless.rules
```

Serverless Rules checks `AWS::Serverless::Function`, `AWS::Serverless::Api`, and `AWS::Serverless::HttpApi` resources directly, including settings inherited from `Globals`, so no `sam build` or deployment step is required.
