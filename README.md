# ByteCache

A high-performance, thread-safe, in-memory key-value database written in Go, featuring raw TCP wire protocol parsing, active background TTL eviction, and zero-allocation hot paths.

## Benchmarks

Benchmarked on Intel Core Ultra 5 125H across 18 cores:

```text
BenchmarkCacheSet-18          37809426   31.89 ns/op   0 B/op   0 allocs/op
BenchmarkCacheGet-18          59172042   20.40 ns/op   0 B/op   0 allocs/op
BenchmarkCacheConcurrency-18  16154098   75.07 ns/op   0 B/op   0 allocs/op

Read throughput: ~49M ops/sec (20.4 ns/op) with zero heap allocations.
Write throughput: ~31M ops/sec (31.8 ns/op) with zero heap allocations.
Data race safety: 100% verified with go test -race.

```
## Features

- Thread-Safe Storage Engine: Sharded read/write locks (sync.RWMutex) maximizing concurrent read throughput.
- Active Background TTL Eviction: Dedicated daemon worker cleans expired keys periodically without blocking reads.
- Raw TCP Wire Protocol: Fast line-delimited socket parser supporting SET, GET, and DEL.
- Production CLI: Built with spf13/cobra supporting custom ports and configurations.


## Quick Start
```
Build and Run
bash
# Clone
git clone https://github.com/PrimeFold/bytecache.git
cd bytecache
# Start the server on port 6379 (Redis default)
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
# Run unit tests and race detection
```
go test -v -race ./...
```
# Run performance benchmarks
```
go test -bench=. -benchmem -run=^$ .
```
## License
Apache License 2.0
