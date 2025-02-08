package internal

import (
	"sync"
	"time"
)

type Cache struct {
  entries map[string]cacheEntry
  mutex *sync.Mutex
}

type cacheEntry struct{
  createdAt time.Time
  val []byte
}

func NewCache(interval time.Duration) *Cache {
  return &Cache{
    entries: make(map[string]cacheEntry),
    mutex: &sync.Mutex{},
  }
}

func (c *Cache) Get(key string, val []byte) ([]byte, bool) {
  c.mutex.Lock()
  defer c.mutex.Unlock()

  entry, exists := c.entries[key]
  if !exists {
    return nil, false
  }
  return entry.val, true
}

func (c *Cache) Add(key string, val []byte) {
  c.mutex.Lock()
  defer c.mutex.Unlock()

  c.entries[key] = cacheEntry{
    createdAt: time.Now(),
    val: val,
  }
}
