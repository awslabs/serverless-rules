# cfn-lint-serverless

Ruleset for [cfn-lint](https://github.com/aws-cloudformation/cfn-lint) to validate CloudFormation templates for serverless applications against recommended best practices.

## Installation

```bash
pip install cfn-lint cfn-lint-serverless
```

## Usage

Run cfn-lint with the serverless rules module:

```bash
cfn-lint template.yaml -a cfn_lint_serverless.rules
```

## Supported Rules

This module provides validation rules for various AWS serverless resources:

- Lambda Functions
- API Gateway
- Step Functions
- SQS
- SNS
- EventBridge
- AppSync

For a detailed list of rules, refer to the [documentation](https://awslabs.github.io/serverless-rules/rules/).

## Examples

Try it with the examples provided in the repository:

```bash
# For SAM templates
cfn-lint examples/sam/template.yaml -a cfn_lint_serverless.rules

# For Serverless Framework templates
cfn-lint examples/serverless-framework/template.yaml -a cfn_lint_serverless.rules
```

## Contributing

Contributions are welcome! Please see the [CONTRIBUTING.md](../CONTRIBUTING.md) file for guidelines.

## License

This project is licensed under the MIT-0 License. See the [LICENSE](../LICENSE) file for details.
