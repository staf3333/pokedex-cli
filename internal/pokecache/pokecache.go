package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	cacheMap map[string]cacheEntry
	mu sync.Mutex
	interval time.Duration
}

func NewCache(interval time.Duration) Cache {
	return Cache{
		cacheMap: make(map[string]cacheEntry),
		interval: interval,
	}
}

func (c *Cache) Add(key string, val []byte) {
	// need to use a mutex to lock the map while doing operation
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheMap[key] = cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.cacheMap[key]
	if !ok {
		return nil, false
	} 
	return entry.val, true
}

func (c *Cache) reapLoop() {
	
}