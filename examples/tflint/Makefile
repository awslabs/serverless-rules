TFLINT_FOLDER = ../../tflint-ruleset-aws-serverless
TFLINT_PLUGINS_FOLDER = $(CURDIR)/.tflint.d/plugins
TFLINT_BINARY = tflint-ruleset-aws-serverless

# Keep plugin installation local to this example directory
export TFLINT_PLUGIN_DIR = $(TFLINT_PLUGINS_FOLDER)

test-example:
	tflint --init
	$(MAKE) -C $(TFLINT_FOLDER) build
	mkdir -p $(TFLINT_PLUGINS_FOLDER)
	mv -f $(TFLINT_FOLDER)/$(TFLINT_BINARY) $(TFLINT_PLUGINS_FOLDER)/github.com/awslabs/serverless-rules/0.3.6/$(TFLINT_BINARY)
	tflint
