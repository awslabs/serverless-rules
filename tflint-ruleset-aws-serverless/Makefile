dev:
	@echo "Nothing to do here"

format:
	go fmt

lint: format

test:
	go test ./...

pr: lint test
ifneq ($(shell grep -oP "^\t+New[A-Za-z0-9]+\(\)," rules/provider.go | wc -l), $(shell grep -oP "^func New[A-Za-z0-9]+\(\)" rules/* | wc -l))
	$(error Mismatch in rule count ($(shell grep -oP "^\t+New[A-Za-z0-9]+\(\)," rules/provider.go | wc -l) vs $(shell grep -oP "^func New[A-Za-z0-9]+\(\)" rules/* | wc -l)) - check rules/provider.go)
else
	$(info Match in rule count)
endif

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