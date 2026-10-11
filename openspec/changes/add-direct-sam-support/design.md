# Design

## Context

See `proposal.md` for motivation. The Python plugin exports 17 rule classes through `cfn_lint_serverless.rules`; every rule uses cfn-lint's existing `CloudFormationLintRule` and `RuleMatch` interfaces.

Relevant observations:

- `lambda_.py` scans native functions, permissions, roles, event source mappings, and event-invoke configurations. `api_gateway.py` scans native REST/HTTP stages. `eventbridge.py` scans native rules.
- Most existing Lambda/API template fixtures actually declare SAM resources. Some permission tests rely on permissions generated from function events; wildcard tests rely on inline SAM policies; async tests include a SAM-managed failure queue declared only by `Type: SQS`.
- Before this upgrade, the locked host is cfn-lint 1.33.2 with an `aws-sam-translator` dependency; the package declares `cfn-lint>=1.33.0` and Python >=3.9.1. CI varies Python 3.9 through 3.13 but uses the same lockfile. These are observed starting points, not support promises for the upgraded project.
- In local cfn-lint 1.57.1, the transform merges SAM Globals but leaves SAM resource types unchanged. All 17 rules load, and native Lambda controls report correctly. The original reproduction exposed `E0002` from `WS1003` and no native functions for the direct Lambda checks. Remote workshop rebuilding and its exact host version are no longer acceptance targets.
- Go declares 1.23, HCL v2.23.0, and TFLint SDK v0.21.0. Tests use SDK `helper.TestRunner`; rules depend on `GetResourceContent`, `EvaluateExpr`, and SDK types, while the entry point uses `plugin.Serve`. CI still includes Go 1.21/1.22, below the module declaration.
- Root documentation requirements are unbounded `mkdocs` and `mkdocs-material`. CI/release configurations have separate Python, Go, TFLint, uv/action, and GoReleaser references; docs deployment currently follows a floating `master` action reference.
- The CDK example pins `aws-cdk.core==1.204.0` and split v1 service packages. Terraform examples pin AWS provider 5.40.0 and random 3.7.1; example runtime/tool declarations must be aligned with current frameworks and providers.
- Native `WS1003` iterates `Statement` as a list even when it is an object, so valid API Gateway/Pipes trust documents cause a crash.
- `tests/test_templates.py` checks target-rule presence and absence through the compatibility core API, not the actual CLI. Its filename convention and rule-inventory check must remain usable.

## Goals / Non-Goals

**Goals:**

- Reuse existing checks through small, source-aware projections of properties and events.
- Keep direct SAM results tied to real resources, not fictitious generated logical IDs.
- Operate on the template supplied by the host without changing it or invoking another transform.
- Preserve native-resource checks alongside direct SAM checks on the newly supported host range.
- Establish one recorded stable dependency/toolchain selection, migrate affected APIs, and keep manifests, local commands, CI, and support documentation coherent.

**Non-Goals:**

- Recreate SAM's full generated-resource graph, generated logical IDs, or deployment behavior.
- Expand managed-policy names or policy templates into IAM permissions.
- Evaluate arbitrary intrinsics, load external OpenAPI definitions, call AWS, or add a translator dependency.
- Add direct support for other SAM resource types such as StateMachine, SimpleTable, or WebSocketApi.
- Apply destination rules to SAM event types that cannot configure the corresponding destination.
- Retain obsolete language/host floors or historical CI environments solely for backward compatibility.
- Rebuild a workshop server, deploy examples, publish docs/packages, or change release credentials and destinations.

## Decisions

### 0. Upgrade the complete stable baseline before layering in SAM changes

The confirmed policy is latest stable, including major migrations and higher minimum versions. At implementation start, resolve targets from official package registries and release metadata, exclude prereleases, and record the resolution date, selected versions, required API migrations, and validated support floors. Do not assume the previously discussed 1.57.2 is still the latest host. If a latest stable component cannot preserve required behavior, stop for an explicit decision instead of silently pinning an older target.

