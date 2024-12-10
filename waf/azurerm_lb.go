package waf

import (
	"github.com/Azure/tflint-ruleset-avm/attrvalue"
)

func (wf WafRules) AzurermLbSku() *attrvalue.SimpleRule[string] {
	return attrvalue.NewSimpleRule[string](
		"azurerm_lb",
		"sku",
		[]string{"Standard"},
		"https://azure.github.io/Azure-Proactive-Resiliency-Library-v2/azure-resources/Network/loadBalancers/#use-standard-load-balancer-sku",
		false,
		"",
	)
}
