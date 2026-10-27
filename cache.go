package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
)

type item struct{
	value []byte
	expiresAt time.Time 
}


type Cache struct{
	mu sync.RWMutex
	data map[string]item
}

func NewCache() *Cache{
	return &Cache{
		data: make(map[string]item),
	}
}

func (c *Cache) Set(key string , value []byte , ttl time.Duration){
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiresAt time.Time

	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	c.data[key] = item{
		value: value,
		expiresAt: expiresAt,
	}
	
}


func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	sanitizedKey := strings.TrimSpace(key)
	if sanitizedKey == "" {
		return nil, false
	}

	item, exists := c.data[sanitizedKey]
	if !exists {
		return nil, false
	}

	// Check TTL expiration
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		return nil, false
	}

	return item.value, true
}



func (c *Cache) Delete(key string){
	c.mu.Lock()
	defer c.mu.Unlock()

	sanitizedKey:= strings.TrimSpace(key);

	if key == ""{
		log.Printf("Invalid key  : [%s]",key)
	}

	delete(c.data,sanitizedKey)	
}


func (c *Cache) evictExpired(){
	c.mu.Lock()
	defer c.mu.Unlock()
	for k,v := range c.data{
		if !v.expiresAt.IsZero() && time.Now().After(v.expiresAt){
			delete(c.data,k)
		}
	}

}

// StartEvictionWorker runs a background loop that actively deletes expired keys
func (c *Cache) StartEvictionWorker(ctx context.Context,interval time.Duration){
	ticker:= time.NewTicker(interval)

	defer ticker.Stop()

	for{
		select{
		case<-ctx.Done():
			return

		case <-ticker.C:
			c.evictExpired()
		}
		
	}


}


func (c *Cache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.data)
}

