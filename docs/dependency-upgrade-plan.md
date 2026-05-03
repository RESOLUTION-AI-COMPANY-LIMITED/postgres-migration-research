# Dependency Upgrade Plan: PocketBase v0.37.5 + PostgreSQL Migration

**Plan Date**: 2026-05-03  
**Analyst**: AGENT-03 (Dependency Analysis Specialist)  
**Purpose**: Phased upgrade strategy with breaking change mitigation and rollback procedures

---

## Executive Summary

### Upgrade Strategy: Phased Approach

```
Phase 1: PocketBase v0.37.5 Base     (2 hours)  ✅ Low risk
Phase 2: PostgreSQL Driver           (1 hour)   ✅ Low risk
Phase 3: Redis Client                (2 hours)  ✅ Medium risk (graceful fallback)
Phase 4: MySQL Driver (Optional)     (1 hour)   ✅ Low risk
```

**Total Estimated Time**: 6 hours  
**Risk Level**: 🟡 MEDIUM (manageable with testing)  
**Rollback Complexity**: ✅ LOW (git-based, per-phase checkpoints)

### Key Decisions

| Decision | Rationale |
|----------|-----------|
| **Phased Rollout** | Isolate breaking changes, enable per-phase rollback |
| **PocketBase First** | Leverage upstream testing, get security fixes early |
| **PostgreSQL Before Redis** | Database is critical, Redis is optional enhancement |
| **Graceful Degradation** | Redis failures don't crash the app |
| **Comprehensive Testing** | Test after each phase before proceeding |

---

## Phase 1: Upgrade to PocketBase v0.37.5 Base (Without PostgreSQL)

### Objective
Upgrade all shared dependencies to PocketBase v0.37.5 versions, resolve breaking changes (especially JWT v4→v5), but keep SQLite as the database.

### Duration
**Estimated**: 2 hours  
**Breakdown**:
- go.mod updates: 15 minutes
- JWT v4→v5 code migration: 60 minutes
- Testing: 30 minutes
- Verification: 15 minutes

---

### Step 1.1: Create Checkpoint
```bash
# Create git checkpoint before changes
git add .
git commit -m "Checkpoint: Before PocketBase v0.37.5 upgrade"
git tag v0.2.0-pre-pocketbase-upgrade
```

---

### Step 1.2: Update go.mod (Shared Dependencies)

#### Commands
```bash
# Update Go version
go mod edit -go=1.25.0

# Update shared dependencies to PocketBase v0.37.5 versions
go get github.com/disintegration/imaging@v1.6.2         # No change
go get github.com/domodwyer/mailyak/v3@v3.6.2           # No change
go get github.com/fatih/color@v1.19.0                   # 1.15.0 → 1.19.0
go get github.com/gabriel-vasile/mimetype@v1.4.13       # 1.4.2 → 1.4.13
go get github.com/ganigeorgiev/fexpr@v0.5.0             # 0.3.0 → 0.5.0
go get github.com/go-ozzo/ozzo-validation/v4@v4.3.0     # No change
go get github.com/spf13/cast@v1.10.0                    # 1.6.0 → 1.10.0
go get github.com/spf13/cobra@v1.10.2                   # 1.7.0 → 1.10.2
go get golang.org/x/crypto@v0.50.0                      # 0.12.0 → 0.50.0 (CRITICAL)
go get golang.org/x/image@v0.39.0                       # 0.11.0 → 0.39.0
go get golang.org/x/net@v0.53.0                         # 0.14.0 → 0.53.0
go get golang.org/x/oauth2@v0.36.0                      # 0.11.0 → 0.36.0
go get golang.org/x/sync@v0.20.0                        # 0.3.0 → 0.20.0
go get modernc.org/sqlite@v1.50.0                       # 1.25.0 → 1.50.0

# BREAKING CHANGE: JWT v4 → v5
go get github.com/golang-jwt/jwt/v5@v5.3.1

# Clean up
go mod tidy
```

