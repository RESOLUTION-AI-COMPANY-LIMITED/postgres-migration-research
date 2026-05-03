# Unit Test Plan (AGENT-10)

**Agent**: AGENT-10 - Unit Test Engineer  
**Date**: 2026-05-03  
**Project**: PocketBase v0.37.5 → PostgreSQL + Redis Migration

---

## 📋 Overview

Comprehensive unit testing plan cho tất cả components mới. Đây là **RESEARCH PROJECT** - documentation ONLY.

### Success Criteria
- ✅ Test cases cho `db_postgresql.go` documented
- ✅ Test cases cho Redis integration documented
- ✅ Test cases cho dbx query builder documented
- ✅ Coverage requirements defined (>80%)

---

## 🎯 Phase 4 - AGENT-10 Tasks

### Dependencies
- 🟡 AGENT-05 (DB Connection) - Complete
- 🟡 AGENT-06 (Redis) - Complete
- 🟡 AGENT-08 (dbx) - Complete

### Estimated Time
8 hours (research documentation)

---

## 🧪 TEST SUITE 1: Database Connection Tests

### File: `core/db_postgresql_test.go`

```go
package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Test Group 1: DSN Parsing & Driver Detection
// =============================================================================

func TestConnectDB_PostgreSQL_WithPrefix(t *testing.T) {
	// Test: postgres:// prefix detected
	dsn := "postgres://user:pass@localhost:5432/testdb?sslmode=disable"
	db, err := connectDB(dsn)
	
	require.NoError(t, err)
	require.NotNil(t, db)
	assert.Equal(t, "postgres", db.DriverName())
	
	defer db.Close()
}

func TestConnectDB_PostgreSQL_AlternativePrefix(t *testing.T) {
	// Test: postgresql:// prefix (alias) detected
	dsn := "postgresql://user:pass@localhost:5432/testdb"
	db, err := connectDB(dsn)
	
	require.NoError(t, err)
	assert.Equal(t, "postgres", db.DriverName())
	
	defer db.Close()
}

func TestConnectDB_MySQL_WithPrefix(t *testing.T) {
	// Test: mysql:// prefix detected
	dsn := "mysql://root:pass@tcp(localhost:3306)/testdb"
	db, err := connectDB(dsn)
	
	require.NoError(t, err)
	assert.Equal(t, "mysql", db.DriverName())
	
	defer db.Close()
}

func TestConnectDB_PostgreSQL_NoPrefix(t *testing.T) {
	// Test: No prefix defaults to PostgreSQL
	dsn := "host=localhost port=5432 user=postgres dbname=testdb sslmode=disable"
	db, err := connectDB(dsn)
	
	require.NoError(t, err)
	assert.Equal(t, "postgres", db.DriverName())
	
	defer db.Close()
}

// =============================================================================
// Test Group 2: Connection Error Handling
// =============================================================================

func TestConnectDB_InvalidDSN(t *testing.T) {
	// Test: Invalid DSN format
	dsn := "invalid://malformed-dsn"
	db, err := connectDB(dsn)
	
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestConnectDB_ConnectionRefused(t *testing.T) {
	// Test: Database not running (connection refused)
	dsn := "postgres://localhost:9999/testdb"  // Wrong port
	db, err := connectDB(dsn)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
}

func TestConnectDB_AuthenticationFailed(t *testing.T) {
	// Test: Wrong credentials
	dsn := "postgres://wronguser:wrongpass@localhost:5432/testdb"
	db, err := connectDB(dsn)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication failed")
}

func TestConnectDB_DatabaseNotExist(t *testing.T) {
	// Test: Database does not exist
	dsn := "postgres://postgres:password@localhost:5432/nonexistent_db"
	db, err := connectDB(dsn)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

// =============================================================================
// Test Group 3: Connection Pool Configuration
// =============================================================================

func TestConnectDB_ConnectionPool_Defaults(t *testing.T) {
	// Test: Connection pool settings applied
	dsn := "postgres://localhost/testdb?sslmode=disable"
	db, err := connectDB(dsn)
	require.NoError(t, err)
	defer db.Close()
	
	// Verify pool settings (if configured)
	stats := db.DB().Stats()
	assert.GreaterOrEqual(t, stats.MaxOpenConnections, 1)
}

// =============================================================================
// Test Group 4: Query Execution
// =============================================================================

func TestConnectDB_SimpleQuery_PostgreSQL(t *testing.T) {
	// Test: Execute simple query on PostgreSQL
	dsn := "postgres://postgres:password@localhost:5432/testdb?sslmode=disable"
	db, err := connectDB(dsn)
	require.NoError(t, err)
	defer db.Close()
	
	var version string
	err = db.NewQuery("SELECT version()").Row(&version)
	assert.NoError(t, err)
	assert.Contains(t, version, "PostgreSQL")
}

func TestConnectDB_SimpleQuery_MySQL(t *testing.T) {
	// Test: Execute simple query on MySQL
	dsn := "mysql://root:password@tcp(localhost:3306)/testdb"
	db, err := connectDB(dsn)
	require.NoError(t, err)
	defer db.Close()
	
	var version string
	err = db.NewQuery("SELECT VERSION()").Row(&version)
	assert.NoError(t, err)
	assert.NotEmpty(t, version)
}
```

