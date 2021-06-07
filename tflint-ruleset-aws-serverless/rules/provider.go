package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var Rules = []tflint.Rule{
	NewAwsAPIGatewayMethodSettingsThrottlingRule(),
	NewAwsAPIGatewayStageLoggingRule(),
	NewAwsAPIGatewayStageTracingRule(),
	NewAwsAPIGatewayStageV2LoggingRule(),
	NewAwsApigatewayV2StageThrottlingRule(),
	NewAwsAppsyncGraphqlAPITracingRule(),
	NewAwsCloudwatchEventTargetNoDlqRule(),
	NewAwsLambdaEventSourceMappingFailureDestinationRule(),
	NewAwsLambdaFunctionTracingRule(),
	NewAwsSfnStateMachineTracingRule(),
	NewAwsSnsTopicSubscriptionRedrivePolicyRule(),
	NewAwsSqsQueueRedrivePolicyRule(),
}
