// Package ws provides WebSocket support for PDA real-time sessions.
package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// PDASession represents a connected PDA device session.
type PDASession struct {
	ID        string    `json:"id"`
	Warehouse string    `json:"warehouse"`
	Operator  string    `json:"operator"`
	Device    string    `json:"device"`
	Status    string    `json:"status"` // online, scanning, idle, error
	Page      string    `json:"page"`
	Area      string    `json:"area"`
	Location  string    `json:"location"`
	LoginAt   time.Time `json:"login_at"`
	LastBeat  time.Time `json:"last_beat"`
	conn      *websocket.Conn
	mu        sync.Mutex
}

// Hub manages all PDA WebSocket connections.
type Hub struct {
	mu       sync.RWMutex
	sessions map[string]*PDASession
}

// NewHub creates a PDA WebSocket hub.
func NewHub() *Hub {
	return &Hub{sessions: make(map[string]*PDASession)}
}

// HandleUpgrade upgrades HTTP to WebSocket for PDA clients.
func (h *Hub) HandleUpgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ws] upgrade error: %v", err)
		return
	}

	sessionID := fmt.Sprintf("pda-%d", time.Now().UnixNano())
	session := &PDASession{
		ID:       sessionID,
		Status:   "online",
		LoginAt:  time.Now(),
		LastBeat: time.Now(),
		conn:     conn,
	}
	h.mu.Lock()
	h.sessions[sessionID] = session
	h.mu.Unlock()

	log.Printf("[ws] PDA connected: %s", sessionID)
	h.broadcastState()
	go h.readPump(session)
}

func (h *Hub) readPump(s *PDASession) {
	defer func() {
		s.conn.Close()
		h.mu.Lock()
		delete(h.sessions, s.ID)
		h.mu.Unlock()
		log.Printf("[ws] PDA disconnected: %s", s.ID)
		h.broadcastState()
	}()

	s.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	s.conn.SetPongHandler(func(string) error {
		s.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		return nil
	})

	for {
		_, msg, err := s.conn.ReadMessage()
		if err != nil {
			break
		}
		var update map[string]any
		if err := json.Unmarshal(msg, &update); err != nil {
			continue
		}
		s.mu.Lock()
		if w, ok := update["warehouse"].(string); ok {
			s.Warehouse = w
		}
		if o, ok := update["operator"].(string); ok {
			s.Operator = o
		}
		if d, ok := update["device"].(string); ok {
			s.Device = d
		}
		if st, ok := update["status"].(string); ok {
			s.Status = st
		}
		if p, ok := update["page"].(string); ok {
			s.Page = p
		}
		if a, ok := update["area"].(string); ok {
			s.Area = a
		}
		if l, ok := update["location"].(string); ok {
			s.Location = l
		}
		s.LastBeat = time.Now()
		s.mu.Unlock()

		h.broadcastState()
	}
}

func (h *Hub) broadcastState() {
	h.mu.RLock()
	sessions := make([]*PDASession, 0, len(h.sessions))
	for _, s := range h.sessions {
		sessions = append(sessions, s)
	}
	h.mu.RUnlock()

	type sessionJSON struct {
		ID        string `json:"id"`
		Warehouse string `json:"warehouse"`
		Operator  string `json:"operator"`
		Device    string `json:"device"`
		Status    string `json:"status"`
		Page      string `json:"page"`
		Area      string `json:"area"`
		Location  string `json:"location"`
		LoginAt   string `json:"login_at"`
		LastBeat  string `json:"last_beat"`
	}

	out := make([]sessionJSON, 0, len(sessions))
	for _, s := range sessions {
		s.mu.Lock()
		out = append(out, sessionJSON{
			ID:        s.ID,
			Warehouse: s.Warehouse,
			Operator:  s.Operator,
			Device:    s.Device,
			Status:    s.Status,
			Page:      s.Page,
			Area:      s.Area,
			Location:  s.Location,
			LoginAt:   s.LoginAt.Format(time.RFC3339),
			LastBeat:  s.LastBeat.Format(time.RFC3339),
		})
		s.mu.Unlock()
	}

	msg, _ := json.Marshal(map[string]any{
		"type":     "pda_sessions",
		"sessions": out,
		"count":    len(out),
	})

	// Broadcast to all admin watchers
	h.mu.RLock()
	for _, s := range h.sessions {
		s.mu.Lock()
		if s.conn != nil {
			s.conn.WriteMessage(websocket.TextMessage, msg)
		}
		s.mu.Unlock()
	}
	h.mu.RUnlock()
}

// HandleAdminWatch serves real-time PDA session data to admin pages via SSE.
func (h *Hub) HandleAdminWatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", 500)
		return
	}

	ctx := r.Context()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	getSessions := func() []any {
		h.mu.RLock()
		defer h.mu.RUnlock()
		out := make([]any, 0, len(h.sessions))
		for _, s := range h.sessions {
			s.mu.Lock()
			out = append(out, map[string]any{
				"id": s.ID, "warehouse": s.Warehouse, "operator": s.Operator,
				"device": s.Device, "status": s.Status, "page": s.Page,
				"area": s.Area, "location": s.Location,
				"login_at": s.LoginAt.Format(time.RFC3339),
				"last_beat": s.LastBeat.Format(time.RFC3339),
			})
			s.mu.Unlock()
		}
		return out
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sessions := getSessions()
			data, _ := json.Marshal(sessions)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