#### Expected go.mod (After Step 1.2)
```go
module github.com/your-org/rai-backend

go 1.25.0

require (
    github.com/disintegration/imaging v1.6.2
    github.com/domodwyer/mailyak/v3 v3.6.2
    github.com/fatih/color v1.19.0               // Updated
    github.com/gabriel-vasile/mimetype v1.4.13   // Updated
    github.com/ganigeorgiev/fexpr v0.5.0         // Updated
    github.com/go-ozzo/ozzo-validation/v4 v4.3.0
    github.com/golang-jwt/jwt/v5 v5.3.1          // BREAKING: v4 → v5
    github.com/spf13/cast v1.10.0                // Updated
    github.com/spf13/cobra v1.10.2               // Updated
    golang.org/x/crypto v0.50.0                  // Updated (security)
    golang.org/x/image v0.39.0                   // Updated
    golang.org/x/net v0.53.0                     // Updated
    golang.org/x/oauth2 v0.36.0                  // Updated
    golang.org/x/sync v0.20.0                    // Updated
    modernc.org/sqlite v1.50.0                   // Updated
)
```

---

### Step 1.3: Migrate JWT v4 → v5 (BREAKING CHANGE)

#### Breaking Changes in jwt/v5

##### Change 1: Import Path
```go
// OLD (jwt/v4)
import "github.com/golang-jwt/jwt/v4"

// NEW (jwt/v5)
import "github.com/golang-jwt/jwt/v5"
```

**Migration Command**:
```bash
# Find all files using jwt/v4
grep -r "jwt/v4" . --include="*.go"

# Update imports
find . -name "*.go" -exec sed -i 's|github.com/golang-jwt/jwt/v4|github.com/golang-jwt/jwt/v5|g' {} \;
```

---

##### Change 2: Parse() Function Signature

**OLD (jwt/v4)**:
```go
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    // Validate algorithm
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    return []byte(secret), nil
})
```

**NEW (jwt/v5)**:
```go
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte(secret), nil
}, jwt.WithValidMethods([]string{"HS256", "HS384", "HS512"}))
```

**Key Changes**:
- Algorithm validation moved to `jwt.WithValidMethods()` option
- More explicit and secure (required algorithms must be specified)
- Prevents algorithm confusion attacks

**Migration Pattern**:
```go
// Find this pattern
token, err := jwt.Parse(tokenString, keyFunc)

// Replace with
token, err := jwt.Parse(tokenString, keyFunc, jwt.WithValidMethods([]string{"HS256"}))
```

---

##### Change 3: ParseWithClaims() Function

**OLD (jwt/v4)**:
```go
token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, keyFunc)
```

**NEW (jwt/v5)**:
```go
token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, keyFunc, 
    jwt.WithValidMethods([]string{"HS256"}))
```

---

##### Change 4: NewWithClaims() (No Changes)

**GOOD NEWS**: Token creation unchanged
```go
// jwt/v4 and jwt/v5 (same)
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, err := token.SignedString([]byte(secret))
```

---

#### Files to Update (Typical PocketBase-based app)

##### File: `core/auth.go` (or similar)
```go
// BEFORE (jwt/v4)
import "github.com/golang-jwt/jwt/v4"

func ValidateToken(tokenString string, secret string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })
    return token, err
}

// AFTER (jwt/v5)
import "github.com/golang-jwt/jwt/v5"

func ValidateToken(tokenString string, secret string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    }, jwt.WithValidMethods([]string{"HS256", "HS384", "HS512"}))
    return token, err
}
```

---

##### File: `apis/auth.go` (if exists)
```go
// BEFORE (jwt/v4)
import "github.com/golang-jwt/jwt/v4"

type AuthClaims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

func ParseAuthToken(tokenString string, secret string) (*AuthClaims, error) {
    claims := &AuthClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return []byte(secret), nil
    })
    // ... error handling ...
}

// AFTER (jwt/v5)
import "github.com/golang-jwt/jwt/v5"

type AuthClaims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

func ParseAuthToken(tokenString string, secret string) (*AuthClaims, error) {
    claims := &AuthClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    }, jwt.WithValidMethods([]string{"HS256"}))
    // ... error handling ...
}
```

---

#### Automated Migration Script

