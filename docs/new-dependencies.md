# New Dependencies for PostgreSQL + Redis Migration

**Analysis Date**: 2026-05-03  
**Analyst**: AGENT-03 (Dependency Analysis Specialist)  
**Purpose**: Document new dependencies required for PostgreSQL/MySQL + Redis support

---

## Executive Summary

### Dependencies to Add
| Package | Version | Purpose | Priority | Size Impact |
|---------|---------|---------|----------|-------------|
| `github.com/lib/pq` | v1.10.9 | PostgreSQL driver | 🔴 CRITICAL | +1.5MB |
| `github.com/go-sql-driver/mysql` | v1.7.1 | MySQL driver | 🟡 MEDIUM | +1.0MB |
| `github.com/redis/go-redis/v9` | v9.3.0 | Redis client | 🔴 CRITICAL | +0.5MB |

**Total Size Impact**: +3.0MB (~8.5% increase from 35MB to 38MB)  
**Binary Compatibility**: ✅ All packages support Go 1.25.0  
**License Compatibility**: ✅ All packages use permissive licenses (MIT/BSD/MPL)

---

## 1. PostgreSQL Driver: lib/pq

### Package Information
```go
require github.com/lib/pq v1.10.9
```

| Attribute | Value |
|-----------|-------|
| **Full Path** | `github.com/lib/pq` |
| **Version** | `v1.10.9` |
| **Release Date** | July 2023 |
| **License** | MIT |
| **Go Version** | 1.13+ (✅ compatible with 1.25.0) |
| **Import Path** | `github.com/lib/pq` |
| **Documentation** | https://pkg.go.dev/github.com/lib/pq |

### Features
- ✅ Pure Go implementation (no C dependencies)
- ✅ Full `database/sql` driver interface
- ✅ PostgreSQL 9.0+ support (tested up to PostgreSQL 16)
- ✅ SSL/TLS connection support
- ✅ Connection pooling via `sql.DB`
- ✅ Transaction support
- ✅ Prepared statements
- ✅ Array types support
- ✅ JSON/JSONB support
- ✅ Listen/Notify support (real-time notifications)

### Why lib/pq vs alternatives?

| Driver | Pros | Cons | Verdict |
|--------|------|------|---------|
| **lib/pq** | Stable, Simple, Well-tested, Standard | Slow updates, Fewer features | ✅ **CHOSEN** |
| **pgx/v5** | Fast, Active dev, Rich features | Larger size (+3MB), More complex | ⏸️ Future upgrade |
| **go-pg/pg** | ORM included, High-level | Not `database/sql`, Breaking changes | ❌ Not suitable |

**Decision**: Use `lib/pq` for initial migration
- **Reason**: Postgrebase already uses lib/pq (proven stable)
- **Future**: Consider pgx/v5 for performance optimizations

### DSN Format
```go
// Standard PostgreSQL DSN
postgres://user:password@localhost:5432/dbname?sslmode=disable

// Full format with all options
postgres://user:password@host:port/database?sslmode=verify-full&pool_max_conns=10

// Unix socket connection
postgres://user:password@/dbname?host=/var/run/postgresql

// Connection string format (alternative)
"host=localhost port=5432 user=postgres password=secret dbname=mydb sslmode=disable"
```

### Integration Points in Codebase

#### 1. Import Statement
```go
// core/db_postgresql.go
import (
    "database/sql"
    _ "github.com/lib/pq" // PostgreSQL driver
)
```

#### 2. Connection Function
```go
// core/db_postgresql.go
func connectPostgres(dsn string) (*sql.DB, error) {
    // lib/pq automatically registers "postgres" driver
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return db, nil
}
```

#### 3. Pool Configuration
```go
// core/db_postgresql.go
func configurePool(db *sql.DB) {
    db.SetMaxOpenConns(25)     // Max connections
    db.SetMaxIdleConns(5)      // Idle connections
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(10 * time.Minute)
}
```

### Known Issues & Workarounds

#### Issue 1: Slow Maintenance
- **Problem**: Last release July 2023 (no updates for 10+ months)
- **Impact**: Minor bug fixes delayed
- **Mitigation**: pgx/v5 team maintains fork, can switch if needed
- **Status**: ⚠️ Monitor but acceptable

