package channels

import (
	"context"
	"fmt"

	"github.com/i56/framework/core/notification"
)

type EmailChannel struct {
	smtpHost string
	smtpPort int
}

func NewEmailChannel(host string, port int) *EmailChannel {
	return &EmailChannel{smtpHost: host, smtpPort: port}
}

func (c *EmailChannel) Name() string { return notification.ChannelEmail }

func (c *EmailChannel) Send(ctx context.Context, msg notification.Message) error {
	fmt.Printf("[EMAIL] To:%v Subject:%s Body:%s\n", msg.To, msg.Title, msg.Body)
	return nil
}
