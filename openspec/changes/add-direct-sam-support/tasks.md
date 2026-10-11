# Tasks

## 1. Stable target selection and Python dependency upgrade

- [x] 1.1 Resolve current stable targets from official metadata for every dependency/tool surface and record the selection date, versions, major migrations, and proposed validated floors in upgrade guidance; verify the inventory covers Python runtime/dev/build, Go/transitives, docs, CI/tools, examples, and providers without defaulting to workshop versions or prereleases.
- [x] 1.2 Upgrade Python runtime/dev/build declarations and the supported Python/cfn-lint minimums, regenerate `uv.lock`, and align local setup commands; verify clean locked dependency installation, selected-version output, and an isolated Python distribution build with matching metadata.
- [x] 1.3 Migrate affected cfn-lint test/runner APIs without weakening rule assertions; verify rule discovery and native-only regression controls on the selected baseline, explicitly retaining known direct-SAM failures until their implementation tasks rather than claiming they pass.
- [x] 1.4 Update Python installation and contributor guidance with new minimums, reproducible setup, and major support changes; verify documented commands in a clean supported environment and consistency with package metadata.

## 2. Go dependencies, SDK migration, and TFLint host

- [x] 2.1 Upgrade the Go toolchain declaration, HCL, TFLint SDK, and compatible transitive graph in `go.mod`/`go.sum`; verify clean module resolution, `go mod tidy` consistency, and the recorded selected versions.
- [x] 2.2 Migrate SDK/HCL consumers and tests as needed, including resource queries, expression evaluation, helper runners, and plugin serving; verify `go test ./...`, `go vet ./...`, plugin build, and existing rule-registration checks without changing rule IDs.
- [x] 2.3 Align the TFLint host and GoReleaser configuration with the selected stable toolchain; verify the compiled plugin loads and runs locally using a temporary plugin directory and that local release-configuration checks pass without publishing.
- [x] 2.4 Update TFLint installation/build guidance and contributor requirements; verify host/tool minimums and local plugin-loading commands agree with the tested dependency selection.

## 3. Documentation and development-tool compatibility

- [x] 3.1 Upgrade root documentation packages with reproducible selected-version constraints and adapt MkDocs/theme configuration if required; verify a local documentation build with the upgraded packages and no deployment.
- [x] 3.2 Validate upgraded formatter, coverage, security, and complexity tooling against their project configurations; verify their scoped local commands execute with recorded versions, report actual findings, and do not silently revert dependencies or scan user-owned workshop content.
- [x] 3.3 Align root development/Make commands and docs/contributor instructions with the upgraded reproducible environments; verify a clean documented setup can invoke Python, Go, and docs tooling without mixing unsupported system interpreters or old dependency baselines.

## 4. Current frameworks, providers, and examples

- [x] 4.1 Migrate the CDK example from split CDK v1/core packages to the current supported library/construct model and selected CLI; verify fresh installation, local synthesis, expected CloudFormation resource shapes, and matching CDK example instructions.
- [x] 4.2 Upgrade SAM and Serverless Framework tool/runtime requirements and adapt example configuration to current syntax; verify applicable local SAM validation and Serverless packaging, recording any auth/licensing blocker explicitly, and update their example instructions without deployment.
- [x] 4.3 Upgrade Terraform/TFLint example tools and AWS/random provider constraints, and update unsupported runtime/config declarations; verify backend-disabled local initialization, Terraform validation, and linting with the newly built plugin, with documented expected rule findings.
- [x] 4.4 Align example Make targets and READMEs with recorded current tool/framework selections; verify the documented synthesize/build/package/validate workflows use the new baseline and never require a workshop rebuild or infrastructure deployment.

## 5. CI actions and support-range alignment

- [x] 5.1 Upgrade test/docs/release/dependency-maintenance action and tool references to selected stable releases, replacing floating development branches; verify workflow/config syntax and local tool requirements while preserving existing publication destinations and credentials without running publish jobs.
- [x] 5.2 Align Python/Go/Node and host/tool selections across tests, docs, and release configurations; verify minimum/current supported combinations, remove obsolete below-floor jobs, and ensure GoReleaser configuration or local snapshot validation does not publish assets.
- [x] 5.3 Add isolated cfn-lint floor/current/latest-stable validation and dependency-update coverage for the upgraded manifests; verify each environment prints and uses its selected version without a later lockfile resync undoing it, and update CI/support guidance to match the new range.

## 6. Source-aware SAM helpers and IAM statement robustness

- [x] 6.1 Add read-only SAM function, API, and event views carrying original owner/path information; verify focused helper tests cover effective property inputs, missing/unsupported shapes, and unchanged host-template contents on the upgraded host.
- [x] 6.2 Normalize single-object/list IAM statements and fix native `WS1003` trust-policy processing; verify tests cover Lambda and non-Lambda roles, both statement shapes, literal wildcard actions, and API Gateway/Pipes trust shapes without `E0002`.
- [x] 6.3 Update `docs/rules/lambda/star_permissions.md` to describe supported inline statements and action-versus-resource wildcards; verify its examples match the passing/failing policy tests.