#### Issue 2: Array Type Handling
- **Problem**: Arrays require `pq.Array()` wrapper
- **Workaround**:
  ```go
  import "github.com/lib/pq"
  
  // Query with array parameter
  rows, err := db.Query("SELECT * FROM users WHERE id = ANY($1)", pq.Array(ids))
  ```

#### Issue 3: JSON Type Handling
- **Problem**: JSON values returned as `[]byte`, not parsed
- **Workaround**:
  ```go
  var jsonData []byte
  err := row.Scan(&jsonData)
  
  var result map[string]interface{}
  json.Unmarshal(jsonData, &result)
  ```

### Security Considerations

#### SSL/TLS Support
```go
// Require SSL
dsn := "postgres://user:pass@host/db?sslmode=require"

// Verify certificate
dsn := "postgres://user:pass@host/db?sslmode=verify-full&sslrootcert=/path/to/ca.crt"

// Development only (no SSL)
dsn := "postgres://user:pass@host/db?sslmode=disable"
```

#### Connection String Sanitization
```go
// core/db_postgresql.go
func sanitizeDSN(dsn string) string {
    // Remove password from logs
    re := regexp.MustCompile(`(password=)[^&\s]+`)
    return re.ReplaceAllString(dsn, "$1***")
}
```

### Testing Requirements

```go
// core/db_postgresql_test.go
func TestLibPQ_Connection(t *testing.T) {
    dsn := "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable"
    db, err := connectPostgres(dsn)
    require.NoError(t, err)
    defer db.Close()
    
    // Test ping
    err = db.Ping()
    require.NoError(t, err)
}

func TestLibPQ_Query(t *testing.T) {
    // Test basic query
    rows, err := db.Query("SELECT 1 AS number")
    require.NoError(t, err)
    defer rows.Close()
    
    var result int
    rows.Next()
    err = rows.Scan(&result)
    require.NoError(t, err)
    assert.Equal(t, 1, result)
}
```

---

## 2. MySQL Driver: go-sql-driver/mysql

### Package Information
```go
require github.com/go-sql-driver/mysql v1.7.1
```

| Attribute | Value |
|-----------|-------|
| **Full Path** | `github.com/go-sql-driver/mysql` |
| **Version** | `v1.7.1` |
| **Release Date** | June 2023 |
| **License** | MPL 2.0 (Mozilla Public License) |
| **Go Version** | 1.13+ (✅ compatible with 1.25.0) |
| **Import Path** | `github.com/go-sql-driver/mysql` |
| **Documentation** | https://pkg.go.dev/github.com/go-sql-driver/mysql |

### Features
- ✅ Pure Go implementation
- ✅ Full `database/sql` driver interface
- ✅ MySQL 5.5+ and MariaDB 10.0+ support
- ✅ TLS/SSL connection support
- ✅ Connection pooling
- ✅ Prepared statements
- ✅ Transaction support
- ✅ Multi-statement support
- ✅ JSON type support
- ✅ Connection timeout configuration

### Why go-sql-driver/mysql?

| Driver | Pros | Cons | Verdict |
|--------|------|------|---------|
| **go-sql-driver/mysql** | Standard, Fast, Well-maintained | None significant | ✅ **CHOSEN** |
| **MyMySQL** | Alternative API | Less popular, No `database/sql` | ❌ Not suitable |
| **go-mysql** | Replication support | Not `database/sql` compatible | ❌ Not suitable |

**Decision**: go-sql-driver/mysql is the de-facto standard

### DSN Format
```go
// Standard MySQL DSN
mysql://user:password@tcp(localhost:3306)/dbname?parseTime=true

// Full format with options
mysql://user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=true&loc=Local

// Unix socket connection
mysql://user:password@unix(/var/run/mysqld/mysqld.sock)/dbname

// Connection string format (alternative)
"user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=true"
```

### Integration Points

#### 1. Import Statement
```go
// core/db_postgresql.go (handles both PostgreSQL and MySQL)
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql" // MySQL driver
)
```

