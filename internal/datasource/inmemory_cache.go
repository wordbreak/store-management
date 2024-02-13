package datasource

import (
	"sync"
	"time"
)

type Item struct {
	value      interface{}
	expiration int64
}

type InMemoryCache struct {
	data sync.Map
	mu   sync.RWMutex
}

func NewInMemoryCache() Cache {
	return &InMemoryCache{
		data: sync.Map{},
		mu:   sync.RWMutex{},
	}
}

func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data.Store(key, Item{
		value:      value,
		expiration: time.Now().Add(ttl).UnixNano(),
	})
	return nil
}

func (c *InMemoryCache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data.Load(key)
	if !ok {
		return nil, nil
	}

	if item.(Item).expiration > time.Now().UnixNano() {
		return item.(Item).value, nil
	}

	c.mu.RUnlock()
	c.mu.Lock()
	c.data.Delete(key)
	c.mu.Unlock()
	c.mu.RLock()

	return nil, nil
}

func (c *InMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data.Range(func(key, value interface{}) bool {
		c.data.Delete(key)
		return true
	})
}

func (c *InMemoryCache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data.Delete(key)
}
