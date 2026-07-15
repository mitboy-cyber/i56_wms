package auth_test

import (
	"fmt"
	"time"

	"github.com/i56/framework/core/auth"
	"github.com/i56/framework/core/config"
)

// ExampleTokenManager demonstrates JWT token issuance and validation.
func ExampleTokenManager() {
	cfg := config.AuthConfig{
		Issuer:          "my-app",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 30 * 24 * time.Hour,
	}
	tm, err := auth.NewTokenManager(cfg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Issue an access token
	token, err := tm.IssueAccessToken("user-42", "tenant-1",
		[]string{"admin", "operator"},
		[]string{"read", "write", "delete"})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Validate it back
	claims, err := tm.ValidateAccessToken(token)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(claims.Subject)
	fmt.Println(claims.TenantID)
	fmt.Println(claims.Roles[0])
	// Output:
	// user-42
	// tenant-1
	// admin
}

// ExampleAPIKeyManager demonstrates API key creation and validation.
func ExampleAPIKeyManager() {
	store := auth.NewInMemAPIKeyStore()
	mgr := auth.NewAPIKeyManager(store, "i56_")

	key, err := mgr.CreateKey(auth.APIKeyIdentity{
		TenantID: "tenant-1",
		UserID:   "user-1",
		Name:     "production-key",
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Validate the key
	id, ok, err := mgr.ValidateKey(key)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(ok)
	fmt.Println(id.Name)
	// Output:
	// true
	// production-key
}

// ExampleHMACSigner demonstrates request signing.
func ExampleHMACSigner() {
	signer := auth.NewHMACSigner("my-shared-secret")

	payload := []byte(`{"order_id": "ORD-1234"}`)
	sig, ts := signer.SignWithTimestamp(payload)

	// Verify the signature is valid
	valid := signer.VerifyWithTimestamp(payload, sig, ts, 5*time.Minute)
	fmt.Println(valid)
	// Output:
	// true
}