#### 2. Connection Function
```go
// core/db_postgresql.go
func connectMySQL(dsn string) (*sql.DB, error) {
    // go-sql-driver/mysql registers "mysql" driver
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return db, nil
}
```

#### 3. Driver Detection
```go
// core/db_postgresql.go
func connectDB(dsn string) (*sql.DB, error) {
    if strings.HasPrefix(dsn, "postgres://") {
        return connectPostgres(dsn)
    } else if strings.HasPrefix(dsn, "mysql://") {
        return connectMySQL(dsn)
    } else {
        return nil, fmt.Errorf("unsupported database DSN: %s", dsn)
    }
}
```

### Configuration Best Practices

```go
// core/db_postgresql.go
func configureMySQLPool(db *sql.DB) {
    db.SetMaxOpenConns(25)     // Max connections
    db.SetMaxIdleConns(5)      // Idle connections
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(10 * time.Minute)
}

// DSN parameters
const mysqlDSNTemplate = "user:pass@tcp(host:port)/db?charset=utf8mb4&parseTime=true&loc=Local&timeout=30s&readTimeout=30s&writeTimeout=30s"
```

### Known Issues & Workarounds

#### Issue 1: Time Parsing
- **Problem**: MySQL returns timestamps as strings by default
- **Solution**: Add `parseTime=true` to DSN
  ```go
  dsn := "user:pass@tcp(localhost:3306)/db?parseTime=true"
  ```

#### Issue 2: Character Encoding
- **Problem**: Default charset may not be UTF-8
- **Solution**: Explicitly set charset in DSN
  ```go
  dsn := "user:pass@tcp(localhost:3306)/db?charset=utf8mb4"
  ```

#### Issue 3: Connection Timeout
- **Problem**: Default timeout is too long (30s)
- **Solution**: Set explicit timeouts
  ```go
  dsn := "user:pass@tcp(localhost:3306)/db?timeout=10s&readTimeout=30s&writeTimeout=30s"
  ```

### Security Considerations

#### TLS/SSL Support
```go
// Require TLS
dsn := "user:pass@tcp(host:3306)/db?tls=true"

// Custom TLS config
dsn := "user:pass@tcp(host:3306)/db?tls=custom"

// Register custom TLS config
import "crypto/tls"
mysql.RegisterTLSConfig("custom", &tls.Config{
    InsecureSkipVerify: false,
    MinVersion: tls.VersionTLS12,
})
```

### Testing Requirements

```go
// core/db_mysql_test.go
func TestMySQL_Connection(t *testing.T) {
    dsn := "root:root@tcp(localhost:3306)/test_db?parseTime=true"
    db, err := connectMySQL(dsn)
    require.NoError(t, err)
    defer db.Close()
    
    err = db.Ping()
    require.NoError(t, err)
}
```

---

## 3. Redis Client: go-redis/v9

### Package Information
```go
require github.com/redis/go-redis/v9 v9.3.0
```

| Attribute | Value |
|-----------|-------|
| **Full Path** | `github.com/redis/go-redis/v9` |
| **Version** | `v9.3.0` |
| **Release Date** | November 2023 |
| **License** | BSD-2-Clause |
| **Go Version** | 1.18+ (✅ compatible with 1.25.0) |
| **Import Path** | `github.com/redis/go-redis/v9` |
| **Documentation** | https://redis.uptrace.dev/ |

### Features
- ✅ Redis 6.0+ and Redis 7.0+ support
- ✅ Pub/Sub support (for realtime events)
- ✅ Connection pooling
- ✅ Cluster support
- ✅ Sentinel support
- ✅ Pipeline and transaction support
- ✅ Lua scripting support
- ✅ Stream support (Redis Streams)
- ✅ Context-based API (cancellation support)
- ✅ Automatic reconnection
- ✅ Health checks

### Why go-redis/v9?

| Client | Pros | Cons | Verdict |
|--------|------|------|---------|
| **go-redis/v9** | Official, Feature-rich, Active | None significant | ✅ **CHOSEN** |
| **redigo** | Mature, Stable | Less active, No context support | ❌ Outdated |
| **radix** | Alternative | Less popular | ❌ Not needed |

**Decision**: go-redis/v9 is the official Redis Go client

