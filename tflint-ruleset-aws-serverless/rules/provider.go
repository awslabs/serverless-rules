package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var Rules = []tflint.Rule{
	// NewAwsAPIGatewayMethodSettingsThrottlingRule(),
	// NewAwsAPIGatewayStageLoggingRule(),
	// NewAwsAPIGatewayStageTracingRule(),
	// NewAwsAPIGatewayStageV2LoggingRule(),
	// NewAwsApigatewayStageStructuredLoggingRule(),
	// NewAwsApigatewayV2StageStructuredLoggingRule(),
	// NewAwsApigatewayV2StageThrottlingRule(),
	// NewAwsAppsyncGraphqlAPITracingRule(),
	// NewAwsCloudwatchEventTargetNoDlqRule(),
	// NewAwsCloudwatchLogGroupLambdaRetentionRule(),
	// NewAwsIamRoleLambdaNoStarRule(),
	// NewAwsLambdaEventInvokeConfigAsyncOnFailureRule(),
	// NewAwsLambdaEventSourceMappingFailureDestinationRule(),
	// NewAwsLambdaFunctionDefaultMemoryRule(),
	// NewAwsLambdaFunctionDefaultTimeoutRule(),
	// NewAwsLambdaFunctionEolRuntimeRule(),
	// NewAwsLambdaFunctionTracingRule(),
	// NewAwsLambdaPermissionMultiplePrincipalsRule(),
	// NewAwsSfnStateMachineTracingRule(),
	// NewAwsSnsTopicSubscriptionRedrivePolicyRule(),
	NewAwsSqsQueueRedrivePolicyRule(),
}
