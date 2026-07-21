// Package pda provides PDA real-time session management via SSE.
package pda

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Session represents a connected PDA device.
type Session struct {
	ID        string    `json:"id"`
	Warehouse string    `json:"warehouse"`
	Operator  string    `json:"operator"`
	Device    string    `json:"device"`
	Status    string    `json:"status"` // online, scanning, idle
	Page      string    `json:"page"`
	Area      string    `json:"area"`
	Location  string    `json:"location"`
	LoginAt   time.Time `json:"login_at"`
	LastBeat  time.Time `json:"last_beat"`
}

// Hub manages all PDA sessions and broadcasts state via SSE.
type Hub struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// NewHub creates a PDA session hub.
func NewHub() *Hub {
	h := &Hub{sessions: make(map[string]*Session)}

	// Seed demo sessions
	now := time.Now()
	h.sessions["pda-01"] = &Session{ID: "pda-01", Warehouse: "厦门仓", Operator: "张三", Device: "TC26", Status: "scanning", Page: "入库上架", Area: "A区-01", Location: "A01-03-05", LoginAt: now.Add(-2 * time.Hour), LastBeat: now}
	h.sessions["pda-02"] = &Session{ID: "pda-02", Warehouse: "厦门仓", Operator: "李四", Device: "TC21", Status: "idle", Page: "拣货任务", Area: "B区-02", Location: "B02-01-01", LoginAt: now.Add(-1 * time.Hour), LastBeat: now}
	h.sessions["pda-03"] = &Session{ID: "pda-03", Warehouse: "台北仓", Operator: "王五", Device: "RFID", Status: "online", Page: "出库扫描", Area: "C区-01", Location: "C01-02-08", LoginAt: now.Add(-30 * time.Minute), LastBeat: now}

	return h
}

// UpdateSession updates a PDA session with new data.
func (h *Hub) UpdateSession(body map[string]any) {
	id, _ := body["id"].(string)
	if id == "" {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	s, ok := h.sessions[id]
	if !ok {
		s = &Session{ID: id, LoginAt: time.Now()}
		h.sessions[id] = s
	}
	if v, ok := body["warehouse"].(string); ok {
		s.Warehouse = v
	}
	if v, ok := body["operator"].(string); ok {
		s.Operator = v
	}
	if v, ok := body["device"].(string); ok {
		s.Device = v
	}
	if v, ok := body["status"].(string); ok {
		s.Status = v
	}
	if v, ok := body["page"].(string); ok {
		s.Page = v
	}
	if v, ok := body["area"].(string); ok {
		s.Area = v
	}
	if v, ok := body["location"].(string); ok {
		s.Location = v
	}
	s.LastBeat = time.Now()
}

// ServeSSE streams PDA session state to admin clients via Server-Sent Events.
func (h *Hub) ServeSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", 500)
		return
	}

	ctx := r.Context()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	send := func() {
		h.mu.RLock()
		sessions := make([]map[string]any, 0, len(h.sessions))
		for _, s := range h.sessions {
			sessions = append(sessions, map[string]any{
				"id":        s.ID,
				"warehouse": s.Warehouse,
				"operator":  s.Operator,
				"device":    s.Device,
				"status":    s.Status,
				"page":      s.Page,
				"area":      s.Area,
				"location":  s.Location,
				"login_at":  s.LoginAt.Format(time.RFC3339),
				"last_beat": s.LastBeat.Format(time.RFC3339),
			})
		}
		h.mu.RUnlock()

		data, _ := json.Marshal(map[string]any{
			"sessions": sessions,
			"count":    len(sessions),
		})
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	// Send immediately, then on ticker
	send()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			send()
		}
	}
}

// HandleUpdate is the HTTP handler for PDA status updates.
func (h *Hub) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	h.UpdateSession(body)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GetSessions returns current session snapshot.
func (h *Hub) GetSessions() []map[string]any {
	h.mu.RLock()
	defer h.mu.RUnlock()
	sessions := make([]map[string]any, 0, len(h.sessions))
	for _, s := range h.sessions {
		sessions = append(sessions, map[string]any{
			"id": s.ID, "warehouse": s.Warehouse, "operator": s.Operator,
			"device": s.Device, "status": s.Status, "page": s.Page,
			"area": s.Area, "location": s.Location,
			"login_at": s.LoginAt.Format(time.RFC3339),
			"last_beat": s.LastBeat.Format(time.RFC3339),
		})
	}
	return sessions
}

// Logf prints a hub log line.
func Logf(format string, args ...any) {
	log.Printf("[pda-hub] "+format, args...)
}
