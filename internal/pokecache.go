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
  newCache := &Cache{
    entries: make(map[string]cacheEntry),
    mutex: &sync.Mutex{},
  }
  go newCache.reapLoop(interval)
  return newCache
}


func (c *Cache) Get(key string) ([]byte, bool) {
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

func (c *Cache) reapLoop(interval time.Duration) {
  ticker := time.NewTicker(interval)
  defer ticker.Stop()

  for range ticker.C {
    c.mutex.Lock()
      for key, entry := range c.entries {
        elapsedTime := time.Since(entry.createdAt)
        if elapsedTime > interval {
          delete(c.entries, key)
        }
      }
      c.mutex.Unlock()
    }
}