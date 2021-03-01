CFN_LINT = cfn-lint-serverless
TFLINT = tflint-ruleset-aws-serverless

dev:
	make -C $(CFN_LINT) dev
	make -C $(TFLINT) dev

pr:
	make -C $(CFN_LINT) pr
	make -C $(TFLINT) pr