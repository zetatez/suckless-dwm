package cache

import (
	"container/list"
	"sync"
	"time"
)

type entry struct {
	value      any
	expiration int64
	key        string
}

type Cache struct {
	data     map[string]*entry
	mu       sync.RWMutex
	ttl      time.Duration
	stop     chan struct{}
	maxItems int
	lruList  *list.List
}

func NewCache(ttl time.Duration, cleanupInterval time.Duration) *Cache {
	return NewCacheWithMax(ttl, cleanupInterval, 0)
}

func NewCacheWithMax(ttl time.Duration, cleanupInterval time.Duration, maxItems int) *Cache {
	c := &Cache{
		data:     make(map[string]*entry),
		ttl:      ttl,
		stop:     make(chan struct{}),
		maxItems: maxItems,
		lruList:  list.New(),
	}
	if cleanupInterval > 0 {
		go c.cleanupLoop(cleanupInterval)
	}
	return c
}

func (c *Cache) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixMilli()
			c.mu.Lock()
			for k, v := range c.data {
				if v.expiration > 0 && v.expiration < now {
					c.removeEntry(k, v)
				}
			}
			c.mu.Unlock()
		case <-c.stop:
			return
		}
	}
}

func (c *Cache) removeEntry(k string, v *entry) {
	delete(c.data, k)
	if c.lruList != nil {
		for e := c.lruList.Front(); e != nil; e = e.Next() {
			if e.Value.(*entry).key == k {
				c.lruList.Remove(e)
				break
			}
		}
	}
}

func (c *Cache) evictLRU() {
	if c.maxItems <= 0 || c.lruList == nil {
		return
	}
	for c.lruList.Len() >= c.maxItems {
		e := c.lruList.Front()
		if e == nil {
			break
		}
		v := e.Value.(*entry)
		delete(c.data, v.key)
		c.lruList.Remove(e)
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

	// 如果 key 已存在，先移除旧的 LRU 结点
	if _, ok := c.data[key]; ok && c.lruList != nil {
		for e := c.lruList.Front(); e != nil; e = e.Next() {
			if e.Value.(*entry).key == key {
				c.lruList.Remove(e)
				break
			}
		}
	}

	c.data[key] = &entry{value: value, expiration: expir, key: key}
	if c.lruList != nil {
		c.lruList.PushBack(c.data[key])
	}

	c.evictLRU()
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
		c.removeEntry(key, e)
		c.mu.Unlock()
		return nil, false
	}

	// LRU: 移到链表尾部
	if c.lruList != nil {
		c.mu.Lock()
		for el := c.lruList.Front(); el != nil; el = el.Next() {
			if el.Value.(*entry).key == key {
				c.lruList.MoveToBack(el)
				break
			}
		}
		c.mu.Unlock()
	}

	return e.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.data[key]; ok {
		c.removeEntry(key, e)
	}
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

	e.value = value
	e.expiration = expir

	// LRU: 移到链表尾部
	if c.lruList != nil {
		for el := c.lruList.Front(); el != nil; el = el.Next() {
			if el.Value.(*entry).key == key {
				c.lruList.MoveToBack(el)
				break
			}
		}
	}

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
	if c.lruList != nil {
		c.lruList.Init()
	}
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
