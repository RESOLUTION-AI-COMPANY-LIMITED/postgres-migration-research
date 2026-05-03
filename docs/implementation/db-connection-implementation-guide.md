# ⚠️ IMPLEMENTATION GUIDE - NOT ACTUAL CODE

**This is a GUIDE on how to implement, NOT working code.**

The actual implementation has NOT been done yet. This document provides:
- ✅ Step-by-step instructions
- ✅ Code examples to copy
- ✅ Testing procedures

**To use this guide**: Follow steps in actual PocketBase fork repository.

---

# Database Connection Layer Implementation Guide (AGENT-05)

**Agent**: AGENT-05 - Database Connection Layer Engineer  
**Date**: 2026-05-03  
**Project**: PocketBase v0.37.5 → PostgreSQL + Redis Migration

---

## 📋 Overview

Hướng dẫn chi tiết để implement PostgreSQL/MySQL connection layer trong PocketBase. Đây là **RESEARCH PROJECT** - tạo documentation, KHÔNG phải actual implementation.

### Success Criteria
- ✅ `core/db_postgresql.go` implementation documented
- ✅ `core/base.go` modifications mapped
- ✅ `go.mod` changes specified
- ✅ Compilation verification steps defined
- ✅ Error handling strategy documented

---

## 🎯 Phase 2 - AGENT-05 Tasks

### Dependencies
- ✅ AGENT-02 (Code Extraction) - Complete
- ✅ AGENT-03 (Dependency Analysis) - Complete
- 🟢 Ready to start

### Estimated Time
6 hours (research documentation)

---

## 📦 STEP 1: Create `core/db_postgresql.go`

### 1.1 File Location
```bash
# Target file path (DOCUMENTATION ONLY - don't create actual file)
rai-backend/core/db_postgresql.go
```

### 1.2 Complete Implementation (Reference)

```go
package core

import (
	"strings"

	// PostgreSQL driver - imported for side effects (registers "postgres" driver)
	_ "github.com/lib/pq"

	// MySQL driver - imported for side effects (registers "mysql" driver)
	_ "github.com/go-sql-driver/mysql"

	// Modified dbx query builder with PostgreSQL/MySQL support
	// NOTE: This is a custom fork, not the original go-ozzo/ozzo-dbx
	"github.com/free/postgresqlbaseapi/dbx"
)

// connectDB creates a database connection from a DSN string
//
// Design Pattern: Factory Pattern with Auto-Detection
// - Detects driver type from DSN prefix
// - Supports: postgres://, postgresql://, mysql://
// - Defaults to postgres if no prefix
//
// Parameters:
//   dsn - Data Source Name (connection string)
//         Examples:
//         - "postgres://user:pass@localhost/dbname"
//         - "postgresql://user:pass@localhost/dbname"
//         - "mysql://user:pass@localhost/dbname"
//
// Returns:
//   *dbx.DB - Database connection with query builder
//   error   - Connection error if any
//
// Usage in PocketBase:
//   This replaces the original SQLite-only connection logic in core/base.go
//   Original: db, err := connectDB(dataDir)
//   Modified: db, err := connectDB(dataDsn) // DSN passed from CLI flag
func connectDB(dsn string) (*dbx.DB, error) {
	// Default driver is PostgreSQL
	driver := "postgres"

	// Parse driver from DSN prefix
	// NOTE: This allows users to specify driver explicitly via DSN
	if strings.HasPrefix(dsn, "mysql://") {
		driver = "mysql"
		// Remove mysql:// prefix for driver compatibility
		dsn = strings.TrimPrefix(dsn, "mysql://")
	} else if strings.HasPrefix(dsn, "postgres://") {
		driver = "postgres"
		// Keep postgres:// prefix - lib/pq expects it
	} else if strings.HasPrefix(dsn, "postgresql://") {
		driver = "postgres"
		// postgresql:// is alias for postgres://
	}
	// If no prefix, assumes PostgreSQL with raw DSN

	// Open database connection using dbx query builder
	// dbx.Open() wraps database/sql Open() and adds query builder
	db, err := dbx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	// TODO: Connection pooling configuration (production recommendation)
	// db.DB().SetMaxOpenConns(100)
	// db.DB().SetMaxIdleConns(10)
	// db.DB().SetConnMaxLifetime(time.Hour)

	return db, nil
}
```

### 1.3 Implementation Notes

