package go_memkvstore

import (
	"sync"
	"time"
)

type Item[V any] struct {
	Value      V
	Expiration time.Time
}

func (m Item[V]) IsExpired(now time.Time) bool {
	return m.Expiration.Before(now)
}

type Store[V any] struct {
	sync.RWMutex
	Expiration time.Duration
	Store      map[string]Item[V]
}

// New creates a new MemoryKV with default expiration
func New[V any](expiration time.Duration) *Store[V] {
	return &Store[V]{
		Store:      make(map[string]Item[V]),
		Expiration: expiration,
	}
}

// Get returns an item from the memory store and a boolean indicating if the item was found
// if the item does not exist or is expired the result will be false
func (m *Store[V]) Get(key string, noValue V) (V, bool) {
	now := time.Now().UTC()
	m.RLock()
	defer m.RUnlock()
	item, ok := m.Store[key]
	if !ok {
		return noValue, false
	}
	if item.IsExpired(now) {
		return noValue, false
	}
	return item.Value, true
}

func (m *Store[V]) GetWithExpiration(key string, noValue V) (V, time.Time, bool) {
	now := time.Now().UTC()
	m.RLock()
	defer m.RUnlock()
	item, ok := m.Store[key]
	if !ok {
		return noValue, time.Time{}, false
	}
	if item.IsExpired(now) {
		return noValue, time.Time{}, false
	}
	return item.Value, item.Expiration, true
}

// Set sets an item in the memory store with default expiration
func (m *Store[V]) Set(key string, value V) {
	now := time.Now().UTC()
	m.Lock()
	defer m.Unlock()
	m.Store[key] = Item[V]{
		Value:      value,
		Expiration: now.Add(m.Expiration),
	}
}

// SetWithExpiration sets an item with specific expiration in the memory store
func (m *Store[V]) SetWithExpiration(key string, value V, expiration time.Duration) {
	expires := time.Now().UTC().Add(expiration)
	m.Lock()
	defer m.Unlock()
	m.Store[key] = Item[V]{
		Value:      value,
		Expiration: expires,
	}
}

// Delete removes an item from the memory store
func (m *Store[V]) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.Store, key)
}

// Cleanup removes expired items from the memory store
func (m *Store[V]) Cleanup() {
	now := time.Now().UTC()
	m.Lock()
	defer m.Unlock()
	for key, item := range m.Store {
		if item.IsExpired(now) {
			delete(m.Store, key)
		}
	}
}
