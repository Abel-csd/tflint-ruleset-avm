package waf_test

import (
	"testing"

	"github.com/prashantv/gostub"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/Azure/tflint-ruleset-avm/waf"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestAzurermApplicationGatewayZones(t *testing.T) {
	wafRules := waf.WafRules{}

	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "correct setting",
			rule: wafRules.AzurermApplicationGatewayZones(),
			content: `
	variable "zones" {
		type    = list(number)
		default = [1, 2, 3]
	}
	resource "azurerm_application_gateway" "example" {
		zones = var.zones
	}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect setting",
			rule: wafRules.AzurermApplicationGatewayZones(),
			content: `
	variable "zones" {
		type    = list(number)
		default = [2, 3]
	}
	resource "azurerm_application_gateway" "example" {
		zones = var.zones
	}`,
			expected: helper.Issues{
				{
					Rule:    wafRules.AzurermApplicationGatewayZones(),
					Message: "\"[2 3]\" is an invalid attribute value of `zones` - expecting (one of) [[1 2 3]]",
				},
			},
		},
		{
			name: "correct with string list",
			rule: wafRules.AzurermApplicationGatewayZones(),
			content: `
	variable "zones" {
		type    = list(string)
		default = ["1", "2", "3"]
	}
	resource "azurerm_application_gateway" "example" {
		zones = var.zones
	}`,
			expected: helper.Issues{},
		},
		{
			name: "correct but different order",
			rule: wafRules.AzurermApplicationGatewayZones(),
			content: `
	variable "zones" {
		type    = list(number)
		default = [2, 3, 1]
	}
	resource "azurerm_application_gateway" "example" {
		zones = var.zones
	}`,
			expected: helper.Issues{},
		},
		{
			name: "variable without default",
			rule: wafRules.AzurermApplicationGatewayZones(),
			content: `
		variable "zones" {
			type    = list(number)
		}
		resource "azurerm_application_gateway" "example" {
			zones = var.zones
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
func TestAzurermApplicationGatewaySku(t *testing.T) {
	wafRules := waf.WafRules{}

	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "correct setting",
			rule: wafRules.AzurermApplicationGatewaySku(),
			content: `
	variable "sku" {
		type = list(string)
		default = ["Standard_v2", "WAF_v2"]
	}
	resource "azurerm_application_gateway" "example" {
		for_each = toset(var.sku)
		sku {
			name = each.value
		}
	}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect setting",
			rule: wafRules.AzurermApplicationGatewaySku(),
			content: `
	variable "sku" {
		type = string
		default = "Standard_v3"
	}
	resource "azurerm_application_gateway" "example" {
		sku {
			name = var.sku
		}
	}`,
			expected: helper.Issues{
				{
					Rule:    wafRules.AzurermApplicationGatewaySku(),
					Message: "Standard_v3 is an invalid attribute value of `name` - expecting (one of) [Standard_v2 WAF_v2]",
				},
			},
		},
		{
			name: "null value",
			rule: wafRules.AzurermApplicationGatewaySku(),
			content: `
	variable "sku" {
		type    = string
		default = null
	}
	resource "azurerm_application_gateway" "example" {
		sku = var.sku
	}`,
			expected: helper.Issues{},
		},
		{
			name: "missing attribute",
			rule: wafRules.AzurermApplicationGatewaySku(),
			content: `
	resource "azurerm_application_gateway" "example" {

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
