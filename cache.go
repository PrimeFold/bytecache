package main

import (
	"context"
	"hash/fnv"
	"log"
	"strings"
	"sync"
	"time"
)

type item struct{
	value []byte
	expiresAt time.Time 
}


type cacheShard struct{
	mu sync.RWMutex
	data map[string]item
}


const shardCount = 32

type cacheShardMap struct{
	shards [shardCount]*cacheShard
	maxCount int
	windowSec time.Duration
}

func NewShardMap(maxCount int , windowsSec int) *cacheShardMap{
	shardMap:= &cacheShardMap{
		maxCount: maxCount,
		windowSec: time.Duration(windowsSec)*time.Second,
	}

	for i := 0; i < shardCount; i++ {
		shardMap.shards[i] = &cacheShard{
			data: make(map[string]item),
		}
	}

	return shardMap

}


func (c *cacheShardMap) getShard(key string) *cacheShard{
	h:= fnv.New32a()
	h.Write([]byte(key))

	shardIndex:= h.Sum32() & (shardCount - 1)
	return c.shards[shardIndex]
}

func (c *cacheShardMap) Set(key string , value []byte , ttl time.Duration){
	
	shard:= c.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	var expiresAt time.Time

	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	shard.data[key] = item{
		value: value,
		expiresAt: expiresAt,
	}
	
}


func (c *cacheShardMap) Get(key string) ([]byte, bool) {
	
	shard:= c.getShard(key)

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	sanitizedKey := strings.TrimSpace(key)
	if sanitizedKey == "" {
		return nil, false
	}

	item, exists := shard.data[sanitizedKey]
	if !exists {
		return nil, false
	}

	// Check TTL expiration
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		return nil, false
	}

	return item.value, true
}



func (c *cacheShardMap) Delete(key string){

	shard:= c.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	sanitizedKey:= strings.TrimSpace(key);

	if key == ""{
		log.Printf("Invalid key  : [%s]",key)
	}

	delete(shard.data,sanitizedKey)	
}


func (c *cacheShardMap) evictExpired(){
	now:= time.Now()

	var wg sync.WaitGroup

	for i := 0; i < shardCount; i++ {
		wg.Add(1)
		go func (shard *cacheShard)  {
			defer wg.Done()

			shard.mu.Lock()
			defer shard.mu.Unlock()

			for key, item := range shard.data {
				// If expiresAt is set and 'now' is past the expiration time:
				if !item.expiresAt.IsZero() && now.After(item.expiresAt) {
					delete(shard.data, key)
				}
			}
			
		}(c.shards[i])
	}
	wg.Wait()
}

// StartEvictionWorker runs a background loop that actively deletes expired keys
func (c *cacheShardMap) StartEvictionWorker(ctx context.Context,interval time.Duration){
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


func (c *cacheShardMap) Count() int {
	total:= 0

	for i := 0; i < shardCount; i++ {
		c.shards[i].mu.RLock()
		total = total + len(c.shards[i].data)
		c.shards[i].mu.Unlock()
	}

	return total
}

