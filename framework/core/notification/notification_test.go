package notification

import (
	"context"
	"errors"
	"testing"
)

type testLogger struct{}

func (l testLogger) Debug(msg string, args ...any) {}
func (l testLogger) Info(msg string, args ...any)  {}
func (l testLogger) Warn(msg string, args ...any)  {}
func (l testLogger) Error(msg string, args ...any) {}
func (l testLogger) With(args ...any) Logger       { return l }
func (l testLogger) WithGroup(name string) Logger  { return l }

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) Logger
	WithGroup(name string) Logger
}

// mockChannel implements Channel for testing.
type mockChannel struct {
	name    string
	sent    []Message
	sendErr error
}

func (c *mockChannel) Name() string { return c.name }

func (c *mockChannel) Send(ctx context.Context, msg Message) error {
	c.sent = append(c.sent, msg)
	return c.sendErr
}

func TestService_RegisterAndSend(t *testing.T) {
	svc := NewService(testLogger{})
	ch := &mockChannel{name: "email"}
	svc.Register(ch)

	msg := Message{Title: "Hello", Body: "World", To: []string{"user@example.com"}}
	err := svc.Send(context.Background(), "email", msg)
	if err != nil {
		t.Fatalf("Send: %v", err)
	}

	if len(ch.sent) != 1 {
		t.Fatalf("expected 1 sent message, got %d", len(ch.sent))
	}
	if ch.sent[0].Title != "Hello" {
		t.Errorf("expected 'Hello', got %q", ch.sent[0].Title)
	}
}

func TestService_SendUnknownChannel(t *testing.T) {
	svc := NewService(testLogger{})

	err := svc.Send(context.Background(), "nonexistent", Message{})
	if err != ErrChannelNotFound {
		t.Errorf("expected ErrChannelNotFound, got %v", err)
	}
}

func TestService_SendAll(t *testing.T) {
	svc := NewService(testLogger{})
	email := &mockChannel{name: "email"}
	sms := &mockChannel{name: "sms"}

	svc.Register(email)
	svc.Register(sms)

	msg := Message{Title: "Alert", Body: "System down"}
	errs := svc.SendAll(context.Background(), msg)

	if len(errs) != 0 {
		t.Errorf("expected no errors, got %d: %v", len(errs), errs)
	}
	if len(email.sent) != 1 {
		t.Errorf("expected email to receive message")
	}
	if len(sms.sent) != 1 {
		t.Errorf("expected sms to receive message")
	}
}

func TestService_SendAllWithError(t *testing.T) {
	svc := NewService(testLogger{})
	email := &mockChannel{name: "email"}
	failing := &mockChannel{name: "failing", sendErr: errors.New("down")}

	svc.Register(email)
	svc.Register(failing)

	msg := Message{Title: "Test"}
	errs := svc.SendAll(context.Background(), msg)

	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}
	if len(email.sent) != 1 {
		t.Errorf("expected email to receive message even though other channel failed")
	}
}

func TestMessage(t *testing.T) {
	msg := Message{
		Title:    "Test",
		Body:     "Hello World",
		To:       []string{"recipient"},
		Template: "welcome",
		Data:     map[string]any{"key": "value"},
	}

	if msg.Title != "Test" {
		t.Errorf("expected 'Test', got %q", msg.Title)
	}
	if len(msg.To) != 1 {
		t.Errorf("expected 1 recipient, got %d", len(msg.To))
	}
}

func TestChannelConstants(t *testing.T) {
	if ChannelEmail != "email" {
		t.Errorf("expected 'email', got %q", ChannelEmail)
	}
	if ChannelSMS != "sms" {
		t.Errorf("expected 'sms', got %q", ChannelSMS)
	}
	if ChannelWebhook != "webhook" {
		t.Errorf("expected 'webhook', got %q", ChannelWebhook)
	}
	if ChannelSlack != "slack" {
		t.Errorf("expected 'slack', got %q", ChannelSlack)
	}
}

func TestErrorType(t *testing.T) {
	if ErrChannelNotFound.Error() != "notification channel not found" {
		t.Errorf("unexpected error message: %q", ErrChannelNotFound)
	}
}
