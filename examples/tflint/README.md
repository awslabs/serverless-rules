TFLint Example
==============

Example on how to use the Serverless Rules ruleset with [TFLint](https://github.com/terraform-linters/tflint) and Terraform.

Requirements
------------

- Terraform 1.13 or later (validated with 1.16.5)
- TFLint v0.64.0 or later
- Providers: `hashicorp/aws` 6.68.0, `hashicorp/random` 3.9.1 (see `main.tf` and `.terraform.lock.hcl`)

Usage
-----

```bash
# Initialize providers locally (no backend, no deployment)
terraform init -backend=false

# Validate the Terraform configuration
terraform validate

# Install the configured TFLint plugins (see .tflint.hcl)
tflint --init

# Lint the configuration
tflint
```

This example intentionally omits `memory_size` and `timeout` on the Lambda function, so `tflint` reports the `aws_lambda_function_default_memory` and `aws_lambda_function_default_timeout` findings from the Serverless Rules plugin.

To lint with a locally built plugin instead of a released one, run `make test-example`, which builds the plugin from `../../tflint-ruleset-aws-serverless` into a local `.tflint.d` plugin directory.
