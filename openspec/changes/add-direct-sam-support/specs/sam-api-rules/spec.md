# Spec Delta

## Purpose

Apply existing API Gateway stage recommendations to explicit and implicit SAM REST and HTTP APIs, preserving the distinction between their supported settings.

## ADDED Requirements

### Requirement: SAM API access logging
The plugin SHALL apply `ES2000` to explicit SAM REST and HTTP APIs using effective `AccessLogSetting` and `AccessLogSettings`, respectively. Missing settings or a missing `DestinationArn` SHALL produce a finding on the original API resource.

#### Scenario: REST access logging is absent
- **WHEN** an `AWS::Serverless::Api` lacks effective access log settings
- **THEN** `ES2000` is reported on that API

#### Scenario: HTTP logging has no destination
- **WHEN** an `AWS::Serverless::HttpApi` has `AccessLogSettings: {}`
- **THEN** `ES2000` is reported on that API

#### Scenario: Logging is inherited
- **WHEN** the applicable SAM Globals settings provide an access log destination to an explicit API
- **THEN** that API does not produce `ES2000`

### Requirement: SAM API structured logging
The plugin SHALL apply `WS2001` to literal access log formats on SAM REST and HTTP APIs using the existing JSON-format check, including API Gateway context substitutions. An absent format SHALL remain the responsibility of other applicable checks. An unresolved intrinsic format SHALL not cause an exception or be claimed to be valid JSON.

#### Scenario: Context substitutions in JSON
- **WHEN** a literal access log format is valid JSON after supported `$context` substitutions
- **THEN** it does not produce `WS2001`

#### Scenario: Unstructured log format
- **WHEN** a literal access log format is not JSON
- **THEN** `WS2001` is reported on the API

#### Scenario: Format is an intrinsic expression
- **WHEN** an access log format is an unresolved `Fn::Sub`, `Fn::If`, or other intrinsic expression
- **THEN** the structured-logging check does not raise a rule-processing exception
- **AND** it does not treat the expression as proof of a passing JSON format

### Requirement: SAM API default throttling
The plugin SHALL apply `ES2003` to SAM REST APIs using wildcard `MethodSettings`, and to SAM HTTP APIs using `DefaultRouteSettings`. Both `ThrottlingBurstLimit` and `ThrottlingRateLimit` SHALL be required for a passing literal configuration, consistent with existing stage checks.

#### Scenario: REST default throttling is incomplete
- **WHEN** the wildcard method setting on a SAM REST API lacks either limit
- **THEN** `ES2003` is reported on that API

#### Scenario: HTTP default throttling is configured
- **WHEN** a SAM HTTP API declares both default route throttling limits
- **THEN** it does not produce `ES2003`

#### Scenario: A specific REST route is not a default
- **WHEN** throttling is declared only for a non-wildcard REST method or resource
- **THEN** `ES2003` is reported

### Requirement: SAM REST API tracing
The plugin SHALL apply `WS2002` to SAM REST APIs using effective `TracingEnabled`. Missing or literal false settings SHALL produce a warning. SAM HTTP APIs SHALL not produce this warning because HTTP APIs do not support X-Ray tracing.

#### Scenario: REST tracing is absent
- **WHEN** a SAM REST API has no effective `TracingEnabled`
- **THEN** `WS2002` is reported on that API

#### Scenario: REST tracing is inherited
- **WHEN** `Globals.Api.TracingEnabled` is true and the API does not override it
- **THEN** that API does not produce `WS2002`

#### Scenario: HTTP API has no tracing setting
- **WHEN** an API is `AWS::Serverless::HttpApi`
- **THEN** it does not produce `WS2002`

### Requirement: Implicit SAM API checks
The plugin SHALL check implicit REST and HTTP APIs created by SAM function `Api` and `HttpApi` events without explicit API references. Applicable Globals SHALL supply their settings. Findings SHALL identify an originating function event; multiple events for the same implicit API on one function SHALL not duplicate a rule's finding.

#### Scenario: Implicit REST API lacks settings
- **WHEN** a function declares an `Api` event without `RestApiId` and the applicable global settings are absent
- **THEN** the implicit API produces `ES2000`, `ES2003`, and `WS2002` findings associated with that function's API event

#### Scenario: Implicit API inherits configured settings
- **WHEN** applicable Globals supply logging, structured log format, throttling, and REST tracing for an implicit API
- **THEN** that implicit API produces none of those rule findings

#### Scenario: Event references an explicit API
- **WHEN** a function event references an explicitly declared SAM API
- **THEN** the plugin checks the API resource rather than creating a duplicate implicit API check

#### Scenario: Multiple routes on one function
- **WHEN** one function has multiple events targeting the same implicit API
- **THEN** each failed API rule produces at most one finding for that function

### Requirement: API checks do not depend on definition expansion
The plugin SHALL check declarative SAM API stage settings without downloading or expanding external API definitions. A nested `AWS::Include` in `DefinitionBody` SHALL not prevent evaluation of separately declared logging, throttling, or tracing settings.

#### Scenario: Workshop REST API includes an external definition
- **WHEN** a SAM API supplies access logging and wildcard throttling, uses `AWS::Include` for its definition, and omits tracing
- **THEN** the API produces `WS2002`
- **AND** its declared logging and throttling satisfy `ES2000` and `ES2003`
- **AND** the plugin does not fetch the included definition
