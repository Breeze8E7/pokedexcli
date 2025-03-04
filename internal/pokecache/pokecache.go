package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	m        map[string]cacheEntry
	mu       sync.Mutex
	ticker   *time.Ticker
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		m:        make(map[string]cacheEntry),
		ticker:   time.NewTicker(interval),
		interval: interval,
	}
	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var entry cacheEntry
	entry.createdAt = time.Now()
	entry.val = val
	c.m[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, exists := c.m[key]
	if !exists {
		return nil, false
	} else {
		return entry.val, true
	}
}

func (c *Cache) reapLoop() {
	for {
		<-c.ticker.C
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.m {
			if now.Sub(entry.createdAt) > c.interval {
				delete(c.m, key)
			}
		}
		c.mu.Unlock()
	}
}