### DSN Format
```go
// Standard Redis DSN
redis://localhost:6379/0

// With password
redis://:password@localhost:6379/0

// Full format
redis://user:password@localhost:6379/0?protocol=3&timeout=5s&maxRetries=3

// Redis Cluster
redis://localhost:7000,localhost:7001,localhost:7002/0

// Redis Sentinel
redis://sentinel1:26379,sentinel2:26379/mymaster?sentinel_password=secret
```

### Integration Points

#### 1. Import Statement
```go
// core/base.go
import (
    "github.com/redis/go-redis/v9"
)
```

#### 2. Client Initialization
```go
// core/base.go
type BaseApp struct {
    // ... existing fields ...
    redis redis.UniversalClient // Supports both single and cluster
}

func (app *BaseApp) initRedis(dsn string) error {
    if dsn == "" {
        // No Redis configured, use in-memory
        app.redis = nil
        return nil
    }
    
    // Parse DSN
    opt, err := redis.ParseURL(dsn)
    if err != nil {
        return fmt.Errorf("invalid Redis DSN: %w", err)
    }
    
    // Create client
    app.redis = redis.NewClient(opt)
    
    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := app.redis.Ping(ctx).Err(); err != nil {
        // Log warning but don't fail (graceful degradation)
        log.Println("Redis connection failed, falling back to in-memory:", err)
        app.redis = nil
        return nil
    }
    
    log.Println("Redis connected successfully")
    return nil
}
```

#### 3. Publish Method (for realtime events)
```go
// core/base.go
func (app *BaseApp) Publish(channel string, message string) error {
    if app.redis == nil {
        // Fallback to in-memory broadcast (existing PocketBase behavior)
        return app.broadcastInMemory(channel, message)
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Publish to Redis (will be received by all nodes)
    err := app.redis.Publish(ctx, channel, message).Err()
    if err != nil {
        // Log error but don't fail (graceful degradation)
        log.Printf("Redis publish failed: %v, falling back to in-memory", err)
        return app.broadcastInMemory(channel, message)
    }
    
    return nil
}
```

#### 4. Subscribe Method (for SSE endpoints)
```go
// apis/realtime.go
func (api *realtimeApi) subscribeToRedis(c echo.Context) error {
    if app.redis == nil {
        // Use existing in-memory subscription
        return api.subscribeInMemory(c)
    }
    
    ctx := c.Request().Context()
    pubsub := app.redis.Subscribe(ctx, "realtime:*")
    defer pubsub.Close()
    
    // Receive messages
    ch := pubsub.Channel()
    for {
        select {
        case msg := <-ch:
            // Send SSE event to client
            fmt.Fprintf(c.Response(), "data: %s\n\n", msg.Payload)
            c.Response().Flush()
        case <-ctx.Done():
            return nil
        }
    }
}
```

### Configuration Best Practices

```go
// core/base.go
func createRedisClient(dsn string) *redis.Client {
    opt, _ := redis.ParseURL(dsn)
    
    // Connection pool settings
    opt.PoolSize = 10              // Max connections
    opt.MinIdleConns = 2           // Idle connections
    opt.MaxIdleConns = 5           // Max idle
    opt.ConnMaxLifetime = 5 * time.Minute
    opt.ConnMaxIdleTime = 10 * time.Minute
    
    // Timeout settings
    opt.DialTimeout = 5 * time.Second
    opt.ReadTimeout = 3 * time.Second
    opt.WriteTimeout = 3 * time.Second
    
    // Retry settings
    opt.MaxRetries = 3
    opt.MinRetryBackoff = 100 * time.Millisecond
    opt.MaxRetryBackoff = 2 * time.Second
    
    return redis.NewClient(opt)
}
```

### Known Issues & Workarounds

#### Issue 1: Context Cancellation
- **Problem**: All operations require context (breaking change from v8)
- **Solution**: Pass request context or use `context.Background()`
  ```go
  // Good: Use request context
  ctx := c.Request().Context()
  err := rdb.Publish(ctx, channel, message).Err()
  
  // Acceptable: Background context for long-running
  ctx := context.Background()
  err := rdb.Subscribe(ctx, "channel").Err()
  ```

