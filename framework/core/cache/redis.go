// Package cache provides multi-level caching with Redis support.
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements Cache using Redis.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis-backed cache.
func NewRedisCache(addr string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("cache: redis connection failed: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// NewRedisCacheWithOptions creates a Redis cache with full options.
func NewRedisCacheWithOptions(addr, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("cache: redis connection failed: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// Get retrieves a value from Redis.
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("cache: redis get failed: %w", err)
	}
	return val, nil
}

// Set stores a value in Redis with TTL.
func (c *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("cache: redis set failed: %w", err)
	}
	return nil
}

// Delete removes a key from Redis.
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("cache: redis delete failed: %w", err)
	}
	return nil
}

// Exists checks if a key exists in Redis.
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("cache: redis exists check failed: %w", err)
	}
	return n > 0, nil
}

// Clear flushes all keys from the current Redis DB.
func (c *RedisCache) Clear(ctx context.Context) error {
	err := c.client.FlushDB(ctx).Err()
	if err != nil {
		return fmt.Errorf("cache: redis flush failed: %w", err)
	}
	return nil
}

// Flush is an alias for Clear (backward compatibility).
func (c *RedisCache) Flush(ctx context.Context) error {
	return c.Clear(ctx)
}

// Close closes the Redis connection.
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// Client returns the underlying Redis client for advanced operations.
func (c *RedisCache) Client() *redis.Client {
	return c.client
}

// Ensure RedisCache implements the Cache interface.
var _ Cache = (*RedisCache)(nil)
