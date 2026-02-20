## 1. Update Version Strings in All Files

- [x] 1.1 Update version from `0.3.4` to `0.3.5` in `cfn-lint-serverless/pyproject.toml`
- [x] 1.2 Update Version constant from `0.3.4` to `0.3.5` in `tflint-ruleset-aws-serverless/main.go`
- [x] 1.3 Update version example from `0.3.4` to `0.3.5` in `README.md` (tflint plugin block)
- [x] 1.4 Update version example from `0.3.4` to `0.3.5` in `docs/tflint.md`
- [x] 1.5 Update version from `0.3.4` to `0.3.5` in `examples/tflint/.tflint.hcl`

## 2. Validate Version Consistency

- [x] 2.1 Run `RELEASE_TAG_VERSION=0.3.5 make release-check` locally to verify all version strings match
- [x] 2.2 Verify all occurrences use format `0.3.5` (not `v0.3.5`)

## 3. Commit and Merge Changes

- [x] 3.1 Create branch for version bump (e.g., `release/v0.3.5`)
- [x] 3.2 Commit version changes with message following conventional commits format
- [ ] 3.3 Push branch and create pull request
- [ ] 3.4 Wait for CI checks to pass on PR
- [ ] 3.5 Merge PR to main branch

## 4. Create Release Tag and Publish

- [ ] 4.1 Pull latest main branch with merged version changes
- [ ] 4.2 Create git tag `v0.3.5` on the merged commit
- [ ] 4.3 Push tag to GitHub (`git push origin v0.3.5`)
- [ ] 4.4 Create GitHub release for tag `v0.3.5` (triggers publish workflow)
- [ ] 4.5 Monitor GitHub Actions publish workflow for success

## 5. Verify Published Artifacts

- [ ] 5.1 Verify Python package published to PyPI at version 0.3.5
- [ ] 5.2 Verify Go plugin binary published to GitHub releases for tag v0.3.5
- [ ] 5.3 Test installing cfn-lint-serverless from PyPI (`pip install cfn-lint-serverless==0.3.5`)
