package repository

import "sync"

// ClientRoutePriceDisplay represents route pricing data for client portal display.
// Fields match the client_route_prices.html template.
type ClientRoutePriceDisplay struct {
	RouteName            string  `json:"route_name"`
	TransportType        string  `json:"transport_type"`         // sea / sea_express / air
	CargoType            string  `json:"cargo_type"`             // general / class1 etc.
	TaxType              string  `json:"tax_type"`               // full_inclusive / tax_excluded
	FirstWeight          string  `json:"first_weight"`           // e.g. "10"
	FirstWeightPrice     string  `json:"first_weight_price"`     // e.g. "32.00"
	AdditionalWeightPrice string `json:"additional_weight_price"` // e.g. "3.20"
	FirstVolume          string  `json:"first_volume"`           // e.g. "1"
	FirstVolumePrice     string  `json:"first_volume_price"`     // e.g. "20.00"
	MinCharge            string  `json:"min_charge"`             // e.g. "50.00"
}

// MemRoutePriceRepo is an in-memory seed repo for client route pricing display.
type MemRoutePriceRepo struct {
	mu     sync.RWMutex
	prices []ClientRoutePriceDisplay
}

func NewMemRoutePriceRepo() *MemRoutePriceRepo {
	return &MemRoutePriceRepo{
		prices: []ClientRoutePriceDisplay{
			{
				RouteName:             "廈門→台灣(海運)",
				TransportType:         "sea",
				CargoType:             "class1",
				TaxType:               "full_inclusive",
				FirstWeight:           "10",
				FirstWeightPrice:      "32.00",
				AdditionalWeightPrice: "3.20",
				FirstVolume:           "1",
				FirstVolumePrice:      "20.00",
				MinCharge:             "50.00",
			},
			{
				RouteName:             "廈門→台灣(海快)",
				TransportType:         "sea_express",
				CargoType:             "general",
				TaxType:               "full_inclusive",
				FirstWeight:           "1",
				FirstWeightPrice:      "15.00",
				AdditionalWeightPrice: "8.00",
				FirstVolume:           "1",
				FirstVolumePrice:      "15.00",
				MinCharge:             "25.00",
			},
			{
				RouteName:             "深圳→台灣(空運普貨)",
				TransportType:         "air",
				CargoType:             "general",
				TaxType:               "tax_excluded",
				FirstWeight:           "0.5",
				FirstWeightPrice:      "25.00",
				AdditionalWeightPrice: "12.00",
				FirstVolume:           "0.5",
				FirstVolumePrice:      "25.00",
				MinCharge:             "50.00",
			},
			{
				RouteName:             "廈門→台灣(商業海快)",
				TransportType:         "sea_express",
				CargoType:             "class3",
				TaxType:               "full_inclusive",
				FirstWeight:           "5",
				FirstWeightPrice:      "20.00",
				AdditionalWeightPrice: "5.00",
				FirstVolume:           "1",
				FirstVolumePrice:      "18.00",
				MinCharge:             "40.00",
			},
		},
	}
}

// List returns all route prices (client portal reads all).
func (r *MemRoutePriceRepo) List() []ClientRoutePriceDisplay {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]ClientRoutePriceDisplay, len(r.prices))
	copy(result, r.prices)
	return result
}

// Add appends a new route price.
func (r *MemRoutePriceRepo) Add(p ClientRoutePriceDisplay) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prices = append(r.prices, p)
}

// Remove deletes a route price by index.
func (r *MemRoutePriceRepo) Remove(idx int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if idx >= 0 && idx < len(r.prices) {
		r.prices = append(r.prices[:idx], r.prices[idx+1:]...)
	}
}
