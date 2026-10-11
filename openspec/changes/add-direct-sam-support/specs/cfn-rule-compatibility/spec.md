# Spec Delta

## Purpose

Keep Serverless Rules discoverable and reliable across supported cfn-lint versions, while preserving native CloudFormation behavior, source attribution, and suppression controls.

## ADDED Requirements

### Requirement: Stable rule discovery and IDs
The existing append-rules module SHALL continue to expose all 17 current serverless rule IDs. Direct SAM support SHALL use existing IDs and severities without requiring a new CLI flag, an upstream cfn-lint modification, or user-side SAM translation.

#### Scenario: Rules are listed
- **WHEN** the user runs `cfn-lint --list-rules --append-rules cfn_lint_serverless.rules` in an environment where the plugin is installed
- **THEN** the output includes all 17 existing serverless rule IDs

#### Scenario: Existing invocation checks SAM
- **WHEN** the user runs `cfn-lint template.yaml -a cfn_lint_serverless.rules` on a supported SAM template
- **THEN** applicable serverless findings are reported using the existing rule IDs

### Requirement: Effective SAM Globals
Direct SAM checks SHALL respect the effective SAM settings supplied by cfn-lint, including primitive overrides, recursive map merging, list inheritance, and supported `IgnoreGlobals` declarations. The plugin SHALL not apply inheritance a second time to already merged resource settings or change the host template.

#### Scenario: Overrides and map inheritance
- **WHEN** global settings and resource settings supply overlapping and distinct nested properties
- **THEN** direct SAM checks observe the effective inherited properties and local overrides

#### Scenario: Ignored global setting
- **WHEN** a supported `IgnoreGlobals` declaration excludes a global property
- **THEN** the excluded value does not satisfy the corresponding rule

#### Scenario: Already merged policy list
- **WHEN** cfn-lint has merged global and local policy lists before invoking the plugin
- **THEN** the plugin evaluates that list without prepending global policies again
- **AND** the host template remains unchanged by the plugin

### Requirement: Native and SAM checks coexist
Existing native CloudFormation checks SHALL remain active in native-only and mixed templates. If a supported host supplies already expanded native resources, the plugin SHALL check those resources without adding duplicate findings for SAM declarations no longer present.

#### Scenario: Mixed function types
- **WHEN** a template contains both a native Lambda function and a SAM function lacking tracing
- **THEN** each function produces its applicable `WS1000` finding

#### Scenario: Already expanded native resources
- **WHEN** a supported host supplies expanded native resources
- **THEN** existing native rule behavior is preserved
- **AND** the plugin does not check an additional synthetic copy of the same expanded resource

### Requirement: Original source attribution and suppression
Direct SAM findings SHALL identify original resource or event paths. Existing command-line ignores and template/resource metadata suppressions SHALL apply to their existing rule IDs. Findings for an implicit API SHALL be associated with its owning function so resource-level suppression remains meaningful.

#### Scenario: Template-wide suppression
- **WHEN** a SAM template ignores `WS1000` in its cfn-lint metadata
- **THEN** direct SAM tracing findings are not reported

#### Scenario: Resource-level event suppression
- **WHEN** a function's metadata suppresses `ES1007` and its asynchronous event lacks a failure destination
- **THEN** the function's direct event finding is suppressed
- **AND** an unsuppressed function's equivalent violation remains reportable

#### Scenario: Source location of a direct finding
- **WHEN** a direct SAM check finds a violation on an explicit resource or inline event
- **THEN** the diagnostic points into that original declaration rather than a fabricated generated logical ID

### Requirement: IAM statement shape robustness
The wildcard-permission check SHALL accept IAM `Statement` as either one object or a list in trust policies and inline permission policies. Non-Lambda trust policies SHALL not cause exceptions or wildcard-permission findings merely because they use the single-object form.

#### Scenario: Non-Lambda single-statement trust policy
- **WHEN** an API Gateway or EventBridge Pipes role uses a single-object trust statement
- **THEN** `WS1003` processing does not produce an `E0002` exception

#### Scenario: Lambda role with one inline statement
- **WHEN** a Lambda-trusted IAM role has a single-object inline policy statement with `Action: "*"`
- **THEN** it produces `WS1003`

### Requirement: Supported-version regression coverage
The plugin SHALL declare a validated cfn-lint support floor for the upgraded stable baseline. Validation SHALL exercise that floor and the selected current/latest stable host on supported Python versions, including rule discovery, native templates, and CLI-level SAM diagnostics. Historical hosts below the new floor SHALL not be required.

#### Scenario: Current stable host reports direct SAM violations
- **WHEN** the selected current stable host checks a self-contained SAM fixture with inherited memory/timeout, missing function tracing/retention, an untraced REST API, and relevant suppressions
- **THEN** it reports at least `WS1000`, `WS1004`, and `WS2002`
- **AND** the inherited memory and timeout do not produce `ES1005` or `ES1006`
- **AND** processing does not produce `E0002`

#### Scenario: New support floor retains fixture outcomes
- **WHEN** the existing template fixtures run against the newly declared support-floor job
- **THEN** their expected target-rule pass/fail outcomes are preserved

#### Scenario: Install metadata matches supported hosts
- **WHEN** users inspect the updated Python package metadata and installation guidance
- **THEN** both declare the validated minimum cfn-lint and Python versions
- **AND** they do not promise support for the previous 1.33.0 baseline unless it remains within the newly validated range

### Requirement: Explicit SAM coverage limits
Documentation SHALL identify supported SAM resources and event declarations. Unexpanded managed-policy names, SAM policy-template permissions, arbitrary intrinsic expressions, external definitions, and SAM resource types outside this change SHALL not be represented as fully analyzed.

#### Scenario: Unexpanded policy template
- **WHEN** users consult the direct SAM support documentation
- **THEN** it explains that named policy templates and managed policies are not expanded into effective permissions
- **AND** it identifies the supported direct resource and event checks
