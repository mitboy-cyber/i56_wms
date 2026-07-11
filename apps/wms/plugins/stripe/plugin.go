// Package stripe provides Stripe payment integration for the I56 WMS.
package stripe

import (
	"fmt"
	"net/http"

	"github.com/i56/framework/plugin"
)

// StripePlugin handles payment processing via Stripe.
type StripePlugin struct {
	config map[string]interface{}
}

func (p *StripePlugin) Name() string    { return "stripe" }
func (p *StripePlugin) Version() string { return "1.0.0" }

func (p *StripePlugin) Init(config map[string]interface{}) error {
	p.config = config
	fmt.Println("[stripe] initialized with key prefix:", p.maskKey(config["api_key"]))
	return nil
}

func (p *StripePlugin) Routes() []plugin.Route {
	return []plugin.Route{
		{Method: http.MethodPost, Path: "/api/plugins/stripe/charge", Handler: p.handleCharge},
		{Method: http.MethodGet, Path: "/api/plugins/stripe/transactions", Handler: p.handleListTransactions},
	}
}

func (p *StripePlugin) MenuItems() []plugin.MenuItem {
	return []plugin.MenuItem{
		{Group: "Integrations", Label: "Stripe Payments", URL: "/plugins/stripe", Icon: "credit-card"},
	}
}

func (p *StripePlugin) Shutdown() error {
	fmt.Println("[stripe] shutdown complete")
	return nil
}

func (p *StripePlugin) maskKey(v interface{}) string {
	s := fmt.Sprint(v)
	if len(s) > 8 {
		return s[:4] + "..." + s[len(s)-4:]
	}
	return "***"
}
