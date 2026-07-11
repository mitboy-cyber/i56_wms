// Package rag implements a RAG (Retrieval-Augmented Generation) pipeline
// backed by pgvector for document chunking, embedding, and similarity search.
package rag

import (
	"context"
	"sync"
)

// Chunk is a text segment with its vector embedding.
type Chunk struct {
	ID        string    `json:"id"`
	DocumentID string   `json:"document_id"`
	Content   string    `json:"content"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
	Metadata  map[string]string `json:"metadata"`
}

// SearchResult is a ranked chunk with a similarity score.
type SearchResult struct {
	Chunk      Chunk   `json:"chunk"`
	Score      float64 `json:"score"`
	DocumentID string  `json:"document_id"`
}

// Embedder generates vector embeddings for text.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
}

// VectorStore persists and searches vector embeddings.
type VectorStore interface {
	Insert(ctx context.Context, chunks []Chunk) error
	Search(ctx context.Context, embedding []float32, topK int) ([]SearchResult, error)
	Delete(ctx context.Context, documentID string) error
}

// Pipeline is the RAG ingestion and retrieval pipeline.
type Pipeline struct {
	mu         sync.RWMutex
	embedder   Embedder
	store      VectorStore
	chunkSize  int
	chunkOverlap int
}

// NewPipeline creates a RAG pipeline.
func NewPipeline(embedder Embedder, store VectorStore) *Pipeline {
	return &Pipeline{
		embedder:     embedder,
		store:        store,
		chunkSize:    512,
		chunkOverlap: 50,
	}
}

// SetChunkConfig configures chunking parameters.
func (p *Pipeline) SetChunkConfig(size, overlap int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.chunkSize = size
	p.chunkOverlap = overlap
}

// Ingest splits a document into chunks, embeds them, and stores them.
func (p *Pipeline) Ingest(ctx context.Context, documentID, content string) error {
	// Split into chunks
	chunks := p.splitContent(documentID, content)
	if len(chunks) == 0 {
		return nil
	}

	// Embed chunks
	texts := make([]string, len(chunks))
	for i, c := range chunks {
		texts[i] = c.Content
	}

	embeddings, err := p.embedder.EmbedBatch(ctx, texts)
	if err != nil {
		return err
	}

	// Attach embeddings to chunks
	for i := range chunks {
		chunks[i].Embedding = embeddings[i]
	}

	// Store in vector database
	return p.store.Insert(ctx, chunks)
}

// Search finds the top-k most relevant chunks for a query.
func (p *Pipeline) Search(ctx context.Context, query string, topK int) ([]SearchResult, error) {
	embedding, err := p.embedder.Embed(ctx, query)
	if err != nil {
		return nil, err
	}
	return p.store.Search(ctx, embedding, topK)
}

// Delete removes all chunks for a document.
func (p *Pipeline) Delete(ctx context.Context, documentID string) error {
	return p.store.Delete(ctx, documentID)
}

// splitContent divides text into overlapping chunks.
func (p *Pipeline) splitContent(documentID, content string) []Chunk {
	p.mu.RLock()
	size := p.chunkSize
	overlap := p.chunkOverlap
	p.mu.RUnlock()

	if len(content) <= size {
		return []Chunk{{
			DocumentID: documentID,
			Content:    content,
			Index:      0,
			Metadata:   make(map[string]string),
		}}
	}

	runes := []rune(content)
	var chunks []Chunk
	step := size - overlap
	if step <= 0 {
		step = size
	}

	for i := 0; i < len(runes); i += step {
		end := i + size
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, Chunk{
			DocumentID: documentID,
			Content:    string(runes[i:end]),
			Index:      len(chunks),
			Metadata:   make(map[string]string),
		})
		if end >= len(runes) {
			break
		}
	}

	return chunks
}
