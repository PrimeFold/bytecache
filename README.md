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
