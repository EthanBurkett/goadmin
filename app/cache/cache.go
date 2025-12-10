package cache

import (
	"sync"
	"time"
)

// CacheItem represents a cached value with expiration
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// Cache is a simple in-memory cache with TTL support
type Cache struct {
	items map[string]*CacheItem
	mu    sync.RWMutex
}

// GlobalCache is the global cache instance
var GlobalCache *Cache

// Init initializes the global cache
func Init() {
	if GlobalCache != nil {
		return // Already initialized
	}

	GlobalCache = &Cache{
		items: make(map[string]*CacheItem),
	}

	// Start cleanup goroutine
	go GlobalCache.cleanup()
}

// cleanup periodically removes expired items
func (c *Cache) cleanup() {
	if c == nil {
		return
	}

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.Expiration) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	if c == nil {
		return nil, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(item.Expiration) {
		return nil, false
	}

	return item.Value, true
}

// Set stores a value in the cache with TTL
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &CacheItem{
		Value:      value,
		Expiration: time.Now().Add(ttl),
	}
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*CacheItem)
}

// InvalidatePattern removes all items with keys matching a pattern
func (c *Cache) InvalidatePattern(pattern string) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.items {
		// Simple pattern matching - contains check
		if len(pattern) > 0 && containsPattern(key, pattern) {
			delete(c.items, key)
		}
	}
}

// containsPattern checks if the key matches the pattern
// Simple implementation - can be enhanced with regex if needed
func containsPattern(key, pattern string) bool {
	return len(key) >= len(pattern) && key[:len(pattern)] == pattern
}

// Size returns the number of items in the cache
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}
