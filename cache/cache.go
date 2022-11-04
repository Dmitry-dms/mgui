package cache

import "sync"

type RamCache[T any] struct {
	m  map[string]T
	mu sync.RWMutex
}

func NewRamCache[T any]() *RamCache[T] {
	c := RamCache[T]{
		mu: sync.RWMutex{},
		m:  make(map[string]T, 100),
	}
	return &c
}

func (c *RamCache[T]) Add(key string, val T) bool {
	_, ok := c.Get(key)
	if ok {
		return false
	}
	c.m[key] = val
	return true
}

func (c *RamCache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[key]
	return v, ok
}
