package jwt

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// Claims represents JWT payload
type Claims struct {
	Sub       string `json:"sub"`        // user id
	TenantID  int64  `json:"tenant_id"`  // multi-tenant
	Role      string `json:"role"`       // admin/operator/client
	Username  string `json:"username"`   // display name
	WarehouseIDs []int64 `json:"wh_ids,omitempty"` // data scope
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

// Token represents a signed JWT
type Token struct {
	Raw       string
	Claims    Claims
	ExpiresAt time.Time
}

// Service handles JWT signing and verification with Ed25519
type Service struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
	issuer     string
}

// NewService creates a JWT service with a new keypair
func NewService(issuer string) (*Service, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ed25519 key: %w", err)
	}
	return &Service{privateKey: priv, publicKey: pub, issuer: issuer}, nil
}

// NewServiceFromSeed creates a JWT service from a base64-encoded seed
func NewServiceFromSeed(seedB64 string, issuer string) (*Service, error) {
	seed, err := base64.StdEncoding.DecodeString(seedB64)
	if err != nil {
		return nil, fmt.Errorf("invalid seed: %w", err)
	}
	priv := ed25519.NewKeyFromSeed(seed)
	return &Service{privateKey: priv, publicKey: priv.Public().(ed25519.PublicKey), issuer: issuer}, nil
}

// PublicKeyBase64 returns the base64-encoded public key
func (s *Service) PublicKeyBase64() string {
	return base64.StdEncoding.EncodeToString(s.publicKey)
}

// Sign creates a signed JWT token with the given claims
func (s *Service) Sign(claims Claims) (string, error) {
	claims.IssuedAt = time.Now().Unix()
	if claims.ExpiresAt == 0 {
		claims.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()
	}

	// Build header + payload
	header := map[string]string{"alg": "EdDSA", "typ": "JWT"}
	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(claims)

	// Encode
	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// Sign
	signingInput := headerB64 + "." + payloadB64
	sig := ed25519.Sign(s.privateKey, []byte(signingInput))
	sigB64 := base64.RawURLEncoding.EncodeToString(sig)

	return signingInput + "." + sigB64, nil
}

// Verify validates a JWT token and returns the claims
func (s *Service) Verify(tokenString string) (*Claims, error) {
	parts := splitToken(tokenString)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	signingInput := parts[0] + "." + parts[1]
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid signature encoding: %w", err)
	}

	// Verify Ed25519 signature
	if !ed25519.Verify(s.publicKey, []byte(signingInput), sig) {
		return nil, fmt.Errorf("invalid signature")
	}

	// Decode payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid payload encoding: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, fmt.Errorf("invalid claims: %w", err)
	}

	// Check expiration
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

func splitToken(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}
