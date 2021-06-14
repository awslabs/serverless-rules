`tflint` plugin
===============

## Installation

This plugin depends on [tflint](https://github.com/terraform-linters/tflint#installation). If you use `tflint` version 0.29 or newer, you can leverage the `tflint --init` command to automatically install the plugin. Otherwise, you will need to download the `tflint-ruleset-aws-serverless` binary corresponding to your system from the [releases page](https://github.com/awslabs/serverless-rules/releases).

You can enable the Serverless Rules plugin by adding a plugin section in the `.tflint.hcl` file in your project:

```terraform
plugin "aws-serverless" {
  enabled = true
  version = "0.1.9"
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

Serverless Rules is a set of recommended practices. 

We recommend you to keep Error-level rules enabled. Non-error rules, for example [Lambda Tracing](rules/lambda/tracing.md), contain detailed scenarios on when it’s safe to ignore them.

When needed, you can ignore any specific rule that doesn’t match your environment.

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
      version = "0.1.9"
      source = "github.com/awslabs/serverless-rules"
    }

    # Disable the aws_lambda_function_tracing_rule rule
    rule "aws_lambda_function_tracing_rule" {
      enabled = false
    }
    ```

## Continuous integration

You can use Serverless Rules and `tflint` with your continuous integration tool to automatically check CloudFormation templates with rules from this project. For example, you can validate on pull requests, merge to your main branch, or before deploying to production.

If there are any issues with your template, `tflint` will return a non-zero error code.

### AWS CodeBuild

Assuming that you are storing your terraform configuration files and a `.tflint.hcl` file at the root of your repository, you can create a [buildspec file](https://docs.aws.amazon.com/codebuild/latest/userguide/build-spec-ref.html) such as this one.

!!! important

    Make sure that you include the __aws-serverless__ plugin into your `.tflint.hcl` configuration file, otherwise tflint will not install this ruleset. See [Installation](#installation) for more information.

=== "Sample buildspec"

    ```yaml
    env:
      variables:
        # TODO: replace "v0.29.0" with the latest version of tflint
        TFLINT_VERSION: "0.29.0"
        TFLINT_OS: "amd64"

    phases:
      install:
        commands:
          # Install tflint
          - wget https://github.com/terraform-linters/tflint/releases/download/v${TFLINT_VERSION}/tflint_linux_${TFLINT_OS}.zip -O tflint.zip
          - unzip tflint.zip
          # Install tflint plugins
          - ./tflint --init
      pre_build:
        commands:
          - ./tflint
    ```

=== "With JUnit report"

    ```yaml
    env:
      variables:
        # TODO: replace "v0.29.0" with the latest version of tflint
        TFLINT_VERSION: "0.29.0"
        TFLINT_OS: "amd64"

    phases:
      install:
        commands:
          # Install tflint
          - wget https://github.com/terraform-linters/tflint/releases/download/v${TFLINT_VERSION}/tflint_linux_${TFLINT_OS}.zip -O tflint.zip
          - unzip tflint.zip
          # Install tflint plugins
          - ./tflint --init
      pre_build:
        commands:
          - ./tflint -f junit > tflint_report.xml

    reports:
      tflint:
        files:
          - tflint_report.xml
    ```

### GitHub Actions

Assuming that you are storing your terraform configuration files and a `.tflint.hcl` file at the root of your repository, and that you are using `main` as your target branch for pull requests, you can create a GitHub actions workflow file such as this one:

!!! important

    Make sure that you include the __aws-serverless__ plugin into your `.tflint.hcl` configuration file, otherwise tflint will not install this ruleset. See [Installation](#installation) for more information.

=== "Sample workflow"

    ```yaml
    name: tflint-serverless

    on:
      pull_request:
        branches:
          # TODO: replace this if you are not using 'main' as your target
          # branch for pull requests.
          - main

    jobs:
      tflint-serverless:
        runs-on: ubuntu-latest
        steps:
          - uses: actions/checkout@v2
          - name: Setup TFLint
            uses: terraform-linters/setup-tflint@v1
            with:
              tflint_version: v0.29.0
          - name: Install Terraform plugins
            run: tflint --init
          - name: Lint Terraform files
            run: tflint
    ```

=== "With JUnit report"

    ```yaml
    name: tflint-serverless

    on:
      pull_request:
        branches:
          # TODO: replace this if you are not using 'main' as your target
          # branch for pull requests.
          - main

    jobs:
      tflint-serverless:
        runs-on: ubuntu-latest
        steps:
          - uses: actions/checkout@v2
          - name: Setup TFLint
            uses: terraform-linters/setup-tflint@v1
            with:
              tflint_version: v0.29.0
          - name: Install Terraform plugins
            run: tflint --init
          - name: Lint Terraform files
            run: tflint -f junit > tflint_report.xml
          - name: Publish test report
            uses: mikepenz/action-junit-report@v2
            # Only run this step on failure
            if: ${{ failure() }}
            with:
              report_paths: cfn_lint_report.xml
    ```

### GitLab

Assuming that you are storing your terraform configuration files and a `.tflint.hcl` file at the root of your repository, you can create a `.gitlab-ci.yml` file such as this one:

=== "Sample file"

    ```yaml
    tflint-serverless:
      variables:
        # TODO: replace "v0.29.0" with the latest version of tflint
        TFLINT_VERSION: "0.29.0"
        TFLINT_OS: "amd64"
      only:
        - merge_requests
      script:
        # Install tflint
        - wget https://github.com/terraform-linters/tflint/releases/download/v${TFLINT_VERSION}/tflint_linux_${TFLINT_OS}.zip -O tflint.zip
        - unzip tflint.zip
        # Install tflint plugins
        - ./tflint --init
        # Run tflint
        - ./tflint
    ```

=== "With JUnit report"

    ```yaml
    tflint-serverless:
      variables:
        # TODO: replace "v0.29.0" with the latest version of tflint
        TFLINT_VERSION: "0.29.0"
        TFLINT_OS: "amd64"
      only:
        - merge_requests
      script:
        # Install tflint
        - wget https://github.com/terraform-linters/tflint/releases/download/v${TFLINT_VERSION}/tflint_linux_${TFLINT_OS}.zip -O tflint.zip
        - unzip tflint.zip
        # Install tflint plugins
        - ./tflint --init
        # Run tflint
        - ./tflint -f junit > tflint_report.xml
      artifacts:
        when: always
        reports:
          junit: cfn_lint_report.xml
    ```