package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type AsyncPermission struct {
	functionName string
	found        bool
	expression   hcl.Expression
}

// AwsLambdaEventInvokeConfigAsyncOnFailure checks if an event invoke config has a destination on failure if the function has permission for an async principal
type AwsLambdaEventInvokeConfigAsyncOnFailureRule struct {
	tflint.DefaultRule
	permissionType        string
	eventInvokeConfigType string
	functionName          string
	principal             string
	block1Name            string
	block2Name            string
	attributeName         string
	enumPrincipal         []string
}

// NewAwsLambdaEventInvokeConfigAsyncOnFailureRule returns new rule with default attributes
func NewAwsLambdaEventInvokeConfigAsyncOnFailureRule() *AwsLambdaEventInvokeConfigAsyncOnFailureRule {
	return &AwsLambdaEventInvokeConfigAsyncOnFailureRule{
		permissionType:        "aws_lambda_permission",
		eventInvokeConfigType: "aws_lambda_function_event_invoke_config",
		functionName:          "function_name",
		principal:             "principal",
		block1Name:            "destination_config",
		block2Name:            "on_failure",
		attributeName:         "destination",
		enumPrincipal: []string{
			"events.amazonaws.com",
			"events.amazonaws.com.cn",
			"iot.amazonaws.com",
			"iot.amazonaws.com.cn",
			"s3.amazonaws.com",
			"s3.amazonaws.com.cn",
			"sns.amazonaws.com",
			"sns.amazonaws.com.cn",
		},
	}
}

// Name returns the rule name
func (r *AwsLambdaEventInvokeConfigAsyncOnFailureRule) Name() string {
	return "aws_lambda_event_invoke_config_async_on_failure"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsLambdaEventInvokeConfigAsyncOnFailureRule) Enabled() bool {
	// TODO: Determine whether the rule is enabled by default
	return false
}

// Severity returns the rule severity
func (r *AwsLambdaEventInvokeConfigAsyncOnFailureRule) Severity() tflint.Severity {
	return tflint.ERROR // TODO: Determine the rule's severity
}

// Link returns the rule reference link
func (r *AwsLambdaEventInvokeConfigAsyncOnFailureRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/lambda/async_failure_destination/"
}

// Check checks if an event invoke config has a destination on failure if the function has permission for an async principal
func (r *AwsLambdaEventInvokeConfigAsyncOnFailureRule) Check(runner tflint.Runner) error {
	var asyncPerms []AsyncPermission

	// Scan permissions
	resources, err := runner.GetResourceContent(r.permissionType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.functionName},
			{Name: r.principal},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		// Get the permission principal
		var principalVal string
		principal, ok := resource.Body.Attributes[r.principal]
		if !ok {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				resource.DefRange,
			)
			continue
		}
		err := runner.EvaluateExpr(principal.Expr, &principalVal, nil)
		if err != nil {
			return err
		}

		// Check if the permission is for a principal that invokes the function asynchronously
		isAsync := false
		for _, asyncPrincipal := range r.enumPrincipal {
			if principalVal == asyncPrincipal {
				isAsync = true
				break
			}
		}
		// Permission is not async, break early
		if !isAsync {
			continue
		}

		// Get the function name
		var functionNameVal string
		functionName, ok := resource.Body.Attributes[r.functionName]
		if !ok {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				resource.DefRange,
			)
			continue
		}
		err = runner.EvaluateExpr(functionName.Expr, &functionNameVal, nil)
		if err != nil {
			return err
		}

		// Add the function name to the list of async functions
		asyncPerms = append(asyncPerms, AsyncPermission{
			functionName: functionNameVal,
			found:        false,
			expression:   functionName.Expr,
		})
	}

	// Scan Event Invoke Config
	configs, err := runner.GetResourceContent(r.eventInvokeConfigType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.functionName},
		},
		Blocks: []hclext.BlockSchema{
			{
				Type: r.block1Name,
				Body: &hclext.BodySchema{
					Blocks: []hclext.BlockSchema{
						{
							Type: r.block2Name,
							Body: &hclext.BodySchema{
								Attributes: []hclext.AttributeSchema{
									{Name: r.attributeName},
								},
							},
						},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, config := range configs.Blocks {
		// Get the function name
		var functionNameVal string
		functionName, ok := config.Body.Attributes[r.functionName]
		if !ok {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				config.DefRange,
			)
			continue
		}
		err = runner.EvaluateExpr(functionName.Expr, &functionNameVal, nil)
		if err != nil {
			return err
		}

		// Check if the function is async
		isAsync := false
		for pos := range asyncPerms {
			if functionNameVal == asyncPerms[pos].functionName {
				asyncPerms[pos].found = true
				isAsync = true
				break
			}
		}

		// Function is not async, break early
		if !isAsync {
			continue
		}

		// Check for block level 1
		blocks := config.Body.Blocks.OfType(r.block1Name)
		if len(blocks) == 0 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.block1Name),
				config.DefRange,
			)
			continue
		}

		// Check for block level 2
		onFailureBlocks := blocks[0].Body.Blocks.OfType(r.block2Name)
		if len(onFailureBlocks) == 0 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.block2Name),
				blocks[0].DefRange,
			)
			continue
		}

		// Check for attribute in block
		_, exists := onFailureBlocks[0].Body.Attributes[r.attributeName]
		if !exists {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				onFailureBlocks[0].DefRange,
			)
		}
	}

	// Scan for missing Event Invoke Config.
	for _, asyncPerm := range asyncPerms {
		if !asyncPerm.found {
			runner.EmitIssue(
				r,
				fmt.Sprintf("missing \"%s\" resource for %s.", r.eventInvokeConfigType, r.functionName),
				asyncPerm.expression.Range(),
			)
		}
	}

	return nil
}