```bash
#!/bin/bash
# migrate-jwt-v5.sh

echo "Migrating JWT v4 → v5..."

# Step 1: Update imports
echo "Updating imports..."
find . -name "*.go" -exec sed -i 's|github.com/golang-jwt/jwt/v4|github.com/golang-jwt/jwt/v5|g' {} \;

# Step 2: Find Parse() calls (manual review required)
echo "Finding jwt.Parse() calls that need jwt.WithValidMethods()..."
grep -rn "jwt.Parse(" . --include="*.go" | grep -v "WithValidMethods"

echo ""
echo "⚠️  MANUAL REVIEW REQUIRED:"
echo "Add jwt.WithValidMethods() to all jwt.Parse() and jwt.ParseWithClaims() calls above"
echo ""
echo "Example:"
echo "  OLD: token, err := jwt.Parse(tokenString, keyFunc)"
echo "  NEW: token, err := jwt.Parse(tokenString, keyFunc, jwt.WithValidMethods([]string{\"HS256\"}))"
```

**Run Script**:
```bash
chmod +x migrate-jwt-v5.sh
./migrate-jwt-v5.sh > jwt-migration-report.txt
```

---

### Step 1.4: Verify Compilation

```bash
# Build all packages
go build -v ./...

# Expected output: Success (no errors)
```

**If Errors Occur**:
```bash
# Common error: jwt.Parse() missing algorithm validation
# Solution: Add jwt.WithValidMethods() as shown above

# Common error: jwt.StandardClaims not found
# Solution: Use jwt.RegisteredClaims instead
#   OLD: jwt.StandardClaims
#   NEW: jwt.RegisteredClaims
```

---

### Step 1.5: Run Unit Tests

```bash
# Run all tests
go test ./... -v

# Run only auth-related tests
go test ./core/... -v -run TestAuth
go test ./apis/... -v -run TestAuth
```

**Expected Results**:
- ✅ All tests pass
- ✅ JWT token validation works
- ✅ Token generation works

**If Tests Fail**:
1. Check JWT algorithm mismatch (ensure HS256 in both sign and verify)
2. Check claims structure (StandardClaims → RegisteredClaims)
3. Review error messages carefully

---

### Step 1.6: Manual Testing Checklist

- [ ] Start server: `go run main.go serve`
- [ ] Test user login (JWT generation)
- [ ] Test API call with JWT token (JWT validation)
- [ ] Test token expiration (should reject expired tokens)
- [ ] Test invalid tokens (should reject)
- [ ] Test tampered tokens (should reject)

---

### Step 1.7: Create Phase 1 Checkpoint

```bash
# Commit Phase 1 changes
git add .
git commit -m "Phase 1 Complete: PocketBase v0.37.5 base + JWT v5 migration"
git tag v0.2.1-phase1-complete

# Verify checkpoint
git log --oneline -3
git tag | grep phase1
```

---

### Step 1.8: Rollback Procedure (if Phase 1 fails)

```bash
# Rollback to pre-upgrade checkpoint
git reset --hard v0.2.0-pre-pocketbase-upgrade

# Verify rollback
go build -v ./...
go test ./...

# Alternative: Rollback specific package
go get github.com/golang-jwt/jwt/v4@v4.5.0
go mod tidy
```

---

## Phase 2: Add PostgreSQL Driver

### Objective
Add `lib/pq` driver and database connection logic WITHOUT changing existing SQLite functionality.

### Duration
**Estimated**: 1 hour  
**Breakdown**:
- Add dependency: 5 minutes
- Create db_postgresql.go: 30 minutes
- Testing: 20 minutes
- Verification: 5 minutes

---

### Step 2.1: Create Checkpoint

```bash
git add .
git commit -m "Checkpoint: Before PostgreSQL driver addition"
git tag v0.2.2-pre-postgres
```

---

### Step 2.2: Add lib/pq Dependency

```bash
# Add PostgreSQL driver
go get github.com/lib/pq@v1.10.9

# Verify go.mod
grep "lib/pq" go.mod
# Expected: github.com/lib/pq v1.10.9

# Clean up
go mod tidy
```

---

### Step 2.3: Create Database Connection Layer

**File**: `core/db_postgresql.go` (new file)

