package main

import (
	"encoding/json"
	"net/http"
)

// apiJSON writes a JSON response with the given status code.
func apiJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
