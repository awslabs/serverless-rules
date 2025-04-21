package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLambdaPermissionMultiplePrincipals checks if there are multiple Lambda permission with different principals for a single function
type AwsLambdaPermissionMultiplePrincipalsRule struct {
	resourceType string
	functionName string
	principal    string
	tflint.DefaultRule
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
func (r *AwsLambdaPermissionMultiplePrincipalsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsLambdaPermissionMultiplePrincipalsRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/lambda/permission_multiple_principals/"
}

// Metadata returns the rule metadata
func (r *AwsLambdaPermissionMultiplePrincipalsRule) Metadata() interface{} {
	return struct {
		Name     string
		Severity tflint.Severity
		Link     string
	}{
		Name:     r.Name(),
		Severity: r.Severity(),
		Link:     r.Link(),
	}
}

// Check checks if there are multiple Lambda permission with different principals for a single function
func (r *AwsLambdaPermissionMultiplePrincipalsRule) Check(runner tflint.Runner) error {
	permissions := make(map[string]map[string][]hcl.Expression)

	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{
				Name: r.principal,
			},
			{
				Name: r.functionName,
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		body := resource.Body
		// Get attribute value
		var principalVal string
		principal, ok := body.Attributes[r.principal]
		if !ok {
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.principal),
				resource.DefRange,
			)
			continue
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
				resource.DefRange,
			)
			continue
		}
		err = runner.EvaluateExpr(functionName.Expr, &functionNameVal, nil)
		if err != nil {
			return err
		}

		if _, ok = permissions[functionNameVal]; !ok {
			permissions[functionNameVal] = make(map[string][]hcl.Expression)
		}

		permissions[functionNameVal][principalVal] = append(permissions[functionNameVal][principalVal], principal.Expr)
	}

	// Parse permissions
	for _, fnPermissions := range permissions {
		// More than one principal type for a given functionName
		if len(fnPermissions) > 1 {
			for _, exprs := range fnPermissions {
				for _, expr := range exprs {
					runner.EmitIssue(
						r,
						fmt.Sprintf("different \"%s\" values for the same %s.", r.principal, r.functionName),
						expr.Range(),
					)
				}
			}
		}
	}

	return nil
}
