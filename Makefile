CFN_LINT = cfn-lint-serverless
TFLINT = tflint-ruleset-aws-serverless
EXAMPLES = examples

dev:
	$(MAKE) -C $(CFN_LINT) dev
	$(MAKE) -C $(TFLINT) dev

pr:
	$(MAKE) -C $(CFN_LINT) pr
	$(MAKE) -C $(TFLINT) pr

test:
	$(MAKE) -C $(CFN_LINT) test
	$(MAKE) -C $(TFLINT) test
	$(MAKE) -C $(EXAMPLES) test

release-check:
	grep "version = \"$$RELEASE_TAG_VERSION\"" cfn-lint-serverless/pyproject.toml
	grep "Version: \"$$RELEASE_TAG_VERSION\"" tflint-ruleset-aws-serverless/main.go
	grep "version = \"$$RELEASE_TAG_VERSION\"" README.md
	grep "version = \"$$RELEASE_TAG_VERSION\"" docs/tflint.md