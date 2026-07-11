package channels

import (
	"context"
	"fmt"

	"github.com/i56/framework/core/notification"
)

type SMSChannel struct {
	apiKey string
}

func NewSMSChannel(apiKey string) *SMSChannel {
	return &SMSChannel{apiKey: apiKey}
}

func (c *SMSChannel) Name() string { return notification.ChannelSMS }

func (c *SMSChannel) Send(ctx context.Context, msg notification.Message) error {
	fmt.Printf("[SMS] To:%v Body:%s\n", msg.To, msg.Body)
	return nil
}
