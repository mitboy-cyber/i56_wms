// Package memory implements user and enterprise memory stores
// with JSONB summaries for persistent conversation context.
package memory

import (
	"sync"
	"time"
)

// Summary is a compressed representation of past interactions.
type Summary struct {
	Key       string    `json:"key"`
	Content   string    `json:"content"`
	Tokens    int       `json:"tokens"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserMemory stores per-user conversation history and preferences.
type UserMemory struct {
	mu       sync.RWMutex
	UserID   string             `json:"user_id"`
	Summaries map[string]*Summary `json:"summaries"`
	Facts    []string           `json:"facts"`    // Extracted facts about the user
	Preferences map[string]string `json:"preferences"` // User preferences
}

// NewUserMemory creates a user memory store.
func NewUserMemory(userID string) *UserMemory {
	return &UserMemory{
		UserID:      userID,
		Summaries:   make(map[string]*Summary),
		Preferences: make(map[string]string),
	}
}

// AddSummary stores a conversation summary.
func (m *UserMemory) AddSummary(s *Summary) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s.UpdatedAt = time.Now()
	m.Summaries[s.Key] = s
}

// GetSummary retrieves a summary by key.
func (m *UserMemory) GetSummary(key string) *Summary {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Summaries[key]
}

// AddFact records an extracted fact about the user.
func (m *UserMemory) AddFact(fact string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Facts = append(m.Facts, fact)
}

// SetPreference stores a user preference.
func (m *UserMemory) SetPreference(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Preferences[key] = value
}

// EnterpriseMemory stores organization-wide knowledge.
type EnterpriseMemory struct {
	mu        sync.RWMutex
	TenantID  string              `json:"tenant_id"`
	Docs      map[string]*Summary `json:"docs"`
	SOPs      map[string]*Summary `json:"sops"`
	FAQs      map[string]*Summary `json:"faqs"`
}

// NewEnterpriseMemory creates an enterprise memory store.
func NewEnterpriseMemory(tenantID string) *EnterpriseMemory {
	return &EnterpriseMemory{
		TenantID: tenantID,
		Docs:     make(map[string]*Summary),
		SOPs:     make(map[string]*Summary),
		FAQs:     make(map[string]*Summary),
	}
}

// AddDoc stores a document summary.
func (m *EnterpriseMemory) AddDoc(key string, s *Summary) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Docs[key] = s
}

// AddSOP stores a standard operating procedure.
func (m *EnterpriseMemory) AddSOP(key string, s *Summary) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SOPs[key] = s
}

// AddFAQ stores a frequently asked question entry.
func (m *EnterpriseMemory) AddFAQ(key string, s *Summary) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.FAQs[key] = s
}

// Search performs a basic keyword search across all memory stores.
// TODO: Replace with vector similarity search.
func (m *EnterpriseMemory) Search(query string) []*Summary {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var results []*Summary
	for _, s := range m.Docs {
		results = append(results, s)
	}
	for _, s := range m.SOPs {
		results = append(results, s)
	}
	for _, s := range m.FAQs {
		results = append(results, s)
	}
	return results
}
