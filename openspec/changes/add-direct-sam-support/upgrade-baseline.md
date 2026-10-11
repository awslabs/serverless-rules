# Upgrade Baseline Selection

Resolved from official release metadata (PyPI, Go module proxy via direct VCS, npm registry, GitHub releases, Terraform Registry, go.dev, HashiCorp checkpoint API) on **2026-10-10**. Prereleases excluded. No workshop-server versions were used as targets.

## Selected stable versions

| Surface | Component | Selected | Notes |
|---|---|---|---|
| Python plugin | Python runtime floor | 3.10 | Required by cfn-lint 1.57.x, hatchling, bandit, coverage, pytest |
| Python plugin | Python supported range | 3.10-3.14 | cfn-lint declares `<3.15`; 3.14.8 is the local/current interpreter |
| Python plugin | cfn-lint | 1.57.2 (floor 1.57.0) | Validated floor 1.57.0, current 1.57.2; replaces `>=1.33.0` |
| Python plugin | hatchling | 1.32.4 | Build backend |
| Python plugin | bandit | 1.9.4 | Dev |
| Python plugin | coverage | 7.16.2 | Dev |
| Python plugin | pytest | 9.1.1 | Dev, major upgrade from 8.x |
| Python plugin | pytest-cov | 7.1.0 | Dev, major upgrade from 6.x |
| Python plugin | radon | 6.0.1 | Dev, unchanged |
| Python plugin | ruff | 0.17.0 | Dev |
| Python plugin | xenon | 0.9.3 | Dev, unchanged (latest available) |
| Go plugin | Go toolchain | 1.27 (latest 1.27.2) | Replaces `go 1.23` |
| Go plugin | github.com/hashicorp/hcl/v2 | v2.25.0 | Replaces v2.23.0 |
| Go plugin | tflint-plugin-sdk | v0.25.0 | Replaces v0.21.0 |
| Go plugin | TFLint host | v0.64.0 | Replaces v0.50.0 in release CI |
| Go plugin | GoReleaser | v2.18.3 | Major upgrade; config must use `version: 2` schema |
| Docs | mkdocs | 1.6.1 | Root `requirements-dev.txt`, now pinned |
| Docs | mkdocs-material | 9.7.7 | Root `requirements-dev.txt`, now pinned |
| CI | actions/checkout | v7 (v7.0.1) | Major upgrade from v4 |
| CI | actions/setup-python | v7 (v7.0.0) | Major upgrade from v5 |
| CI | actions/setup-go | v7 (v7.0.0) | Major upgrade from v5 |
| CI | actions/setup-node | v7 (v7.1.0) | For Node-dependent example jobs |
| CI | astral-sh/setup-uv | v10 (v10.3.0) | Major upgrade from v5/v6 |
| CI | goreleaser/goreleaser-action | v7 (v7.2.3) | Major upgrade from v6 |
| CI | release-drafter/release-drafter | v7 (v7.10.0) | Major upgrade from v6 |
| CI | mhausenblas/mkdocs-deploy-gh-pages | 1.26 | Replaces floating `@master` |
| Examples | CDK CLI (npm `aws-cdk`) | 2.1145.0 | Replaces CDK v1 |
| Examples | aws-cdk-lib (Python) | 2.273.0 | Replaces split CDK v1 packages |
| Examples | SAM CLI | 1.168.0 | Local 1.167.0 acceptable within 1.x |
| Examples | Serverless Framework | 4.43.0 | Stays on framework v4 |
| Examples | Terraform | 1.16.5 (local 1.13.3 usable) | `required_version` allows current 1.x |
| Examples | hashicorp/aws provider | 6.68.0 | Major upgrade from 5.40.0 |
| Examples | hashicorp/random provider | 3.9.1 | Replaces 3.7.1 |

## Major migrations required

