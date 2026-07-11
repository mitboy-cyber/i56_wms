package repository

import (
	"sync"
	"time"
)

// ApiCredentialDisplay represents API credential data for client portal display.
type ApiCredentialDisplay struct {
	AppKey         string `json:"app_key"`
	AppSecret      string `json:"app_secret"`
	MaskedSecret   string `json:"masked_secret"`
	SecretVisible  bool   `json:"secret_visible"`
	Active         bool   `json:"active"`
	CreatedAt      string `json:"created_at"`
	Scopes         string `json:"scopes"`
	Timestamp      int64  `json:"timestamp"`
	Nonce          string `json:"nonce"`
	SignatureExample string `json:"signature_example"`
}

// MemApiCredentialRepo is an in-memory seed repo for API credentials.
type MemApiCredentialRepo struct {
	mu          sync.RWMutex
	credentials []ApiCredentialDisplay
}

func NewMemApiCredentialRepo() *MemApiCredentialRepo {
	return &MemApiCredentialRepo{
		credentials: []ApiCredentialDisplay{
			{
				AppKey:         "i5k_live_a1b2c3d4e5f6",
				AppSecret:      "i5s_live_z9y8x7w6v5u4t3s2r1q0",
				MaskedSecret:   "i5s_****r1q0",
				SecretVisible:  false,
				Active:         true,
				CreatedAt:      "2026-07-01",
				Scopes:         "read+write",
				Timestamp:      time.Now().Unix(),
				Nonce:          "a1b2c3d4e5f6g7h8",
				SignatureExample: "5k+Yma...8f3g",
			},
		},
	}
}

// List returns all API credentials.
func (r *MemApiCredentialRepo) List() []ApiCredentialDisplay {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]ApiCredentialDisplay, len(r.credentials))
	copy(result, r.credentials)
	return result
}

// GetByAppKey returns a credential by AppKey.
func (r *MemApiCredentialRepo) GetByAppKey(appKey string) *ApiCredentialDisplay {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.credentials {
		if r.credentials[i].AppKey == appKey {
			c := r.credentials[i]
			return &c
		}
	}
	return nil
}

// Create adds a new API credential.
func (r *MemApiCredentialRepo) Create(appKey, appSecret string) *ApiCredentialDisplay {
	r.mu.Lock()
	defer r.mu.Unlock()
	cred := ApiCredentialDisplay{
		AppKey:        appKey,
		AppSecret:     appSecret,
		MaskedSecret:  "i5s_****" + appSecret[len(appSecret)-4:],
		SecretVisible: false,
		Active:        true,
		CreatedAt:     time.Now().Format("2006-01-02"),
		Scopes:        "read+write",
		Timestamp:     time.Now().Unix(),
		Nonce:         "n" + time.Now().Format("20060102150405"),
		SignatureExample: "xkYma...8f3g",
	}
	r.credentials = append(r.credentials, cred)
	return &cred
}

// ResetSecret generates a new secret for an existing credential.
func (r *MemApiCredentialRepo) ResetSecret(appKey, newSecret string) *ApiCredentialDisplay {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.credentials {
		if r.credentials[i].AppKey == appKey {
			r.credentials[i].AppSecret = newSecret
			r.credentials[i].MaskedSecret = "i5s_****" + newSecret[len(newSecret)-4:]
			r.credentials[i].SecretVisible = false
			return &r.credentials[i]
		}
	}
	return nil
}