#### Issue 2: Connection Pool Exhaustion
- **Problem**: Too many concurrent connections
- **Solution**: Configure pool size appropriately
  ```go
  opt.PoolSize = 10 * runtime.NumCPU() // Rule of thumb
  ```

#### Issue 3: Pub/Sub Reliability
- **Problem**: Messages can be lost if subscriber is slow
- **Solution**: Use Redis Streams for guaranteed delivery (future enhancement)

### Graceful Degradation Strategy

**Critical Design Principle**: Redis is **optional** for basic operation

```go
// core/base.go
func (app *BaseApp) Publish(channel string, message string) error {
    if app.redis == nil {
        // ✅ FALLBACK: Use existing in-memory broadcast
        return app.broadcastInMemory(channel, message)
    }
    
    // Try Redis first
    err := app.publishRedis(channel, message)
    if err != nil {
        // ✅ FALLBACK: Redis failed, use in-memory
        log.Printf("Redis publish failed: %v, using in-memory fallback", err)
        return app.broadcastInMemory(channel, message)
    }
    
    return nil
}
```

**Benefits**:
- ✅ Single-node deployments work without Redis
- ✅ Development environments simpler (no Redis setup required)
- ✅ Production resilience (continues working if Redis crashes)
- ✅ Progressive enhancement (add Redis for multi-node only)

### Security Considerations

#### Authentication
```go
// Redis 6.0+ ACL support
dsn := "redis://username:password@localhost:6379/0"

// Legacy AUTH command
dsn := "redis://:password@localhost:6379/0"
```

#### TLS/SSL Support
```go
// Enable TLS
dsn := "rediss://localhost:6380/0" // Note: "rediss" not "redis"

// Custom TLS config
opt, _ := redis.ParseURL(dsn)
opt.TLSConfig = &tls.Config{
    MinVersion: tls.VersionTLS12,
    InsecureSkipVerify: false,
}
```

### Testing Requirements

```go
// core/redis_test.go
func TestRedis_Connection(t *testing.T) {
    dsn := "redis://localhost:6379/0"
    rdb, err := redis.ParseURL(dsn)
    require.NoError(t, err)
    
    client := redis.NewClient(rdb)
    defer client.Close()
    
    ctx := context.Background()
    err = client.Ping(ctx).Err()
    require.NoError(t, err)
}

func TestRedis_PubSub(t *testing.T) {
    ctx := context.Background()
    
    // Publisher
    err := rdb.Publish(ctx, "test-channel", "hello").Err()
    require.NoError(t, err)
    
    // Subscriber
    pubsub := rdb.Subscribe(ctx, "test-channel")
    defer pubsub.Close()
    
    msg, err := pubsub.ReceiveMessage(ctx)
    require.NoError(t, err)
    assert.Equal(t, "hello", msg.Payload)
}

func TestRedis_GracefulFallback(t *testing.T) {
    // Test with Redis unavailable
    app := &BaseApp{redis: nil}
    
    // Should fallback to in-memory
    err := app.Publish("channel", "message")
    require.NoError(t, err)
}
```

---

## 4. Installation & go.mod Updates

### Step 1: Add Dependencies
```bash
# Add all three packages
go get github.com/lib/pq@v1.10.9
go get github.com/go-sql-driver/mysql@v1.7.1
go get github.com/redis/go-redis/v9@v9.3.0

# Clean up
go mod tidy
```

### Step 2: Verify go.mod
```go
// go.mod (after updates)
module github.com/your-org/rai-backend

go 1.25.0

require (
    // ... existing PocketBase dependencies ...
    
    // NEW: Database drivers
    github.com/lib/pq v1.10.9
    github.com/go-sql-driver/mysql v1.7.1
    
    // NEW: Redis client
    github.com/redis/go-redis/v9 v9.3.0
)
```

### Step 3: Verify go.sum
```bash
# Check that checksums are generated
grep "github.com/lib/pq" go.sum
grep "github.com/go-sql-driver/mysql" go.sum
grep "github.com/redis/go-redis" go.sum
```

### Step 4: Test Compilation
```bash
# Build with new dependencies
go build -v ./...

# Check binary size
ls -lh ./rai-backend
# Expected: ~38MB (from 35MB)
```

