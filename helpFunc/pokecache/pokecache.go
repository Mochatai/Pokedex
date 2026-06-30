package pokecache

import (
	"sync"
	"time"
)

type chacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu       sync.Mutex
	data     map[string]chacheEntry
	interval time.Duration
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = chacheEntry{time.Now(), val}

}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if val, ok := c.data[key]; ok {
		return val.val, true
	}

	return nil, false
}

func (c *Cache) reapLoop(inte time.Duration) {
	tiker := time.NewTicker(c.interval)

	for {
		select {
		case <-tiker.C:
			c.clearCache()
		}
	}

}

func (c *Cache) clearCache() {
	currTime := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, _ := range c.data {
		if currTime.After(c.data[key].createdAt) {
			delete(c.data, key)
		}

	}
}

func NewCache(inte time.Duration) *Cache {
	cache := &Cache{
		data:     make(map[string]chacheEntry),
		interval: inte,
	}

	go cache.reapLoop(inte)

	return cache

}
