dev:
	@echo "Nothing to do here"

format:
	go fmt

lint: format

test:
	go test ./...

pr: lint test

build:
	go build

install: build
	mkdir -p ~/.tflint.d/plugins
	mv ./tflint-ruleset-template ~/.tflint.d/plugins