dev:
	@echo "Nothing to do here"

format:
	go fmt

lint: format

test:
	go test ./...

pr: lint test check-rules

build:
	go build

install: build
	mkdir -p ~/.tflint.d/plugins
	mv ./tflint-ruleset-template ~/.tflint.d/plugins

# Check if the number of rules in provider.go corresponds to the total number of rules
check-rules:
ifneq ($(shell grep -oP "^\t+New[A-Za-z0-9]+\(\)," rules/provider.go | wc -l), $(shell grep -oP "^func New[A-Za-z0-9]+\(\)" rules/* | wc -l))
	$(error Mismatch in rule count ($(shell grep -oP "^\t+New[A-Za-z0-9]+\(\)," rules/provider.go | wc -l) vs $(shell grep -oP "^func New[A-Za-z0-9]+\(\)" rules/* | wc -l)) - check rules/provider.go)
else
	$(info Match in rule count)
endif

# Create a new rule
add-rule:
ifeq ($(RULE_NAME),)
	$(error Missing RULE_NAME environment variable)
endif
ifeq ($(RULE_NAME_CC),)
	$(error Missing RULE_NAME_CC environment variable)
endif
	gomplate -f templates/rule.go.tmpl -o rules/$(RULE_NAME).go
	gomplate -f templates/rule_test.go.tmpl -o rules/$(RULE_NAME)_test.go

# List all the rules for provider.go
list-rules:
	grep -ohE 'func New[^\(]+\(\)' rules/* | sed 's/func //' | sed 's/$$/,/' | sort