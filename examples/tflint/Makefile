TFLINT_FOLDER = "../../tflint-ruleset-aws-serverless"
TFLINT_PLUGINS_FOLDER = ".tflint.d/plugins"
TFLINT_BINARY = "tflint-ruleset-aws-serverless"

test-example:
	$(MAKE) -C $(TFLINT_FOLDER) build
	mkdir -p $(TFLINT_PLUGINS_FOLDER)
	mv -f $(TFLINT_FOLDER)/$(TFLINT_BINARY) $(TFLINT_PLUGINS_FOLDER)/$(TFLINT_BINARY)
	tflint