- CDK example: v1 split packages (`aws-cdk.core`, `aws-cdk.aws-*`) to CDK v2 `aws-cdk-lib` + `constructs`.
- Terraform AWS provider 5.x to 6.x in the TFLint example.
- TFLint plugin SDK v0.21 to v0.25; verify `helper.TestRunner`, `GetResourceContent`, `EvaluateExpr`, and `plugin.Serve` APIs.
- GoReleaser configuration to the v2 schema (`version: 2`).
- pytest 8 to 9 and pytest-cov 6 to 7 in the Python dev group.
- CI actions to their current major versions; docs deploy action pinned to tag `1.26` instead of `@master`.
- Python floor raised from 3.9.1 to 3.10; cfn-lint floor raised from 1.33.0 to 1.57.0; Go floor raised from 1.23 to 1.27.

## Proposed validated floors

- Python `>=3.10,<3.15`; CI tests 3.10 (floor) and 3.14 (current), with intermediate versions as practical.
- cfn-lint `>=1.57.0`; isolated validation at 1.57.0 (floor) and 1.57.2 (current/latest at selection date).
- Go `1.27`; CI tests the declared toolchain.
- TFLint host `v0.64.0` validated against SDK `v0.25.0` plugin builds.
- A latest-stable drift job may resolve newer releases and must print the versions it used.

## Validation outcomes

Recorded here as implementation progresses; see `tasks.md` groups 1-5 for per-area checks and group 11 for the final integration record.

### Blocked checks

- **Serverless Framework v4 packaging** (`examples/serverless-framework`, `sls package`): Serverless Framework v4.43.0 requires interactive login, `SERVERLESS_ACCESS_KEY`, or `SERVERLESS_LICENSE_KEY`. No account or credential was created during this upgrade, so packaging was not executed. The example configuration was updated to current v4 syntax (`provider.runtime: python3.13`) but its end-to-end packaging check remains blocked on user-provided licensing.

### Completed validation results (2026-10-10)

- **Python plugin**: `uv sync --locked` clean install on Python 3.14.8; locked cfn-lint 1.57.2; `uv build` isolated distribution with `Requires-Python: >=3.10,<3.15` and `Requires-Dist: cfn-lint>=1.57.0,<2`.
- **Python tests**: 112 passed (`make pr`: pytest + coverage, ruff 0.17.0 format/check, bandit 1.9.4, radon, xenon) on the locked baseline; the full suite also passes on the floor combination Python 3.10.22 + cfn-lint 1.57.0 in an isolated pip environment.
- **CLI regression**: all 17 rule IDs listed by `cfn-lint --list-rules -a cfn_lint_serverless.rules`; the self-contained `tests/cli/sam_app.yaml` fixture reports WS1000/WS1004/WS2002 without ES1005/ES1006/ES2000/ES2003 false positives or E0002; command-line ignores and resource metadata suppressions honored.
- **Go plugin**: HCL v2.25.0, tflint-plugin-sdk v0.25.0, Go 1.27; `go mod tidy` consistent, `go vet ./...` clean, `go test ./...` passes, rule IDs unchanged; `goreleaser check` and a local `--snapshot` release pass with GoReleaser v2.18.3 (no publish).
- **TFLint host**: plugin built from source loads and reports expected `aws_lambda_function_default_memory`/`aws_lambda_function_default_timeout` findings under TFLint v0.64.0 from a temporary plugin directory.
- **Docs**: `mkdocs build --strict` passes with mkdocs 1.6.1 and mkdocs-material 9.7.7 from pinned `requirements-dev.txt` (no deployment).
- **Examples**: CDK example migrated to aws-cdk-lib 2.273.0/constructs 10.8.1 and synthesizes with CDK CLI 2.1145.0 (expected ES1005/ES1006 findings lint cleanly through the plugin); SAM example passes `sam validate --lint` (SAM CLI 1.167.0) and produces only the expected ES1005/ES1006 findings; Terraform example pins aws 6.68.0/random 3.9.1, passes `terraform init -backend=false` + `terraform validate` (Terraform 1.13.3) and TFLint with the locally built plugin.
- **CI**: workflows updated to checkout v7, setup-python v7, setup-go v7, setup-uv v10, setup-tflint v6, goreleaser-action v7, release-drafter v7; docs deploy uses pinned actions (peaceiris/actions-gh-pages v4) instead of the floating `master` reference; cfn-lint floor/latest compat jobs install explicit versions without lockfile resync.
- **Untouched as required**: no rule IDs/severities changed, `unicorn_contracts/` unchanged, no publishing, no deployments.
