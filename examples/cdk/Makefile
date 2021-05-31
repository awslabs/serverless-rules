test-example:
	cdk synth > template.yaml
	PYTHONPATH="../../cfn-lint-serverless/:$$PYTHONPATH" cfn-lint template.yaml -a cfn_lint_serverless.rules -i W3005,W2001