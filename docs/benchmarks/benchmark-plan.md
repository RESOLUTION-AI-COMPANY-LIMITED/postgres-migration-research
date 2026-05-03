# ⚠️ BENCHMARK PLAN - NOT ACTUAL BENCHMARKS

**This is a PLAN for benchmarking, NOT actual benchmark results.**

No benchmarks have been run yet. This document provides:
- ✅ Benchmark scenarios
- ✅ Metrics to measure
- ✅ Expected baselines

**To use this plan**: Run actual benchmarks after implementing PocketBase fork.

---

# Performance Benchmark Plan (AGENT-12)

**Agent**: AGENT-12 - Performance Benchmark Engineer  
**Date**: 2026-05-03  
**Project**: PocketBase v0.37.5 → PostgreSQL + Redis Migration

---

## 📋 Overview

Comprehensive performance benchmarking plan để measure SQLite vs PostgreSQL performance và Redis Pub/Sub latency. Đây là **RESEARCH PROJECT** - documentation ONLY.

### Success Criteria
- ✅ Database query benchmarks documented
- ✅ Redis Pub/Sub benchmarks documented
- ✅ Performance baselines established
- ✅ Optimization recommendations provided

---

## 🎯 Phase 4 - AGENT-12 Tasks

### Dependencies
- 🟡 AGENT-11 (Integration Tests) - Complete

### Estimated Time
6 hours (research documentation)

---

## 📊 BENCHMARK SUITE 1: Database Queries

### File: `benchmarks/db_bench_test.go`

