package cache

import (
	"sync"
	"time"
)

type Cache interface {
	Set(key string, value any, ttl time.Duration) error
	Get(key string) (any, bool)
}

type cacheItem struct {
	value     any
	expiresAt time.Time
}

type inMemoryCache struct {
	items map[string]*cacheItem
	mutex sync.RWMutex
}

func NewInMemoryCache(defaultTTL time.Duration) Cache {
	cache := &inMemoryCache{
		items: make(map[string]*cacheItem),
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

func (c *inMemoryCache) Set(key string, value any, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = &cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}

	return nil
}

func (c *inMemoryCache) Get(key string) (any, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.expiresAt) {
		delete(c.items, key)
		return nil, false
	}

	return item.value, true
}

func (c *inMemoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.expiresAt) {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}