```go
package core

import (
    "database/sql"
    "fmt"
    "strings"
    
    _ "github.com/lib/pq" // PostgreSQL driver
)

// connectDB creates a database connection based on DSN prefix
func connectDB(dsn string) (*sql.DB, error) {
    if dsn == "" {
        return nil, fmt.Errorf("DSN is empty")
    }
    
    // Detect database type from DSN
    if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
        return connectPostgres(dsn)
    }
    
    // Default: assume SQLite (existing behavior)
    return nil, fmt.Errorf("unsupported DSN format: %s", dsn)
}

// connectPostgres creates a PostgreSQL connection
func connectPostgres(dsn string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open postgres connection: %w", err)
    }
    
    // Test connection
    if err := db.Ping(); err != nil {
        db.Close()
        return nil, fmt.Errorf("failed to ping postgres: %w", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    
    return db, nil
}
```

---

### Step 2.4: Verify Compilation

```bash
# Build with new file
go build -v ./core/...

# Expected: Success
```

---

### Step 2.5: Unit Test

**File**: `core/db_postgresql_test.go` (new file)

```go
package core

import (
    "os"
    "testing"
)

func TestConnectPostgres(t *testing.T) {
    // Skip if no PostgreSQL available
    dsn := os.Getenv("TEST_POSTGRES_DSN")
    if dsn == "" {
        t.Skip("TEST_POSTGRES_DSN not set")
    }
    
    db, err := connectPostgres(dsn)
    if err != nil {
        t.Fatalf("connectPostgres() error = %v", err)
    }
    defer db.Close()
    
    // Test query
    var result int
    err = db.QueryRow("SELECT 1").Scan(&result)
    if err != nil {
        t.Errorf("Query failed: %v", err)
    }
    
    if result != 1 {
        t.Errorf("Expected 1, got %d", result)
    }
}
```

**Run Test**:
```bash
# Set test DSN (requires PostgreSQL running)
export TEST_POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable"

# Run test
go test ./core/... -v -run TestConnectPostgres
```

---

### Step 2.6: Create Phase 2 Checkpoint

```bash
git add .
git commit -m "Phase 2 Complete: PostgreSQL driver added (lib/pq)"
git tag v0.2.3-phase2-complete
```

---

### Step 2.7: Rollback Procedure (if Phase 2 fails)

```bash
# Rollback to Phase 1 checkpoint
git reset --hard v0.2.1-phase1-complete

# Or remove just PostgreSQL
go mod edit -droprequire github.com/lib/pq
rm core/db_postgresql.go core/db_postgresql_test.go
go mod tidy
```

---

## Phase 3: Add Redis Client

### Objective
Add Redis client for realtime Pub/Sub, with graceful fallback to in-memory when Redis unavailable.

### Duration
**Estimated**: 2 hours  
**Breakdown**:
- Add dependency: 5 minutes
- Implement initRedis(): 30 minutes
- Implement Publish/Subscribe: 45 minutes
- Testing: 30 minutes
- Verification: 10 minutes

---

### Step 3.1: Create Checkpoint

```bash
git add .
git commit -m "Checkpoint: Before Redis client addition"
git tag v0.2.4-pre-redis
```

---

### Step 3.2: Add go-redis/v9 Dependency

```bash
# Add Redis client
go get github.com/redis/go-redis/v9@v9.3.0

# Verify go.mod
grep "redis/go-redis" go.mod
# Expected: github.com/redis/go-redis/v9 v9.3.0

# Clean up
go mod tidy
```

---

### Step 3.3: Add Redis Client to BaseApp

**File**: `core/base.go` (modify)

```go
package core

import (
    // ... existing imports ...
    "github.com/redis/go-redis/v9"
)

type BaseApp struct {
    // ... existing fields ...
    
    // NEW: Redis client (optional)
    redis redis.UniversalClient
}

// initRedis initializes Redis client (graceful fallback)
func (app *BaseApp) initRedis(dsn string) error {
    if dsn == "" {
        app.Logger().Info("Redis DSN not provided, using in-memory realtime")
        app.redis = nil
        return nil
    }
    
    // Parse DSN
    opt, err := redis.ParseURL(dsn)
    if err != nil {
        app.Logger().Warn("Invalid Redis DSN, using in-memory realtime", "error", err)
        app.redis = nil
        return nil
    }
    
    // Create client
    client := redis.NewClient(opt)
    
    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := client.Ping(ctx).Err(); err != nil {
        app.Logger().Warn("Redis connection failed, using in-memory realtime", "error", err)
        app.redis = nil
        return nil
    }
    
    app.redis = client
    app.Logger().Info("Redis connected successfully")
    return nil
}

// Publish publishes a realtime event (Redis or in-memory)
func (app *BaseApp) Publish(channel string, message string) error {
    if app.redis == nil {
        // Fallback to existing in-memory broadcast
        return app.publishInMemory(channel, message)
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    err := app.redis.Publish(ctx, channel, message).Err()
    if err != nil {
        // Graceful fallback
        app.Logger().Warn("Redis publish failed, using in-memory", "error", err)
        return app.publishInMemory(channel, message)
    }
    
    return nil
}
```

