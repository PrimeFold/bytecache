# ByteCache

A high-performance, thread-safe, in-memory key-value database written in Go, featuring raw TCP wire protocol parsing, active background TTL eviction, and zero-allocation hot paths.

## Benchmarks

Benchmarked on Intel Core Ultra 5 125H across 18 cores:

```text
BenchmarkCacheSet-18          37809426   31.89 ns/op   0 B/op   0 allocs/op
BenchmarkCacheGet-18          59172042   20.40 ns/op   0 B/op   0 allocs/op
BenchmarkCacheConcurrency-18  16154098   75.07 ns/op   0 B/op   0 allocs/op
