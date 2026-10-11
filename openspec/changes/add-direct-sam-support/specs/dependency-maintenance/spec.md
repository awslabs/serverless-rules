# Spec Delta

## Purpose

Maintain a current stable dependency and toolchain baseline across the Python and Go plugins, documentation, CI, and examples, with reproducible environments and explicit supported-version boundaries.

## ADDED Requirements

### Requirement: Project-wide stable dependency selection
The project SHALL upgrade its first-party dependency surfaces to latest stable releases selected from official release metadata at implementation start. Major-version migrations and higher minimum versions SHALL be permitted. Prereleases SHALL not be selected by default, and an incompatible latest release SHALL be reported rather than silently replaced by an older target.

#### Scenario: All dependency surfaces are considered
- **WHEN** the upgrade baseline is selected
- **THEN** it covers Python runtime/dev/build packages, Go dependencies, documentation packages, CI actions/tools, example frameworks/CLIs, and Terraform providers
- **AND** the selected versions and required migrations are recorded

#### Scenario: Latest release requires a migration
- **WHEN** a latest stable dependency changes its supported API or runtime requirements
- **THEN** the project migrates to those requirements and validates the affected behavior
- **AND** it does not retain the old release merely to preserve the previous minimum version

#### Scenario: Latest release cannot be integrated
- **WHEN** the selected latest stable release cannot satisfy required project behavior
- **THEN** the incompatibility is reported as a blocker requiring an explicit decision
- **AND** the upgrade is not represented as complete with an undisclosed older pin

### Requirement: Coherent supported-version boundaries
The project SHALL declare validated minimum Python, Go, cfn-lint, TFLint, and other required tool versions appropriate to the upgraded baseline. Package metadata, toolchain declarations, installation guidance, examples, and CI SHALL agree on those boundaries. The workshop server's versions SHALL not determine the supported range.

#### Scenario: Runtime floor is raised
- **WHEN** an upgraded dependency requires a newer language or host version
- **THEN** manifests, relevant CI environments, and documentation reflect the new validated minimum
- **AND** obsolete compatibility claims and unsupported CI jobs are removed

#### Scenario: Below-floor installation
- **WHEN** a user attempts a Python package installation below its declared Python or cfn-lint minimum
- **THEN** package metadata communicates the incompatibility rather than promising an unsupported installation

### Requirement: Reproducible Python development and builds
Updated Python dependency declarations and lock/constraint data SHALL resolve coherently for the new supported range. Clean development, test, documentation, and package-build environments SHALL use recorded stable versions without unintended reversion to the old lockfile baseline.

#### Scenario: Fresh development environment
- **WHEN** the documented clean setup and test commands run
- **THEN** they select the upgraded validated baseline
- **AND** dependency resolution succeeds with runtime, development, and build requirements consistent

#### Scenario: Distribution build
- **WHEN** the Python distribution is built locally in an isolated build environment
- **THEN** the build succeeds
- **AND** its declared host/runtime requirements match the new support boundary

### Requirement: Current Go plugin integration
The Go plugin SHALL build and preserve its rule behavior with the selected stable HCL and TFLint SDK releases, compatible updated transitive dependencies, and the declared Go toolchain. Its binary SHALL load in the selected supported TFLint host without changing existing rule identifiers.

#### Scenario: Upgraded Go test and build
- **WHEN** the Go plugin is tested and built with the selected dependency baseline
- **THEN** its unit tests and rule-registration checks pass
- **AND** module/checksum data remain consistent

#### Scenario: TFLint host loads the plugin
- **WHEN** the newly built plugin is loaded locally by the selected TFLint host
- **THEN** its existing rules are discoverable and execute on a representative Terraform example
- **AND** no deployment or publication is needed to verify integration

### Requirement: Current documentation and CI tooling
Documentation builds and CI validation SHALL use stable, recorded dependency/tool versions consistent with the declared runtime range. Dependency references to floating development branches SHALL be replaced with stable release references or immutable references to those releases. Local configuration checks SHALL not publish docs or release artifacts.

#### Scenario: Documentation is built locally
- **WHEN** documentation is built with the upgraded docs packages
- **THEN** the site's configured features and build succeed without a deployment

#### Scenario: CI toolchain is checked
- **WHEN** updated test, documentation, and release configurations are validated
- **THEN** their action/tool references and language versions are consistent with the selected baseline
- **AND** validation does not invoke publishing steps

### Requirement: Maintained examples on upgraded dependencies
Project examples SHALL use current stable framework/CLI and provider dependencies with the migrations required to keep their documented local workflows functional. Deprecated CDK v1 packages SHALL be replaced by the current supported CDK library model. Example checks SHALL synthesize, build, package, or validate locally without deploying infrastructure.

#### Scenario: CDK example is migrated
- **WHEN** the CDK example is installed and synthesized using the upgraded dependencies
- **THEN** it no longer depends on CDK v1 packages
- **AND** it produces CloudFormation for the expected example resources

#### Scenario: Other examples are checked
- **WHEN** SAM, Serverless Framework, and Terraform/TFLint examples run their applicable local checks
- **THEN** updated runtime declarations, tool requirements, provider constraints, and configuration syntax are valid for the selected baseline
- **AND** no infrastructure is deployed

### Requirement: Upgrades preserve rule behavior
Dependency/API migrations SHALL preserve existing Python and Go rule identifiers, severities, and applicable native-template outcomes. Direct SAM functionality SHALL be validated on the upgraded baseline rather than assumed to follow from successful dependency resolution.

#### Scenario: Regression suites run after migration
- **WHEN** Python and Go rule tests execute with upgraded dependencies
- **THEN** their existing target-rule expectations remain valid
- **AND** the new direct-SAM requirements are separately verified

### Requirement: Transparent upgrade validation
The project SHALL document selected dependency/tool versions, new minimums, major migrations, and validation outcomes. Required checks blocked by credentials, licensing, unavailable tools, or network access SHALL be identified as blocked, not reported as passing. Workshop rebuilds, release publication, and remote deployment SHALL remain outside the upgrade validation.

#### Scenario: A check requires credentials
- **WHEN** an example's local packaging check requires unavailable authentication or licensing
- **THEN** that limitation is recorded explicitly
- **AND** no account creation, remote deployment, or unrequested publishing is attempted

#### Scenario: Upgrade is ready for review
- **WHEN** the upgrade is presented for review
- **THEN** dependency selections, support changes, migrations, and actual test/build results are available
- **AND** user-owned workshop files remain unchanged
