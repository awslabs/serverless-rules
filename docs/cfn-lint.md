`cfn-lint` module
=================

## Installation

Serverless Rules for `cfn-lint` is available through the `cfn-lint-serverless` Python module in PyPi. You can use pip to install the module:

```bash
pip install cfn-lint cfn-lint-serverless
```

## Usage

You can now instruct `cfn-lint` to use Serverless Rules installed previous via `--append-rules`, or `-a` for short:

```bash
cfn-lint my_template.yaml -a cfn_lint_serverless.rules
```

You can also try with a [Serverless Application Model (SAM)](https://github.com/awslabs/serverless-rules/tree/main/examples/sam) example provided in this repository. In the root folder of the repository, you can run:

```bash
cfn-lint examples/sam/template.yaml -a cfn_lint_serverless.rules
```

### Sample outputs

If the template fulfills the requirements for all rules, `cfn-lint` will return an empty output. Otherwise, `cfn-lint` will output recommendations.

=== "Matching template"

    ```bash
    $ cfn-lint template.yaml -a cfn_lint_serverless.rules
    $
    ```

=== "With recommendations"

    ```bash
    $ cfn-lint template.yaml -a cfn_lint_serverless.rules
    ES4000 EventBridge rule ConsumerFunctionTrigger should have a DeadLetterConfig.Arn property for all its Targets.
    template.yaml:5:1

    WS1000 Lambda function ConsumerFunction should have TracingConfig.Mode set to 'Active'.
    template.yaml:7:3

    WS1004 Lambda function ConsumerFunction does not have a corresponding log group with a Retention property
    template.yaml:7:3
    $
    ```

### Ignoring rules

Serverless Rules is a set of recommended practices. 

We recommend you to keep Error-level rules enabled. Non-error rules, for example [Lambda Tracing](rules/lambda/tracing.md), contain detailed scenarios on when it’s safe to ignore them.

When needed, you can ignore any specific rule that doesn’t match your environment.

Rules in `cfn-lint` can be disabled either through the `--ignore-checks`/`-i` command line argument, or with the `.cfnlintrc` configuration file in the current working directory. See the [`cfn-lint` documentation](https://github.com/aws-cloudformation/cfn-lint#configuration) for more information.

=== "Command line"

    ```bash
    # Disable rule WS1000
    cfn-lint my_template.yaml -a cfn_lint_serverless.rules -i WS1000
    ```

=== ".cfnlintrc"

    ```yaml
    ignore_checks:
      # Disable rule WS1000
      - WS1000
    ```

## Other frameworks

You can use the `cfn-lint` module with all frameworks that synthesize to a CloudFormation template, such as the [AWS Cloud Development Kit (CDK)](https://docs.aws.amazon.com/cdk/latest/guide/home.html) and the [Serverless Framework](https://www.serverless.com/).

### CDK

With the __AWS CDK__, you define your cloud resources using a familiar programming language such as TypeScript, Python, Java, C#/.Net, or Go. It will then use CloudFormation to provision those resources in your AWS environments.

Under the hood, CDK will generate a CloudFormation template before initiating a deployment on AWS. You can use the `cdk synth` command to generate that template manually. You can then run `cfn-lint` for inspecting that template.

```bash
cdk synth > template.yaml
cfn-lint template.yaml -a cfn_lint_serverless.rules
```

You can also try with a [CDK example](https://github.com/awslabs/serverless-rules/tree/main/examples/cdk) provided in this repository. In the root folder of the repository, you can run:

```bash
cd examples/cdk/
cdk synth > template.yaml
cfn-lint template.yaml -a cfn_lint_serverless.rules
```

### Serverless Framework

The __Serverless Framework__ is a specialized framework for Serverless applications that uses CloudFormation under the hood when deploying resources to AWS. You can manually create a package using the `sls package` command, which will generate the CloudFormation templates. With the `--package` option, you can specify in which folder it will store your package. In that folder, you can find the CloudFormation templates as JSON files starting with `cloudformation-template-`. You can then run `cfn-lint` for inspecting those templates.

```bash
sls package --package output/
cfn-lint output/cloudformation-template-*.json -a cfn_lint_serverless.rules
```

## Continuous integration

You can use Serverless Rules and `cfn-lint` with your continuous integration tool to automatically check CloudFormation templates with rules from this project. For example, you can validate on pull requests, merge to your main branch, or before deploying to production.

If there are any issues with your template, `cfn-lint` will return a non-zero error code. You can find more information about `cfn-lint` return codes in [its documentation](https://github.com/aws-cloudformation/cfn-lint).

### AWS CodeBuild

Assuming that you are storing your template as `template.yaml` at the root of your repository, you can create a [buildspec file](https://docs.aws.amazon.com/codebuild/latest/userguide/build-spec-ref.html) such as this one:

=== "Sample buildspec"

    ```yaml
    version: 0.2

    phases:
      install:
        runtime-versions:
          python: "3.8"
        commands:
          # Install cfn-lint-serverless
          - pip install cfn-lint cfn-lint-serverless
      pre_build:
        commands:
          # TODO: replace here with your template name if you are not
          # using 'template.yaml'.
          - cfn-lint template.yaml -a cfn_lint_serverless.rules
    ```

If you want to run `cfn-lint` with other frameworks, see how you can generate CloudFormation templates in the [Other frameworks](#other-frameworks) section of this documentation.

### GitHub Actions

Assuming that your template is stored as `template.yaml` at the root of your repository and that you are using `main` as your target branch for pull requests, you can create a GitHub actions workflow file such as this one:

=== "GitHub Actions workflow"

    ```yaml
    name: cfn-lint-serverless

    on:
      pull_request:
        branches:
          # TODO: replace this if you are not using 'main' as your target
          # branch for pull requests.
          - main

    jobs:
      cfn-lint-serverless:
        runs-on: ubuntu-latest
        steps:
          - uses: actions/checkout@v2
          - name: Set up Python 3.8
            uses: actions/setup-python@v2
            with:
              python-version: "3.8"
          - name: Install cfn-lint-serverless
            # Install cfn-lint-serverless
            run: pip install cfn-lint cfn-lint-serverless
          - name: Lint CloudFormation template
            # TODO: replace here with your template name if you are not
            # using 'template.yaml'.
            run: cfn-lint template.yaml -a cfn_lint_serverless.rules
    ```

If you want to run `cfn-lint` with other frameworks, see how you can generate CloudFormation templates in the [Other frameworks](#other-frameworks) section of this documentation.

### GitLab

Assuming that your template is stored as `template.yaml` at the root of your repository, you can create a `.gitlab-ci.yml` file such as this one:

=== ".gitlab-ci.yml"

    ```yaml
    cfn-lint-serverless:
      image: python:latest
      only:
        - merge_requests
      script:
        # Install cfn-lint-serverless
        - pip install cfn-lint cfn-lint-serverless
        # TODO: replace here with your template name if you are not
        # using 'template.yaml'.
        - cfn-lint template.yaml -a cfn_lint_serverless.rules
    ```

If you want to run `cfn-lint` with other frameworks, see how you can generate CloudFormation templates in the [Other frameworks](#other-frameworks) section of this documentation.

## IDE integration

### Visual Studio Code

![Screenshot of VS Code using CloudFormation Linter and Serverless Rules](images/cfn_lint_vscode.png)

For Visual Studio Code, you can add the [CloudFormation Linter](https://marketplace.visualstudio.com/items?itemName=kddejong.vscode-cfn-lint) extension, which will automatically run `cfn-lint` on your CloudFormation templates. In the extension's `settings.json` file, you can add additional rules like so:

=== "Extension settings"

    ```json
    {
      // ... other settings omitted

      "cfnLint.appendRules": [
        "cfn_lint_serverless.rules"
      ]
    }
    ```