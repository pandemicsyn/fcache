fcache
======

It's a simple cache with TTL support that does what *I need*. It probably doesn't do what you need, but thats okay...


[![GoDoc](https://godoc.org/github.com/pandemicsyn/fcache?status.svg)](https://godoc.org/github.com/pandemicsyn/fcache)

benchmark
=========

Mid 2012 Macbook Pro Retina 15 - 2.3 GHz Intel Core i7 (4 cores) - 16GB Ram

```
PASS
BenchmarkFCacheGet  30000000            46.8 ns/op         0 B/op          0 allocs/op
BenchmarkFCacheConcurrentGet    30000000            47.8 ns/op         0 B/op          0 allocs/op
BenchmarkFCacheIncrementInt 10000000           120 ns/op           8 B/op          1 allocs/op
BenchmarkFCacheIncrementFloat64  5000000           293 ns/op          32 B/op          2 allocs/op
BenchmarkFCacheSet   5000000           301 ns/op          48 B/op          2 allocs/op
BenchmarkRWMutexMapSet  20000000            67.6 ns/op         0 B/op          0 allocs/op
BenchmarkFCacheSetDelete     5000000           389 ns/op          48 B/op          2 allocs/op
BenchmarkRWMutexMapSetDelete    10000000           144 ns/op           0 B/op          0 allocs/op
BenchmarkFCacheConcurrentSet     5000000           326 ns/op          48 B/op          2 allocs/op
BenchmarkFCacheSetGrowing    1000000          1192 ns/op         173 B/op          3 allocs/op
BenchmarkFCacheSetDeleteGrowing  2000000           680 ns/op          56 B/op          4 allocs/op
ok      github.com/pandemicsyn/fcache   18.457s
```

supporting cast
===============
Inspired by go-cache and others (that extra stuff or where missing bits I wanted)
