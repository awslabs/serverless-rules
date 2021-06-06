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