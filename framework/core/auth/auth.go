// Package auth provides ed25519 JWT token management with refresh rotation,
// API key authentication, and HMAC request signing.
package auth

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/i56/framework/core/config"
	"github.com/i56/framework/core/errors"
)

// ---------------------------------------------------------------------------
// JWT Token Manager
// ---------------------------------------------------------------------------

// TokenManager handles JWT token issuance, validation, and refresh.
type TokenManager struct {
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
	signKey    ed25519.PrivateKey
	verifyKey  ed25519.PublicKey
}

// Claims represents JWT standard + custom claims.
type Claims struct {
	Subject   string   `json:"sub"`
	TenantID  string   `json:"tid,omitempty"`
	UserID    string   `json:"uid,omitempty"`
	Roles     []string `json:"roles,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
	Issuer    string   `json:"iss"`
	Audience  string   `json:"aud,omitempty"`
	IssuedAt  int64    `json:"iat"`
	ExpiresAt int64    `json:"exp"`
	NotBefore int64    `json:"nbf,omitempty"`
	JTI       string   `json:"jti"`
	Type      string   `json:"type,omitempty"` // "access" | "refresh"
}

// NewTokenManager creates a TokenManager with generated ed25519 keys.
func NewTokenManager(cfg config.AuthConfig) (*TokenManager, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("auth: generate key: %w", err)
	}
	return &TokenManager{
		issuer:     cfg.Issuer,
		accessTTL:  cfg.AccessTokenTTL,
		refreshTTL: cfg.RefreshTokenTTL,
		signKey:    priv,
		verifyKey:  pub,
	}, nil
}

// NewTokenManagerWithKeys creates a TokenManager with existing keys.
func NewTokenManagerWithKeys(cfg config.AuthConfig, priv ed25519.PrivateKey) *TokenManager {
	return &TokenManager{
		issuer:     cfg.Issuer,
		accessTTL:  cfg.AccessTokenTTL,
		refreshTTL: cfg.RefreshTokenTTL,
		signKey:    priv,
		verifyKey:  priv.Public().(ed25519.PublicKey),
	}
}

// IssueAccessToken creates a signed JWT access token.
func (tm *TokenManager) IssueAccessToken(subject, tenantID string, roles, scopes []string) (string, error) {
	now := time.Now()
	claims := Claims{
		Subject:   subject,
		TenantID:  tenantID,
		Roles:     roles,
		Scopes:    scopes,
		Issuer:    tm.issuer,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(tm.accessTTL).Unix(),
		NotBefore: now.Unix(),
		JTI:       generateJTI(),
		Type:      "access",
	}
	return tm.sign(claims)
}

// IssueRefreshToken creates a signed JWT refresh token.
func (tm *TokenManager) IssueRefreshToken(subject, tenantID string) (string, error) {
	now := time.Now()
	claims := Claims{
		Subject:   subject,
		TenantID:  tenantID,
		Issuer:    tm.issuer,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(tm.refreshTTL).Unix(),
		JTI:       generateJTI(),
		Type:      "refresh",
	}
	return tm.sign(claims)
}

// ValidateAccessToken validates and parses an access token.
func (tm *TokenManager) ValidateAccessToken(tokenStr string) (*Claims, error) {
	claims, err := tm.verify(tokenStr)
	if err != nil {
		return nil, err
	}
	if claims.Type != "access" {
		return nil, errors.NewUnauthorized("token is not an access token")
	}
	return claims, nil
}

// ValidateRefreshToken validates and parses a refresh token.
func (tm *TokenManager) ValidateRefreshToken(tokenStr string) (*Claims, error) {
	claims, err := tm.verify(tokenStr)
	if err != nil {
		return nil, err
	}
	if claims.Type != "refresh" {
		return nil, errors.NewUnauthorized("token is not a refresh token")
	}
	return claims, nil
}

// sign creates a signed JWT string using ed25519.
func (tm *TokenManager) sign(claims Claims) (string, error) {
	header := map[string]string{"alg": "EdDSA", "typ": "JWT"}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("auth: marshal header: %w", err)
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("auth: marshal claims: %w", err)
	}

	headerB64 := base64URLEncode(headerJSON)
	claimsB64 := base64URLEncode(claimsJSON)
	signingInput := headerB64 + "." + claimsB64

	signature := ed25519.Sign(tm.signKey, []byte(signingInput))
	sigB64 := base64URLEncode(signature)

	return signingInput + "." + sigB64, nil
}

// verify parses and validates a JWT token.
func (tm *TokenManager) verify(tokenStr string) (*Claims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, errors.NewUnauthorized("invalid token format")
	}

	signingInput := parts[0] + "." + parts[1]
	signature, err := base64URLDecode(parts[2])
	if err != nil {
		return nil, errors.NewUnauthorized("invalid token signature encoding")
	}

	claimsJSON, err := base64URLDecode(parts[1])
	if err != nil {
		return nil, errors.NewUnauthorized("invalid token claims encoding")
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, errors.NewUnauthorized("invalid token claims")
	}

	// Verify signature
	if !ed25519.Verify(tm.verifyKey, []byte(signingInput), signature) {
		return nil, errors.NewUnauthorized("invalid token signature")
	}

	// Check expiry
	now := time.Now().Unix()
	if claims.ExpiresAt > 0 && now > claims.ExpiresAt {
		return nil, errors.NewUnauthorized("token has expired")
	}

	// Check not-before
	if claims.NotBefore > 0 && now < claims.NotBefore {
		return nil, errors.NewUnauthorized("token not yet valid")
	}

	// Check issuer
	if claims.Issuer != tm.issuer {
		return nil, errors.NewUnauthorized("invalid token issuer")
	}

	return &claims, nil
}

// PublicKey returns the public key for JWKS endpoints.
func (tm *TokenManager) PublicKey() ed25519.PublicKey {
	return tm.verifyKey
}

// AccessTTL returns the access token TTL.
func (tm *TokenManager) AccessTTL() time.Duration {
	return tm.accessTTL
}

// RefreshTTL returns the refresh token TTL.
func (tm *TokenManager) RefreshTTL() time.Duration {
	return tm.refreshTTL
}

// ---------------------------------------------------------------------------
// API Key Authentication
// ---------------------------------------------------------------------------

// APIKeyStore is the interface for persisting and validating API keys.
type APIKeyStore interface {
	// ValidateKey checks an API key and returns the associated identity.
	ValidateKey(key string) (*APIKeyIdentity, bool, error)
	// GenerateKey creates a new API key for a given identity.
	GenerateKey(identity APIKeyIdentity) (string, error)
	// RevokeKey invalidates an API key.
	RevokeKey(key string) error
}

// APIKeyIdentity holds the identity associated with an API key.
type APIKeyIdentity struct {
	KeyID       string   `json:"key_id"`
	TenantID    string   `json:"tenant_id"`
	UserID      string   `json:"user_id,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Name        string   `json:"name"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// APIKeyManager manages API key generation and validation.
type APIKeyManager struct {
	store  APIKeyStore
	prefix string // e.g. "i56_" — all keys start with this
}

// NewAPIKeyManager creates an APIKeyManager.
func NewAPIKeyManager(store APIKeyStore, prefix string) *APIKeyManager {
	if prefix == "" {
		prefix = "i56_"
	}
	return &APIKeyManager{store: store, prefix: prefix}
}

// CreateKey generates a new API key for the given identity.
func (m *APIKeyManager) CreateKey(identity APIKeyIdentity) (string, error) {
	return m.store.GenerateKey(identity)
}

// ValidateKey checks an API key and returns the identity.
func (m *APIKeyManager) ValidateKey(key string) (*APIKeyIdentity, bool, error) {
	return m.store.ValidateKey(key)
}

// RevokeKey invalidates an API key.
func (m *APIKeyManager) RevokeKey(key string) error {
	return m.store.RevokeKey(key)
}

// GenerateSecureKey creates a cryptographically random API key.
func GenerateSecureKey(prefix string) (string, error) {
	if prefix == "" {
		prefix = "i56_"
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("auth: generate key: %w", err)
	}
	return prefix + base64.RawURLEncoding.EncodeToString(raw), nil
}

// InMemAPIKeyStore is an in-memory implementation of APIKeyStore.
type InMemAPIKeyStore struct {
	keys  map[string]*APIKeyIdentity // key → identity
}

// NewInMemAPIKeyStore creates an in-memory API key store.
func NewInMemAPIKeyStore() *InMemAPIKeyStore {
	return &InMemAPIKeyStore{keys: make(map[string]*APIKeyIdentity)}
}

func (s *InMemAPIKeyStore) ValidateKey(key string) (*APIKeyIdentity, bool, error) {
	id, ok := s.keys[key]
	if !ok {
		return nil, false, nil
	}
	if id.ExpiresAt != nil && time.Now().After(*id.ExpiresAt) {
		delete(s.keys, key)
		return nil, false, nil
	}
	return id, true, nil
}

func (s *InMemAPIKeyStore) GenerateKey(identity APIKeyIdentity) (string, error) {
	key, err := GenerateSecureKey("i56_")
	if err != nil {
		return "", err
	}
	id := identity
	id.KeyID = key
	s.keys[key] = &id
	return key, nil
}

func (s *InMemAPIKeyStore) RevokeKey(key string) error {
	delete(s.keys, key)
	return nil
}

// ---------------------------------------------------------------------------
// HMAC Request Signing
// ---------------------------------------------------------------------------

// HMACSigner creates and verifies HMAC-SHA256 request signatures.
type HMACSigner struct {
	secret []byte
}

// NewHMACSigner creates a new HMACSigner with the given secret.
func NewHMACSigner(secret string) *HMACSigner {
	return &HMACSigner{secret: []byte(secret)}
}

// Sign computes an HMAC-SHA256 signature for the given payload.
// Returns a hex-encoded signature.
func (s *HMACSigner) Sign(payload []byte) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// SignWithTimestamp includes a timestamp in the payload to prevent replay attacks.
// Returns signature and the timestamp string suitable for a header.
func (s *HMACSigner) SignWithTimestamp(payload []byte) (signature, timestamp string) {
	ts := fmt.Sprintf("%d", time.Now().Unix())
	combined := append([]byte(ts+":"), payload...)
	return s.Sign(combined), ts
}

// Verify checks whether the given signature matches the payload.
func (s *HMACSigner) Verify(payload []byte, signatureHex string) bool {
	sig, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(payload)
	expected := mac.Sum(nil)
	return hmac.Equal(sig, expected)
}

// VerifyWithTimestamp checks a signature that was created with SignWithTimestamp.
func (s *HMACSigner) VerifyWithTimestamp(payload []byte, signatureHex, timestampStr string, maxAge time.Duration) bool {
	ts, err := time.Parse("2006-01-02T15:04:05Z", timestampStr)
	if err != nil {
		// Try Unix timestamp
		var unix int64
		if _, err := fmt.Sscanf(timestampStr, "%d", &unix); err != nil {
			return false
		}
		ts = time.Unix(unix, 0)
	}
	if time.Since(ts) > maxAge || time.Since(ts) < -maxAge {
		return false
	}
	combined := append([]byte(timestampStr+":"), payload...)
	return s.Verify(combined, signatureHex)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func generateJTI() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func base64URLDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}
