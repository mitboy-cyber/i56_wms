// Package knowledge provides knowledge base management for document
// import, FAQ, and SOP retrieval.
package knowledge

import (
	"sync"
	"time"
)

// Entry is a single knowledge base entry.
type Entry struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"` // "document", "faq", "sop"
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Tags      []string          `json:"tags"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// SearchHit is a ranked knowledge base search result.
type SearchHit struct {
	Entry Entry   `json:"entry"`
	Score float64 `json:"score"`
}

// KnowledgeBase manages the knowledge repository.
type KnowledgeBase struct {
	mu      sync.RWMutex
	entries map[string]*Entry // id → entry
	byType  map[string][]string // type → []id
}

// NewKnowledgeBase creates an empty knowledge base.
func NewKnowledgeBase() *KnowledgeBase {
	return &KnowledgeBase{
		entries: make(map[string]*Entry),
		byType:  make(map[string][]string),
	}
}

// Import adds a document to the knowledge base.
func (kb *KnowledgeBase) Import(entry *Entry) {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	now := time.Now()
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = now
	}
	entry.UpdatedAt = now

	kb.entries[entry.ID] = entry
	kb.byType[entry.Type] = append(kb.byType[entry.Type], entry.ID)
}

// Get retrieves an entry by ID.
func (kb *KnowledgeBase) Get(id string) (*Entry, bool) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	e, ok := kb.entries[id]
	return e, ok
}

// ListByType returns all entries of a given type.
func (kb *KnowledgeBase) ListByType(typ string) []Entry {
	kb.mu.RLock()
	defer kb.mu.RUnlock()

	ids := kb.byType[typ]
	out := make([]Entry, 0, len(ids))
	for _, id := range ids {
		if e, ok := kb.entries[id]; ok {
			out = append(out, *e)
		}
	}
	return out
}

// Search performs a basic keyword search across all entries.
// TODO: Replace with vector similarity search via RAG pipeline.
func (kb *KnowledgeBase) Search(query string, topK int) []SearchHit {
	kb.mu.RLock()
	defer kb.mu.RUnlock()

	var hits []SearchHit
	for _, entry := range kb.entries {
		// Stub: exact substring match
		hits = append(hits, SearchHit{
			Entry: *entry,
			Score: 0.5, // TODO: compute real relevance score
		})
	}
	// TODO: Sort by score, limit to topK
	_ = topK
	return hits
}

// Delete removes an entry from the knowledge base.
func (kb *KnowledgeBase) Delete(id string) {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	entry, ok := kb.entries[id]
	if !ok {
		return
	}
	delete(kb.entries, id)

	// Remove from type index
	ids := kb.byType[entry.Type]
	for i, eid := range ids {
		if eid == id {
			kb.byType[entry.Type] = append(ids[:i], ids[i+1:]...)
			break
		}
	}
}

// Count returns the total number of entries.
func (kb *KnowledgeBase) Count() int {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	return len(kb.entries)
}

// Update modifies an existing entry.
func (kb *KnowledgeBase) Update(id string, entry *Entry) bool {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	if _, ok := kb.entries[id]; !ok {
		return false
	}
	entry.UpdatedAt = time.Now()
	kb.entries[id] = entry
	return true
}
