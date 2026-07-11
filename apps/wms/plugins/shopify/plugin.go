// Package shopify provides Shopify order import integration for the I56 WMS.
package shopify

import (
	"fmt"
	"net/http"

	"github.com/i56/framework/plugin"
)

// ShopifyPlugin imports orders from Shopify into the WMS.
type ShopifyPlugin struct {
	config map[string]interface{}
}

func (p *ShopifyPlugin) Name() string    { return "shopify" }
func (p *ShopifyPlugin) Version() string { return "1.0.0" }

func (p *ShopifyPlugin) Init(config map[string]interface{}) error {
	p.config = config
	fmt.Println("[shopify] initialized with store:", config["store"])
	return nil
}

func (p *ShopifyPlugin) Routes() []plugin.Route {
	return []plugin.Route{
		{Method: http.MethodGet, Path: "/api/plugins/shopify/orders", Handler: p.handleListOrders},
		{Method: http.MethodPost, Path: "/api/plugins/shopify/import", Handler: p.handleImportOrders},
	}
}

func (p *ShopifyPlugin) MenuItems() []plugin.MenuItem {
	return []plugin.MenuItem{
		{Group: "Integrations", Label: "Shopify Orders", URL: "/plugins/shopify", Icon: "shopping-cart"},
	}
}

func (p *ShopifyPlugin) Shutdown() error {
	fmt.Println("[shopify] shutdown complete")
	return nil
}