### Coverage Target
- **Target:** >90% code coverage
- **Critical Paths:** All driver detection branches, error handling

---

## 🧪 TEST SUITE 2: Redis Integration Tests

### File: `core/redis_test.go`

```go
package core

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Test Group 1: Redis Initialization
// =============================================================================

func TestInitRedis_NoDSN_SingleNodeMode(t *testing.T) {
	// Test: No Redis DSN → single-node mode (no error)
	app := NewBaseApp(BaseAppConfig{
		DataDir:  "./test_data",
		RedisDsn: "",  // Empty DSN
	})
	
	err := app.initRedis()
	
	assert.NoError(t, err)
	assert.Nil(t, app.redisCache)  // Redis should be nil
}

func TestInitRedis_InvalidDSN_FailFast(t *testing.T) {
	// Test: Invalid DSN → fail fast with error
	app := NewBaseApp(BaseAppConfig{
		DataDir:  "./test_data",
		RedisDsn: "invalid://bad-format",
	})
	
	err := app.initRedis()
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid Redis DSN")
}

func TestInitRedis_ConnectionFailed_GracefulDegradation(t *testing.T) {
	// Test: Redis not running → graceful degradation (no error)
	app := NewBaseApp(BaseAppConfig{
		DataDir:  "./test_data",
		RedisDsn: "redis://localhost:9999/0",  // Wrong port
	})
	
	err := app.initRedis()
	
	// Should NOT error (graceful degradation)
	assert.NoError(t, err)
	assert.Nil(t, app.redisCache)
}

func TestInitRedis_Success_ConnectionEstablished(t *testing.T) {
	// Test: Valid DSN + running Redis → connection established
	// Prerequisite: Redis must be running on localhost:6379
	
	app := NewBaseApp(BaseAppConfig{
		DataDir:  "./test_data",
		RedisDsn: "redis://localhost:6379/0",
	})
	
	err := app.initRedis()
	
	require.NoError(t, err)
	require.NotNil(t, app.redisCache)
	
	// Verify PING works
	pong, err := app.redisCache.Ping(app.redisContext).Result()
	assert.NoError(t, err)
	assert.Equal(t, "PONG", pong)
}

// =============================================================================
// Test Group 2: Publish Method
// =============================================================================

func TestPublish_NoRedis_LocalBroadcastFallback(t *testing.T) {
	// Test: No Redis → fallback to local broadcast
	app := NewBaseApp(BaseAppConfig{
		DataDir:  "./test_data",
		RedisDsn: "",  // No Redis
	})
	app.initRedis()
	
	// Register local listener
	received := false
	app.OnRealtimeBroadcast().Add(func(e *RealtimeBroadcastEvent) error {
		received = true
		return nil
	})
	
	// Publish message
	err := app.Publish("realtime", map[string]any{
		"action": "test",
	})
	
	assert.NoError(t, err)
	assert.True(t, received, "Local listener should receive message")
}

func TestPublish_WithRedis_MessagePublished(t *testing.T) {
	// Test: With Redis → message published to Redis channel
	// Prerequisite: Redis running
	
	app := NewBaseApp(BaseAppConfig{
		DataDir:  "./test_data",
		RedisDsn: "redis://localhost:6379/0",
	})
	err := app.initRedis()
	require.NoError(t, err)
	require.NotNil(t, app.redisCache)
	
	// Publish message
	err = app.Publish("realtime", map[string]any{
		"action": "test",
		"data":   "hello",
	})
	
	assert.NoError(t, err)
}

func TestPublish_InvalidData_MarshalError(t *testing.T) {
	// Test: Invalid data (cannot marshal to JSON) → error
	app := NewBaseApp(BaseAppConfig{
		DataDir:  "./test_data",
		RedisDsn: "redis://localhost:6379/0",
	})
	app.initRedis()
	
	// Try to publish non-serializable data
	invalidData := make(chan int)  // Channels cannot be marshaled
	err := app.Publish("realtime", invalidData)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marshal")
}

// =============================================================================
// Test Group 3: Subscriber (Background Goroutine)
// =============================================================================

func TestRedisSubscriber_ReceivesMessages(t *testing.T) {
	// Test: Subscriber goroutine receives published messages
	// Prerequisite: Redis running
	
	app := NewBaseApp(BaseAppConfig{
		DataDir:  "./test_data",
		RedisDsn: "redis://localhost:6379/0",
	})
	err := app.initRedis()
	require.NoError(t, err)
	
	// Register listener
	received := false
	app.OnRealtimeBroadcast().Add(func(e *RealtimeBroadcastEvent) error {
		received = true
		t.Logf("Received: %s", string(e.Payload))
		return nil
	})
	
	// Publish message (will be received by subscriber goroutine)
	time.Sleep(100 * time.Millisecond)  // Let subscriber start
	err = app.Publish("realtime", map[string]any{"test": "data"})
	require.NoError(t, err)
	
	// Wait for message propagation
	time.Sleep(200 * time.Millisecond)
	
	assert.True(t, received, "Subscriber should receive published message")
}
```

