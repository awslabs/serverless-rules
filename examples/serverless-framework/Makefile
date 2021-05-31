test-example:
	sls package
	PYTHONPATH="../../cfn-lint-serverless/:$$PYTHONPATH" cfn-lint .serverless/cloudformation-template-update-stack.json -a cfn_lint_serverless.rules