---

### Step 3.4: Verify Compilation

```bash
# Build with Redis support
go build -v ./core/...

# Expected: Success
```

---

### Step 3.5: Unit Test

**File**: `core/redis_test.go` (new file)

```go
package core

import (
    "context"
    "os"
    "testing"
    "time"
    
    "github.com/redis/go-redis/v9"
)

func TestInitRedis_Success(t *testing.T) {
    dsn := os.Getenv("TEST_REDIS_DSN")
    if dsn == "" {
        t.Skip("TEST_REDIS_DSN not set")
    }
    
    app := &BaseApp{}
    err := app.initRedis(dsn)
    if err != nil {
        t.Errorf("initRedis() error = %v", err)
    }
    
    if app.redis == nil {
        t.Error("Expected Redis client, got nil")
    }
}

func TestInitRedis_Fallback(t *testing.T) {
    // Invalid DSN should fallback gracefully
    app := &BaseApp{}
    err := app.initRedis("redis://invalid:9999/0")
    
    // Should not error (graceful fallback)
    if err != nil {
        t.Errorf("initRedis() should not error on invalid DSN, got: %v", err)
    }
    
    // Should fallback to nil (in-memory)
    if app.redis != nil {
        t.Error("Expected nil Redis client (fallback), got non-nil")
    }
}

func TestPublish_Redis(t *testing.T) {
    dsn := os.Getenv("TEST_REDIS_DSN")
    if dsn == "" {
        t.Skip("TEST_REDIS_DSN not set")
    }
    
    app := &BaseApp{}
    app.initRedis(dsn)
    
    err := app.Publish("test-channel", "test-message")
    if err != nil {
        t.Errorf("Publish() error = %v", err)
    }
}
```

**Run Test**:
```bash
# Set test DSN (requires Redis running)
export TEST_REDIS_DSN="redis://localhost:6379/0"

# Run test
go test ./core/... -v -run TestRedis
```

---

### Step 3.6: Integration Test (Optional)

```bash
# Start Redis with Docker
docker run -d --name test-redis -p 6379:6379 redis:7-alpine

# Run integration test
export TEST_REDIS_DSN="redis://localhost:6379/0"
go test ./core/... -v -run TestPublish

# Cleanup
docker stop test-redis
docker rm test-redis
```

---

### Step 3.7: Create Phase 3 Checkpoint

```bash
git add .
git commit -m "Phase 3 Complete: Redis client added with graceful fallback"
git tag v0.2.5-phase3-complete
```

---

### Step 3.8: Rollback Procedure (if Phase 3 fails)

```bash
# Rollback to Phase 2 checkpoint
git reset --hard v0.2.3-phase2-complete

# Or remove just Redis
go mod edit -droprequire github.com/redis/go-redis/v9
git checkout core/base.go  # Restore original
rm core/redis_test.go
go mod tidy
```

---

## Phase 4: Add MySQL Driver (Optional)

### Objective
Add MySQL driver for multi-database support (PostgreSQL + MySQL + SQLite).

### Duration
**Estimated**: 1 hour  
**Breakdown**:
- Add dependency: 5 minutes
- Update db_postgresql.go: 30 minutes
- Testing: 20 minutes
- Verification: 5 minutes

---

### Step 4.1: Create Checkpoint

```bash
git add .
git commit -m "Checkpoint: Before MySQL driver addition"
git tag v0.2.6-pre-mysql
```

---

### Step 4.2: Add go-sql-driver/mysql Dependency

