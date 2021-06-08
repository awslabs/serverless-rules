---
title: Homepage
description: Serverless Rules
---

__Serverless Rules__ are a compilation of rules to validate infrastructure-as-code template against recommended practices. This currently provides a module for [cfn-lint](https://github.com/aws-cloudformation/cfn-python-lint) and a plugin for [tflint](https://github.com/terraform-linters/tflint).

You can use those rules to get quick feedback on recommended practices while building a serverless application, as part of an automated code review process, or as guardrails before deploying to production.

!!! danger "Public preview"
    This project is currently in __public preview__ to get feedback from the serverless community. APIs, tools, and rules might change between the beginning of public preview and version 1.

## Current modules and plugins

* The `cfn-lint` module supports checking CloudFormation and [Serverless Application Model (SAM)](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/what-is-sam.html) templates. It also supports templates defined with the [AWS Cloud Development Kit (CDK)](https://docs.aws.amazon.com/cdk/latest/guide/home.html) and the [Serverless Framework](https://www.serverless.com/) by synthesizing to CloudFormation.
* The `tflint` plugin supports checking Terraform configuration files.

If you would like native support for other frameworks, please consider adding a 👍 reaction to the corresponding comment in [this GitHub issue](https://github.com/awslabs/serverless-rules/issues/23).

## Usage

To learn how to use Serverless Rules, see the detailed usage guide for each plugin:

* [With `cfn-lint`](cfn-lint.md)
* [With `tflint`](tflint.md)

## Frequently asked questions

### How is this different from `cfn-lint` or `tflint-ruleset-aws`?

`cfn-lint` and `tflint` main goals are to find possible errors in templates and configuration files before you try to deploy resources to AWS. By comparison, __Serverless Rules__ goes one step further by providing prescriptive guidance based on the [AWS Well-Architected pillars](https://aws.amazon.com/architecture/well-architected/?wa-lens-whitepapers.sort-by=item.additionalFields.sortDate&wa-lens-whitepapers.sort-order=desc).

For example, you can define AWS Lambda functions without tracing enabled, which is valid for both CloudFormation and Terraform. This project adds [a rule on Lambda tracing](rules/lambda.md#tracing) as a recommended practice for Operational Excellence.

### Why use `cfn-lint` instead of `cfn-guard` or `cfn_nag`?

[`cfn-guard`](https://github.com/aws-cloudformation/cloudformation-guard) provides developers with a simplified language to define polices and validate JSON- or YAML- formatted documents. You can use that tool to create your own rules that match your compliance needs or internal recommended practices. [`cfn_nag`](https://github.com/stelligent/cfn_nag) provides a set of default rules with a strong focus on security. You can [extend it](https://stelligent.com/2020/02/27/custom-rule-distribution-enhancements-for-cfn_nag/) in ruby to create custom rules.

By comparison, [`cfn-lint`](https://github.com/aws-cloudformation/cfn-lint) provides a set of default rules focused on validating templates against the [CloudFormation resource specification](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-resource-specification.html) and you can create your own rule with Python modules.

For __Serverless Rules__, using a programming language like Python or Ruby gives more flexibility when defining complex rules integrating multiple resources, such as checking if [a Lambda function has a log group with retention configured](rules/lambda.md#log-retention).