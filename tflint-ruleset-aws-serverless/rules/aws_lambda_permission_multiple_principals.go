package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLambdaPermissionMultiplePrincipals checks if there are multiple Lambda permission with different principals for a single function
type AwsLambdaPermissionMultiplePrincipalsRule struct {
	resourceType string
	functionName string
	principal    string
}

// NewAwsLambdaPermissionMultiplePrincipalsRule returns new rule with default attributes
func NewAwsLambdaPermissionMultiplePrincipalsRule() *AwsLambdaPermissionMultiplePrincipalsRule {
	return &AwsLambdaPermissionMultiplePrincipalsRule{
		resourceType: "aws_lambda_permission",
		functionName: "function_name",
		principal:    "principal",
	}
}

// Name returns the rule name
func (r *AwsLambdaPermissionMultiplePrincipalsRule) Name() string {
	return "aws_lambda_permission_multiple_principals"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsLambdaPermissionMultiplePrincipalsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsLambdaPermissionMultiplePrincipalsRule) Severity() string {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsLambdaPermissionMultiplePrincipalsRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/lambda/permission_multiple_principals/"
}

// Check checks if there are multiple Lambda permission with different principals for a single function
func (r *AwsLambdaPermissionMultiplePrincipalsRule) Check(runner tflint.Runner) error {
	permissions := make(map[string]map[string][]hcl.Expression)

	err := runner.WalkResources(r.resourceType, func(resource *configs.Resource) error {
		// Attribute
		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name: r.principal,
				},
				{
					Name: r.functionName,
				},
			},
		})

		if diags.HasErrors() {
			return diags
		}

		// Get attribute value
		var principalVal string
		principal, ok := body.Attributes[r.principal]
		if !ok {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.principal),
				body.MissingItemRange,
			)

			return nil
		}
		err := runner.EvaluateExpr(principal.Expr, &principalVal, nil)
		if err != nil {
			return err
		}

		// Get reference attribute value
		var functionNameVal string
		functionName, ok := body.Attributes[r.functionName]
		if !ok {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.functionName),
				body.MissingItemRange,
			)

			return nil
		}
		err = runner.EvaluateExpr(functionName.Expr, &functionNameVal, nil)
		if err != nil {
			return err
		}

		if _, ok = permissions[functionNameVal]; !ok {
			permissions[functionNameVal] = make(map[string][]hcl.Expression)
		}

		permissions[functionNameVal][principalVal] = append(permissions[functionNameVal][principalVal], principal.Expr)

		return nil
	})

	// Parse permissions
	for _, fnPermissions := range permissions {
		// More than one principal type for a given functionName
		if len(fnPermissions) > 1 {
			for _, exprs := range fnPermissions {
				for _, expr := range exprs {
					runner.EmitIssueOnExpr(
						r,
						fmt.Sprintf("different \"%s\" values for the same %s.", r.principal, r.functionName),
						expr,
					)
				}
			}
		}
	}

	return err
}
