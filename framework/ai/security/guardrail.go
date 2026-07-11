// Package security implements AI-specific security controls including
// PII masking, prompt injection detection, and RBAC circuit breakers.
package security

import (
	"context"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Guardrail enforces security policies on all AI interactions.
type Guardrail struct {
	mu sync.RWMutex

	// PII patterns for masking sensitive data.
	piiPatterns []*regexp.Regexp

	// Injection patterns for detecting prompt injection attempts.
	injectionPatterns []*regexp.Regexp

	// RBAC rules: permission → allowed model tiers.
	rbacRules map[string]string

	// CircuitBreaker tracks consecutive failures to open the breaker.
	consecutiveFailures int
	maxFailures         int
	breakerOpen         bool
	breakerOpenedAt     time.Time
	cooldownDuration    time.Duration

	// Stats for observability.
	blockedCount int
	piiMaskCount int
}

// Config holds security configuration.
type Config struct {
	MaxFailures      int           `json:"max_failures"`
	CooldownDuration time.Duration `json:"cooldown_duration"`
	CustomPIIPatterns []string     `json:"custom_pii_patterns"`
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxFailures:      5,
		CooldownDuration: 30 * time.Second,
	}
}

// New creates a Guardrail with the given config.
func New(cfg Config) *Guardrail {
	g := &Guardrail{
		rbacRules:       make(map[string]string),
		maxFailures:     cfg.MaxFailures,
		cooldownDuration: cfg.CooldownDuration,
	}

	// Built-in PII patterns
	g.piiPatterns = []*regexp.Regexp{
		regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),                              // SSN
		regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`),                             // Credit card
		regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`), // Email
		regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`),             // IPv4
	}

	// Built-in injection patterns
	g.injectionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)ignore (all |previous )?(instructions?|prompts?)`),
		regexp.MustCompile(`(?i)you are now (DAN|jailbroken|unrestricted)`),
		regexp.MustCompile(`(?i)pretend (you are|to be)`),
		regexp.MustCompile(`(?i)system:\s*`),
	}

	// Add custom patterns
	for _, p := range cfg.CustomPIIPatterns {
		if re, err := regexp.Compile(p); err == nil {
			g.piiPatterns = append(g.piiPatterns, re)
		}
	}

	return g
}

// MaskPII scans text for PII patterns and replaces them with placeholders.
func (g *Guardrail) MaskPII(text string) (string, bool) {
	masked := false
	result := text
	for _, re := range g.piiPatterns {
		if re.MatchString(result) {
			result = re.ReplaceAllString(result, "[REDACTED]")
			masked = true
		}
	}
	if masked {
		g.mu.Lock()
		g.piiMaskCount++
		g.mu.Unlock()
	}
	return result, masked
}

// DetectInjection checks whether the input contains prompt injection patterns.
// Returns true if injection is detected, along with a description of the match.
func (g *Guardrail) DetectInjection(input string) (bool, string) {
	for _, re := range g.injectionPatterns {
		if match := re.FindString(input); match != "" {
			g.mu.Lock()
			g.blockedCount++
			g.mu.Unlock()
			return true, match
		}
	}
	return false, ""
}

// CheckRBAC verifies whether a user role is authorized for the requested action.
func (g *Guardrail) CheckRBAC(role, permission string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	allowedRole, ok := g.rbacRules[permission]
	if !ok {
		return true // No rule defined → default allow
	}
	return strings.EqualFold(role, allowedRole) || allowedRole == "*"
}

// SetRBACRule configures a permission-to-role mapping.
func (g *Guardrail) SetRBACRule(permission, requiredRole string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.rbacRules[permission] = requiredRole
}

// AllowRequest checks if the circuit breaker permits requests.
func (g *Guardrail) AllowRequest() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if !g.breakerOpen {
		return true
	}
	if time.Since(g.breakerOpenedAt) > g.cooldownDuration {
		return true
	}
	return false
}

// RecordSuccess resets the circuit breaker on success.
func (g *Guardrail) RecordSuccess() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.consecutiveFailures = 0
	g.breakerOpen = false
}

// RecordFailure increments the failure counter and opens the breaker if threshold exceeded.
func (g *Guardrail) RecordFailure() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.consecutiveFailures++
	if g.consecutiveFailures >= g.maxFailures {
		g.breakerOpen = true
		g.breakerOpenedAt = time.Now()
	}
}

// PreFlight runs all security checks before an AI request.
// It returns nil if the request is safe to process.
func (g *Guardrail) PreFlight(ctx context.Context, input string, role string, permission string) error {
	// Check circuit breaker
	if !g.AllowRequest() {
		return &SecurityError{Code: "circuit_open", Message: "circuit breaker is open"}
	}

	// Check RBAC
	if !g.CheckRBAC(role, permission) {
		return &SecurityError{Code: "rbac_denied", Message: "insufficient permissions"}
	}

	// Check prompt injection
	if detected, match := g.DetectInjection(input); detected {
		return &SecurityError{Code: "injection_detected", Message: "prompt injection detected: " + match}
	}

	return nil
}

// Stats returns current security metrics.
type Stats struct {
	BlockedCount    int  `json:"blocked_count"`
	PIIMaskCount    int  `json:"pii_mask_count"`
	BreakerOpen     bool `json:"breaker_open"`
	ConsecutiveFails int `json:"consecutive_fails"`
}

// Stats returns the current guardrail statistics.
func (g *Guardrail) Stats() Stats {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return Stats{
		BlockedCount:     g.blockedCount,
		PIIMaskCount:     g.piiMaskCount,
		BreakerOpen:      g.breakerOpen,
		ConsecutiveFails: g.consecutiveFailures,
	}
}

// SecurityError is a structured error from security checks.
type SecurityError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *SecurityError) Error() string {
	return e.Code + ": " + e.Message
}