**Driver Detection Logic:**
| DSN Prefix | Detected Driver | DSN Processing |
|------------|----------------|----------------|
| `mysql://` | `mysql` | Remove prefix |
| `postgres://` | `postgres` | Keep prefix |
| `postgresql://` | `postgres` | Keep prefix |
| (no prefix) | `postgres` | Use raw DSN |

**Connection Pool Settings (Production):**
```go
// Add to connectDB() function after db, err := dbx.Open(...)
db.DB().SetMaxOpenConns(100)      // Max simultaneous connections
db.DB().SetMaxIdleConns(10)       // Keep idle connections ready
db.DB().SetConnMaxLifetime(time.Hour) // Recycle connections every hour
```

### 1.4 Testing Strategy

```go
// File: core/db_postgresql_test.go (future test)
package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConnectDB_PostgreSQL(t *testing.T) {
	dsn := "postgres://localhost/test?sslmode=disable"
	db, err := connectDB(dsn)
	
	assert.NoError(t, err)
	assert.NotNil(t, db)
	
	// Verify driver
	assert.Equal(t, "postgres", db.DriverName())
}

func TestConnectDB_MySQL(t *testing.T) {
	dsn := "mysql://root:pass@tcp(localhost:3306)/test"
	db, err := connectDB(dsn)
	
	assert.NoError(t, err)
	assert.NotNil(t, db)
	
	// Verify driver
	assert.Equal(t, "mysql", db.DriverName())
}

func TestConnectDB_InvalidDSN(t *testing.T) {
	dsn := "invalid://localhost/test"
	db, err := connectDB(dsn)
	
	assert.Error(t, err)
	assert.Nil(t, db)
}
```

---

## 🔧 STEP 2: Update `core/base.go` Struct

### 2.1 Locate Existing BaseApp Struct

**File:** `rai-backend/core/base.go` (existing file)

**Original Struct (typical PocketBase structure):**
```go
type BaseApp struct {
	// Existing fields (keep as-is)
	dataDir           string
	encryptionEnv     string
	isDebug           bool
	db                *dbx.DB
	dao               *daos.Dao
	store             *store.Store<any>
	settings          *settings.Settings
	cache             *store.Store<any>
	
	// Event hooks
	onBeforeServe     *hook.Hook[*ServeEvent]
	onRealtimeConnect *hook.Hook[*RealtimeConnectEvent]
	// ... more hooks ...
}
```

### 2.2 Add New Fields for PostgreSQL + Redis

**Modified Struct:**
```go
type BaseApp struct {
	// Existing fields (keep as-is)
	dataDir           string
	encryptionEnv     string
	isDebug           bool
	db                *dbx.DB
	dao               *daos.Dao
	store             *store.Store<any>
	settings          *settings.Settings
	cache             *store.Store<any>
	
	// Event hooks
	onBeforeServe     *hook.Hook[*ServeEvent]
	onRealtimeConnect *hook.Hook[*RealtimeConnectEvent]
	// ... more hooks ...
	
	// =========================================================================
	// NEW: PostgreSQL + Redis Support
	// =========================================================================
	
	// Database DSN field
	dataDsn           string         // PostgreSQL/MySQL connection string
	
	// Redis integration
	redisDsn          string         // Redis connection string (optional)
	redisCache        *redis.Client  // Redis client instance (nil if disabled)
	redisContext      context.Context // Context for Redis operations
}
```

### 2.3 Update BaseAppConfig Struct

**Location:** Same file (`core/base.go`)

**Original Config:**
```go
type BaseAppConfig struct {
	DataDir       string
	EncryptionEnv string
	IsDebug       bool
}
```

**Modified Config:**
```go
type BaseAppConfig struct {
	// Existing fields (keep as-is)
	DataDir       string
	EncryptionEnv string
	IsDebug       bool
	
	// =========================================================================
	// NEW: PostgreSQL + Redis Configuration
	// =========================================================================
	
	DataDsn       string // PostgreSQL/MySQL DSN (e.g., "postgres://user:pass@localhost/db")
	RedisDsn      string // Redis DSN (optional, e.g., "redis://localhost:6379/0")
}
```

### 2.4 Update NewBaseApp() Constructor

