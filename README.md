Serverless Rules
================

The __Serverless Rules__ are a compilation of rules to validate infrastructure as code template against recommended practices. This currently provides a module for [cfn-lint](https://github.com/aws-cloudformation/cfn-python-lint) and a plugin for [tflint](https://github.com/terraform-linters/tflint).

You can use those rules to get quick feedback on recommended practices while building a serverless application, as part of automated code review process, or as guardrails before deploying to production.

<p align="center">

**[📜Documentation](https://awslabs.github.io/serverless-rules/)** | **[🐍PyPi](https://pypi.org/project/cfn-lint-serverless/)**

</p>

__PUBLIC PREVIEW__: this project is currently in __public preview__ to get feedback from the serverless community. APIs, tools, and rules might change between the beginning of public preview and version 1.

You can find a list of currently supported rules [in the documentation](https://awslabs.github.io/serverless-rules/rules/).

## Usage guide
### cfn-lint

To get started with Serverless Rules and [cfn-lint](https://github.com/aws-cloudformation/cfn-lint), install `cfn-lint-serverless` module: `pip install cfn-lint cfn-lint-serverless`

You can now instruct `cfn-lint` to use Serverless Rules module installed previously via `--append-rules` or `-a` for short:

```bash
cfn-lint my_template.yaml -a cfn_lint_serverless.rules
```

You can try with a Serverless Application Model (SAM) example provided in this repository by running:

```bash
cfn-lint examples/sam/template.yaml -a cfn_lint_serverless.rules
```

### tflint

This plugin depends on [tflint](https://github.com/terraform-linters/tflint#installation). If you use `tflint` version 0.29 or newer, you can leverage the `tflint --init` command to automatically install the plugin. Otherwise, you will need to download the `tflint-ruleset-aws-serverless` binary corresponding to your system from the [releases page](https://github.com/awslabs/serverless-rules/releases).

You can enable the Serverless Rules plugin by adding a plugin section in the `.tflint.hcl` file in your project:

```terraform
plugin "aws-serverless" {
  enabled = true
  version = "0.2.1"
  source = "github.com/awslabs/serverless-rules"
}
```

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) to learn how to contribute to this project.

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This library is licensed under the MIT-0 License. See the [LICENSE](./LICENSE) file.
