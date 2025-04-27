package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Data      []byte
	Timestamp time.Time
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]map[string]CacheItem // propertyID -> (subject+params) -> data
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]map[string]CacheItem),
	}
}

func (c *Cache) Get(propertyID, key string) (CacheItem, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	subMap, ok := c.items[propertyID]
	if !ok {
		return CacheItem{}, false
	}
	item, ok := subMap[key]
	return item, ok
}

func (c *Cache) Set(propertyID, key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.items[propertyID]; !ok {
		c.items[propertyID] = make(map[string]CacheItem)
	}
	c.items[propertyID][key] = CacheItem{
		Data:      data,
		Timestamp: time.Now(),
	}
}

func (c *Cache) Invalidate(propertyID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, propertyID)
}
