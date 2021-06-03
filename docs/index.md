---
title: Homepage
description: Serverless Rules
---

__Serverless Rules__ are a compilation of rules to validate infrastructure as code template against recommended practices. This currently provides a module for [cfn-lint](https://github.com/aws-cloudformation/cfn-python-lint) and a plugin for [tflint](https://github.com/terraform-linters/tflint).

!!! danger "Public preview"
    This project is currently in __public preview__ to get feedback from the serverless community. APIs, tools, and rules might change between the beginning of public preview and version 1.

### Current modules and plugins

* The `cfn-lint` module supports checking CloudFormation and [Serverless Application Model (SAM)](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/what-is-sam.html) templates. It also supports templates defined with the [AWS Cloud Development Kit (CDK)](https://docs.aws.amazon.com/cdk/latest/guide/home.html) and the [Serverless Framework](https://www.serverless.com/) by synthesizing to CloudFormation.
* The `tflint` plugin supports checking Terraform configuration files.

If you would like native support for other frameworks, please consider adding a 👍 reaction to the corresponding comment in [this GitHub issue](https://github.com/awslabs/serverless-rules/issues/23).

## Using the `cfn-lint` module

### Installation

Serverless Rules for `cfn-lint` is available through the `cfn-lint-serverless` Python module in PyPi. You can use pip to install the module:

```bash
pip install cfn-lint cfn-lint-serverless
```

### Usage

You can now instruct `cfn-lint` to use Serverless Rules installed previous via `--append-rules`, or `-a` for short:

```bash
cfn-lint my_template.yaml -a cfn_lint_serverless.rules
```

You can also try with a Serverless Application Model (SAM) example provided in this repository. In the root folder of the repository, you can run:

```bash
cfn-lint examples/sam/template.yaml -a cfn_lint_serverless.rules
```

#### Sample outputs

If the template fulfill the requirements for all rules, `cfn-lint` will return an empty output:

```
$ cfn-lint template.yaml -a cfn_lint_serverless.rules
$
```

If there are potential improvements to your templates, `cfn-lint` will output recommendations:

```
$ cfn-lint template.yaml -a cfn_lint_serverless.rules
ES4000 EventBridge rule ConsumerFunctionTrigger should have a DeadLetterConfig.Arn property for all its Targets.
template.yaml:5:1

WS1000 Lambda function ConsumerFunction should have TracingConfig.Mode set to 'Active'.
template.yaml:7:3

WS1004 Lambda function ConsumerFunction does not have a corresponding log group with a Retention property
template.yaml:7:3
$
```

## Using the `tflint` plugin

### Installation

This plugin depends on [tflint being installed](https://github.com/terraform-linters/tflint#installation). If you are using `tflint` version 0.29 or newer, you can leverage the `tflint --init` command to automatically install the plugin. Otherwise, you will need to download the `tflint-ruleset-aws-serverless` binary corresponding to your system from the [releases page](https://github.com/awslabs/serverless-rules/releases).

You can enable the Serverless Rules plugin by adding a plugin section in the `.tflint.hcl` file in your project:

```terraform
plugin "aws-serverless" {
  enabled = true
  # Replace this with the latest version
  version = "0.1.4"
  source = "github.com/aws-samples/serverless-rules/tflint-ruleset-aws-serverless"
}
```

### Usage

You can now run the `tflint` command, which will automatically use the Serverless Rules plugin:

```bash
tflint
```

You can also try with a Terraform example provided in this repository. From the root folder of the repository, you can run:

```bash
cd examples/tflint/
tflint 
```

#### Sample outputs

If the Terraform configuration files fulfill the requirements for all the rules, `tflint` will return an empty output:

```
$ tflint
$
```

If there are potential improvements to your templates, `tflint` will output recommendations:

```
$ tflint
1 issue(s) found:

Warning: "tracing_config" is not present. (aws_lambda_function_tracing_rule)

    on main.tf line 20:

$
```