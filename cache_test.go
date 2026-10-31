package main

import (
	"context"
	"testing"
	"time"
)


const maxCount = 32
const windowsSec = 60
func TestActiveEvictionWorker(t *testing.T){
	
	cache:= NewShardMap(maxCount,windowsSec)

	ctx,cancel:= context.WithCancel(context.Background())
	defer cancel()


	//The eviction worker SHOULD remove the cache set for just 50 ms , meaning that worker would run
	// around 5 times in order to remove the key on its expiry..
	go cache.StartEvictionWorker(ctx,10*time.Millisecond)


	//Setting the Key for 50ms
	cache.Set("her",[]byte("impossible"),50*time.Millisecond)


	//Just a check to make sure the count is 1
	if cache.Count() != 1 {
		t.Fatalf("expected count 1, got %d", cache.Count())
	}

	//Waiting for 200 ms..
	time.Sleep(200*time.Millisecond)


	if count:= cache.Count() ; count != 0{
		t.Errorf("expected background worker to evict expired key, but count is still %d", count)
	}

}


func BenchmarkCacheSet(b *testing.B){
	c:= NewShardMap(maxCount,windowsSec)
	val:= []byte("benchmark-value")
	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop(){
		c.Set("bench:key",val,0)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	c := NewShardMap(maxCount,windowsSec)
	key := "bench:key"
	c.Set(key, []byte("benchmark-value"), 0)

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, _ = c.Get(key)
	}
}

func BenchmarkCacheConcurrency(b *testing.B){
	c := NewShardMap(maxCount,windowsSec)
	c.Set("key:read",[]byte("read-value"),0)
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB){
		for pb.Next(){
			_,_=c.Get("key:read")
		}
	})

}



