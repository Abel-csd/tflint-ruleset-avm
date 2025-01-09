package waf_test

import (
	"testing"

	"github.com/prashantv/gostub"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/Azure/tflint-ruleset-avm/waf"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestAzurermStorageAccountAccountReplicationType(t *testing.T) {
	wafRules := waf.WafRules{}

	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "correct setting",
			rule: wafRules.AzurermStorageAccountAccountReplicationType(),
			content: `
	variable "account_replication_type" {
		type    = list(string)
		default = ["GRS","ZRS"]
	}
	resource "azurerm_storage_account" "example" {
		for_each = toset(var.account_replication_type)
		account_replication_type = each.value
	}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect setting",
			rule: wafRules.AzurermStorageAccountAccountReplicationType(),
			content: `
	resource "azurerm_storage_account" "example" {
		account_replication_type = "LRS"
	}`,
			expected: helper.Issues{
				{
					Rule:    wafRules.AzurermStorageAccountAccountReplicationType(),
					Message: "LRS is an invalid attribute value of `account_replication_type` - expecting (one of) [GRS ZRS]",
				},
			},
		},
		{
			name: "null value",
			rule: wafRules.AzurermStorageAccountAccountReplicationType(),
			content: `
	variable "account_replication_type" {
		type    = string
		default = null
	}
	resource "azurerm_storage_account" "example" {
		account_replication_type = var.account_replication_type
	}`,
			expected: helper.Issues{},
		},
		{
			name: "missing attribute",
			rule: wafRules.AzurermStorageAccountAccountReplicationType(),
			content: `
	resource "azurerm_storage_account" "example" {
	}`,
			expected: helper.Issues{},
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