**Original Constructor:**
```go
func NewBaseApp(config BaseAppConfig) *BaseApp {
	app := &BaseApp{
		dataDir:       config.DataDir,
		encryptionEnv: config.EncryptionEnv,
		isDebug:       config.IsDebug,
		store:         store.New[any](nil),
		cache:         store.New[any](nil),
	}
	
	// Initialize hooks
	app.initHooks()
	
	return app
}
```

**Modified Constructor:**
```go
func NewBaseApp(config BaseAppConfig) *BaseApp {
	app := &BaseApp{
		// Existing fields (keep as-is)
		dataDir:       config.DataDir,
		encryptionEnv: config.EncryptionEnv,
		isDebug:       config.IsDebug,
		store:         store.New[any](nil),
		cache:         store.New[any](nil),
		
		// =========================================================================
		// NEW: PostgreSQL + Redis Configuration
		// =========================================================================
		
		dataDsn:       config.DataDsn,      // PostgreSQL/MySQL DSN
		redisDsn:      config.RedisDsn,     // Redis DSN (optional)
		redisContext:  context.Background(), // Redis operation context
	}
	
	// Initialize hooks
	app.initHooks()
	
	return app
}
```

### 2.5 Update Bootstrap() Method

**Original Bootstrap:**
```go
func (app *BaseApp) Bootstrap() error {
	// Connect to SQLite database
	db, err := connectDB(app.dataDir + "/pb_data.db")
	if err != nil {
		return err
	}
	app.db = db
	
	// Initialize DAO, settings, etc.
	// ...
	
	return nil
}
```

**Modified Bootstrap:**
```go
func (app *BaseApp) Bootstrap() error {
	// =========================================================================
	// STEP 1: Connect to Database (PostgreSQL/MySQL/SQLite)
	// =========================================================================
	
	var db *dbx.DB
	var err error
	
	// If dataDsn provided, use PostgreSQL/MySQL
	if app.dataDsn != "" {
		// Use new connectDB() from db_postgresql.go
		db, err = connectDB(app.dataDsn)
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		
		if app.IsDebug() {
			color.Green("PostgreSQL/MySQL connected: %s", app.dataDsn)
		}
	} else {
		// Fallback to SQLite (original behavior)
		sqlitePath := app.dataDir + "/pb_data.db"
		db, err = connectDB("sqlite://" + sqlitePath)
		if err != nil {
			return fmt.Errorf("failed to connect to SQLite: %w", err)
		}
		
		if app.IsDebug() {
			color.Yellow("SQLite mode: %s", sqlitePath)
		}
	}
	
	app.db = db
	
	// =========================================================================
	// STEP 2: Initialize Redis (optional, for clustering)
	// =========================================================================
	
	if err := app.initRedis(); err != nil {
		// Log warning but don't fail (graceful degradation)
		if app.IsDebug() {
			color.Yellow("Redis initialization failed: %v", err)
			color.Yellow("Continuing in single-node mode (no clustering)")
		}
	}
	
	// =========================================================================
	// STEP 3: Initialize DAO, settings, etc. (existing logic)
	// =========================================================================
	
	// Initialize DAO
	app.dao = daos.New(app.db)
	
	// Initialize settings
	app.settings, err = settings.New(app.dao)
	if err != nil {
		return err
	}
	
	// ... rest of existing Bootstrap() logic ...
	
	return nil
}
```

---

## 📦 STEP 3: Update `go.mod`

### 3.1 Add PostgreSQL Driver Dependency

**Command (documentation only):**
```bash
go get github.com/lib/pq@v1.10.9
```

**Resulting `go.mod` changes:**
```go
module github.com/pocketbase/pocketbase

go 1.25.0

require (
	// Existing dependencies...
	
	// NEW: PostgreSQL driver
	github.com/lib/pq v1.10.9
)
```

### 3.2 Add MySQL Driver Dependency

**Command (documentation only):**
```bash
go get github.com/go-sql-driver/mysql@v1.7.1
```

**Resulting `go.mod` changes:**
```go
require (
	// Existing dependencies...
	
	// NEW: Database drivers
	github.com/lib/pq v1.10.9
	github.com/go-sql-driver/mysql v1.7.1
)
```

### 3.3 Add Redis Client Dependency (for AGENT-06)

**Command (documentation only):**
```bash
go get github.com/redis/go-redis/v9@v9.3.0
```

