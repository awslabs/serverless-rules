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

add-rule:
ifeq ($(RULE_NAME),)
	$(error Missing RULE_NAME environment variable)
endif
ifeq ($(RULE_NAME_CC),)
	$(error Missing RULE_NAME_CC environment variable)
endif
	gomplate -f templates/rule.go.tmpl -o rules/$(RULE_NAME).go
	gomplate -f templates/rule_test.go.tmpl -o rules/$(RULE_NAME)_test.go