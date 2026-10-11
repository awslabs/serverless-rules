# Proposal

## Why

Recent cfn-lint versions validate SAM resources directly instead of expanding them into native CloudFormation resources, so Serverless Rules silently skips checks that depend on generated Lambda functions, API stages, permissions, and roles. The project also has outdated Python, Go, documentation, CI, and example dependencies; upgrading the complete toolchain alongside direct SAM support will establish a current, tested baseline rather than preserve the workshop server or old dependency floors.

## What Changes

- Upgrade all first-party project dependency surfaces to latest stable releases: Python runtime/dev/build packages, Go direct and transitive dependencies, documentation packages, CI actions and tools, and maintained examples including their Terraform providers.
- Refresh manifests, constraints, lockfiles/checksums, tool declarations, examples, and documentation together; resolve exact versions from official metadata at implementation start and record the validated selection.
- **BREAKING**: Allow major-version migrations and raise Python, Go, cfn-lint, TFLint, or other minimum versions as needed by the upgraded baseline. Replace unsupported/deprecated dependencies, including CDK v1 packages, and document the resulting support boundaries.
- Add direct support for `AWS::Serverless::Function`, `AWS::Serverless::Api`, and `AWS::Serverless::HttpApi`, using the existing rule IDs and severities.
- Restore the existing SAM-tested Lambda behavior: tracing, explicit memory and timeout, log retention, invocation principals, inline wildcard permissions, and asynchronous failure destinations.
- Check explicit SAM APIs and implicit APIs created by function events for access logging, structured logging, default throttling, and REST API tracing.
- Check the supported SAM event declarations corresponding to event source mappings and EventBridge rules without requiring generated resources.
- Respect effective SAM `Globals`, resource overrides, metadata suppressions, and original source paths without modifying the host template or translating the whole template.
- Normalize IAM statements supplied as either a single object or a list, fixing the `WS1003` crash exposed by the workshop template.
- Preserve native CloudFormation checks and Python/Go rule behavior while adapting to current dependency APIs. Test the newly declared support floor and selected current/latest releases rather than requiring historical cfn-lint versions.
- Document supported SAM declarations, dependency/tool requirements, migrations, and explicit limitations. Do not rebuild the workshop server, change user-owned workshop files, change rule IDs/severities, add unrelated rules, or implement a general SAM transformer, policy-template expander, or intrinsic-function evaluator.

## Capabilities

### New Capabilities

- `sam-function-rules`: Direct checks of SAM functions, their inline policies, invocation sources, event source mappings, and associated log groups.
- `sam-api-rules`: Direct checks of explicit and implicit SAM REST/HTTP APIs with their effective stage settings.
- `cfn-rule-compatibility`: Stable rule discovery, native/SAM coexistence, source-aware suppression, coverage of the upgraded host support range, and robust IAM statement handling.
- `dependency-maintenance`: A current stable dependency/toolchain baseline with coherent manifests, reproducible development/CI environments, maintained examples, and documented support boundaries.

### Modified Capabilities

None. The project has no existing OpenSpec capabilities.

## Impact

- Python manifests and lockfiles, build/dev tooling, SAM-aware rule implementations, tests, and any adaptations required by current cfn-lint APIs.
- Go module manifests/checksums, SDK/HCL integrations, Go toolchain settings, TFLint host compatibility, tests, and local GoReleaser configuration.
- Root docs dependencies, GitHub Actions and CI/release tool versions, dependency-update coverage, Make targets, and contributor instructions.
- CDK, SAM, Serverless Framework, and Terraform/TFLint examples: dependency declarations, framework/CLI requirements, provider constraints, and compatibility migrations. Example validation remains local and does not deploy infrastructure.
- Existing rule documentation plus installation, support-range, example, and upgrade guidance. No upstream cfn-lint changes and no new whole-template SAM translator dependency.
- The user-owned `unicorn_contracts/` tree remains untouched and is not an implementation dependency. Keep a generic self-contained SAM regression case; remote workshop rebuild/validation is deferred.