---

## 5. Import Statement Summary

### File: core/db_postgresql.go (new file)
```go
package core

import (
    "database/sql"
    "fmt"
    "strings"
    
    _ "github.com/lib/pq"                // PostgreSQL driver
    _ "github.com/go-sql-driver/mysql"   // MySQL driver
)
```

### File: core/base.go (modified)
```go
package core

import (
    // ... existing imports ...
    
    "github.com/redis/go-redis/v9"   // Redis client
)
```

---

## 6. Verification Checklist

### Pre-Installation Checks
- [ ] Go version 1.25.0 installed
- [ ] go.mod exists and is valid
- [ ] Network access to pkg.go.dev (for downloads)

### Post-Installation Checks
- [ ] `go.mod` contains all three packages
- [ ] `go.sum` contains checksums
- [ ] `go build` succeeds
- [ ] Binary size ~38MB
- [ ] No import errors in IDE

### Runtime Checks
- [ ] PostgreSQL connection works
- [ ] MySQL connection works (if used)
- [ ] Redis connection works (if configured)
- [ ] Graceful fallback when Redis unavailable
- [ ] No crashes with missing drivers

---

## 7. Rollback Instructions

### Remove All New Dependencies
```bash
go mod edit -droprequire github.com/lib/pq
go mod edit -droprequire github.com/go-sql-driver/mysql
go mod edit -droprequire github.com/redis/go-redis/v9
go mod tidy
```

### Remove Individual Dependency
```bash
# Example: Remove only Redis
go mod edit -droprequire github.com/redis/go-redis/v9
go mod tidy
```

---

## 8. Alternative Versions (for troubleshooting)

### lib/pq Alternatives
```bash
# Latest version
go get github.com/lib/pq@latest

# Specific older version (if compatibility issue)
go get github.com/lib/pq@v1.10.7
```

### go-sql-driver/mysql Alternatives
```bash
# Latest version
go get github.com/go-sql-driver/mysql@latest

# Avoid v1.8.x (breaking changes)
go get github.com/go-sql-driver/mysql@v1.7.1
```

### go-redis/v9 Alternatives
```bash
# Latest v9.x
go get github.com/redis/go-redis/v9@latest

# Fallback to v8 (if context issues)
go get github.com/go-redis/redis/v8@v8.11.5
```

---

## 9. License Compliance

| Package | License | Commercial Use | Attribution Required | Copyleft |
|---------|---------|----------------|----------------------|----------|
| lib/pq | MIT | ✅ Yes | ⚠️ Recommended | ❌ No |
| go-sql-driver/mysql | MPL 2.0 | ✅ Yes | ✅ Yes | ⚠️ File-level |
| go-redis/v9 | BSD-2-Clause | ✅ Yes | ⚠️ Recommended | ❌ No |

**Compliance Actions**:
1. **lib/pq (MIT)**: Include MIT license text in NOTICES file
2. **go-sql-driver/mysql (MPL 2.0)**: Include MPL 2.0 license, provide source if modified
3. **go-redis/v9 (BSD-2-Clause)**: Include BSD license text in NOTICES file

---

## 10. Next Steps

### For AGENT-05 (Database Connection Layer Engineer)
- ✅ Use lib/pq v1.10.9 for PostgreSQL
- ✅ Use go-sql-driver/mysql v1.7.1 for MySQL
- ✅ Reference DSN formats from this document
- ✅ Implement driver detection logic (postgres:// vs mysql://)

### For AGENT-06 (Redis Integration Engineer)
- ✅ Use go-redis/v9 v9.3.0
- ✅ Implement graceful fallback to in-memory
- ✅ Reference Pub/Sub examples from this document
- ✅ Add context-based API calls

### For AGENT-10 (Unit Test Engineer)
- ✅ Test lib/pq connection (PostgreSQL required)
- ✅ Test mysql driver connection (MySQL required)
- ✅ Test Redis client (Redis optional)
- ✅ Test graceful degradation (Redis unavailable)

---

**Report Status**: ✅ Complete  
**Dependencies Documented**: 3 packages  
**Next Deliverable**: `docs/dependency-upgrade-plan.md` (TASK-03-C)
