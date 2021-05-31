Serverless Rules
================

The __Serverless Rules__ are a compilation of rules to validate infrastructure as code template against recommended practices. This currently provides a module for [cfn-lint](https://github.com/aws-cloudformation/cfn-python-lint) and a plugin for [tflint](https://github.com/terraform-linters/tflint).

You can find a list of currently supported rules [in the documentation](docs/rules.md).

## Usage guide

__For `cfn-lint`__

```bash
# To use this project with cfn-lint, you need to install the cfn-lint-serverless python module
pip install cfn-lint cfn-lint-serverless

# You can then run cfn-lint with the additional rules
cfn-lint my_template.yaml -a cfn_lint_serverless.rules
```

__For `tflint`__

You can add the `aws-serverless` plugin to your .tflint.hcl configuration file:

```hcl
plugin "aws-serverless" {
  enabled = true
  version = "0.1.0"
  source = "github.com/aws-samples/serverless-rules/tflint-ruleset-aws-serverless"
}
```

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) to learn how to contribute to this project.

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This library is licensed under the MIT-0 License. See the [LICENSE](./LICENSE) file.

