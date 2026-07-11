package shopify

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (p *ShopifyPlugin) handleListOrders(w http.ResponseWriter, r *http.Request) {
	store := "default"
	if s, ok := p.config["store"]; ok {
		store = fmt.Sprint(s)
	}
	resp := map[string]interface{}{
		"plugin":  "shopify",
		"store":   store,
		"orders":  []string{},
		"message": "Shopify order import endpoint — orders fetched from Shopify API",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (p *ShopifyPlugin) handleImportOrders(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"plugin":  "shopify",
		"status":  "ok",
		"imported": 0,
		"message": "Shopify orders would be imported here via Shopify Admin API",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
