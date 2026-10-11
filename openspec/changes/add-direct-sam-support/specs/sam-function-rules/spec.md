# Spec Delta

## Purpose

Apply existing serverless recommendations directly to SAM functions and their supported inline declarations, without depending on generated CloudFormation resources.

## ADDED Requirements

### Requirement: SAM function tracing
The plugin SHALL apply `WS1000` to `AWS::Serverless::Function` using its effective `Tracing` setting. Literal `Active` SHALL pass; an absent setting or literal `PassThrough` SHALL produce a warning associated with the original function.

#### Scenario: Tracing is missing
- **WHEN** a SAM function has no local or inherited `Tracing` setting
- **THEN** `WS1000` is reported for that function

#### Scenario: Tracing is active
- **WHEN** a SAM function has effective `Tracing: Active`
- **THEN** the function does not produce `WS1000`

#### Scenario: Resource overrides inherited tracing
- **WHEN** `Globals.Function.Tracing` is `Active` and a function declares `Tracing: PassThrough`
- **THEN** that function produces `WS1000`

### Requirement: SAM function memory and timeout
The plugin SHALL apply `ES1005` and `ES1006` to SAM functions when effective `MemorySize` and `Timeout`, respectively, are absent. Inherited values SHALL satisfy the corresponding checks.

#### Scenario: Both settings are missing
- **WHEN** a SAM function has neither local nor inherited memory or timeout settings
- **THEN** it produces `ES1005` and `ES1006`

#### Scenario: Settings come from Globals
- **WHEN** `Globals.Function` supplies `MemorySize: 128` and `Timeout: 15`
- **THEN** an inheriting function produces neither `ES1005` nor `ES1006`

### Requirement: SAM function log retention
The plugin SHALL apply `WS1004` to SAM functions using the existing log-group association checks: a matching log group must declare `RetentionInDays`. Matching SHALL support references to the SAM function logical ID and an explicitly declared function name.

#### Scenario: Log group is missing
- **WHEN** a SAM function has no associated log group declaring retention
- **THEN** `WS1004` is reported for the original function

#### Scenario: Log group references the function
- **WHEN** a log group has `LogGroupName: !Sub "/aws/lambda/${Function}"`, declares `RetentionInDays`, and `Function` is a SAM function
- **THEN** that function does not produce `WS1004`

#### Scenario: Existing supported name expressions
- **WHEN** a matching log group uses a supported `Fn::Join`, `Fn::Sub` variable map, or explicit function name and declares retention
- **THEN** the equivalent existing log-retention fixtures retain their passing results

### Requirement: SAM invocation principals
The plugin SHALL apply `WS1002` to the union of identifiable principals from explicit Lambda permissions and supported SAM function events. Events with the same service principal SHALL count as one principal; polling events SHALL not be treated as resource-based invocation permissions.

#### Scenario: Different SAM event services
- **WHEN** a function declares SNS and EventBridgeRule events
- **THEN** `WS1002` is reported for the original function

#### Scenario: SAM and explicit permissions are combined
- **WHEN** a function declares an SNS event and an explicit permission references that function with an EventBridge principal
- **THEN** `WS1002` is reported

#### Scenario: Repeated service and polling event
- **WHEN** a function has two SNS invocation sources and an SQS polling event
- **THEN** those sources alone do not produce `WS1002`

### Requirement: SAM inline wildcard permissions
The plugin SHALL apply `WS1003` to literal actions in inline IAM policy documents on SAM functions when the function uses a SAM-managed role. A wildcard resource ARN alone SHALL not constitute a wildcard action. Named managed policies and SAM policy-template names SHALL not be expanded into inferred permissions.

#### Scenario: Inline wildcard action
- **WHEN** a SAM function without an explicit `Role` has an inline policy with `Action: "s3:*"`
- **THEN** `WS1003` is reported for that function

#### Scenario: Scoped action with wildcard resource
- **WHEN** an inline policy uses `Action: s3:GetObject` and a resource ARN ending in `/*`
- **THEN** that policy does not produce `WS1003`

#### Scenario: Policy template is not an IAM document
- **WHEN** a function declares `DynamoDBWritePolicy` and a scoped inline policy
- **THEN** the plugin checks the inline document without crashing or inventing wildcard actions for the policy template

#### Scenario: Explicit role supersedes SAM-managed policies
- **WHEN** a SAM function declares an explicit `Role`
- **THEN** `WS1003` evaluates any applicable explicit IAM role rather than treating ignored function `Policies` as deployed permissions

### Requirement: SAM asynchronous failure destinations
The plugin SHALL apply `ES1007` to functions with supported asynchronous SAM event sources using their effective `EventInvokeConfig.DestinationConfig.OnFailure`. Supported sources SHALL include CloudWatchEvent, EventBridgeRule, IoTRule, S3, Schedule, and SNS. A valid SAM declaration requesting an automatically created failure destination SHALL satisfy the check.

#### Scenario: Async event lacks failure destination
- **WHEN** a function declares an SNS event without a local or inherited failure destination
- **THEN** `ES1007` is reported at the original event declaration

#### Scenario: Inherited auto-created destination
- **WHEN** `Globals.Function.EventInvokeConfig.DestinationConfig.OnFailure` declares `Type: SQS`
- **THEN** inheriting functions with supported asynchronous events do not produce `ES1007`

#### Scenario: Explicit destination
- **WHEN** a function has a supported asynchronous event and an on-failure destination referencing an existing queue
- **THEN** that event does not produce `ES1007`

### Requirement: SAM stream event failure destinations
The plugin SHALL apply `ES1001` to SAM Kinesis and DynamoDB event declarations using their on-failure destination configuration. It SHALL NOT recommend unsupported destination properties on SQS polling declarations. Existing checks for explicit native event source mappings SHALL remain unchanged.

#### Scenario: Stream event lacks destination
- **WHEN** a SAM function declares a Kinesis or DynamoDB event without an on-failure destination
- **THEN** `ES1001` is reported at that event declaration

#### Scenario: Stream event declares destination
- **WHEN** a SAM stream event declares a valid on-failure destination
- **THEN** the event does not produce `ES1001`

#### Scenario: SQS polling is not a stream destination check
- **WHEN** a SAM function declares only an SQS polling event
- **THEN** that declaration does not produce `ES1001`

### Requirement: SAM EventBridge target dead letter queues
The plugin SHALL apply `ES4000` to CloudWatchEvent, EventBridgeRule, and Schedule event declarations whose generated EventBridge target lacks a dead letter queue. Both an explicit DLQ ARN and a valid SAM declaration requesting an automatically created DLQ SHALL satisfy this check.

#### Scenario: EventBridge event lacks DLQ
- **WHEN** a function declares an EventBridgeRule event without `DeadLetterConfig`
- **THEN** `ES4000` is reported at the event declaration

#### Scenario: EventBridge event requests a DLQ
- **WHEN** the event declares a valid explicit DLQ or a valid SAM-managed DLQ configuration
- **THEN** the event does not produce `ES4000`
