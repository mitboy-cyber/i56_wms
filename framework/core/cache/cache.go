// Package cache provides multi-level caching.
package cache

import (
	"context"
	"sync"
	"time"
)

// Cache is the abstract cache interface.
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Clear(ctx context.Context) error
}

// MemCache is an in-memory cache implementation.
type MemCache struct {
	mu    sync.RWMutex
	items map[string]*cacheItem
}

type cacheItem struct {
	value   string
	expires time.Time
}

// NewMemCache creates a new in-memory cache.
func NewMemCache() *MemCache {
	mc := &MemCache{items: make(map[string]*cacheItem)}
	go mc.cleaner()
	return mc
}

func (c *MemCache) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return "", ErrNotFound
	}
	if time.Now().After(item.expires) {
		return "", ErrNotFound
	}
	return item.value, nil
}

func (c *MemCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &cacheItem{
		value:   value,
		expires: time.Now().Add(ttl),
	}
	return nil
}

func (c *MemCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

func (c *MemCache) Exists(ctx context.Context, key string) (bool, error) {
	_, err := c.Get(ctx, key)
	if err == ErrNotFound {
		return false, nil
	}
	return err == nil, err
}

func (c *MemCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*cacheItem)
	return nil
}

func (c *MemCache) cleaner() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for k, v := range c.items {
			if now.After(v.expires) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}

// MultiLevelCache combines L1 (memory) and L2 (distributed) caches.
type MultiLevelCache struct {
	L1 Cache
	L2 Cache
}

// NewMultiLevelCache creates a two-level cache.
func NewMultiLevelCache(l1, l2 Cache) *MultiLevelCache {
	return &MultiLevelCache{L1: l1, L2: l2}
}

func (c *MultiLevelCache) Get(ctx context.Context, key string) (string, error) {
	// Try L1 first
	val, err := c.L1.Get(ctx, key)
	if err == nil {
		return val, nil
	}
	// Fall back to L2
	val, err = c.L2.Get(ctx, key)
	if err != nil {
		return "", err
	}
	// Populate L1
	_ = c.L1.Set(ctx, key, val, 5*time.Minute)
	return val, nil
}

func (c *MultiLevelCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if err := c.L2.Set(ctx, key, value, ttl); err != nil {
		return err
	}
	return c.L1.Set(ctx, key, value, min(ttl, 5*time.Minute))
}

func (c *MultiLevelCache) Delete(ctx context.Context, key string) error {
	_ = c.L1.Delete(ctx, key)
	return c.L2.Delete(ctx, key)
}

func (c *MultiLevelCache) Exists(ctx context.Context, key string) (bool, error) {
	ok, _ := c.L1.Exists(ctx, key)
	if ok {
		return true, nil
	}
	return c.L2.Exists(ctx, key)
}

func (c *MultiLevelCache) Clear(ctx context.Context) error {
	_ = c.L1.Clear(ctx)
	return c.L2.Clear(ctx)
}

// ErrNotFound is returned when a cache key is not found.
var ErrNotFound = &cacheError{"cache: key not found"}

type cacheError struct{ msg string }

func (e *cacheError) Error() string { return e.msg }
