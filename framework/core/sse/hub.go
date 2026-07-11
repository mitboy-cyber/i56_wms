package sse

import (
	"fmt"
	"log"
	"sync"
)

// Event represents a server-sent event
type Event struct {
	Channel string
	Type    string // "parcel_received", "parcel_stored", "order_created", etc.
	Data    string // JSON payload
}

// Client represents a single SSE connection
type Client struct {
	ID       string
	Channel  string
	Events   chan Event
	Done     chan struct{}
}

// Hub manages SSE connections and event broadcasting
type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client // key: channel name
	// Multichannel: one client per channel, broadcast to all
	subscribers map[string]map[string]*Client // channel → clientID → client
	nextID      int
}

// NewHub creates a new SSE hub
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[string]*Client),
		subscribers: make(map[string]map[string]*Client),
	}
}

// Subscribe creates a new client for a channel
func (h *Hub) Subscribe(channel string) *Client {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.nextID++
	c := &Client{
		ID:      fmt.Sprintf("sse-%d", h.nextID),
		Channel: channel,
		Events:  make(chan Event, 32),
		Done:    make(chan struct{}),
	}
	
	if h.subscribers[channel] == nil {
		h.subscribers[channel] = make(map[string]*Client)
	}
	h.subscribers[channel][c.ID] = c
	log.Printf("[SSE] Client %s subscribed to channel %s (total: %d)", c.ID, channel, len(h.subscribers[channel]))
	return c
}

// Unsubscribe removes a client
func (h *Hub) Unsubscribe(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if ch, ok := h.subscribers[c.Channel]; ok {
		delete(ch, c.ID)
		if len(ch) == 0 {
			delete(h.subscribers, c.Channel)
		}
	}
	close(c.Events)
	log.Printf("[SSE] Client %s unsubscribed from channel %s", c.ID, c.Channel)
}

// Publish sends an event to all subscribers of a channel
func (h *Hub) Publish(channel string, ev Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	ch, ok := h.subscribers[channel]
	if !ok {
		return
	}
	
	for _, c := range ch {
		select {
		case c.Events <- ev:
		default:
			// Client buffer full, skip (client too slow)
			log.Printf("[SSE] Client %s buffer full, dropping event", c.ID)
		}
	}
}

// BroadCast sends an event to ALL channels
func (h *Hub) BroadCast(ev Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	for _, ch := range h.subscribers {
		for _, c := range ch {
			select {
			case c.Events <- ev:
			default:
			}
		}
	}
}

// ActiveChannels returns list of channels with active subscribers
func (h *Hub) ActiveChannels() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	channels := make([]string, 0, len(h.subscribers))
	for ch := range h.subscribers {
		channels = append(channels, ch)
	}
	return channels
}

// SubscriberCount returns the number of active subscribers for a channel
func (h *Hub) SubscriberCount(channel string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	if ch, ok := h.subscribers[channel]; ok {
		return len(ch)
	}
	return 0
}