**Resulting `go.mod` changes:**
```go
require (
	// Existing dependencies...
	
	// NEW: Database drivers
	github.com/lib/pq v1.10.9
	github.com/go-sql-driver/mysql v1.7.1
	
	// NEW: Redis client
	github.com/redis/go-redis/v9 v9.3.0
)
```

### 3.4 Update dbx Dependency (for AGENT-08)

**See AGENT-08 documentation for dbx vendoring strategy**

**Option 1: Use replace directive (recommended):**
```go
replace github.com/go-ozzo/ozzo-dbx => github.com/free/postgresqlbaseapi/dbx v1.0.0
```

**Option 2: Vendor dbx package (alternative):**
```bash
go mod vendor
# Manual copy of dbx fork to vendor/
```

### 3.5 Run go mod tidy

**Command (documentation only):**
```bash
go mod tidy
```

**Expected Output:**
```
go: downloading github.com/lib/pq v1.10.9
go: downloading github.com/go-sql-driver/mysql v1.7.1
go: downloading github.com/redis/go-redis/v9 v9.3.0
```

---

## 🧪 STEP 4: Compilation Testing

### 4.1 Initial Compilation Test

**Command (documentation only):**
```bash
cd rai-backend
go build -v ./... 2>&1 | tee ../docs/compilation-report.md
```

### 4.2 Expected Compilation Errors (First Run)

**Error 1: Missing Redis import in base.go**
```
core/base.go:15:2: undefined: redis
```

**Fix:**
```go
import (
	// ... existing imports ...
	"github.com/redis/go-redis/v9" // ADD THIS
)
```

**Error 2: Missing context import in base.go**
```
core/base.go:42:18: undefined: context
```

**Fix:**
```go
import (
	"context" // ADD THIS
	// ... existing imports ...
)
```

**Error 3: Missing dbx package**
```
core/db_postgresql.go:14:2: cannot find package "github.com/free/postgresqlbaseapi/dbx"
```

**Fix:** Wait for AGENT-08 (dbx vendoring)

### 4.3 Compilation Verification Checklist

**After all fixes applied:**

```bash
# Step 1: Clean build
go clean -cache

# Step 2: Verify dependencies
go mod verify

# Step 3: Build with verbose output
go build -v ./rai-backend/...

# Step 4: Run go vet
go vet ./rai-backend/...

# Step 5: Check imports
go list -f '{{.Imports}}' ./rai-backend/core

# Expected imports (partial list):
# [context github.com/lib/pq github.com/go-sql-driver/mysql github.com/redis/go-redis/v9 ...]
```

### 4.4 Success Criteria

✅ **Compilation successful:**
```
# go build output (no errors)
```

✅ **All imports resolved:**
```bash
$ go list -f '{{.Imports}}' ./rai-backend/core | grep -E 'pq|mysql|redis'
github.com/lib/pq
github.com/go-sql-driver/mysql
github.com/redis/go-redis/v9
```

✅ **No vet warnings:**
```bash
$ go vet ./rai-backend/core
(no output = success)
```

---

## 🧪 STEP 5: Manual Testing Strategy

### 5.1 Test PostgreSQL Connection

**Create test script:** `scripts/test-postgres-connection.go`

```go
package main

import (
	"fmt"
	"log"
	
	_ "github.com/lib/pq"
	"github.com/free/postgresqlbaseapi/dbx"
)

func main() {
	// Test DSN (update with your credentials)
	dsn := "postgres://postgres:password@localhost:5432/testdb?sslmode=disable"
	
	db, err := dbx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer db.Close()
	
	// Test query
	var version string
	err = db.NewQuery("SELECT version()").Row(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	
	fmt.Printf("✅ PostgreSQL connected successfully\n")
	fmt.Printf("Version: %s\n", version)
}
```

**Run test:**
```bash
# Start PostgreSQL (Docker)
docker run -d --name postgres-test \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=testdb \
  -p 5432:5432 \
  postgres:15

# Run test script
go run scripts/test-postgres-connection.go
```

**Expected output:**
```
✅ PostgreSQL connected successfully
Version: PostgreSQL 15.4 on x86_64-pc-linux-gnu...
```

### 5.2 Test MySQL Connection

**Create test script:** `scripts/test-mysql-connection.go`

