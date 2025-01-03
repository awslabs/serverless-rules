# Development environment setup
dev:
	uv sync --all-groups

format:
	ruff format

lint: format
	ruff check

test:
	PYTHONPATH=".:$$PYTHONPATH" uv run pytest -vvv --cov=./ --cov-report=xml

test-template:
	PYTHONPATH=".:$$PYTHONPATH" cfn-lint $(TEMPLATE) -a cfn_lint_serverless.rules

security-baseline:
	uv run bandit -f json -r cfn_lint_serverless -o bandit.baseline || true

complexity-baseline:
	$(info Maintenability index)
	uv run radon mi cfn_lint_serverless
	$(info Cyclomatic complexity index)
	uv run xenon --max-absolute C --max-modules B --max-average A cfn_lint_serverless

pr: lint test security-baseline complexity-baseline

release:
	uv build
	uv publish --token ${PYPI_TOKEN25}