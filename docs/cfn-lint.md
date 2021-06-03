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

## Other frameworks

You can use the `cfn-lint` module with all frameworks that synthesizes to a CloudFormation template, such as the [AWS Cloud Development Kit (CDK)](https://docs.aws.amazon.com/cdk/latest/guide/home.html) and the [Serverless Framework](https://www.serverless.com/).

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

The __Serverless Framework__ is a specialized framework for Serverless applications which uses CloudFormation under the hood when deploying resources to AWS. You can manually create a package using the `sls package` command, which will generate the CloudFormation templates. With the `--package` option, you can specify in which folder it will store your package. In that folder, you can find the CloudFormation templates as JSON files starting with `cloudformation-template-`. You can then run `cfn-lint` for inspecting those template.

```bash
sls package --package output/
cfn-lint output/cloudformation-template-*.json -a cfn_lint_serverless.rules
```

## IDE integration

### Visual Studio Code

![Screenshot of VS Code using CloudFormation Linter and Serverless Rules](images/cfn_lint_vscode.png)

If you use Visual Studio Code, you can add the [CloudFormation Linter](https://marketplace.visualstudio.com/items?itemName=kddejong.vscode-cfn-lint) extension, which will automatically run `cfn-lint` on your CloudFormation templates. In the extension's `settings.json` file, you can add additional rules like so:

```json
{
  // ... other settings omitted

  "cfnLint.appendRules": [
    "cfn_lint_serverless.rules"
  ]
}
```