```go
package main

import (
	"fmt"
	"log"
	
	_ "github.com/go-sql-driver/mysql"
	"github.com/free/postgresqlbaseapi/dbx"
)

func main() {
	// Test DSN (update with your credentials)
	dsn := "root:password@tcp(localhost:3306)/testdb?parseTime=true"
	
	db, err := dbx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer db.Close()
	
	// Test query
	var version string
	err = db.NewQuery("SELECT VERSION()").Row(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	
	fmt.Printf("✅ MySQL connected successfully\n")
	fmt.Printf("Version: %s\n", version)
}
```

**Run test:**
```bash
# Start MySQL (Docker)
docker run -d --name mysql-test \
  -e MYSQL_ROOT_PASSWORD=password \
  -e MYSQL_DATABASE=testdb \
  -p 3306:3306 \
  mysql:8.0

# Run test script
go run scripts/test-mysql-connection.go
```

**Expected output:**
```
✅ MySQL connected successfully
Version: 8.0.35
```

### 5.3 Test Driver Auto-Detection

**Create test script:** `scripts/test-driver-detection.go`

```go
package main

import (
	"fmt"
	"testing"
)

func TestDriverDetection(t *testing.T) {
	tests := []struct {
		name        string
		dsn         string
		wantDriver  string
		wantDsnStrip string
	}{
		{
			name:        "PostgreSQL with postgres:// prefix",
			dsn:         "postgres://localhost/db",
			wantDriver:  "postgres",
			wantDsnStrip: "postgres://localhost/db", // Keep prefix
		},
		{
			name:        "PostgreSQL with postgresql:// prefix",
			dsn:         "postgresql://localhost/db",
			wantDriver:  "postgres",
			wantDsnStrip: "postgresql://localhost/db", // Keep prefix
		},
		{
			name:        "MySQL with mysql:// prefix",
			dsn:         "mysql://root@localhost/db",
			wantDriver:  "mysql",
			wantDsnStrip: "root@localhost/db", // Strip prefix
		},
		{
			name:        "PostgreSQL with no prefix",
			dsn:         "host=localhost dbname=db",
			wantDriver:  "postgres",
			wantDsnStrip: "host=localhost dbname=db", // Keep raw
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver, dsn := detectDriver(tt.dsn)
			
			if driver != tt.wantDriver {
				t.Errorf("driver = %v, want %v", driver, tt.wantDriver)
			}
			
			if dsn != tt.wantDsnStrip {
				t.Errorf("dsn = %v, want %v", dsn, tt.wantDsnStrip)
			}
		})
	}
}

func detectDriver(dsn string) (driver, strippedDsn string) {
	driver = "postgres"
	strippedDsn = dsn
	
	if strings.HasPrefix(dsn, "mysql://") {
		driver = "mysql"
		strippedDsn = strings.TrimPrefix(dsn, "mysql://")
	} else if strings.HasPrefix(dsn, "postgres://") {
		driver = "postgres"
	} else if strings.HasPrefix(dsn, "postgresql://") {
		driver = "postgres"
	}
	
	return driver, strippedDsn
}
```

---

## 🐛 Error Handling Strategy

### Common Connection Errors

**Error 1: Invalid DSN format**
```
Error: pq: invalid DSN
```

**Root Cause:** Malformed connection string

**Solution:**
```go
// Valid PostgreSQL DSN formats
"postgres://user:pass@localhost:5432/db?sslmode=disable"
"host=localhost port=5432 user=user password=pass dbname=db sslmode=disable"

// Valid MySQL DSN formats
"mysql://user:pass@tcp(localhost:3306)/db?parseTime=true"
"user:pass@tcp(localhost:3306)/db"
```

**Error 2: Connection refused**
```
Error: dial tcp [::1]:5432: connect: connection refused
```

**Root Cause:** Database not running or wrong port

**Solution:**
```bash
# Check PostgreSQL status
docker ps | grep postgres

# Start PostgreSQL
docker start postgres-container

# Verify port
netstat -an | grep 5432
```

**Error 3: Authentication failed**
```
Error: pq: password authentication failed for user "postgres"
```

**Root Cause:** Wrong credentials in DSN

**Solution:**
```go
// Double-check DSN credentials
dsn := "postgres://correct_user:correct_password@localhost/db"
```

**Error 4: Database does not exist**
```
Error: pq: database "mydb" does not exist
```

**Root Cause:** Target database not created

**Solution:**
```bash
# Create database manually
docker exec -it postgres-container psql -U postgres -c "CREATE DATABASE mydb"
```

