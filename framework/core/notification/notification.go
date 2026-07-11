// Package notification provides multi-channel message delivery.
package notification

import (
	"context"

	"github.com/i56/framework/core/logger"
)

// Message represents a notification to be sent.
type Message struct {
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	To       []string          `json:"to"`
	Template string            `json:"template,omitempty"`
	Data     map[string]any    `json:"data,omitempty"`
}

// Channel is the interface that all notification channels implement.
type Channel interface {
	Name() string
	Send(ctx context.Context, msg Message) error
}

// Service orchestrates notification delivery across channels.
type Service struct {
	channels map[string]Channel
	log      logger.Logger
}

// NewService creates a notification service.
func NewService(log logger.Logger) *Service {
	return &Service{
		channels: make(map[string]Channel),
		log:      log,
	}
}

// Register adds a notification channel.
func (s *Service) Register(ch Channel) {
	s.channels[ch.Name()] = ch
	s.log.Info("notification channel registered", "channel", ch.Name())
}

// Send delivers a message through the specified channel.
func (s *Service) Send(ctx context.Context, channelName string, msg Message) error {
	ch, ok := s.channels[channelName]
	if !ok {
		return ErrChannelNotFound
	}
	return ch.Send(ctx, msg)
}

// SendAll delivers a message through all registered channels.
func (s *Service) SendAll(ctx context.Context, msg Message) map[string]error {
	errs := make(map[string]error)
	for name, ch := range s.channels {
		if err := ch.Send(ctx, msg); err != nil {
			errs[name] = err
			s.log.Error("notification send failed", "channel", name, "error", err)
		}
	}
	return errs
}

// Channel names.
const (
	ChannelEmail    = "email"
	ChannelSMS      = "sms"
	ChannelLINE     = "line"
	ChannelTelegram = "telegram"
	ChannelWebhook  = "webhook"
	ChannelSlack    = "slack"
	ChannelInApp    = "in_app"
)

// Predefined errors.
var (
	ErrChannelNotFound = &notifError{"notification channel not found"}
)

type notifError struct{ msg string }

func (e *notifError) Error() string { return e.msg }
