package pokecache

import (
	"sync"
	"time"
//	"fmt"
)

//
type Cache struct {
	entries map[string]cacheEntry // map of cachEntries
	mu sync.Mutex // protect the map across goroutines
	duration time.Duration
}

type cacheEntry struct {
	createdAt time.Time // represents when the entry was created 
	val []byte // represents the raw data we're caching
}

func NewCache(interval time.Duration) *Cache {
	// initialise a map and mutex
	// store interval
	//fmt.Print("Creating new cache")
	c := Cache{
		entries : make(map[string]cacheEntry),
		duration : interval,
	}
	// start reap loop
	//fmt.Print("starting reap loop")
	go c.reapLoop()	
	// return pointer to new cache
	return &c

}

func (c* Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if entry, exists := c.entries[key]; exists {
		return entry.val, true
	}

	return nil, false
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry {
		createdAt: time.Now(),
		val: val,
	}
	//fmt.Println("Added!")
	return
}

func (c *Cache) reapLoop() {
	// should remove any entries that are older than the interval
	// loop through all entries
	//if time.now - entry.createdTime >= duration, then delete entry
	ticker := time.NewTicker(c.duration)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			//fmt.Println("tick at", t)
			c.mu.Lock()
			for k,entry := range c.entries {
				if time.Since(entry.createdAt) >= c.duration {
					//fmt.Println("Delete!")
					delete(c.entries, k)
				}
			}
			c.mu.Unlock()
		}
	}
	return
}
