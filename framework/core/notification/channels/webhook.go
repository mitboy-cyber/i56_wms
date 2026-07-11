package channels

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/i56/framework/core/notification"
)

type WebhookChannel struct {
	client *http.Client
}

func NewWebhookChannel() *WebhookChannel {
	return &WebhookChannel{client: &http.Client{Timeout: 10 * time.Second}}
}

func (c *WebhookChannel) Name() string { return notification.ChannelWebhook }

func (c *WebhookChannel) Send(ctx context.Context, msg notification.Message) error {
	for _, url := range msg.To {
		body, _ := json.Marshal(msg.Data)
		resp, err := c.client.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			return err
		}
		resp.Body.Close()
		fmt.Printf("[WEBHOOK] %s -> %d\n", url, resp.StatusCode)
	}
	return nil
}
