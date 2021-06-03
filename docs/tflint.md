`tflint` plugin
===============

## Installation

This plugin depends on [tflint being installed](https://github.com/terraform-linters/tflint#installation). If you are using `tflint` version 0.29 or newer, you can leverage the `tflint --init` command to automatically install the plugin. Otherwise, you will need to download the `tflint-ruleset-aws-serverless` binary corresponding to your system from the [releases page](https://github.com/awslabs/serverless-rules/releases).

You can enable the Serverless Rules plugin by adding a plugin section in the `.tflint.hcl` file in your project:

```terraform
plugin "aws-serverless" {
  enabled = true
  # Replace this with the latest version
  version = "0.1.5"
  source = "github.com/awslabs/serverless-rules/tflint-ruleset-aws-serverless"
}
```

## Usage

You can now run the `tflint` command, which will automatically use the Serverless Rules plugin:

```bash
tflint
```

You can also try with a Terraform example provided in this repository. From the root folder of the repository, you can run:

```bash
cd examples/tflint/
tflint 
```

### Sample outputs

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