```bash
# Add MySQL driver
go get github.com/go-sql-driver/mysql@v1.7.1

# Verify go.mod
grep "go-sql-driver/mysql" go.mod
# Expected: github.com/go-sql-driver/mysql v1.7.1

# Clean up
go mod tidy
```

---

### Step 4.3: Update Database Connection Layer

**File**: `core/db_postgresql.go` (modify)

```go
package core

import (
    "database/sql"
    "fmt"
    "strings"
    
    _ "github.com/lib/pq"                // PostgreSQL driver
    _ "github.com/go-sql-driver/mysql"   // MySQL driver
)

// connectDB creates a database connection based on DSN prefix
func connectDB(dsn string) (*sql.DB, error) {
    if dsn == "" {
        return nil, fmt.Errorf("DSN is empty")
    }
    
    // Detect database type from DSN
    if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
        return connectPostgres(dsn)
    } else if strings.HasPrefix(dsn, "mysql://") {
        return connectMySQL(dsn)
    }
    
    return nil, fmt.Errorf("unsupported DSN format: %s", dsn)
}

// connectMySQL creates a MySQL connection
func connectMySQL(dsn string) (*sql.DB, error) {
    // Strip mysql:// prefix (driver expects "user:pass@tcp(...)" format)
    dsn = strings.TrimPrefix(dsn, "mysql://")
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open mysql connection: %w", err)
    }
    
    // Test connection
    if err := db.Ping(); err != nil {
        db.Close()
        return nil, fmt.Errorf("failed to ping mysql: %w", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    
    return db, nil
}
```

---

### Step 4.4: Create Phase 4 Checkpoint

```bash
git add .
git commit -m "Phase 4 Complete: MySQL driver added (optional)"
git tag v0.2.7-phase4-complete
```

---

## Overall Rollback Strategy

### Checkpoint Summary

| Checkpoint | Tag | Description | Rollback Command |
|------------|-----|-------------|------------------|
| Pre-Migration | `v0.1.0-pre-migration` | Original codebase | `git reset --hard v0.1.0-pre-migration` |
| Pre-PocketBase | `v0.2.0-pre-pocketbase-upgrade` | Before any changes | `git reset --hard v0.2.0-pre-pocketbase-upgrade` |
| Phase 1 Done | `v0.2.1-phase1-complete` | PocketBase v0.37.5 + JWT v5 | `git reset --hard v0.2.1-phase1-complete` |
| Phase 2 Done | `v0.2.3-phase2-complete` | + PostgreSQL driver | `git reset --hard v0.2.3-phase2-complete` |
| Phase 3 Done | `v0.2.5-phase3-complete` | + Redis client | `git reset --hard v0.2.5-phase3-complete` |
| Phase 4 Done | `v0.2.7-phase4-complete` | + MySQL driver | `git reset --hard v0.2.7-phase4-complete` |

---

### Rollback Decision Matrix

| Failure Scenario | Rollback To | Command |
|------------------|-------------|---------|
| JWT v5 breaks authentication | `v0.2.0-pre-pocketbase-upgrade` | `git reset --hard v0.2.0-pre-pocketbase-upgrade` |
| PostgreSQL driver doesn't compile | `v0.2.1-phase1-complete` | `git reset --hard v0.2.1-phase1-complete` |
| Redis client causes crashes | `v0.2.3-phase2-complete` | `git reset --hard v0.2.3-phase2-complete` |
| MySQL driver conflicts | `v0.2.5-phase3-complete` | `git reset --hard v0.2.5-phase3-complete` |

---

### Partial Rollback (Remove Specific Dependency)

```bash
# Remove only PostgreSQL
go mod edit -droprequire github.com/lib/pq
rm core/db_postgresql.go core/db_postgresql_test.go
go mod tidy

# Remove only Redis
go mod edit -droprequire github.com/redis/go-redis/v9
git checkout core/base.go  # Restore original
rm core/redis_test.go
go mod tidy

# Remove only MySQL
go mod edit -droprequire github.com/go-sql-driver/mysql
# Update core/db_postgresql.go to remove MySQL support
go mod tidy

# Rollback JWT v5 → v4
go get github.com/golang-jwt/jwt/v4@v4.5.0
find . -name "*.go" -exec sed -i 's|jwt/v5|jwt/v4|g' {} \;
# Manually remove jwt.WithValidMethods() calls
go mod tidy
```

