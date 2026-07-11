// Package notification provides a multi-channel notification center with
// built-in message templates for warehouse operations.
package notification

import (
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------
// Notifier interface
// ---------------------------------------------------------------------------

// Notifier is the interface that every notification channel must implement.
type Notifier interface {
	Send(to, subject, body string) error
}

// ---------------------------------------------------------------------------
// Built-in templates
// ---------------------------------------------------------------------------

// Predefined message templates keyed by event name.
var DefaultTemplates = map[string]string{
	"order_created":     "您的订单 {order_no} 已创建，金额 ¥{amount}",
	"parcel_received":   "包裹 {tracking} 已入库",
	"shipment_departed": "您的包裹已从 {warehouse} 发出",
}

// RenderTemplate substitutes {placeholders} in the template string with values
// from data.
func RenderTemplate(tmpl string, data map[string]interface{}) string {
	result := tmpl
	for k, v := range data {
		result = strings.ReplaceAll(result, "{"+k+"}", fmt.Sprint(v))
	}
	return result
}

// ---------------------------------------------------------------------------
// Channel implementations
// ---------------------------------------------------------------------------

// EmailNotifier logs email delivery (smtp integration placeholder).
type EmailNotifier struct {
	SMTPHost string
	Username string
	Password string
}

func (e *EmailNotifier) Send(to, subject, body string) error {
	fmt.Printf("[EMAIL] To: %s | Subject: %s | Body: %s\n", to, subject, body)
	return nil
}

// SMSNotifier logs SMS delivery (SMS gateway placeholder).
type SMSNotifier struct {
	APIKey string
}

func (s *SMSNotifier) Send(to, subject, body string) error {
	fmt.Printf("[SMS] To: %s | Body: %s\n", to, body)
	return nil
}

// WebhookNotifier logs webhook delivery.
type WebhookNotifier struct {
	URL string
}

func (w *WebhookNotifier) Send(to, subject, body string) error {
	fmt.Printf("[WEBHOOK] URL: %s | Subject: %s | Body: %s\n", w.URL, subject, body)
	return nil
}

// ---------------------------------------------------------------------------
// NotificationCenter
// ---------------------------------------------------------------------------

// NotificationCenter manages multiple notifier channels and template-based
// message delivery.
type NotificationCenter struct {
	notifiers map[string]Notifier
	templates map[string]string
}

// NewCenter creates a NotificationCenter pre-loaded with default templates.
func NewCenter() *NotificationCenter {
	// Copy default templates so callers can mutate safely
	tmpl := make(map[string]string, len(DefaultTemplates))
	for k, v := range DefaultTemplates {
		tmpl[k] = v
	}
	return &NotificationCenter{
		notifiers: make(map[string]Notifier),
		templates: tmpl,
	}
}

// Register adds a notifier under the given channel name.
func (nc *NotificationCenter) Register(channel string, n Notifier) {
	nc.notifiers[channel] = n
}

// RegisterTemplate adds or overwrites a named template.
func (nc *NotificationCenter) RegisterTemplate(name, tmpl string) {
	nc.templates[name] = tmpl
}

// Notify sends a message through the given channel after rendering the named
// template with the supplied data.
func (nc *NotificationCenter) Notify(channel, to, template string, data map[string]interface{}) error {
	n, ok := nc.notifiers[channel]
	if !ok {
		return fmt.Errorf("notification: channel %q not registered", channel)
	}

	tmpl, ok := nc.templates[template]
	if !ok {
		return fmt.Errorf("notification: template %q not found", template)
	}

	body := RenderTemplate(tmpl, data)
	subject := template // use template name as subject fallback
	return n.Send(to, subject, body)
}

// ListChannels returns all registered channel names.
func (nc *NotificationCenter) ListChannels() []string {
	names := make([]string, 0, len(nc.notifiers))
	for k := range nc.notifiers {
		names = append(names, k)
	}
	return names
}

// ListTemplates returns all registered template names.
func (nc *NotificationCenter) ListTemplates() []string {
	names := make([]string, 0, len(nc.templates))
	for k := range nc.templates {
		names = append(names, k)
	}
	return names
}
