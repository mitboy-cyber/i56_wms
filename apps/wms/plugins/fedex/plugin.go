// Package fedex provides FedEx tracking integration for the I56 WMS.
package fedex

import (
	"fmt"
	"net/http"

	"github.com/i56/framework/plugin"
)

// FedExPlugin provides tracking integration with FedEx APIs.
type FedExPlugin struct {
	config map[string]interface{}
}

func (p *FedExPlugin) Name() string    { return "fedex" }
func (p *FedExPlugin) Version() string { return "1.0.0" }

func (p *FedExPlugin) Init(config map[string]interface{}) error {
	p.config = config
	fmt.Println("[fedex] initialized with account:", config["account_number"])
	return nil
}

func (p *FedExPlugin) Routes() []plugin.Route {
	return []plugin.Route{
		{Method: http.MethodGet, Path: "/api/plugins/fedex/track/:tracking", Handler: p.handleTrack},
		{Method: http.MethodPost, Path: "/api/plugins/fedex/ship", Handler: p.handleCreateShipment},
	}
}

func (p *FedExPlugin) MenuItems() []plugin.MenuItem {
	return []plugin.MenuItem{
		{Group: "Integrations", Label: "FedEx Tracking", URL: "/plugins/fedex", Icon: "truck"},
	}
}

func (p *FedExPlugin) Shutdown() error {
	fmt.Println("[fedex] shutdown complete")
	return nil
}