### Coverage Target
- **Target:** >85% code coverage
- **Critical Paths:** All Redis connection states, publish/subscribe mechanics

---

## 🧪 TEST SUITE 3: Query Builder Tests

### File: `core/query_builder_test.go`

```go
package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/free/postgresqlbaseapi/dbx"
)

// =============================================================================
// Test Group 1: Placeholder Conversion
// =============================================================================

func TestQueryBuilder_PlaceholderConversion_PostgreSQL(t *testing.T) {
	// Test: ? → $1, $2, $3 for PostgreSQL
	db, _ := dbx.Open("postgres", "postgres://localhost/test?sslmode=disable")
	defer db.Close()
	
	query := db.NewQuery("SELECT * FROM users WHERE email = ? AND age > ?")
	query.Bind("test@example.com", 18)
	
	sql, params := query.Build()
	
	// Verify placeholders converted
	assert.Contains(t, sql, "$1")
	assert.Contains(t, sql, "$2")
	assert.NotContains(t, sql, "?")
	
	// Verify params
	assert.Equal(t, 2, len(params))
	assert.Equal(t, "test@example.com", params[0])
	assert.Equal(t, 18, params[1])
}

func TestQueryBuilder_PlaceholderConversion_MySQL(t *testing.T) {
	// Test: ? remains as ? for MySQL
	db, _ := dbx.Open("mysql", "root:pass@tcp(localhost:3306)/test")
	defer db.Close()
	
	query := db.NewQuery("SELECT * FROM users WHERE email = ? AND age > ?")
	query.Bind("test@example.com", 18)
	
	sql, params := query.Build()
	
	// Verify placeholders remain as ?
	assert.Contains(t, sql, "?")
	assert.NotContains(t, sql, "$1")
}

// =============================================================================
// Test Group 2: CRUD Operations
// =============================================================================

func TestQueryBuilder_INSERT(t *testing.T) {
	db, _ := dbx.Open("postgres", "postgres://localhost/test?sslmode=disable")
	defer db.Close()
	
	// Test INSERT query building
	_, err := db.Insert("users", dbx.Params{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
	}).Execute()
	
	assert.NoError(t, err)
}

func TestQueryBuilder_SELECT(t *testing.T) {
	db, _ := dbx.Open("postgres", "postgres://localhost/test?sslmode=disable")
	defer db.Close()
	
	var users []struct {
		ID    int    `db:"id"`
		Name  string `db:"name"`
		Email string `db:"email"`
	}
	
	err := db.Select("*").From("users").Where(dbx.HashExp{"age": 30}).All(&users)
	assert.NoError(t, err)
}

func TestQueryBuilder_UPDATE(t *testing.T) {
	db, _ := dbx.Open("postgres", "postgres://localhost/test?sslmode=disable")
	defer db.Close()
	
	_, err := db.Update("users", dbx.Params{
		"age": 31,
	}, dbx.HashExp{"email": "john@example.com"}).Execute()
	
	assert.NoError(t, err)
}

func TestQueryBuilder_DELETE(t *testing.T) {
	db, _ := dbx.Open("postgres", "postgres://localhost/test?sslmode=disable")
	defer db.Close()
	
	_, err := db.Delete("users", dbx.HashExp{"email": "john@example.com"}).Execute()
	assert.NoError(t, err)
}
```

