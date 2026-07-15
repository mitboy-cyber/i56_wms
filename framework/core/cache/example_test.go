package cache_test

import (
	"context"
	"fmt"
	"time"

	"github.com/i56/framework/core/cache"
)

// ExampleMemCache demonstrates basic cache operations.
func ExampleMemCache() {
	c := cache.NewMemCache()
	ctx := context.Background()

	// Set a value with TTL
	c.Set(ctx, "user:42", `{"name":"Alice"}`, 5*time.Minute)

	// Get the value
	val, err := c.Get(ctx, "user:42")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(val)

	// Check existence
	exists, _ := c.Exists(ctx, "user:42")
	fmt.Println("Exists:", exists)

	// Delete
	c.Delete(ctx, "user:42")
	_, err = c.Get(ctx, "user:42")
	fmt.Println("After delete:", err)
	// Output:
	// {"name":"Alice"}
	// Exists: true
	// After delete: cache: key not found
}

// ExampleMultiLevelCache demonstrates two-level caching.
func ExampleMultiLevelCache() {
	l1 := cache.NewMemCache() // fast in-memory
	l2 := cache.NewMemCache() // could be Redis in production
	mc := cache.NewMultiLevelCache(l1, l2)
	ctx := context.Background()

	mc.Set(ctx, "config:app", "my-value", time.Hour)

	// Read from L1
	val, _ := mc.Get(ctx, "config:app")
	fmt.Println(val)
	// Output:
	// my-value
}
