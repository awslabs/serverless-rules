# Lambda Star Permissions

__Level__: Warning
{: class="badge badge-yellow" }

__Initial version__: 0.1.3
{: class="badge badge-blue" }

__cfn-lint__: WS1003
{: class="badge" }

__tflint__: aws_iam_role_lambda_no_star
{: class="badge" }

 With Lambda functions, you should follow least-privileged access and only allow the access needed to perform a given operation. Attaching a role with more permissions than necessary can open up your systems for abuse.

??? warning "Limitations on policies"

    This rule only works with inline policies defined as part of the IAM role resource. It will not check managed policies or policies defined as separate resources.

    === "CloudFormation"

        ```yaml
        Resources:
          MyRole:
            Type: AWS::IAM::Role
            Properties:
              AssumeRolePolicyDocument:
                Version: "2012-10-17"
                Statement:
                  - Effect: Allow
                    Principal:
                      Service:
                        - lambda.amazonaws.com
                    Action:
                      - sts:AssumeRole
              # The rule will check this policy
              Policies:
                - PolicyName: root
                  PolicyDocument:
                    Version: "2012-10-17"
                    Statement:
                      - Effect: Allow
                        Action: dynamodb:Query
                        Resource: "*"

          # It will not check this policy
          MyPolicy:
            Type: AWS::IAM::Policy
            Properties:
              PolicyDocument:
                Version: "2012-10-17"
                    Statement:
                      - Effect: Allow
                        Action: dynamodb:*
                        Resource: "*"
              Roles:
                - !GetAtt MyRole.Arn
        ```

    === "Terraform"
    
        ```tf
        resource "aws_iam_role" "this" {
          name = "my-function-role"
          assume_role_policy = data.aws_iam_policy_document.assume.json

          # The rule will check this policy
          inline_policy {
            name = "FunctionPolicy"
            policy = data.aws_iam_policy_document.valid.json
          }
        }

        # It will not check this policy
        resource "aws_iam_policy" "this" {
          policy = data.aws_iam_policy_document.invalid.json
        }

        resource "aws_iam_policy_attachment" "this" {
          roles = [aws_iam_role.this.name]
          policy_arn = aws_iam_policy.this.arn
        }

        data "aws_iam_policy_document" "assume" {
          statement {
            actions = ["sts:AssumeRole"]
            principals {
              type       = "Service"
              identifiers = ["lambda.amazonaws.com"]
            }
          }
        }

        data "aws_iam_policy_document" "valid" {
          statement {
            actions = ["dynamodb:Query"]
            resources = ["arn:aws:dynamodb:eu-west-1:111122223333:table/my-table"]
          }
        }

        data "aws_iam_policy_document" "invalid" {
          statement {
            actions = ["dynamodb:*"]
            resources = ["*"]
          }
        }
        ```

## Why is this a warning?

If your Lambda function needs a broad range of permissions, you do not know ahead of time which permissions you will need, and you have evaluated the risks of using broad permissions for this function, you might ignore this rule.

## Implementations

=== "CDK"

    ```typescript
    import { AttributeType, Table } from '@aws-cdk/aws-dynamodb';
    import { Code, Function, Runtime } from '@aws-cdk/aws-lambda';

    export class MyStack extends cdk.Stack {
      constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const myTable = new Table(
          scope, 'MyTable',
          {
            partitionKey: {
              name: 'id',
              type: AttributeType.STRING,
            }
          },
        );

        const myFunction = new Function(
          scope, 'MyFunction',
          {
            code: Code.fromAsset('src/hello/'),
            handler: 'main.handler',
            runtime: Runtime.PYTHON_3_8,
          }
        );

        // Grant read access to the DynamoDB table
        table.grantReadData(myFunction);
      }
    }
    ```

=== "CloudFormation (JSON)"

    ```json
    {
      "Resources": {
        "MyFunction": {
          "Type": "AWS::Serverless::Function",
          "Properties": {
            "CodeUri": ".",
            "Runtime": "python3.8",
            "Handler": "main.handler",

            "Policies": [{
              "Version": "2012-10-17",
              "Statement": [{
                "Effect": "Allow",
                // Tightly scoped permissions to just 's3:GetObject'
                // instead of 's3:*' or '*'
                "Action": "s3:GetObject",
                "Resource": "arn:aws:s3:::my-bucket/*"
              }]
            }]
          }
        }
      }
    }
    ```

=== "CloudFormation (YAML)"

    ```yaml
    Resources:
      MyFunction:
        Type: AWS::Serverless::Function
        Properties:
          CodeUri: .
          Runtime: python3.12
          Handler: main.handler

          Policies:
            - Version: "2012-10-17"
              Statement:
                - Effect: Allow
                  # Tightly scoped permissions to just 's3:GetObject'
                  # instead of 's3:*' or '*'
                  Action: s3:GetObject
                  Resource: "arn:aws:s3:::my-bucket/*"
    ```

=== "Serverless Framework"

    ```yaml
    provider:
      name: aws
      iam:
        role:
          name: my-function-role
          statements:
            - Effect: Allow
              # Tightly scoped permissions to just 's3:GetObject'
              # instead of 's3:*' or '*'
              Action: s3:GetObject
              Resource: "arn:aws:s3:::my-bucket/*"

    functions:
      hello:
        handler: handler.hello
    ```

=== "Terraform"

    ```tf
    resource "aws_iam_role" "this" {
      name = "my-function-role"
      assume_role_policy = data.aws_iam_policy_document.assume.json

      inline_policy {
        name = "FunctionPolicy"
        policy = data.aws_iam_policy_document.this.json
      }
    }

    data "aws_iam_policy_document" "assume" {
      statement {
        actions = ["sts:AssumeRole"]
        principals {
          type       = "Service"
          identifiers = ["lambda.amazonaws.com"]
        }
      }
    }

    data "aws_iam_policy_document" "this" {
      statement {
        # Tightly scoped permissions to just 'dynamodb:Query'
        # instead of 'dynamodb:*' or '*'
        actions = ["dynamodb:Query"]
        resources = ["arn:aws:dynamodb:eu-west-1:111122223333:table/my-table"]
      }
    }

    resource "aws_lambda_function" "this" {
      function_name = "my-function"
      handler       = "main.handler"
      runtime       = "python3.8"
      filename      = "function.zip"
      role          = aws_iam_role.this.arn
    }
    ```

## See also
* [Serverless Lens: Identity and Access Management](https://docs.aws.amazon.com/wellarchitected/latest/serverless-applications-lens/identity-and-access-management.html)
* [AWS Lambda execution role](https://docs.aws.amazon.com/lambda/latest/dg/lambda-intro-execution-role.html)