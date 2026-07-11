package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemCache_SetAndGet(t *testing.T) {
	c := NewMemCache()
	ctx := context.Background()

	err := c.Set(ctx, "key1", "value1", time.Minute)
	if err != nil {
		t.Fatalf("Set: %v", err)
	}

	val, err := c.Get(ctx, "key1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if val != "value1" {
		t.Errorf("expected 'value1', got %q", val)
	}
}

func TestMemCache_NotFound(t *testing.T) {
	c := NewMemCache()
	ctx := context.Background()

	_, err := c.Get(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemCache_Expiry(t *testing.T) {
	c := NewMemCache()
	ctx := context.Background()

	c.Set(ctx, "temp", "val", 50*time.Millisecond)
	time.Sleep(100 * time.Millisecond)

	_, err := c.Get(ctx, "temp")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound for expired key, got %v", err)
	}
}

func TestMemCache_Delete(t *testing.T) {
	c := NewMemCache()
	ctx := context.Background()

	c.Set(ctx, "key", "val", time.Minute)
	c.Delete(ctx, "key")

	_, err := c.Get(ctx, "key")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestMemCache_Exists(t *testing.T) {
	c := NewMemCache()
	ctx := context.Background()

	ok, err := c.Exists(ctx, "key")
	if err != nil {
		t.Fatalf("Exists (missing): %v", err)
	}
	if ok {
		t.Error("expected false for missing key")
	}

	c.Set(ctx, "key", "val", time.Minute)

	ok, err = c.Exists(ctx, "key")
	if err != nil {
		t.Fatalf("Exists (present): %v", err)
	}
	if !ok {
		t.Error("expected true for present key")
	}
}

func TestMemCache_Clear(t *testing.T) {
	c := NewMemCache()
	ctx := context.Background()

	c.Set(ctx, "a", "1", time.Minute)
	c.Set(ctx, "b", "2", time.Minute)
	c.Clear(ctx)

	ok, _ := c.Exists(ctx, "a")
	if ok {
		t.Error("expected false after clear")
	}
}

func TestMultiLevelCache_L1Hit(t *testing.T) {
	l1 := NewMemCache()
	l2 := NewMemCache()
	mc := NewMultiLevelCache(l1, l2)
	ctx := context.Background()

	mc.Set(ctx, "shared", "value", time.Minute)

	val, err := mc.Get(ctx, "shared")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if val != "value" {
		t.Errorf("expected 'value', got %q", val)
	}
}

func TestMultiLevelCache_L2Fallback(t *testing.T) {
	l1 := NewMemCache()
	l2 := NewMemCache()
	mc := NewMultiLevelCache(l1, l2)
	ctx := context.Background()

	// Set only in L2
	l2.Set(ctx, "only_l2", "from_l2", time.Minute)

	val, err := mc.Get(ctx, "only_l2")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if val != "from_l2" {
		t.Errorf("expected 'from_l2', got %q", val)
	}

	// Should now be in L1 too
	l1Val, _ := l1.Get(ctx, "only_l2")
	if l1Val != "from_l2" {
		t.Errorf("expected L1 to be populated, got %q", l1Val)
	}
}

func TestMultiLevelCache_Delete(t *testing.T) {
	l1 := NewMemCache()
	l2 := NewMemCache()
	mc := NewMultiLevelCache(l1, l2)
	ctx := context.Background()

	mc.Set(ctx, "key", "val", time.Minute)
	mc.Delete(ctx, "key")

	_, err := l1.Get(ctx, "key")
	if err != ErrNotFound {
		t.Error("L1 should not have key")
	}
	_, err = l2.Get(ctx, "key")
	if err != ErrNotFound {
		t.Error("L2 should not have key")
	}
}

func TestCacheError(t *testing.T) {
	if ErrNotFound.Error() != "cache: key not found" {
		t.Errorf("unexpected error message: %q", ErrNotFound)
	}
}