---

## 📊 Performance Considerations

### Connection Pool Tuning

**Production Recommendations:**
```go
// In connectDB() function, after db, err := dbx.Open(...)

// Set connection pool limits
db.DB().SetMaxOpenConns(100)      // Maximum open connections (adjust based on load)
db.DB().SetMaxIdleConns(10)       // Idle connections to keep ready
db.DB().SetConnMaxLifetime(time.Hour) // Recycle connections every hour

// Log pool stats (debug mode)
if app.IsDebug() {
	stats := db.DB().Stats()
	log.Printf("DB Pool Stats: OpenConns=%d, InUse=%d, Idle=%d",
		stats.OpenConnections, stats.InUse, stats.Idle)
}
```

**Load Testing Script:**
```bash
# Test concurrent connections
for i in {1..100}; do
  curl http://localhost:8090/api/collections &
done
wait

# Check pool stats in logs
grep "DB Pool Stats" logs/pocketbase.log
```

---

## 📁 Deliverables Summary

### Files Created (Documentation)
1. ✅ `docs/implementation/db-connection-implementation-guide.md` (this file)
2. 📄 `core/db_postgresql.go` (implementation reference)
3. 📄 `core/base.go` (modification guide)
4. 📄 `go.mod` (dependency changes)
5. 📄 `docs/compilation-report.md` (compilation verification)

### Code Changes Required (Summary)
1. **New File:** `core/db_postgresql.go` (~90 lines)
2. **Modified File:** `core/base.go` (~50 lines changed)
   - Add fields: `dataDsn`, `redisDsn`, `redisCache`, `redisContext`
   - Update `BaseAppConfig` struct
   - Update `NewBaseApp()` constructor
   - Update `Bootstrap()` method
3. **Modified File:** `go.mod`
   - Add `github.com/lib/pq v1.10.9`
   - Add `github.com/go-sql-driver/mysql v1.7.1`

### Testing Scripts (Future)
1. `scripts/test-postgres-connection.go`
2. `scripts/test-mysql-connection.go`
3. `scripts/test-driver-detection.go`

---

## ✅ Verification Checklist

**Before marking AGENT-05 complete:**

- [ ] `db_postgresql.go` implementation documented
- [ ] `base.go` modification guide complete
- [ ] `go.mod` changes specified
- [ ] Compilation error fixes documented
- [ ] Error handling strategy defined
- [ ] Performance tuning recommendations provided
- [ ] Testing scripts documented
- [ ] Integration points with AGENT-06 (Redis) identified
- [ ] Integration points with AGENT-07 (CLI flags) identified
- [ ] Integration points with AGENT-08 (dbx) identified

---

## 🔗 Dependencies & Handoffs

### Upstream Dependencies (Completed)
- ✅ **AGENT-02** - Extracted `db_postgresql.go` code
- ✅ **AGENT-03** - Dependency version analysis (`lib/pq` v1.10.9, `mysql` v1.7.1)

### Downstream Dependencies (Blocked Until Complete)
- 🟡 **AGENT-06** - Redis integration (needs `base.go` struct changes)
- 🟡 **AGENT-07** - CLI flags (needs `BaseAppConfig` struct changes)
- 🟡 **AGENT-08** - dbx vendoring (needs `connectDB()` implementation)

### Parallel Work (Can Start Simultaneously)
- 🟢 **AGENT-06** - Can document Redis integration (overlapping work)
- 🟢 **AGENT-07** - Can document CLI flags (overlapping work)

---

## 📌 Next Steps

**After AGENT-05 completion:**

1. **Handoff to AGENT-06:** Redis Integration Guide
   - Use `redisCache`, `redisDsn`, `redisContext` fields from this guide
   - Implement `initRedis()`, `Publish()`, `Subscribe()` methods

2. **Handoff to AGENT-07:** CLI Flags Guide
   - Add `--dataDsn` flag mapping to `BaseAppConfig.DataDsn`
   - Add `--redisDsn` flag mapping to `BaseAppConfig.RedisDsn`

3. **Handoff to AGENT-08:** dbx Vendoring Guide
   - Replace `connectDB()` dbx import path
   - Vendor custom dbx fork
   - Update import statements across codebase

---

**AGENT-05 Status:** 🟢 Ready for Review  
**Document Version:** 1.0  
**Last Updated:** 2026-05-03