| Surface | Upgrade and migration scope |
|---|---|
| Python plugin | Runtime dependency, dev group, build backend, declared Python/host floors, lockfile, test harness and affected APIs |
| Go plugin | HCL, TFLint SDK, compatible transitive graph, Go toolchain, SDK/rule/test migrations, TFLint host and GoReleaser config |
| Documentation | Root docs packages, recorded reproducible constraints, MkDocs/theme/config compatibility |
| CI and maintenance | Actions, uv, language/host/framework tools, coherent test/release versions, stable references instead of development branches, update coverage |
| Examples | Current CDK library/construct model and CLI, SAM/Serverless Framework requirements, Terraform/TFLint tools/providers, valid runtime/config declarations |

Set minimum versions from the selected dependencies' requirements and the runtime/tool combinations actually validated. Raising floors is authorized; rule IDs, severities, and expected native behavior remain stable. Add major dependency replacements and support changes to upgrade guidance. Update declarations as well as resolution artifacts: regenerating `uv.lock` alone or only updating `go.sum` is not sufficient.

Upgrade in independently testable sections: Python/build tooling, Go SDK/host, docs tooling, and examples, then finish the cross-project CI alignment and SAM acceptance checks. Preserve meaningful existing tests during API migration. Native baselines can pass before SAM checks are implemented; classify known SAM failures explicitly rather than weakening assertions to manufacture a green upgrade.

Refresh transitive dependencies through the resolver/module graph, within the selected direct dependencies' valid constraints. Pin or constrain local/CI tool resolutions sufficiently to reproduce the recorded baseline; a distinct latest-release drift job may resolve newer stable versions and must print what it used.

**Alternatives:** Keeping the original floors defeats the requested major upgrades and preserves unsupported environments. A lockfile-only refresh leaves incompatible manifests, SDK APIs, examples, and release jobs. Blindly replacing version strings without migrations and real host tests does not establish a usable baseline.

### 1. Use source-aware, read-only SAM helpers

Add a small helper module such as `cfn_lint_serverless/sam.py`, separate from the existing intrinsic-value utilities. It exposes iterators/property views for supported SAM functions, APIs, and function events. Each view carries the owning logical ID, a real `RuleMatch` source path, effective properties, and the relevant resource/event kind.

Reuse existing predicate logic where property shapes match. Map SAM `Tracing` to the Lambda tracing predicate and SAM API properties to the REST/HTTP stage predicates without inserting synthetic resources into `cfn.template`. Immutable source paths or newly constructed lists prevent accidental sharing between findings.

Native checks remain active. SAM helpers inspect only SAM declarations currently present. Already expanded native resources therefore use the native path rather than being projected a second time. This avoids version-number branching and monkey-patching cfn-lint; it does not require supporting historical hosts below the upgraded floor.

**Alternatives:** Simply adding SAM types to every `get_resources` call misses property differences and generated-resource checks. A whole-template translator would restore some prior behavior but adds a deployment-shaped dependency, generated IDs, and source/suppression problems. A custom runner would break ordinary append-rules usage.

### 2. Consume effective Globals, do not merge them twice

For explicit resources, consume properties after cfn-lint's normal template processing. The inspected current host already merges globals, including list concatenation and `IgnoreGlobals`; confirm that contract on the newly selected stable baseline during dependency migration. Do not import private cfn-lint merge helpers or prepend original Globals to already effective properties. Native inputs use their native properties.

For an implicit API with no explicit API resource, read the applicable `Globals.Api` or `Globals.HttpApi` settings into a separate read-only API view. There are no local API overrides to merge in that case. Explicit API references always use the referenced resource's effective settings.

Unit tests exercise the helper views after the normal host transform, including globals, overrides, ignored globals, and inherited policy lists. Assert that invoking plugin checks does not change the supplied template.

**Alternative:** A plugin-owned generic merge adds duplicate inherited list entries on current hosts and unnecessarily couples this change to all SAM global properties.

### 3. Restore Lambda behavior without reconstructing IAM or permissions

Extend the simple function checks and log-group association to SAM function views. Reuse supported log-group name/reference handling from `Value`; do not add general intrinsic evaluation.

