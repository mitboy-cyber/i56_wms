package gateway

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
)

// SSEHub manages Server-Sent Events for AI streaming and real-time updates
type SSEHub struct {
	mu      sync.RWMutex
	clients map[string]chan string // channel → clients
}

// SSEClientInfo holds client metadata
type SSEClientInfo struct {
	Channel string `json:"channel"`
}

// NewSSEHub creates a new SSE hub
func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients: make(map[string]chan string),
	}
}

// Subscribe adds a client channel for SSE
func (h *SSEHub) Subscribe(channel string) chan string {
	h.mu.Lock()
	defer h.mu.Unlock()
	ch := make(chan string, 32)
	h.clients[channel] = ch
	return ch
}

// Unsubscribe removes a client channel
func (h *SSEHub) Unsubscribe(channel string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if ch, ok := h.clients[channel]; ok {
		close(ch)
		delete(h.clients, channel)
	}
}

// Broadcast sends an event to a specific SSE channel
func (h *SSEHub) Broadcast(channel, event, data string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if ch, ok := h.clients[channel]; ok {
		msg := fmt.Sprintf("event: %s\ndata: %s\n\n", event, data)
		select {
		case ch <- msg:
		default:
		}
	}
}

// HandleSSE serves SSE connections
func (gw *Gateway) sseHandler(c *gin.Context) {
	channel := c.Query("channel")
	if channel == "" {
		channel = "default"
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	ch := gw.sseHub.Subscribe(channel)

	// Send initial connection event
	fmt.Fprintf(c.Writer, "event: connected\ndata: {\"channel\":\"%s\"}\n\n", channel)
	c.Writer.Flush()

	ctx := c.Request.Context()
	for {
		select {
		case <-ctx.Done():
			gw.sseHub.Unsubscribe(channel)
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			io.WriteString(c.Writer, msg)
			c.Writer.Flush()
		}
	}
}

// Ensure embedding works with Gin
func init() {
	var _ = log.Default
}
