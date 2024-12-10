package waf_test

import (
	"testing"

	"github.com/prashantv/gostub"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/Azure/tflint-ruleset-avm/waf"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestAzurermKubernetesClusterZones(t *testing.T) {
	wafRules := waf.WafRules{}

	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "correct setting",
			rule: wafRules.AzurermKubernetesClusterZones(),
			content: `
	resource "azurerm_kubernetes_cluster" "example" {
		zones = [1, 2, 3]
	}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect setting",
			rule: wafRules.AzurermKubernetesClusterZones(),
			content: `
	resource "azurerm_kubernetes_cluster" "example" {
		zones = [1, 2]
	}`,
			expected: helper.Issues{
				{
					Rule:    wafRules.AzurermKubernetesClusterZones(),
					Message: "\"[1 2]\" is an invalid attribute value of `zones` - expecting (one of) [[1 2 3]]",
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
