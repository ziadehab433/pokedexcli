package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu sync.Mutex
	cm map[string]CacheEntry
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	newC := Cache{
		cm: make(map[string]CacheEntry),
	}

	ticker := time.NewTicker(interval)
	go func() {
		for {
			<-ticker.C
			newC.reapLoop(interval)
		}
	}()

	return &newC
}

func (c *Cache) Add(k string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := CacheEntry{
		createdAt: time.Now(),
		val:       value,
	}

	c.cm[k] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, e := range c.cm {
		if key == k {
			return e.val, true
		}
	}

	return []byte(""), false
}

func (c *Cache) reapLoop(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, e := range c.cm {
		if (time.Now().UnixMilli() - e.createdAt.UnixMilli()) > interval.Milliseconds() {
			delete(c.cm, k)
		}
	}
}