```go
package benchmarks

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/free/postgresqlbaseapi/dbx"
)

// =============================================================================
// Benchmark 1: SELECT Query Performance
// =============================================================================

func BenchmarkDB_SELECT_SQLite(b *testing.B) {
	db, err := dbx.Open("sqlite3", "./test_data/benchmark.db")
	require.NoError(b, err)
	defer db.Close()
	
	// Setup: Create test table with 10k records
	setupBenchmarkData(b, db, 10000)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []User
		db.Select("*").From("users").Limit(100).All(&users)
	}
}

func BenchmarkDB_SELECT_PostgreSQL(b *testing.B) {
	db, err := dbx.Open("postgres", "postgres://postgres:password@localhost:5432/benchmarkdb?sslmode=disable")
	require.NoError(b, err)
	defer db.Close()
	
	// Setup: Create test table with 10k records
	setupBenchmarkData(b, db, 10000)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []User
		db.Select("*").From("users").Limit(100).All(&users)
	}
}

func BenchmarkDB_SELECT_MySQL(b *testing.B) {
	db, err := dbx.Open("mysql", "root:password@tcp(localhost:3306)/benchmarkdb")
	require.NoError(b, err)
	defer db.Close()
	
	// Setup: Create test table with 10k records
	setupBenchmarkData(b, db, 10000)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []User
		db.Select("*").From("users").Limit(100).All(&users)
	}
}

// =============================================================================
// Benchmark 2: INSERT Query Performance
// =============================================================================

func BenchmarkDB_INSERT_SQLite(b *testing.B) {
	db, err := dbx.Open("sqlite3", "./test_data/benchmark.db")
	require.NoError(b, err)
	defer db.Close()
	
	setupBenchmarkTable(b, db)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Insert("users", dbx.Params{
			"name":  "Benchmark User",
			"email": fmt.Sprintf("bench%d@example.com", i),
			"age":   30,
		}).Execute()
	}
}

func BenchmarkDB_INSERT_PostgreSQL(b *testing.B) {
	db, err := dbx.Open("postgres", "postgres://postgres:password@localhost:5432/benchmarkdb?sslmode=disable")
	require.NoError(b, err)
	defer db.Close()
	
	setupBenchmarkTable(b, db)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Insert("users", dbx.Params{
			"name":  "Benchmark User",
			"email": fmt.Sprintf("bench%d@example.com", i),
			"age":   30,
		}).Execute()
	}
}

// =============================================================================
// Benchmark 3: JOIN Query Performance
// =============================================================================

func BenchmarkDB_JOIN_SQLite(b *testing.B) {
	db, err := dbx.Open("sqlite3", "./test_data/benchmark.db")
	require.NoError(b, err)
	defer db.Close()
	
	setupBenchmarkRelations(b, db)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var results []JoinResult
		db.Select("users.*, orders.total").
			From("users").
			Join("orders", dbx.NewExp("users.id = orders.user_id")).
			Limit(100).
			All(&results)
	}
}

func BenchmarkDB_JOIN_PostgreSQL(b *testing.B) {
	db, err := dbx.Open("postgres", "postgres://postgres:password@localhost:5432/benchmarkdb?sslmode=disable")
	require.NoError(b, err)
	defer db.Close()
	
	setupBenchmarkRelations(b, db)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var results []JoinResult
		db.Select("users.*, orders.total").
			From("users").
			Join("orders", dbx.NewExp("users.id = orders.user_id")).
			Limit(100).
			All(&results)
	}
}

// =============================================================================
// Benchmark 4: Transaction Performance
// =============================================================================

func BenchmarkDB_Transaction_SQLite(b *testing.B) {
	db, err := dbx.Open("sqlite3", "./test_data/benchmark.db")
	require.NoError(b, err)
	defer db.Close()
	
	setupBenchmarkTable(b, db)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx, _ := db.Begin()
		
		for j := 0; j < 10; j++ {
			tx.Insert("users", dbx.Params{
				"name":  "Batch User",
				"email": fmt.Sprintf("batch%d_%d@example.com", i, j),
				"age":   25,
			}).Execute()
		}
		
		tx.Commit()
	}
}

func BenchmarkDB_Transaction_PostgreSQL(b *testing.B) {
	db, err := dbx.Open("postgres", "postgres://postgres:password@localhost:5432/benchmarkdb?sslmode=disable")
	require.NoError(b, err)
	defer db.Close()
	
	setupBenchmarkTable(b, db)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx, _ := db.Begin()
		
		for j := 0; j < 10; j++ {
			tx.Insert("users", dbx.Params{
				"name":  "Batch User",
				"email": fmt.Sprintf("batch%d_%d@example.com", i, j),
				"age":   25,
			}).Execute()
		}
		
		tx.Commit()
	}
}

// =============================================================================
// Helper Functions
// =============================================================================

type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Age   int    `db:"age"`
}

type JoinResult struct {
	User
	Total float64 `db:"total"`
}

func setupBenchmarkTable(b *testing.B, db *dbx.DB) {
	db.NewQuery("DROP TABLE IF EXISTS users").Execute()
	
	if db.DriverName() == "postgres" {
		db.NewQuery(`
			CREATE TABLE users (
				id SERIAL PRIMARY KEY,
				name TEXT,
				email TEXT UNIQUE,
				age INTEGER
			)
		`).Execute()
	} else if db.DriverName() == "sqlite3" {
		db.NewQuery(`
			CREATE TABLE users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT,
				email TEXT UNIQUE,
				age INTEGER
			)
		`).Execute()
	}
}

func setupBenchmarkData(b *testing.B, db *dbx.DB, count int) {
	setupBenchmarkTable(b, db)
	
	for i := 0; i < count; i++ {
		db.Insert("users", dbx.Params{
			"name":  fmt.Sprintf("User %d", i),
			"email": fmt.Sprintf("user%d@example.com", i),
			"age":   20 + (i % 50),
		}).Execute()
	}
}

func setupBenchmarkRelations(b *testing.B, db *dbx.DB) {
	// Create users table
	setupBenchmarkData(b, db, 1000)
	
	// Create orders table
	db.NewQuery("DROP TABLE IF EXISTS orders").Execute()
	
	if db.DriverName() == "postgres" {
		db.NewQuery(`
			CREATE TABLE orders (
				id SERIAL PRIMARY KEY,
				user_id INTEGER,
				total NUMERIC(10, 2)
			)
		`).Execute()
	} else {
		db.NewQuery(`
			CREATE TABLE orders (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER,
				total REAL
			)
		`).Execute()
	}
	
	// Insert orders
	for i := 0; i < 5000; i++ {
		db.Insert("orders", dbx.Params{
			"user_id": (i % 1000) + 1,
			"total":   100.0 + float64(i%500),
		}).Execute()
	}
}
```

---

## 📊 BENCHMARK SUITE 2: Redis Pub/Sub

### File: `benchmarks/redis_bench_test.go`

