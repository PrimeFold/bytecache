# ByteCache

A high-performance, thread-safe, in-memory key-value database written in Go, featuring raw TCP wire protocol parsing, active background TTL eviction, and zero-allocation hot paths.

## 📊 Sharded Performance Benchmarks

Benchmarked on an **Intel® Core™ Ultra 5 125H (18 cores)** under Linux/WSL:

| Benchmark | Operations | Latency | Memory / Op | Allocs / Op |
| :--- | :--- | :--- | :--- | :--- |
| **`BenchmarkCacheGet` (Reads)** | **56,932,548** | **19.72 ns/op** | **0 B/op** | **0 allocs/op** |
| **`BenchmarkCacheSet` (Writes)** | **41,421,339** | **29.56 ns/op** | **0 B/op** | **0 allocs/op** |
| **`BenchmarkCacheConcurrency`** | **14,578,346** | **75.57 ns/op** | **0 B/op** | **0 allocs/op** |

* **32-Shard Lock Striping:** Slices the in-memory map into 32 partitions via FNV-1a hashing, eliminating mutex lock contention under high multi-core concurrency.
* **Parallel Active TTL Eviction:** Background eviction worker sweeps all 32 shards concurrently using `sync.WaitGroup` worker routines on ticker intervals.
* **Zero Allocations:** Hot paths run at **0 B/op and 0 allocs/op**, sustaining **>56M reads/sec** and **>41M writes/sec**.


## Features

- Thread-Safe Storage Engine: Sharded read/write locks (sync.RWMutex) maximizing concurrent read throughput.
- Active Background TTL Eviction: Dedicated daemon worker cleans expired keys periodically without blocking reads.
- Raw TCP Wire Protocol: Fast line-delimited socket parser supporting SET, GET, and DEL.
- Production CLI: Built with spf13/cobra supporting custom ports and configurations.


## Quick Start
```
Build and Run
bash

```

### Clone
```
git clone https://github.com/PrimeFold/bytecache.git
cd bytecache
```

### Start the server on port 6379 (Redis default)
```
go run . start --port 6379

```

## Connect via TCP
```
bash
nc localhost 6379
SET user:101 Aditya
OK
GET user:101
Aditya
SET token secret123 60s
OK
DEL user:101
OK
```


## Running Tests & Benchmarks
```
bash
```
## Run unit tests and race detection
```
go test -v -race ./...
```
## Run performance benchmarks
```
go test -bench=. -benchmem -run=^$ .
```
### License
Apache License 2.0
