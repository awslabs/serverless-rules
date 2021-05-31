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