## 7. Direct SAM function settings and log retention

- [x] 7.1 Extend `WS1000`, `ES1005`, and `ES1006` to SAM functions with the proper SAM property shapes; verify existing target-rule fixtures and new missing/active/override cases pass on the upgraded supported range.
- [x] 7.2 Extend `WS1004` log-group association to SAM function logical IDs and explicit names; verify existing `Fn::Sub`, variable-map, `Fn::Join`, hardcoded-name, and missing-retention fixtures plus a mixed native/SAM case.
- [x] 7.3 Add host-processed Globals regression cases for inherited memory/timeout, tracing overrides, nested maps, supported `IgnoreGlobals`, and already merged lists; verify plugin checks leave the template unchanged and do not duplicate inherited entries.
- [x] 7.4 Update tracing, memory, timeout, and retention rule documentation with direct SAM examples; verify each example's expected rule result against the corresponding fixtures.

## 8. SAM invocation sources, policies, and async destinations

- [x] 8.1 Extend `WS1002` using explicit permissions plus the bounded SAM event-to-principal mapping; verify SNS/EventBridge mixed principals, repeated principals, explicit/SAM combinations, and polling-event exclusions, including the existing permission fixtures.
- [x] 8.2 Extend `WS1003` to effective inline SAM policy documents when SAM manages the role; verify wildcard/scoped actions, single/list statements, policy-template and managed-policy exclusions, and explicit-role behavior without exceptions.
- [x] 8.3 Extend `ES1007` to the six specified asynchronous event types and effective SAM failure destinations; verify existing async fixtures, inherited `Type: SQS`, explicit destinations, absent destinations, and event-level source paths.
- [x] 8.4 Document the supported invocation events, inline-policy coverage limits, and SAM-managed async destinations in the affected Lambda rule pages; verify examples and the stated event list agree with the tests and specs.

## 9. SAM stream and EventBridge destination checks

- [x] 9.1 Extend `ES1001` to SAM Kinesis/DynamoDB events without recommending unsupported destination properties on SQS events; verify missing/present stream destination cases, the SQS exclusion, native mapping fixtures, and original event paths.
- [x] 9.2 Extend `ES4000` to CloudWatchEvent/EventBridgeRule/Schedule events; verify absent DLQs, explicit ARNs, valid SAM-managed DLQs, the ScheduleV2 exclusion, and unchanged native EventBridge fixture outcomes.
- [x] 9.3 Update the event-source-mapping and EventBridge DLQ rule pages with direct SAM examples and event-specific limitations; verify the examples against the corresponding regression cases.

## 10. Explicit and implicit SAM API checks

- [x] 10.1 Extend `ES2000`, `ES2003`, and `WS2002` to explicit SAM REST/HTTP API settings; verify existing API fixtures, inherited settings, missing destinations, incomplete default throttling, local overrides, and the HTTP tracing exclusion.
- [x] 10.2 Extend `WS2001` to direct SAM literal log formats and guard unresolved intrinsic formats; verify existing JSON/non-JSON fixtures, supported context substitutions, and intrinsic inputs without `E0002` or an assertion that they are valid JSON.
- [x] 10.3 Add implicit API views from function Api/HttpApi events and applicable Globals; verify configured/missing settings, explicit and unresolved references, repeated routes per owner, and intentional per-owner findings for shared APIs.
- [x] 10.4 Add direct API regression coverage for nested `AWS::Include` without fetching definitions; verify separately declared logging/throttling still pass and missing REST tracing produces `WS2002`.
- [x] 10.5 Update the four API Gateway rule pages and `docs/cfn-lint.md` with direct/implicit SAM examples, per-owner reporting, and coverage limits; verify all documented property names and expectations against the API tests.

## 11. Whole-project integration and upgrade review

- [x] 11.1 Add a minimal self-contained direct-SAM fixture outside the filename-driven template directory and test it through the actual CLI with JSON output; verify `WS1000`, `WS1004`, and `WS2002`, no memory/timeout/logging/throttling false positives or `E0002`, and preserved suppressions without depending on or modifying `unicorn_contracts/`.
- [x] 11.2 Add CLI discovery, native-only/mixed templates, original source locations, command-line ignores, template/resource metadata suppression, and repeated-route checks; verify all 17 rule IDs remain discoverable and ignored owner findings do not suppress other owners.
- [x] 11.3 Run full Python and Go suites, CLI/plugin integrations, docs and distribution builds, scoped non-mutating lint checks, and applicable local example validations on the selected baseline; verify actual results and record specific blocked checks rather than treating them as passed.
- [x] 11.4 Exercise clean setup and the newly declared minimum/current host/runtime combinations; verify reproducible resolution, coherent metadata/CI/documentation, and selected-version reporting without restoring historical workshop-specific pins.
- [x] 11.5 Review the complete implementation against all four specs and deliver upgrade/migration guidance with selected versions, new minimums, and validation outcomes; verify preserved rule IDs/severities and native behavior, no new whole-template translator, no changed user-owned workshop files, and no deployments or publications.
