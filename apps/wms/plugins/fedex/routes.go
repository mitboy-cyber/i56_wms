package fedex

import (
	"encoding/json"
	"net/http"
)

func (p *FedExPlugin) handleTrack(w http.ResponseWriter, r *http.Request) {
	tracking := r.PathValue("tracking")
	resp := map[string]interface{}{
		"plugin":   "fedex",
		"tracking": tracking,
		"status":   "in_transit",
		"message":  "FedEx tracking lookup — would call FedEx Track API",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (p *FedExPlugin) handleCreateShipment(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"plugin":  "fedex",
		"status":  "ok",
		"message": "FedEx shipment would be created via FedEx Ship API",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
