package waf_test

import (
	"testing"

	"github.com/prashantv/gostub"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/Azure/tflint-ruleset-avm/waf"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestAzurermLbSku(t *testing.T) {
	wafRules := waf.WafRules{}

	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "correct setting",
			rule: wafRules.AzurermLbSku(),
			content: `
	resource "azurerm_lb" "example" {
		sku = "Standard"
	}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect setting",
			rule: wafRules.AzurermLbSku(),
			content: `
	resource "azurerm_lb" "example" {
		sku = "Basic"
	}`,
			expected: helper.Issues{
				{
					Rule:    wafRules.AzurermLbSku(),
					Message: "Basic is an invalid attribute value of `sku` - expecting (one of) [Standard]",
				},
			},
		},
	}

	filename := "main.tf"
	for _, c := range testCases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{filename: tc.content})
			stub := gostub.Stub(&attrvalue.AppFs, mockFs(tc.content))
			defer stub.Reset()
			if err := tc.rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssuesWithoutRange(t, tc.expected, runner.Issues)
		})
	}
}
