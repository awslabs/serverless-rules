package rules

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type awsIamStatement struct {
	Action    interface{}            `json:"Action"`
	Effect    string                 `json:"Effect"`
	Principal map[string]interface{} `json:"Principal"`
}

type awsIamAssumeRole struct {
	Version   string            `json:"Version"`
	Statement []awsIamStatement `json:"Statement"`
}

// AwsIamRoleLambdaNoStar checks if an IAM role with a Lambda principal has broad permissions
type AwsIamRoleLambdaNoStarRule struct {
	tflint.DefaultRule
	resourceType    string
	principalNames  []string
	assumeAttrName  string
	inlineBlockName string
	policyName      string
}

// NewAwsIamRoleLambdaNoStarRule returns new rule with default attributes
func NewAwsIamRoleLambdaNoStarRule() *AwsIamRoleLambdaNoStarRule {
	return &AwsIamRoleLambdaNoStarRule{
		// TODO: Write resource type and attribute name here
		resourceType: "aws_iam_role",
		principalNames: []string{
			"lambda.amazonaws.com",
			"lambda.amazonaws.com.cn",
		},
		assumeAttrName:  "assume_role_policy",
		inlineBlockName: "inline_policy",
		policyName:      "policy",
	}
}

// Name returns the rule name
func (r *AwsIamRoleLambdaNoStarRule) Name() string {
	return "aws_iam_role_lambda_no_star"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIamRoleLambdaNoStarRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsIamRoleLambdaNoStarRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsIamRoleLambdaNoStarRule) Link() string {
	return "https://awslabs.github.io/serverless-rules/rules/lambda/star_permissions.html"
}

// matchPrincipal returns true if the policy has a matching Principal
func (r *AwsIamRoleLambdaNoStarRule) matchPrincipal(runner tflint.Runner, policy *hclext.Attribute) (bool, error) {
	var assumeAttrValue string
	err := runner.EvaluateExpr(policy.Expr, &assumeAttrValue, nil)
	if err != nil {
		return false, err
	}

	assumeRolePolicy := awsIamAssumeRole{}
	err = json.Unmarshal([]byte(assumeAttrValue), &assumeRolePolicy)
	if err != nil {
		return false, err
	}

	for _, principalName := range r.principalNames {
		for _, statement := range assumeRolePolicy.Statement {
			if principalService, ok := statement.Principal["Service"]; ok {
				switch principalService := principalService.(type) {
				case string:
					if principalService == principalName {
						return true, nil
					}
				case []string:
					for i := range principalService {
						if principalService[i] == principalName {
							return true, nil
						}
					}
				}
			}

		}
	}

	return false, nil
}

// matchStarAction returns true if the policy has a broad action in one of its statement
func (r *AwsIamRoleLambdaNoStarRule) matchStarAction(runner tflint.Runner, policy *hclext.Attribute) (bool, error) {
	var policyAttrValue string
	err := runner.EvaluateExpr(policy.Expr, &policyAttrValue, nil)
	if err != nil {
		return false, err
	}

	rolePolicy := awsIamAssumeRole{}
	err = json.Unmarshal([]byte(policyAttrValue), &rolePolicy)
	if err != nil {
		return false, err
	}

	for _, statement := range rolePolicy.Statement {
		switch action := reflect.ValueOf(statement.Action); action.Kind() {
		case reflect.String:
			if action.String() == "*" || strings.Contains(action.String(), ":*") {
				return true, nil
			}
		case reflect.Slice:
			for i := 0; i < action.Len(); i++ {
				v := action.Index(i).Interface().(string)
				if v == "*" || strings.Contains(v, ":*") {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// Check checks if an IAM role with a Lambda principal has broad permissions
func (r *AwsIamRoleLambdaNoStarRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.assumeAttrName},
		},
		Blocks: []hclext.BlockSchema{
			{
				Type: r.inlineBlockName,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: r.policyName},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		// Load assume role policy
		assumeAttr, ok := resource.Body.Attributes[r.assumeAttrName]
		if !ok {
			// This is a mandatory attribute
			runner.EmitIssue(
				r,
				fmt.Sprintf("\"%s\" is not present.", r.assumeAttrName),
				resource.DefRange,
			)
			continue
		}

		// Check if it contains the right principal
		hasLambda, err := r.matchPrincipal(runner, assumeAttr)
		if err != nil {
			return err
		}
		if !hasLambda {
			continue
		}

		// Load inline policy
		inlineBlocks := resource.Body.Blocks.OfType(r.inlineBlockName)
		for _, inlineBlock := range inlineBlocks {
			policyAttr, ok := inlineBlock.Body.Attributes[r.policyName]
			if !ok {
				// This is a mandatory attribute
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is not present.", r.policyName),
					inlineBlock.DefRange,
				)
				continue
			}

			// Check if policy contains stars
			hasStar, err := r.matchStarAction(runner, policyAttr)
			if err != nil {
				return err
			}

			if hasStar {
				runner.EmitIssue(
					r,
					"Inline policy for role with Lambda as principal has policy actions with stars.",
					policyAttr.Expr.Range(),
				)
			}
		}
	}

	return nil
}
