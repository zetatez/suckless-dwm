package cache

import (
	"sync"
	"time"
)

type entry struct {
	value      any
	expiration int64
}

type Cache struct {
	data map[string]*entry
	mu   sync.RWMutex
	ttl  time.Duration
	stop chan struct{}
}

func NewCache(ttl time.Duration, cleanupInterval time.Duration) *Cache {
	c := &Cache{
		data: make(map[string]*entry),
		ttl:  ttl,
		stop: make(chan struct{}),
	}
	go c.cleanupLoop(cleanupInterval)
	return c
}

func (c *Cache) cleanupLoop(interval time.Duration) {
	if interval <= 0 {
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixMilli()
			c.mu.Lock()
			for k, v := range c.data {
				if v.expiration > 0 && v.expiration < now {
					delete(c.data, k)
				}
			}
			c.mu.Unlock()
		case <-c.stop:
			return
		}
	}
}

func (c *Cache) Close() {
	close(c.stop)
}

func (c *Cache) Set(key string, value any, ttl ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expir := int64(0)
	duration := c.ttl
	if len(ttl) > 0 {
		duration = ttl[0]
	}
	if duration > 0 {
		expir = time.Now().Add(duration).UnixMilli()
	}

	c.data[key] = &entry{value: value, expiration: expir}
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	e, ok := c.data[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	if e.expiration > 0 && e.expiration < time.Now().UnixMilli() {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		return nil, false
	}

	return e.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *Cache) Update(key string, value any, ttl ...time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.data[key]
	if !ok {
		return false
	}

	expir := e.expiration
	if len(ttl) > 0 {
		expir = time.Now().Add(ttl[0]).UnixMilli()
	}

	c.data[key] = &entry{value: value, expiration: expir}
	return true
}

func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]*entry)
}

func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.data))
	for k := range c.data {
		keys = append(keys, k)
	}
	return keys
}

func (c *Cache) Values() []any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	values := make([]any, 0, len(c.data))
	for _, v := range c.data {
		values = append(values, v.value)
	}
	return values
}

func (c *Cache) Items() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	items := make(map[string]any, len(c.data))
	for k, v := range c.data {
		items[k] = v.value
	}
	return items
}

func (c *Cache) Range(f func(key string, value any) bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for k, v := range c.data {
		if !f(k, v.value) {
			break
		}
	}
}
