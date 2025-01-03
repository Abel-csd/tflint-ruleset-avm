package waf_test

import (
	"testing"

	"github.com/prashantv/gostub"

	"github.com/Azure/tflint-ruleset-avm/attrvalue"
	"github.com/Azure/tflint-ruleset-avm/waf"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestAzurermVirtualNetworkGatewaySku(t *testing.T) {
	wafRules := waf.WafRules{}

	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "correct setting",
			rule: wafRules.AzurermVirtualNetworkGatewaySku(),
			content: `
	variable "sku_type" {
		type    = list(string)
		default = ["ErGw1AZ", "ErGw2AZ", "ErGw3AZ", "VpnGw1AZ", "VpnGw2AZ", "VpnGw3AZ", "VpnGw4AZ", "VpnGw5AZ"]
	}
	resource "azurerm_virtual_network_gateway" "example" {
		for_each = toset(var.sku_type)
		sku = each.value
	}`,
			expected: helper.Issues{},
		},
		{
			name: "incorrect setting",
			rule: wafRules.AzurermVirtualNetworkGatewaySku(),
			content: `
	resource "azurerm_virtual_network_gateway" "example" {
		sku = "ErGw4AZ"
	}`,
			expected: helper.Issues{
				{
					Rule:    wafRules.AzurermVirtualNetworkGatewaySku(),
					Message: "ErGw4AZ is an invalid attribute value of `sku` - expecting (one of) [ErGw1AZ ErGw2AZ ErGw3AZ VpnGw1AZ VpnGw2AZ VpnGw3AZ VpnGw4AZ VpnGw5AZ]",
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

func TestAzurermVirtualNetworkGatewayVpnActiveActive(t *testing.T) {
	wafRules := waf.WafRules{}

	testCases := []struct {
		name     string
		rule     tflint.Rule
		content  string
		expected helper.Issues
	}{
		{
			name: "correct setting",
			rule: wafRules.AzurermVirtualNetworkGatewayVpnActiveActive(),
			content: `
	resource "azurerm_virtual_network_gateway" "example" {
		active_active = true
	}`,

			expected: helper.Issues{},
		},
		{
			name: "incorrect setting",
			rule: wafRules.AzurermVirtualNetworkGatewayVpnActiveActive(),
			content: `
	resource "azurerm_virtual_network_gateway" "example" {
		active_active = false
	}`,
			expected: helper.Issues{
				{
					Rule:    wafRules.AzurermVirtualNetworkGatewayVpnActiveActive(),
					Message: "false is an invalid attribute value of `active_active` - expecting (one of) [true]",
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
