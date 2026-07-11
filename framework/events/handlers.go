package events

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
)

// ─── Publish Helpers ───────────────────────────────────────────────────────

// PublishOrderCreated publishes an order.created event.
func PublishOrderCreated(orderID int64, orderNo string, clientID int64, price float64) {
	if Bus == nil {
		log.Println("[events] Bus not initialized, skipping PublishOrderCreated")
		return
	}
	ev := OrderCreated{
		OrderID:    orderID,
		OrderNo:    orderNo,
		ClientID:   clientID,
		TotalPrice: price,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}
	publish("order.created", ev)
}

// PublishParcelReceived publishes a parcel.received event.
func PublishParcelReceived(parcelID int64, trackingNo string, warehouseID int64, productName string) {
	if Bus == nil {
		return
	}
	ev := ParcelReceived{
		ParcelID:    parcelID,
		TrackingNo:  trackingNo,
		WarehouseID: warehouseID,
		ProductName: productName,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	publish("parcel.received", ev)
}

// PublishParcelCreated publishes a parcel.created event.
func PublishParcelCreated(parcelID int64, trackingNo string, warehouseID, clientID int64, productName string) {
	if Bus == nil {
		return
	}
	ev := ParcelCreated{
		ParcelID:    parcelID,
		TrackingNo:  trackingNo,
		WarehouseID: warehouseID,
		ClientID:    clientID,
		ProductName: productName,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	publish("parcel.created", ev)
}

// PublishStatusChanged publishes a status.changed event.
func PublishStatusChanged(entityType string, entityID int64, oldStatus, newStatus string) {
	if Bus == nil {
		return
	}
	ev := StatusChanged{
		EntityType: entityType,
		EntityID:   entityID,
		OldStatus:  oldStatus,
		NewStatus:  newStatus,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}
	publish("status.changed", ev)
}

// PublishPaymentReceived publishes a payment.received event.
func PublishPaymentReceived(clientID int64, amount float64, txID string) {
	if Bus == nil {
		return
	}
	ev := PaymentReceived{
		ClientID:  clientID,
		Amount:    amount,
		TxID:      txID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	publish("payment.received", ev)
}

// PublishContainerClosed publishes a container.closed event.
func PublishContainerClosed(containerID int64, parcelCount int) {
	if Bus == nil {
		return
	}
	ev := ContainerClosed{
		ContainerID: containerID,
		ParcelCount: parcelCount,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	publish("container.closed", ev)
}

// ─── Internal publish helper ────────────────────────────────────────────────

func publish(topic string, ev interface{}) {
	data, err := json.Marshal(ev)
	if err != nil {
		log.Printf("[events] Failed to marshal event %s: %v", topic, err)
		return
	}
	msg := message.NewMessage(watermillUUID(), data)
	if err := Bus.Publish(topic, msg); err != nil {
		log.Printf("[events] Failed to publish event %s: %v", topic, err)
		return
	}
	log.Printf("[events] Published %s: %s", topic, string(data))
}

// watermillUUID generates a new Watermill UUID using Google UUID.
func watermillUUID() string {
	// Use a simple counter-based ID since we don't want to pull in the full
	// watermill UUID dep just for this.
	return time.Now().Format("20060102150405.000000")
}

// ─── Subscriber Helpers ─────────────────────────────────────────────────────

// SubscribeHandler subscribes a handler function to a topic on the bus.
// The handler receives the raw JSON payload as []byte.
func SubscribeHandler(topic string, handler func(ctx context.Context, payload []byte) error) {
	if Bus == nil {
		log.Printf("[events] Bus not initialized, cannot subscribe to %s", topic)
		return
	}
	messages, err := Bus.Subscribe(context.Background(), topic)
	if err != nil {
		log.Printf("[events] Failed to subscribe to %s: %v", topic, err)
		return
	}
	go func() {
		for msg := range messages {
			if err := handler(msg.Context(), msg.Payload); err != nil {
				log.Printf("[events] Handler error for %s: %v", topic, err)
			}
			msg.Ack()
		}
	}()
	log.Printf("[events] Subscribed to %s", topic)
}
