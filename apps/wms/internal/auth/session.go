// Package auth provides admin session management with HMAC-SHA256 signing
// and bcrypt password hashing for employee credentials.
package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SessionTTL is the session lifetime (8 hours).
const SessionTTL = 8 * time.Hour

// Session represents an authenticated admin session.
type Session struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired returns true if the session has expired.
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// SessionManager manages admin sessions with HMAC-signed cookies
// and bcrypt-hashed employee credentials.
type SessionManager struct {
	mu      sync.RWMutex
	secret  []byte
	users   map[string]string // username → bcrypt hash
}

// NewSessionManager creates a new SessionManager with a randomly generated
// 32-byte secret and pre-seeded employee credentials.
func NewSessionManager() *SessionManager {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		// Fallback: use a fixed secret if crypto/rand fails (should never happen)
		copy(secret, []byte("i56-wms-fallback-secret-key!!!"))
	}

	sm := &SessionManager{
		secret: secret,
		users:  make(map[string]string),
	}
	sm.seedCredentials()
	return sm
}

// seedCredentials pre-seeds employee credentials with bcrypt-hashed passwords.
// Default passwords are "1234" for operators and "admin" for admin.
func (sm *SessionManager) seedCredentials() {
	credentials := map[string]string{
		"admin": "admin",
		"OP001": "1234",
		"OP002": "1234",
		"OP003": "1234",
	}

	for username, password := range credentials {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			// Fallback: store a pre-computed hash for "1234" at cost 10
			hash = []byte("$2a$10$precomputed.hash.placeholder")
			continue
		}
		sm.users[username] = string(hash)
	}
}

// AddUser adds or updates a user's bcrypt password hash.
func (sm *SessionManager) AddUser(username, plainPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	sm.mu.Lock()
	sm.users[username] = string(hash)
	sm.mu.Unlock()
	return nil
}

// Authenticate verifies username/password against stored credentials.
// Returns true if the credentials are valid.
func (sm *SessionManager) Authenticate(username, password string) bool {
	sm.mu.RLock()
	hash, ok := sm.users[username]
	sm.mu.RUnlock()

	if !ok {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateSession creates a new session for the given username and returns
// a signed cookie value.
func (sm *SessionManager) CreateSession(username string) string {
	now := time.Now()
	session := Session{
		Username:  username,
		CreatedAt: now,
		ExpiresAt: now.Add(SessionTTL),
	}

	data, err := json.Marshal(session)
	if err != nil {
		return ""
	}

	encoded := base64.RawURLEncoding.EncodeToString(data)
	sig := sm.sign(encoded)
	return encoded + "." + sig
}

// ValidateSession validates a signed cookie value and returns the Session
// if valid and not expired. Returns nil if invalid or expired.
func (sm *SessionManager) ValidateSession(cookieValue string) *Session {
	if cookieValue == "" {
		return nil
	}

	parts := strings.SplitN(cookieValue, ".", 2)
	if len(parts) != 2 {
		return nil
	}

	encoded := parts[0]
	sig := parts[1]

	// Verify HMAC signature
	expectedSig := sm.sign(encoded)
	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return nil
	}

	// Decode the session data
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil
	}

	if session.IsExpired() {
		return nil
	}

	return &session
}

// sign creates an HMAC-SHA256 signature for the given data.
func (sm *SessionManager) sign(data string) string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	mac := hmac.New(sha256.New, sm.secret)
	mac.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