---

## Testing Strategy

### Per-Phase Testing

#### Phase 1: PocketBase v0.37.5 Base
```bash
# Unit tests
go test ./core/... -v
go test ./apis/... -v

# Integration tests
go test ./tests/integration/... -v

# Manual testing
go run main.go serve
# - Test login (JWT generation)
# - Test API calls (JWT validation)
# - Test all existing features
```

#### Phase 2: PostgreSQL Driver
```bash
# Unit test (requires PostgreSQL)
export TEST_POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable"
go test ./core/... -v -run TestPostgres

# Integration test
docker run -d --name test-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:16-alpine
export TEST_POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
go test ./core/... -v
docker stop test-postgres && docker rm test-postgres
```

#### Phase 3: Redis Client
```bash
# Unit test (requires Redis)
export TEST_REDIS_DSN="redis://localhost:6379/0"
go test ./core/... -v -run TestRedis

# Integration test
docker run -d --name test-redis -p 6379:6379 redis:7-alpine
export TEST_REDIS_DSN="redis://localhost:6379/0"
go test ./core/... -v
docker stop test-redis && docker rm test-redis

# Fallback test (no Redis)
unset TEST_REDIS_DSN
go test ./core/... -v -run TestRedis_Fallback
```

#### Phase 4: MySQL Driver
```bash
# Unit test (requires MySQL)
export TEST_MYSQL_DSN="root:root@tcp(localhost:3306)/test_db?parseTime=true"
go test ./core/... -v -run TestMySQL

# Integration test
docker run -d --name test-mysql -e MYSQL_ROOT_PASSWORD=root -p 3306:3306 mysql:8
sleep 10  # Wait for MySQL to start
export TEST_MYSQL_DSN="root:root@tcp(localhost:3306)/mysql?parseTime=true"
go test ./core/... -v
docker stop test-mysql && docker rm test-mysql
```

---

### Regression Testing (After All Phases)

```bash
# Full test suite
go test ./... -v -cover

# Coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Performance benchmarks
go test ./benchmarks/... -bench=. -benchmem
```

---

## Breaking Changes Summary

### JWT v4 → v5 (BREAKING)

| Change | Impact | Migration Effort |
|--------|--------|------------------|
| Import path | Update all imports | 5 minutes (automated) |
| `Parse()` signature | Add `WithValidMethods()` | 30 minutes (manual) |
| `StandardClaims` → `RegisteredClaims` | Update struct fields | 10 minutes (search/replace) |

**Total Effort**: ~45 minutes

---

### golang.org/x/crypto v0.12 → v0.50 (Minor Breaking)

| Change | Impact | Migration Effort |
|--------|--------|------------------|
| Internal API changes | Minimal (mostly internal) | 0 minutes |
| Deprecated functions | Compile warnings (not errors) | 0 minutes |

**Total Effort**: 0 minutes (no code changes needed)

---

### Other Dependencies (No Breaking Changes)

| Package | Change | Breaking? |
|---------|--------|-----------|
| fatih/color | v1.15.0 → v1.19.0 | ❌ No (backward compatible) |
| spf13/cobra | v1.7.0 → v1.10.2 | ❌ No (new features only) |
| modernc.org/sqlite | v1.25.0 → v1.50.0 | ❌ No (bug fixes + perf) |

---

## Risk Mitigation Strategies

### Risk 1: JWT v5 Migration Breaks Authentication

**Mitigation**:
- ✅ Comprehensive unit tests for auth flow
- ✅ Integration tests for token validation
- ✅ Manual testing of login/logout
- ✅ Git checkpoint before JWT changes
- ✅ Quick rollback script

**Rollback Plan**:
```bash
git reset --hard v0.2.0-pre-pocketbase-upgrade
go build && go test ./...
```

---

### Risk 2: PostgreSQL Driver Connection Failures

