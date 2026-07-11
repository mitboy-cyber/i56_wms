package stripe

import (
	"encoding/json"
	"net/http"
)

func (p *StripePlugin) handleCharge(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"plugin":  "stripe",
		"status":  "ok",
		"charge":  "pi_mock_12345",
		"message": "Stripe payment charge would be processed via Stripe Charges API",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (p *StripePlugin) handleListTransactions(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"plugin":       "stripe",
		"transactions": []string{},
		"message":      "Stripe transactions would be fetched via Stripe API",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
