package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/i56/framework/core/router"
	weightDomain "github.com/i56/modules/weight/domain"
)

func registerWeightAPI(r *router.Router, repo *weightDomain.MemWeightRepo) {
	// POST /api/v1/weight-records — create
	r.POST("/api/v1/weight-records", func(w http.ResponseWriter, req *http.Request) {
		var input weightDomain.CreateWeightRecordRequest
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			http.Error(w, `{"error":"invalid json"}`, 400)
			return
		}
		if input.TrackingNumber == "" {
			http.Error(w, `{"error":"tracking_number required"}`, 400)
			return
		}
		if input.Weight <= 0 {
			input.Weight = 0.1
		}
		if input.Length <= 0 {
			input.Length = 1
		}
		if input.Width <= 0 {
			input.Width = 1
		}
		if input.Height <= 0 {
			input.Height = 1
		}
		if input.ParcelCount <= 0 {
			input.ParcelCount = 1
		}
		if input.Platform == "" {
			input.Platform = "EZ集运通"
		}

		rec := repo.Create(1, input) // tenantID=1 for now
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": rec})
	})

	// GET /api/v1/weight-records — list with pagination
	r.GET("/api/v1/weight-records", func(w http.ResponseWriter, req *http.Request) {
		offset, _ := strconv.Atoi(req.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
		if limit <= 0 {
			limit = 50
		}

		records, total := repo.List(1, offset, limit)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data":  records,
			"total": total,
		})
	})

	// GET /api/v1/weight-records/search?q= — search
	r.GET("/api/v1/weight-records/search", func(w http.ResponseWriter, req *http.Request) {
		q := req.URL.Query().Get("q")
		records := repo.Search(1, q, 100)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data":  records,
			"total": len(records),
		})
	})

	// GET /api/v1/weight-records/{id} — get single record (handled via query param)
	r.GET("/api/v1/weight-records/:id", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": nil, "error": "use ?id= param"})
	})
}