```go
package benchmarks

import (
	"context"
	"testing"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Benchmark 1: Redis Publish Performance
// =============================================================================

func BenchmarkRedis_Publish(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()
	
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Publish(ctx, "benchmark", "test message").Err()
	}
}

// =============================================================================
// Benchmark 2: Redis Publish + Subscribe Latency
// =============================================================================

func BenchmarkRedis_PubSub_Latency(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()
	
	ctx := context.Background()
	
	// Setup subscriber
	pubsub := client.Subscribe(ctx, "benchmark")
	defer pubsub.Close()
	
	ch := pubsub.Channel()
	
	// Start receiving goroutine
	received := make(chan struct{}, b.N)
	go func() {
		for range ch {
			received <- struct{}{}
		}
	}()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Publish message
		client.Publish(ctx, "benchmark", "test message")
		
		// Wait for receive (measure latency)
		<-received
	}
}

// =============================================================================
// Benchmark 3: Redis Cache Hit Performance
// =============================================================================

func BenchmarkRedis_CacheHit(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()
	
	ctx := context.Background()
	
	// Setup: Populate cache
	for i := 0; i < 1000; i++ {
		client.Set(ctx, fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i), 0)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Get(ctx, fmt.Sprintf("key_%d", i%1000)).Result()
	}
}

// =============================================================================
// Benchmark 4: Redis Cache Miss + DB Fallback
// =============================================================================

func BenchmarkRedis_CacheMiss_DBFallback(b *testing.B) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()
	
	db, err := dbx.Open("postgres", "postgres://postgres:password@localhost:5432/benchmarkdb?sslmode=disable")
	require.NoError(b, err)
	defer db.Close()
	
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("user_%d", i%1000)
		
		// Try cache
		cached, err := redisClient.Get(ctx, key).Result()
		if err == redis.Nil {
			// Cache miss - fetch from DB
			var user User
			db.Select("*").From("users").Where(dbx.HashExp{"id": i % 1000}).One(&user)
			
			// Store in cache
			redisClient.Set(ctx, key, user.Name, 0)
		} else {
			_ = cached // Cache hit
		}
	}
}

// =============================================================================
// Benchmark 5: In-Memory vs Redis Pub/Sub
// =============================================================================

func BenchmarkPubSub_InMemory(b *testing.B) {
	// Simulate in-memory pub/sub (Go channels)
	ch := make(chan string, 100)
	
	// Subscriber goroutine
	received := make(chan struct{}, b.N)
	go func() {
		for range ch {
			received <- struct{}{}
		}
	}()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch <- "test message"
		<-received
	}
}

func BenchmarkPubSub_Redis(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()
	
	ctx := context.Background()
	pubsub := client.Subscribe(ctx, "benchmark")
	defer pubsub.Close()
	
	ch := pubsub.Channel()
	
	// Subscriber goroutine
	received := make(chan struct{}, b.N)
	go func() {
		for range ch {
			received <- struct{}{}
		}
	}()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Publish(ctx, "benchmark", "test message")
		<-received
	}
}
```

---

## 📈 Expected Performance Baselines

### Database Query Benchmarks

| Operation | SQLite | PostgreSQL | MySQL | Notes |
|-----------|--------|------------|-------|-------|
| **SELECT (100 rows)** | 50-100 µs | 80-150 µs | 80-150 µs | PostgreSQL slightly slower (network overhead) |
| **INSERT (single)** | 10-20 µs | 50-100 µs | 50-100 µs | SQLite faster (no network) |
| **JOIN (100 rows)** | 200-300 µs | 150-250 µs | 150-250 µs | PostgreSQL faster (better optimizer) |
| **Transaction (10 INSERTs)** | 50-100 µs | 200-400 µs | 200-400 µs | SQLite faster (no network) |

**Key Insights:**
- ✅ **SQLite faster for single-node**: No network latency
- ✅ **PostgreSQL faster for JOINs**: Better query optimizer
- ✅ **Trade-off acceptable**: +2-5x latency for clustering capability

### Redis Pub/Sub Benchmarks

| Operation | In-Memory | Redis | Notes |
|-----------|-----------|-------|-------|
| **Publish** | 1-2 µs | 50-100 µs | Redis has network overhead |
| **Pub+Sub Latency** | 5-10 µs | 100-200 µs | Redis adds ~100 µs latency |
| **Cache Hit** | N/A | 50-100 µs | Very fast |
| **Cache Miss + DB** | N/A | 150-300 µs | DB query + cache write |