For `WS1002`, combine identifiable explicit permission principals with a bounded event-to-service-principal mapping. Support `Api`, `HttpApi`, `CloudWatchEvent`, `EventBridgeRule`, `Schedule`, `IoTRule`, `S3`, and `SNS`. Polling events such as SQS, Kinesis, and DynamoDB do not imply Lambda resource-based permissions. Deduplicate by service principal, not event name. Preserve existing explicit-permission checks and reference matching.

For `WS1003`, inspect inline IAM documents in effective SAM `Policies` only when SAM manages the execution role. Skip named managed policies and named policy templates; an explicit `Role` means function `Policies` are not treated as deployed inline permissions. Native IAM role checks remain responsible for explicit roles.

Normalize trust-policy and permission-policy statements through a shared object-or-list helper. Check Lambda trust before evaluating a role's permission actions, and keep current action semantics: literal `*` and service wildcards fail, while resource-ARN wildcards alone do not.

For `ES1007`, inspect supported asynchronous event types already exercised by the fixtures. Check effective SAM `EventInvokeConfig.DestinationConfig.OnFailure`; accept an explicit destination or a valid SAM-managed SQS/SNS destination declaration. Keep native permission/event-invoke checks intact.

### 4. Project only supported event destination checks

For `ES1001`, inspect Kinesis and DynamoDB SAM event settings, which expose the relevant stream failure destination. Do not extend the check to SAM SQS declarations that cannot set that property; preserve existing explicit native mapping behavior.

For `ES4000`, inspect `CloudWatchEvent`, `EventBridgeRule`, and `Schedule` events and their target `DeadLetterConfig`. Accept either the explicit ARN shape or valid SAM-managed DLQ configuration. Do not reinterpret `ScheduleV2` as an EventBridge rule.

Report event-level findings on `Resources/<function>/Properties/Events/<event>`. This retains resource metadata filtering and makes the recommended edit visible.

### 5. Represent explicit and implicit API settings separately

Extend the four API rule classes using the following property views:

| Source | Logging property | Throttling property | Tracing |
|---|---|---|---|
| Native REST stage / SAM Api | AccessLogSetting | Wildcard MethodSettings | TracingEnabled |
| Native HTTP stage / SAM HttpApi | AccessLogSettings | DefaultRouteSettings | Not applicable |

Check explicit SAM API settings independently of `DefinitionBody`, including `AWS::Include`.

An `Api` event without `RestApiId` and an `HttpApi` event without `ApiId` produce implicit API views from applicable Globals. Group multiple routes by implicit API kind and owning function, and report each failed API rule at most once per owner. When multiple functions share an implicit API, each owning function can receive a finding; this is intentional so resource-level suppressions do not accidentally suppress another function's findings. An explicit or unresolved external API reference does not create an implicit API.

Retain the literal structured-log JSON check. Guard expression/container inputs: unresolved intrinsic formats are not parsed as strings, do not raise exceptions, and are not asserted to be valid JSON. Document this limitation instead of adding an evaluator.

### 6. Validate through both compatibility APIs and the CLI

Keep the existing filename-driven fixtures and rule-inventory assertion. Add focused direct-SAM unit/integration cases and a small self-contained SAM regression under a separate fixture directory so it does not violate `get_templates()` filename parsing.

The generic direct-SAM fixture retains globals, an untraced function, missing Lambda log group, an untraced REST API with independently configured logging/throttling, and single-object non-Lambda trust documents. Keep relevant ignore metadata and use inline/mock API definitions or a standalone unevaluated include shape. It must not reference the untracked workshop project or require a remote server.

Invoke the actual `cfn-lint` executable from the active test environment with JSON output; do not assume `python -m cfnlint` exists. Adapt the older core-API fixture harness if the selected current host has changed it, retaining observable rule assertions. Treat nonzero exit statuses caused by findings as expected. Assert:

- Rule discovery includes all 17 IDs.
- The direct-SAM regression includes `WS1000`, `WS1004`, and `WS2002`.
- Inherited memory/timeout and declared API logging/throttling do not produce their corresponding findings.
- Suppressed rule findings remain suppressed and `E0002` is absent.
- Native-only and mixed templates retain their applicable results.

