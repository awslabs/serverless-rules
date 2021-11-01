package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type AsyncPermission struct {
	functionName string
	found        bool
	expression   hcl.Expression
}

// AwsLambdaEventInvokeConfigAsyncOnFailure checks if an event invoke config has a destination on failure if the function has permission for an async principal
type AwsLambdaEventInvokeConfigAsyncOnFailureRule struct {
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
func (r *AwsLambdaEventInvokeConfigAsyncOnFailureRule) Severity() string {
	// TODO: Determine the rule's severiry
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsLambdaEventInvokeConfigAsyncOnFailureRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/lambda/async_failure_destination/"
}

// Check checks if an event invoke config has a destination on failure if the function has permission for an async principal
func (r *AwsLambdaEventInvokeConfigAsyncOnFailureRule) Check(runner tflint.Runner) error {
	var asyncPerms []AsyncPermission

	// Scan permissions
	err := runner.WalkResources(r.permissionType, func(resource *configs.Resource) error {
		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name: r.functionName,
				},
				{
					Name: r.principal,
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		// Get the permission principal
		var principalVal string
		principal, ok := body.Attributes[r.principal]
		if !ok {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				body.MissingItemRange,
			)
			return nil
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
			return nil
		}

		// Get the function name
		var functionNameVal string
		functionName, ok := body.Attributes[r.functionName]
		if !ok {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				body.MissingItemRange,
			)
			return nil
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

		return nil
	})

	if err != nil {
		return err
	}

	// Scan Event Invoke Config
	err = runner.WalkResources(r.eventInvokeConfigType, func(resource *configs.Resource) error {
		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name: r.functionName,
				},
			},
			Blocks: []hcl.BlockHeaderSchema{
				{
					Type: r.block1Name,
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		// Get the function name
		var functionNameVal string
		functionName, ok := body.Attributes[r.functionName]
		if !ok {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				body.MissingItemRange,
			)
			return nil
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
			return nil
		}

		// Check for block level 1
		blocks := body.Blocks.OfType(r.block1Name)
		if len(blocks) != 1 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.block1Name),
				body.MissingItemRange,
			)

			return nil
		}

		// Check for block level 2
		body, _, diags = blocks[0].Body.PartialContent(&hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{
				{
					Type: r.block2Name,
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		blocks = body.Blocks.OfType(r.block2Name)
		if len(blocks) != 1 {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.block2Name),
				body.MissingItemRange,
			)

			return nil
		}

		// Check for attribute in block
		body, _, diags = blocks[0].Body.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name: r.attributeName,
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		if _, ok = body.Attributes[r.attributeName]; !ok {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.attributeName),
				body.MissingItemRange,
			)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Scan for missing Event Invoke Config.
	for _, asyncPerm := range asyncPerms {
		if !asyncPerm.found {
			runner.EmitIssueOnExpr(
				r,
				fmt.Sprintf("missing \"%s\" resource for %s.", r.eventInvokeConfigType, r.functionName),
				asyncPerm.expression,
			)
		}
	}

	return nil
}