**Key Insights:**
- ✅ **Redis Pub/Sub adds ~100 µs latency**: Acceptable for realtime use
- ✅ **Cache hit performance**: Very fast (~50 µs)
- ✅ **Trade-off acceptable**: ~10x latency for multi-node capability

---

## 🔧 Running Benchmarks

### Commands

```bash
# Run all benchmarks
go test ./benchmarks/... -bench=. -benchmem -benchtime=10s

# Run specific benchmark
go test ./benchmarks/... -bench=BenchmarkDB_SELECT -benchmem

# Run with CPU profiling
go test ./benchmarks/... -bench=. -cpuprofile=cpu.prof

# Run with memory profiling
go test ./benchmarks/... -bench=. -memprofile=mem.prof

# Compare benchmarks (before/after)
go test ./benchmarks/... -bench=. -benchmem > before.txt
# (make changes)
go test ./benchmarks/... -bench=. -benchmem > after.txt
benchcmp before.txt after.txt
```

### Expected Output Format

```
BenchmarkDB_SELECT_SQLite-8             20000    75000 ns/op    2048 B/op    10 allocs/op
BenchmarkDB_SELECT_PostgreSQL-8         10000   120000 ns/op    2048 B/op    10 allocs/op
BenchmarkDB_SELECT_MySQL-8              10000   115000 ns/op    2048 B/op    10 allocs/op

BenchmarkDB_INSERT_SQLite-8            100000    15000 ns/op     512 B/op     5 allocs/op
BenchmarkDB_INSERT_PostgreSQL-8         20000    75000 ns/op     512 B/op     5 allocs/op

BenchmarkRedis_Publish-8                20000    75000 ns/op     256 B/op     3 allocs/op
BenchmarkRedis_PubSub_Latency-8         10000   150000 ns/op     512 B/op     5 allocs/op
```

---

## 🎯 Optimization Recommendations

### Database Optimizations

**1. Connection Pool Tuning**
```go
db.DB().SetMaxOpenConns(100)      // Increase for high load
db.DB().SetMaxIdleConns(25)       // Keep more idle connections
db.DB().SetConnMaxLifetime(time.Hour) // Recycle connections
```

**2. PostgreSQL Indexes**
```sql
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_orders_user_id ON orders(user_id);
```

**3. Query Optimization**
```go
// BAD: N+1 query problem
for _, user := range users {
    db.Select("*").From("orders").Where(dbx.HashExp{"user_id": user.ID}).All(&orders)
}

// GOOD: Single JOIN query
db.Select("users.*, orders.*").
    From("users").
    Join("orders", dbx.NewExp("users.id = orders.user_id")).
    All(&results)
```

### Redis Optimizations

**1. Connection Pool**
```go
redis.NewClient(&redis.Options{
    Addr:         "localhost:6379",
    PoolSize:     100,   // Increase pool size
    MinIdleConns: 10,    // Keep idle connections ready
})
```

**2. Pub/Sub vs Streams**
```go
// For high-volume events, consider Redis Streams instead of Pub/Sub
// Streams provide persistence and delivery guarantees
client.XAdd(ctx, &redis.XAddArgs{
    Stream: "events",
    Values: map[string]any{"event": "create"},
})
```

**3. Cache Expiry**
```go
// Set reasonable TTL to avoid memory bloat
client.Set(ctx, key, value, 1*time.Hour)
```

---

## 📁 Deliverables Summary

### Benchmark Files
1. ✅ `benchmarks/db_bench_test.go` - Database query benchmarks
2. ✅ `benchmarks/redis_bench_test.go` - Redis Pub/Sub benchmarks

### Performance Reports
1. `docs/performance-report.md` - Benchmark results analysis
2. `docs/optimization-recommendations.md` - Performance tuning guide

---

## ✅ Verification Checklist

**Before marking AGENT-12 complete:**

- [ ] Database benchmarks documented (SELECT, INSERT, JOIN, Transaction)
- [ ] Redis benchmarks documented (Publish, Pub/Sub, Cache)
- [ ] Expected baselines documented
- [ ] Benchmark execution commands documented
- [ ] Optimization recommendations provided (DB + Redis)
- [ ] Performance comparison (SQLite vs PostgreSQL vs MySQL)
- [ ] Performance report generated

---

**AGENT-12 Status:** 🟢 Ready for Review  
**Document Version:** 1.0  
**Last Updated:** 2026-05-03
