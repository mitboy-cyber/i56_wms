package unit

import (
	"testing"
	"time"

	"github.com/i56/framework/core/auth"
	"github.com/i56/framework/core/config"
)

func newTestTokenManager(t *testing.T) *auth.TokenManager {
	cfg := config.AuthConfig{
		Issuer:          "test-issuer",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 30 * 24 * time.Hour,
	}
	tm, err := auth.NewTokenManager(cfg)
	if err != nil {
		t.Fatalf("failed to create token manager: %v", err)
	}
	return tm
}

func TestIssueAccessToken(t *testing.T) {
	tm := newTestTokenManager(t)
	token, err := tm.IssueAccessToken("user-123", "tenant-1", []string{"admin"}, []string{"read", "write"})
	if err != nil {
		t.Fatalf("failed to issue token: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
	// Token should have 3 parts (header.payload.signature)
	parts := 0
	for _, c := range token {
		if c == '.' {
			parts++
		}
	}
	if parts != 2 {
		t.Errorf("expected 3-part JWT, got %d separators", parts)
	}
}

func TestValidateAccessToken(t *testing.T) {
	tm := newTestTokenManager(t)
	token, _ := tm.IssueAccessToken("user-123", "tenant-1", []string{"admin"}, nil)

	claims, err := tm.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.Subject != "user-123" {
		t.Errorf("expected subject 'user-123', got '%s'", claims.Subject)
	}
	if claims.TenantID != "tenant-1" {
		t.Errorf("expected tenant 'tenant-1', got '%s'", claims.TenantID)
	}
	if claims.Type != "access" {
		t.Errorf("expected type 'access', got '%s'", claims.Type)
	}
	if claims.Issuer != "test-issuer" {
		t.Errorf("expected issuer 'test-issuer', got '%s'", claims.Issuer)
	}
}

func TestValidateTokenInvalidSignature(t *testing.T) {
	tm := newTestTokenManager(t)
	token, _ := tm.IssueAccessToken("user-123", "tenant-1", nil, nil)

	// Tamper with the token
	tampered := token[:len(token)-4] + "xxxx"

	if _, err := tm.ValidateAccessToken(tampered); err == nil {
		t.Error("expected error for tampered token")
	}
}

func TestValidateTokenWrongIssuer(t *testing.T) {
	tm := newTestTokenManager(t)
	token, _ := tm.IssueAccessToken("user-123", "tenant-1", nil, nil)

	// Create a second token manager with different issuer
	cfg2 := config.AuthConfig{
		Issuer:          "different-issuer",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 30 * 24 * time.Hour,
	}
	tm2, _ := auth.NewTokenManager(cfg2)

	if _, err := tm2.ValidateAccessToken(token); err == nil {
		t.Error("expected error when validating with wrong issuer")
	}
}

func TestRefreshTokenCannotBeUsedAsAccessToken(t *testing.T) {
	tm := newTestTokenManager(t)
	refreshToken, _ := tm.IssueRefreshToken("user-123", "tenant-1")

	if _, err := tm.ValidateAccessToken(refreshToken); err == nil {
		t.Error("expected error when using refresh token as access token")
	}
}

func TestRefreshTokenValidation(t *testing.T) {
	tm := newTestTokenManager(t)
	refreshToken, _ := tm.IssueRefreshToken("user-123", "tenant-1")

	claims, err := tm.ValidateRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("failed to validate refresh token: %v", err)
	}
	if claims.Type != "refresh" {
		t.Errorf("expected type 'refresh', got '%s'", claims.Type)
	}
}

func TestTokenContainsRoles(t *testing.T) {
	tm := newTestTokenManager(t)
	token, _ := tm.IssueAccessToken("user-123", "tenant-1", []string{"admin", "editor"}, nil)

	claims, _ := tm.ValidateAccessToken(token)
	if len(claims.Roles) != 2 {
		t.Errorf("expected 2 roles, got %d", len(claims.Roles))
	}
}

func TestTokenContainsScopes(t *testing.T) {
	tm := newTestTokenManager(t)
	token, _ := tm.IssueAccessToken("user-123", "tenant-1", nil, []string{"read:orders", "write:orders"})

	claims, _ := tm.ValidateAccessToken(token)
	if len(claims.Scopes) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(claims.Scopes))
	}
}