**Mitigation**:
- ✅ Test with real PostgreSQL instance
- ✅ Graceful error handling (don't crash on connection failure)
- ✅ Fallback to SQLite if PostgreSQL unavailable (future enhancement)
- ✅ Connection pool limits to prevent resource exhaustion

**Rollback Plan**:
```bash
git reset --hard v0.2.1-phase1-complete
go mod edit -droprequire github.com/lib/pq
go mod tidy
```

---

### Risk 3: Redis Client Causes Instability

**Mitigation**:
- ✅ Graceful fallback to in-memory (Redis optional)
- ✅ Timeout on all Redis operations (5 seconds)
- ✅ Error logging (warn, not error)
- ✅ Continue operation even if Redis fails

**Rollback Plan**:
```bash
git reset --hard v0.2.3-phase2-complete
# Or just disable Redis in config (set redisDsn="")
```

---

## Success Criteria

### Phase 1 Success Criteria
- [ ] All dependencies updated to PocketBase v0.37.5 versions
- [ ] JWT v4 → v5 migration complete
- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] Manual testing confirms auth works
- [ ] No performance regression

### Phase 2 Success Criteria
- [ ] lib/pq dependency added
- [ ] PostgreSQL connection logic implemented
- [ ] Unit tests pass (with PostgreSQL)
- [ ] DSN parsing works correctly
- [ ] Connection pool configured

### Phase 3 Success Criteria
- [ ] go-redis/v9 dependency added
- [ ] Redis client initialization implemented
- [ ] Graceful fallback to in-memory works
- [ ] Publish/Subscribe logic implemented
- [ ] Unit tests pass (with and without Redis)

### Phase 4 Success Criteria
- [ ] go-sql-driver/mysql dependency added
- [ ] MySQL connection logic implemented
- [ ] Unit tests pass (with MySQL)
- [ ] Multi-database support works (Postgres + MySQL + SQLite)

---

## Timeline

| Phase | Duration | Start | End | Status |
|-------|----------|-------|-----|--------|
| Phase 0: Planning | - | - | Complete | ✅ |
| Phase 1: PocketBase v0.37.5 | 2 hours | Day 1 | Day 1 | 🟢 Ready |
| Phase 2: PostgreSQL Driver | 1 hour | Day 1 | Day 1 | 🟡 Blocked (wait for Phase 1) |
| Phase 3: Redis Client | 2 hours | Day 2 | Day 2 | 🟡 Blocked (wait for Phase 2) |
| Phase 4: MySQL Driver | 1 hour | Day 2 | Day 2 | 🟡 Blocked (wait for Phase 3) |
| **Total** | **6 hours** | - | - | - |

---

## Next Steps

### Immediate Actions (AGENT-03 Complete)
- [x] Document version conflicts
- [x] Identify breaking changes
- [x] Create phased upgrade plan
- [x] Document rollback procedures
- [ ] Hand off to AGENT-05 (Database Connection Layer Engineer)

### Dependencies for Other Agents
- **AGENT-05** (DB Layer): Can start Phase 2 after Phase 1 complete
- **AGENT-06** (Redis): Can start Phase 3 after Phase 2 complete
- **AGENT-07** (CLI Flags): Can start after Phase 3 complete
- **AGENT-10** (Testing): Needs all phases complete for comprehensive testing

---

## Appendix: Quick Reference Commands

### Checkpoints
```bash
# Create checkpoint
git tag v0.X.Y-checkpoint-name

# List checkpoints
git tag | grep checkpoint

# Rollback to checkpoint
git reset --hard v0.X.Y-checkpoint-name
```

### Dependencies
```bash
# Add dependency
go get <package>@<version>

# Remove dependency
go mod edit -droprequire <package>

# Update dependency
go get -u <package>

# Clean up
go mod tidy

# Verify dependencies
go mod verify
go mod graph
```

### Testing
```bash
# Run all tests
go test ./... -v

# Run specific test
go test ./core/... -v -run TestName

# With coverage
go test ./... -cover -coverprofile=coverage.out

# View coverage
go tool cover -html=coverage.out
```

---

**Report Status**: ✅ Complete  
**Total Phases**: 4 (phased rollout)  
**Total Estimated Time**: 6 hours  
**Risk Level**: 🟡 MEDIUM (manageable with checkpoints)  
**Next Agent**: AGENT-05 (Database Connection Layer Engineer)
