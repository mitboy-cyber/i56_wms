package auth

import (
	"testing"
	"time"

	"github.com/i56/framework/core/config"
)

func TestTokenManager_IssueAndValidate(t *testing.T) {
	cfg := config.AuthConfig{
		Issuer:          "test",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour,
	}
	tm, err := NewTokenManager(cfg)
	if err != nil {
		t.Fatalf("NewTokenManager: %v", err)
	}

	// Issue access token
	token, err := tm.IssueAccessToken("user-1", "tenant-1", []string{"admin"}, []string{"read", "write"})
	if err != nil {
		t.Fatalf("IssueAccessToken: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}

	// Validate
	claims, err := tm.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken: %v", err)
	}
	if claims.Subject != "user-1" {
		t.Errorf("expected subject 'user-1', got %q", claims.Subject)
	}
	if claims.TenantID != "tenant-1" {
		t.Errorf("expected tenant 'tenant-1', got %q", claims.TenantID)
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != "admin" {
		t.Errorf("unexpected roles: %v", claims.Roles)
	}
}

func TestTokenManager_RefreshToken(t *testing.T) {
	cfg := config.AuthConfig{
		Issuer:          "test",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour,
	}
	tm, _ := NewTokenManager(cfg)

	refresh, err := tm.IssueRefreshToken("user-1", "tenant-1")
	if err != nil {
		t.Fatalf("IssueRefreshToken: %v", err)
	}

	claims, err := tm.ValidateRefreshToken(refresh)
	if err != nil {
		t.Fatalf("ValidateRefreshToken: %v", err)
	}
	if claims.Type != "refresh" {
		t.Errorf("expected type 'refresh', got %q", claims.Type)
	}

	// Access token validation should fail for refresh token
	_, err = tm.ValidateAccessToken(refresh)
	if err == nil {
		t.Error("expected error when validating refresh as access token")
	}
}

func TestAPIKeyManager_CreateAndValidate(t *testing.T) {
	store := NewInMemAPIKeyStore()
	mgr := NewAPIKeyManager(store, "i56_")

	identity := APIKeyIdentity{
		TenantID: "tenant-1",
		UserID:   "user-1",
		Name:     "test-key",
	}

	key, err := mgr.CreateKey(identity)
	if err != nil {
		t.Fatalf("CreateKey: %v", err)
	}
	if key == "" {
		t.Error("expected non-empty key")
	}
	if key[:4] != "i56_" {
		t.Errorf("expected key prefix 'i56_', got %q", key[:4])
	}

	// Validate
	id, ok, err := mgr.ValidateKey(key)
	if err != nil {
		t.Fatalf("ValidateKey: %v", err)
	}
	if !ok {
		t.Error("expected key to be valid")
	}
	if id.TenantID != "tenant-1" {
		t.Errorf("expected tenant 'tenant-1', got %q", id.TenantID)
	}
}

func TestAPIKeyManager_InvalidKey(t *testing.T) {
	store := NewInMemAPIKeyStore()
	mgr := NewAPIKeyManager(store, "")

	_, ok, err := mgr.ValidateKey("invalid-key")
	if err != nil {
		t.Fatalf("ValidateKey: %v", err)
	}
	if ok {
		t.Error("expected invalid key to return ok=false")
	}
}

func TestAPIKeyManager_Revoke(t *testing.T) {
	store := NewInMemAPIKeyStore()
	mgr := NewAPIKeyManager(store, "i56_")

	key, _ := mgr.CreateKey(APIKeyIdentity{TenantID: "t1", Name: "key"})

	_, ok, _ := mgr.ValidateKey(key)
	if !ok {
		t.Error("expected valid key before revoke")
	}

	mgr.RevokeKey(key)

	_, ok, _ = mgr.ValidateKey(key)
	if ok {
		t.Error("expected invalid key after revoke")
	}
}

func TestAPIKeyManager_ExpiredKey(t *testing.T) {
	store := NewInMemAPIKeyStore()
	mgr := NewAPIKeyManager(store, "i56_")

	expiresAt := time.Now().Add(-1 * time.Hour)
	key, _ := store.GenerateKey(APIKeyIdentity{TenantID: "t1", Name: "expired", ExpiresAt: &expiresAt})

	_, ok, _ := mgr.ValidateKey(key)
	if ok {
		t.Error("expected expired key to be invalid")
	}
}

func TestHMACSigner_SignAndVerify(t *testing.T) {
	signer := NewHMACSigner("my-secret-key")
	payload := []byte("hello world")

	sig := signer.Sign(payload)
	if sig == "" {
		t.Error("expected non-empty signature")
	}

	if !signer.Verify(payload, sig) {
		t.Error("expected valid signature")
	}

	if signer.Verify([]byte("tampered"), sig) {
		t.Error("expected invalid signature for tampered payload")
	}
}

func TestHMACSigner_SignWithTimestamp(t *testing.T) {
	signer := NewHMACSigner("secret")
	payload := []byte(`{"order":"123"}`)

	sig, ts := signer.SignWithTimestamp(payload)
	if sig == "" || ts == "" {
		t.Error("expected non-empty signature and timestamp")
	}

	if !signer.VerifyWithTimestamp(payload, sig, ts, 5*time.Minute) {
		t.Error("expected valid signature with timestamp")
	}

	// Old timestamp should fail
	if signer.VerifyWithTimestamp(payload, sig, "0", 5*time.Minute) {
		t.Error("expected invalid signature for old timestamp")
	}
}

func TestGenerateSecureKey(t *testing.T) {
	key1, err := GenerateSecureKey("")
	if err != nil {
		t.Fatalf("GenerateSecureKey: %v", err)
	}
	key2, _ := GenerateSecureKey("")

	if key1 == key2 {
		t.Error("expected unique keys for each generation")
	}
	if len(key1) < 40 {
		t.Errorf("expected key length >= 40, got %d", len(key1))
	}
}
