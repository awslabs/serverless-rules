CFN_LINT = cfn-lint-serverless
TFLINT = tflint-ruleset-aws-serverless
EXAMPLES = examples

dev: requirements-dev
	$(MAKE) -C $(CFN_LINT) dev
	$(MAKE) -C $(TFLINT) dev

dev-requirements:
	pip install -r requirements-dev.txt

pr:
	$(MAKE) -C $(CFN_LINT) pr
	$(MAKE) -C $(TFLINT) pr

test:
	$(MAKE) -C $(CFN_LINT) test
	$(MAKE) -C $(TFLINT) test

docs-serve:
	mkdocs serve

release-check:
	grep "version = \"$$RELEASE_TAG_VERSION\"" cfn-lint-serverless/pyproject.toml
	grep "Version: \"$$RELEASE_TAG_VERSION\"" tflint-ruleset-aws-serverless/main.go
	grep "version = \"$$RELEASE_TAG_VERSION\"" README.md
	grep "version = \"$$RELEASE_TAG_VERSION\"" docs/tflint.md
	grep "version = \"$$RELEASE_TAG_VERSION\"" examples/tflint/.tflint.hcl