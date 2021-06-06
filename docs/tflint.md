`tflint` plugin
===============

## Installation

This plugin depends on [tflint](https://github.com/terraform-linters/tflint#installation). If you use `tflint` version 0.29 or newer, you can leverage the `tflint --init` command to automatically install the plugin. Otherwise, you will need to download the `tflint-ruleset-aws-serverless` binary corresponding to your system from the [releases page](https://github.com/awslabs/serverless-rules/releases).

You can enable the Serverless Rules plugin by adding a plugin section in the `.tflint.hcl` file in your project:

```terraform
plugin "aws-serverless" {
  enabled = true
  version = "0.1.6"
  source = "github.com/awslabs/serverless-rules"
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

If the Terraform configuration files fulfill the requirements for all the rules, `tflint` will return an empty output. Otherwise, `tflint` will output recommendations.

=== "Matching template"

    ```bash
    $ tflint
    $
    ```

=== "With recommendations"

    ```bash
    $ tflint
    1 issue(s) found:

    Warning: "tracing_config" is not present. (aws_lambda_function_tracing_rule)

        on main.tf line 20:

    $
    ```

### Ignoring rules

Serverless Rules is a compilation of __recommended practices__ and you might have a valid reason to ignore specific rules. While we recommend that you keep Error-level rules enabled, all other rules contain explanations on when you can safely ignore those rules. See [the Lambda Tracing rule](rules/lambda/#tracing) for an example of such explanation.

Rules in `tflint` can be disabled either through the `--disable-rule` command-line argument or with the `.tflint.hcl` configuration file in the current working directory. See the [`tflint` user guide](https://github.com/terraform-linters/tflint/blob/master/docs/user-guide/config.md) for more information.

=== "Command line"

    ```bash
    # Disable the aws_lambda_function_tracing_rule rule
    tflint --disable-rule aws_lambda_function_tracing_rule
    ```

=== ".tflint.hcl"

    ```terraform
    plugin "aws-serverless" {
      enabled = true
      version = "0.1.6"
      source = "github.com/awslabs/serverless-rules"
    }

    # Disable the aws_lambda_function_tracing_rule rule
    rule "aws_lambda_function_tracing_rule" {
      enabled = false
    }
    ```