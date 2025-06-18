package session

import (
	"context"
	"sync"
	"time"

	"billsplitter-monolith/internal/domain/auth"
)

type Cache struct {
	// key - session id
	kv map[string]cacheEntry
	mu sync.RWMutex
}

type cacheEntry struct {
	value    *auth.Session
	expireAt time.Time
}

func NewMemCache() *Cache {
	return &Cache{
		kv: make(map[string]cacheEntry),
	}
}

func (c *Cache) Set(ctx context.Context, sessionID string, value *auth.Session, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.kv[sessionID] = cacheEntry{
		value:    value,
		expireAt: time.Now().Add(ttl),
	}

	return nil
}

func (c *Cache) Get(ctx context.Context, sessionID string) (*auth.Session, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.kv[sessionID]
	if !ok {
		return nil, nil
	}

	if time.Now().After(entry.expireAt) {
		return nil, nil
	}

	return entry.value, nil
}
