from aws_cdk import core as cdk
from aws_cdk import aws_apigateway as apigw
from aws_cdk import aws_lambda as lambda_
from aws_cdk import aws_logs as logs


class CdkStack(cdk.Stack):

    def __init__(self, scope: cdk.Construct, construct_id: str, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)

        # Lambda function

        hello_function = lambda_.Function(
            self, "hello-function",
            code=lambda_.Code.from_asset("src/hello/"),
            handler="main.handler",
            runtime=lambda_.Runtime.PYTHON_3_8,
            tracing=lambda_.Tracing.ACTIVE
        )

        logs.LogGroup(
            self, "hello-logs",
            log_group_name=f"/aws/lambda/{hello_function.function_name}",
            retention=logs.RetentionDays.ONE_WEEK
        )

        # API Gateway

        api_logs = logs.LogGroup(
            self, "hello-api-logs",
            retention=logs.RetentionDays.ONE_WEEK
        )

        api = apigw.RestApi(
            self, "hello-api",
            deploy_options=apigw.StageOptions(
                access_log_destination=apigw.LogGroupLogDestination(api_logs),
                access_log_format=apigw.AccessLogFormat.json_with_standard_fields(
                    caller=True,
                    http_method=True,
                    ip=True,
                    protocol=True,
                    request_time=True,
                    resource_path=True,
                    response_length=True,
                    status=True,
                    user=True
                ),
                throttling_burst_limit=1000,
                throttling_rate_limit=10,
                tracing_enabled=True
            )
        )

        hello_integration = apigw.LambdaIntegration(hello_function, proxy=True)
        api.root.add_method(
            "GET",
            hello_integration
        )