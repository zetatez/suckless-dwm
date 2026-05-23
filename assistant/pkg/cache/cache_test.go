package cache

import (
	"sync"
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	c := NewCache(time.Hour, 0)
	defer c.Close()

	c.Set("key1", "value1")
	v, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected to get value")
	}
	if v != "value1" {
		t.Fatalf("expected value1, got %v", v)
	}
}

func TestCache_GetNotFound(t *testing.T) {
	c := NewCache(time.Hour, 0)
	defer c.Close()

	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestCache_Delete(t *testing.T) {
	c := NewCache(time.Hour, 0)
	defer c.Close()

	c.Set("key1", "value1")
	c.Delete("key1")
	_, ok := c.Get("key1")
	if ok {
		t.Fatal("expected not found after delete")
	}
}

func TestCache_Update(t *testing.T) {
	c := NewCache(time.Hour, 0)
	defer c.Close()

	c.Set("key1", "value1")
	ok := c.Update("key1", "value2")
	if !ok {
		t.Fatal("expected update to return true")
	}
	v, _ := c.Get("key1")
	if v != "value2" {
		t.Fatalf("expected value2, got %v", v)
	}
}

func TestCache_UpdateNotFound(t *testing.T) {
	c := NewCache(time.Hour, 0)
	defer c.Close()

	ok := c.Update("nonexistent", "value")
	if ok {
		t.Fatal("expected update to return false for nonexistent key")
	}
}

func TestCache_WithTTL(t *testing.T) {
	c := NewCache(time.Hour, 0)
	defer c.Close()

	c.Set("key1", "value1", 50*time.Millisecond)
	time.Sleep(100 * time.Millisecond)

	_, ok := c.Get("key1")
	if ok {
		t.Fatal("expected expired key to be not found")
	}
}

func TestCache_DefaultTTL(t *testing.T) {
	c := NewCache(50*time.Millisecond, 0)
	defer c.Close()

	c.Set("key1", "value1")
	time.Sleep(100 * time.Millisecond)

	_, ok := c.Get("key1")
	if ok {
		t.Fatal("expected expired key to be not found")
	}
}

func TestCache_Clear(t *testing.T) {
	c := NewCache(time.Hour, 0)
	defer c.Close()

	c.Set("key1", "value1")
	c.Set("key2", "value2")
	c.Clear()

	if c.Len() != 0 {
		t.Fatalf("expected len 0, got %d", c.Len())
	}
}

func TestCache_Keys(t *testing.T) {
	c := NewCache(time.Hour, 0)
	defer c.Close()

	c.Set("key1", "value1")
	c.Set("key2", "value2")

	keys := c.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestCache_Values(t *testing.T) {
	c := NewCache(time.Hour, 0)
	defer c.Close()

	c.Set("key1", "value1")
	c.Set("key2", "value2")

	values := c.Values()
	if len(values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(values))
	}
}

func TestCache_Items(t *testing.T) {
	c := NewCache(time.Hour, 0)
	defer c.Close()

	c.Set("key1", "value1")
	c.Set("key2", "value2")

	items := c.Items()
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items["key1"] != "value1" {
		t.Fatal("unexpected item value")
	}
}

func TestCache_Range(t *testing.T) {
	c := NewCache(time.Hour, 0)
	defer c.Close()

	c.Set("key1", "value1")
	c.Set("key2", "value2")

	count := 0
	c.Range(func(key string, value any) bool {
		count++
		return true
	})
	if count != 2 {
		t.Fatalf("expected range count 2, got %d", count)
	}
}

func TestCache_RangeBreak(t *testing.T) {
	c := NewCache(time.Hour, 0)
	defer c.Close()

	c.Set("key1", "value1")
	c.Set("key2", "value2")

	count := 0
	c.Range(func(key string, value any) bool {
		count++
		return false
	})
	if count != 1 {
		t.Fatalf("expected range count 1, got %d", count)
	}
}

func TestCache_Len(t *testing.T) {
	c := NewCache(time.Hour, 0)
	defer c.Close()

	if c.Len() != 0 {
		t.Fatalf("expected len 0, got %d", c.Len())
	}

	c.Set("key1", "value1")
	if c.Len() != 1 {
		t.Fatalf("expected len 1, got %d", c.Len())
	}

	c.Delete("key1")
	if c.Len() != 0 {
		t.Fatalf("expected len 0, got %d", c.Len())
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	c := NewCache(time.Hour, 0)
	defer c.Close()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Set("key", i)
			c.Get("key")
			c.Len()
		}(i)
	}
	wg.Wait()
}
