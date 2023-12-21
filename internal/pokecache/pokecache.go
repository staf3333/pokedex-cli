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

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cacheMap: make(map[string]cacheEntry),
		interval: interval,
	}
	go c.reapLoop()
	return c
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
	// time.NewTicker returns a new Ticker containing a channel that will send the current time on 
	// channel after each tick. Period of ticks is specified by duration arg
	// have to specify time interval in seconds b/c c.interval doesn't convert to second automatically
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	// use for range loop for idiomatic handling of the ticker channel
	// automatically reads values from the channel until it is closed -> don't have to add
	// extra logic to check if channel is closed
	for t := range ticker.C {
		// here we want to range through the entries in the cache map, and compare the createdAt 
		// time to the t recieved from the ticker
		c.mu.Lock()
		for k, v := range c.cacheMap {
			if t.Sub(v.createdAt) > c.interval {
				delete(c.cacheMap, k)
			}
		}
		c.mu.Unlock()
	}
}