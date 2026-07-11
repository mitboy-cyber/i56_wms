// Package webhook provides a Watermill event subscriber that forwards
// domain events to registered webhook URLs.
package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/i56/framework/events"
	whDomain "github.com/i56/modules/webhook/domain"
	whRepo "github.com/i56/modules/webhook/repository"
)

// Dispatcher listens on the Watermill event bus and forwards matching
// events to registered webhook subscription URLs.
type Dispatcher struct {
	repo   *whRepo.MemWebhookRepo
	client *http.Client
}

// NewDispatcher creates a new webhook event dispatcher.
func NewDispatcher(repo *whRepo.MemWebhookRepo) *Dispatcher {
	return &Dispatcher{
		repo: repo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Start subscribes to all supported domain events on the Watermill bus
// and forwards them to matching webhook URLs.
func (d *Dispatcher) Start() {
	topics := []string{
		"order.created",
		"parcel.created",
		"parcel.received",
		"status.changed",
		"payment.received",
		"container.closed",
	}

	for _, topic := range topics {
		topic := topic // capture
		events.SubscribeHandler(topic, func(ctx context.Context, payload []byte) error {
			d.dispatch(topic, payload)
			return nil
		})
	}

	log.Println("[webhook-dispatcher] Started — listening for domain events")
}

// dispatch finds matching webhook subscriptions and POSTs the event payload.
func (d *Dispatcher) dispatch(eventType string, payload []byte) {
	subs := d.repo.FindByEvent(context.Background(), eventType)
	if len(subs) == 0 {
		return
	}

	for _, sub := range subs {
		go d.deliver(sub, eventType, payload)
	}
}

// deliver sends an event payload to a single webhook subscription URL.
func (d *Dispatcher) deliver(sub whDomain.WebhookSubscription, eventType string, payload []byte) {
	req, err := http.NewRequest("POST", sub.URL, bytes.NewReader(payload))
	if err != nil {
		log.Printf("[webhook-dispatcher] Failed to create request for %s: %v", sub.URL, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-I56-Event", eventType)

	resp, err := d.client.Do(req)
	statusCode := 0
	errStr := ""
	if err != nil {
		errStr = err.Error()
	} else {
		statusCode = resp.StatusCode
		resp.Body.Close()
	}

	// Log delivery
	logEntry := &whDomain.WebhookDeliveryLog{
		SubscriptionID: sub.ID,
		Event:          eventType,
		Payload:        compactJSON(payload),
		StatusCode:     statusCode,
		Error:          errStr,
	}
	d.repo.LogDelivery(context.Background(), logEntry)

	if err != nil {
		log.Printf("[webhook-dispatcher] Delivery failed to %s: %v", sub.URL, err)
	} else {
		log.Printf("[webhook-dispatcher] Delivered %s → %s (HTTP %d)", eventType, sub.URL, statusCode)
	}
}

// compactJSON trims whitespace from a JSON payload for logging.
func compactJSON(data []byte) string {
	var buf bytes.Buffer
	if err := json.Compact(&buf, data); err != nil {
		return string(data)
	}
	return buf.String()
}