Use explicit host/runtime environments for the newly declared cfn-lint floor and selected current release; test the minimum supported Python and selected current stable Python, with other supported versions as practical. Add a latest-stable host drift job on a compatible Python. Remove historical 1.33.0/1.33.2/workshop-specific jobs and obsolete Python combinations when they are below the new support boundary.

The current `uv run` test target may resynchronize the lockfile and silently undo a selected host version. Compatibility jobs must install the selected host in isolated environments and invoke pytest without a subsequent lockfile resync, or use an explicit tested override mechanism. Ordinary development tests can retain the existing Make target.

### 7. Verify real tooling locally, without invoking outward-facing workflows

Validate the Go SDK migration with `go test ./...`, rule-registration checks, `go vet`, and a build. Load that built plugin with the selected TFLint host in a temporary plugin directory, not a user's existing installation. Run against a locally initialized example and assert plugin rule discovery/execution; unit tests alone do not validate the host protocol.

Migrate CDK v1 imports/declarations to the current supported CDK library and constructs model, then synthesize locally and inspect the expected CloudFormation resources. Validate SAM locally, package the Serverless Framework example if its current licensing/auth requirements are available, and initialize Terraform with its backend disabled before validating and running TFLint. Update provider constraints and example runtime/config syntax together.

Build docs locally with the upgraded MkDocs/theme, validate the installed formatter/security/coverage tools, build the Python distribution in isolation, and run GoReleaser configuration checks or a local snapshot build appropriate to its selected version. Update CI actions/tool versions and ensure tests, docs, and release jobs use coherent declarations, including any Node runtime requirements introduced by action/framework upgrades.

Do not execute `uv publish`, GoReleaser publishing, docs deployment, Terraform apply, or SAM/Serverless deployment. No workflow may create credentials/accounts or submit a release as part of validation. If packaging needs unavailable auth/licensing or a tool cannot be installed, record the specific blocked check and do not represent integration as verified.

## Risks / Trade-offs

- **Partial SAM coverage can look like full validation:** Publish a resource/event coverage table and name policy-template, intrinsic, and external-definition limitations.
- **Duplicated or missing inherited settings:** Test host-processed inputs and avoid a second global merge.
- **New dependency APIs and minimums are breaking:** Migrate SDK/test/config consumers, publish the validated support boundaries, and preserve rule-level behavior rather than obsolete runtime support.
- **Native and SAM source shapes differ:** Retain native paths and validate native/SAM/mixed fixtures on the upgraded range rather than inspecting private pre-transform state.
- **Implicit APIs are shared:** Group per owner deliberately, test repeated routes, and explain why different owners may each receive a finding.
- **Destination semantics differ by event type:** Use explicit supported event sets and test both automatic and existing destinations; do not recommend properties invalid for that source.
- **Current-only tests hide regressions at the new floor:** Select host versions explicitly and verify the version actually used before running tests.
- **Transitive resolution or mutable references undo upgrades:** Record resolved versions, regenerate graph/lock data, pin stable references, and inspect selected versions in clean environments.
- **Framework licensing/auth or unavailable tools block validation:** Name those limitations, separate static checks from real execution, and never claim an unrun workflow passed.
- **Workshop data is user-owned and untracked:** Use a minimal standalone fixture; do not edit, move, commit, or depend on `unicorn_contracts/`.

## Migration Plan

First select and record the stable targets, update language/host requirements and dependency manifests, migrate the Python/Go integrations and examples, and regenerate lock/checksum data. Validate each upgraded area locally before relying on it in the new CI environments. Implement direct SAM behavior on that baseline, complete floor/current integration checks, and update installation, contributor, migration, and rule documentation.

User SAM templates do not need a new invocation flag. Users may need newer language/host tools because changed minimums are intentional breaking support changes. Release through the existing package process only when separately requested; the workshop server can be rebuilt separately and is not a gate.

Rollback restores the preceding plugin release and its recorded dependency/tool requirements as a coherent set. Do not mix old lockfiles with new manifests or silently change rule IDs/severities. No dependency downgrade that departs from the latest-stable policy is an unannounced completion path; report and resolve a concrete blocker explicitly.
