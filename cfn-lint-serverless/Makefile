# Development environment setup
dev:
	pip install --upgrade pip poetry
	poetry install

format:
	poetry run isort cfn_lint_serverless tests
	poetry run black cfn_lint_serverless tests

lint: format
	poetry run flake8 cfn_lint_serverless/* tests/*

test:
	poetry run pytest -vvv --cov=./ --cov-report=xml

security-baseline:
	poetry run bandit --baseline bandit.baseline -r cfn_lint_serverless

complexity-baseline:
	$(info Maintenability index)
	poetry run radon mi cfn_lint_serverless
	$(info Cyclomatic complexity index)
	poetry run xenon --max-absolute C --max-modules A --max-average A cfn_lint_serverless

pr: lint test security-baseline complexity-baseline

release:
	poetry build
	poetry publish