### Coverage Target
- **Target:** >90% code coverage
- **Critical Paths:** Placeholder conversion, all CRUD operations

---

## 📊 Coverage Requirements

### Overall Coverage Targets

| Component | Target Coverage | Priority |
|-----------|----------------|----------|
| `db_postgresql.go` | >90% | Critical |
| `redis integration` | >85% | High |
| `query builder` | >90% | Critical |
| `CLI flags` | >70% | Medium |

### Coverage Measurement

**Command:**
```bash
go test ./core/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o docs/coverage.html
```

**Expected Output:**
```
core/db_postgresql.go:      92.5%
core/redis.go:              87.3%
core/query_builder.go:      93.1%
---
TOTAL:                      89.8%
```

---

## 🔧 Test Infrastructure Setup

### Test Database Setup

**Docker Compose:** `docker-compose.test.yml`
```yaml
version: '3.8'

services:
  postgres-test:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: testdb
    ports:
      - "5432:5432"

  mysql-test:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: testdb
    ports:
      - "3306:3306"

  redis-test:
    image: redis:7-alpine
    ports:
      - "6379:6379"
```

**Start test infrastructure:**
```bash
docker-compose -f docker-compose.test.yml up -d
```

### Test Execution Commands

```bash
# Run all unit tests
go test ./core/... -v

# Run with coverage
go test ./core/... -cover -coverprofile=coverage.out

# Run specific test
go test ./core/... -run TestConnectDB_PostgreSQL -v

# Run tests in parallel
go test ./core/... -parallel=4

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html
```

---

## ✅ Verification Checklist

**Before marking AGENT-10 complete:**

- [ ] All test files created (`*_test.go`)
- [ ] Database connection tests (10+ test cases)
- [ ] Redis integration tests (8+ test cases)
- [ ] Query builder tests (5+ test cases)
- [ ] All tests passing
- [ ] Coverage >80% overall
- [ ] Coverage >90% for critical components
- [ ] Test infrastructure documented (Docker Compose)
- [ ] Test execution commands documented

---

**AGENT-10 Status:** 🟢 Ready for Review  
**Document Version:** 1.0  
**Last Updated:** 2026-05-03
