package cache

import (
	"sync"
	"wb-tech/internal/model"
)

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]*model.Order
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]*model.Order),
	}
}

func (m *MemoryCache) Get(uid string) *model.Order {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.items[uid]
}

func (m *MemoryCache) Set(uid string, order *model.Order) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[uid] = order
}
