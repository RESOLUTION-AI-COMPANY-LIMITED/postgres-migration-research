# PocketBase v0.37.5 → PostgreSQL + Redis Migration Plan

**Document Version**: 1.0  
**Created**: 2026-05-03  
**Target**: Migrate PocketBase v0.37.5 to PostgreSQL/MySQL + Redis  
**Based On**: Postgrebase methodology analysis  
**Timeline**: 2 tuần (10 ngày làm việc)

---

## Executive Summary

Kế hoạch này mô tả các bước chi tiết để chuyển đổi PocketBase v0.37.5 (latest stable) sang PostgreSQL + Redis, dựa trên cách postgrebase đã thực hiện thành công.

**Phương pháp**: Selective backporting - Chỉ thay đổi tối thiểu (~200 dòng code core) để đạt được PostgreSQL/MySQL support.

**Timeline**: 2 tuần  
**Risk Level**: 🟡 Trung bình (có thể quản lý được)  
**Expected Outcome**: Production-ready PostgreSQL/MySQL backend với Redis clustering

---

## Phase 1: Preparation & Analysis (Ngày 1-2)

### Day 1 Morning: Repository Setup

**Tasks**:
1. Backup current codebase
2. Clone PocketBase v0.37.5 source
3. Compare với rai-backend hiện tại
4. Document differences

**Commands**:
```bash
# Backup current code
cd ~/Documents/postgrebase-super-backend
git add -A
git commit -m "Backup before v0.37.5 PostgreSQL migration"
git tag v0.1.0-pre-migration

# Clone PocketBase v0.37.5 for reference
cd /tmp
git clone --depth 1 --branch v0.37.5 https://github.com/pocketbase/pocketbase.git pocketbase-v0.37.5

# Compare structure
diff -qr rai-backend/ /tmp/pocketbase-v0.37.5/ | grep -E "core|cmd" > ~/Documents/postgrebase-super-backend/docs/v0375-diff.txt
```

**Deliverables**:
- [ ] Backup commit created
- [ ] PocketBase v0.37.5 cloned
- [ ] Diff report generated
- [ ] Version differences documented

### Day 1 Afternoon: Postgrebase Code Extraction

**Tasks**:
1. Extract `core/db_postgresql.go` from postgrebase
2. Extract Redis logic from `core/base.go`
3. Extract cmd flags from `cmd/serve.go`
4. Document all extracted code

**Commands**:
```bash
# Extract key files from postgrebase
cd /tmp/postgrebase

# Copy db_postgresql.go
cp core/db_postgresql.go ~/Documents/postgrebase-super-backend/docs/extracted/

# Extract Redis functions from base.go
grep -A 50 "func.*initRedis" core/base.go > ~/Documents/postgrebase-super-backend/docs/extracted/redis_functions.go

# Extract Publish function
grep -A 20 "func.*Publish" core/base.go >> ~/Documents/postgrebase-super-backend/docs/extracted/redis_functions.go
```

**Deliverables**:
- [ ] `docs/extracted/db_postgresql.go` (29 lines)
- [ ] `docs/extracted/redis_functions.go` (~150 lines)
- [ ] `docs/extracted/cmd_serve_flags.go` (DSN flag definitions)
- [ ] Code annotation comments added

### Day 2 Morning: Dependency Analysis

**Tasks**:
1. Check postgrebase's `go.mod` dependencies
2. Compare với PocketBase v0.37.5 `go.mod`
3. List all new dependencies needed
4. Check for version conflicts

**Commands**:
```bash
# Compare go.mod files
cd /tmp/postgrebase
grep -E "lib/pq|go-sql-driver|redis" go.mod > /tmp/postgrebase-deps.txt

cd /tmp/pocketbase-v0.37.5
grep -E "lib/pq|go-sql-driver|redis" go.mod > /tmp/pocketbase-deps.txt

# Check our current dependencies
cd ~/Documents/postgrebase-super-backend/rai-backend
grep -E "lib/pq|go-sql-driver|redis" go.mod > /tmp/our-deps.txt

# Compare all three
diff /tmp/postgrebase-deps.txt /tmp/our-deps.txt
```

**Expected New Dependencies**:
```go
require (
	github.com/lib/pq v1.10.9              // PostgreSQL driver
	github.com/go-sql-driver/mysql v1.7.1  // MySQL driver (optional)
	github.com/redis/go-redis/v9 v9.0.5    // Redis client
)
```

**Deliverables**:
- [ ] Dependency comparison matrix
- [ ] Version conflict report
- [ ] go.mod update plan

### Day 2 Afternoon: dbx Package Analysis

**Tasks**:
1. Identify dbx package location in postgrebase
2. Compare với original dbx package
3. Document all modifications
4. Plan vendoring strategy

**Commands**:
```bash
# Find dbx package
cd /tmp/postgrebase
find vendor/ -name "dbx" -type d

# Compare with original (if available)
# Note: postgrebase uses custom fork at github.com/free/postgresqlbaseapi/dbx

# Check import paths
grep -r "postgresqlbaseapi/dbx" . --include="*.go" | head -20
```

**Key Questions**:
- ❓ Có phải vendor toàn bộ dbx package?
- ❓ Có thể dùng original dbx được không?
- ❓ Cần modify import paths ở bao nhiêu files?

**Deliverables**:
- [ ] dbx modification report
- [ ] Vendoring strategy document
- [ ] Import path change plan

---

## Phase 2: Core Implementation (Ngày 3-5)

### Day 3: Database Connection Layer

**Tasks**:
1. Add `core/db_postgresql.go` to rai-backend
2. Update `core/base.go` struct definitions
3. Update `go.mod` with new dependencies
4. Test compilation

**Implementation Steps**:

#### Step 3.1: Create `rai-backend/core/db_postgresql.go`
```go
package core

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/free/postgresqlbaseapi/dbx"  // TODO: Update import path
)

func connectDB(dsn string) (*dbx.DB, error) {
	driver := "postgres"
	// Parse driver from DSN prefix
	if strings.HasPrefix(dsn, "mysql://") {
		driver = "mysql"
		dsn = strings.TrimPrefix(dsn, "mysql://")
	} else if strings.HasPrefix(dsn, "postgres://") {
		driver = "postgres"
	} else if strings.HasPrefix(dsn, "postgresql://") {
		driver = "postgres"
	}

	db, err := dbx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
```

#### Step 3.2: Update `rai-backend/core/base.go`
```go
// Add to BaseApp struct (line ~44):
type BaseApp struct {
	// ... existing fields ...
	
	dataDsn          string          // ← ADD THIS
	redisDsn         string          // ← ADD THIS
	redisCache       *redis.Client   // ← ADD THIS
	redisContext     context.Context // ← ADD THIS
	
	// ... rest unchanged ...
}

// Add to BaseAppConfig struct (line ~176):
type BaseAppConfig struct {
	// ... existing fields ...
	
	DataDsn          string   // ← ADD THIS
	RedisDsn         string   // ← ADD THIS
	
	// ... rest unchanged ...
}

// Update NewBaseApp function (line ~193):
func NewBaseApp(config BaseAppConfig) *BaseApp {
	app := &BaseApp{
		// ... existing fields ...
		
		dataDsn:     config.DataDsn,   // ← ADD THIS
		redisDsn:    config.RedisDsn,  // ← ADD THIS
		
		// ... rest unchanged ...
	}
	// ... rest unchanged ...
}
```

#### Step 3.3: Update `go.mod`
```bash
cd rai-backend
go get github.com/lib/pq@v1.10.9
go get github.com/go-sql-driver/mysql@v1.7.1
go get github.com/redis/go-redis/v9@v9.0.5
go mod tidy
```

#### Step 3.4: Test Compilation
```bash
cd rai-backend
go build -o pb ./
# Expected: Compilation successful
```

**Deliverables**:
- [ ] `core/db_postgresql.go` created (29 lines)
- [ ] `core/base.go` updated (struct changes)
- [ ] `go.mod` updated (3 new dependencies)
- [ ] Code compiles successfully

**Rollback Plan**: `git reset --hard HEAD` if compilation fails

### Day 4: Redis Integration

**Tasks**:
1. Add `initRedis()` function to `core/base.go`
2. Add `Publish()` function for Redis Pub/Sub
3. Update `Bootstrap()` to initialize Redis
4. Test Redis connection

**Implementation Steps**:

#### Step 4.1: Add `initRedis()` function
```go
// Add to core/base.go (after initDataDB function, ~line 1061)
func (app *BaseApp) initRedis() error {
	if app.redisDsn == "" {
		return nil  // Skip if Redis not configured
	}

	opt, err := redis.ParseURL(app.redisDsn)
	if err != nil {
		return err
	}

	app.redisCache = redis.NewClient(opt)

	// Test connection (5 second timeout)
	ctx, cancel := context.WithTimeout(app.redisContext, 5*time.Second)
	defer cancel()
	if _, err := app.redisCache.Ping(ctx).Result(); err != nil {
		app.redisCache = nil // Disable Redis if connection fails
		if app.IsDebug() {
			color.Red("Redis connection failed: %v", err)
		}
	} else {
		if app.IsDebug() {
			color.Green("Redis connected successfully")
		}

		// Subscribe to realtime channel (background goroutine)
		go func() {
			pubsub := app.redisCache.Subscribe(app.redisContext, "realtime")
			defer pubsub.Close()

			ch := pubsub.Channel()
			for msg := range ch {
				event := &RealtimeBroadcastEvent{
					App:     app,
					Channel: msg.Channel,
					Payload: []byte(msg.Payload),
				}
				if err := app.OnRealtimeBroadcast().Trigger(event); err != nil && app.IsDebug() {
					log.Println("Realtime broadcast error:", err)
				}
			}
		}()
	}

	return nil
}
```

#### Step 4.2: Add `Publish()` function
```go
// Add to core/base.go (after initRedis function, ~line 1126)
func (app *BaseApp) Publish(channel string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if app.redisCache != nil {
		// Use Redis Pub/Sub if available
		return app.redisCache.Publish(app.redisContext, channel, payload).Err()
	}

	// Fallback to local broadcast if Redis is not enabled
	event := &RealtimeBroadcastEvent{
		App:     app,
		Channel: channel,
		Payload: payload,
	}
	return app.OnRealtimeBroadcast().Trigger(event)
}
```

#### Step 4.3: Update `Bootstrap()` function
```go
// Modify Bootstrap() in core/base.go (line ~335)
func (app *BaseApp) Bootstrap() error {
	// ... existing code ...

	if err := app.initDataDB(); err != nil {
		return err
	}

	if err := app.initRedis(); err != nil {  // ← ADD THIS LINE
		return err
	}

	// ... rest unchanged ...
}
```

#### Step 4.4: Update `initDataDB()` to use `connectDB()`
```go
// Modify initDataDB() in core/base.go (line ~1018)
func (app *BaseApp) initDataDB() error {
	maxOpenConns := DefaultDataMaxOpenConns
	maxIdleConns := DefaultDataMaxIdleConns
	// ... existing config code ...

	// REPLACE: concurrentDB, err := openDB(...)
	// WITH:
	concurrentDB, err := connectDB(app.dataDsn)  // ← USE NEW FUNCTION
	if err != nil {
		return err
	}
	// ... rest unchanged ...

	// REPLACE: nonconcurrentDB, err := openDB(...)
	// WITH:
	nonconcurrentDB, err := connectDB(app.dataDsn)  // ← USE NEW FUNCTION
	if err != nil {
		return err
	}
	// ... rest unchanged ...
}
```

**Deliverables**:
- [ ] `initRedis()` function added (~50 lines)
- [ ] `Publish()` function added (~20 lines)
- [ ] `Bootstrap()` updated (1 line)
- [ ] `initDataDB()` updated (use `connectDB()`)
- [ ] Code compiles successfully

### Day 5: Command-Line Flags

**Tasks**:
1. Add `--dataDsn` flag to `cmd/serve.go`
2. Add `--redisDsn` flag to `cmd/serve.go`
3. Pass flags to `BaseAppConfig`
4. Test with dummy DSN strings

**Implementation Steps**:

#### Step 5.1: Update `cmd/serve.go`
```go
// Find the flags section in cmd/serve.go
// Add these flags (after existing --dir flag):

app.RootCmd.PersistentFlags().StringVar(
	&dataDsn,
	"dataDsn",
	"",
	"Database connection string (postgresql://, mysql://, or sqlite file path)",
)

app.RootCmd.PersistentFlags().StringVar(
	&redisDsn,
	"redisDsn",
	"",
	"Redis connection string (optional, enables distributed caching)",
)

// Update NewBaseApp call to include new flags:
app := core.NewBaseApp(core.BaseAppConfig{
	DataDir:          dataDir,
	DataDsn:          dataDsn,   // ← ADD THIS
	RedisDsn:         redisDsn,  // ← ADD THIS
	EncryptionEnv:    encryptionEnv,
	IsDebug:          isDebug,
	DataMaxOpenConns: dataMaxOpenConns,
	DataMaxIdleConns: dataMaxIdleConns,
	// ... rest unchanged ...
})
```

#### Step 5.2: Test compilation and help text
```bash
cd rai-backend
go build -o pb ./

# Test help text
./pb serve --help | grep -E "dataDsn|redisDsn"
# Expected output:
#   --dataDsn string    Database connection string (postgresql://, ...)
#   --redisDsn string   Redis connection string (optional, ...)
```

**Deliverables**:
- [ ] `--dataDsn` flag added
- [ ] `--redisDsn` flag added
- [ ] Flags passed to `BaseAppConfig`
- [ ] Help text displays correctly
- [ ] Binary compiles successfully

---

## Phase 3: Query Builder Integration (Ngày 6-7)

### Day 6: dbx Package Vendoring

**Tasks**:
1. Vendor postgrebase's dbx package
2. Update import paths in all files
3. Test basic queries
4. Document SQL dialect differences

**Implementation Steps**:

#### Step 6.1: Vendor dbx package
```bash
# Option A: Copy from postgrebase
cd /tmp/postgrebase
cp -r vendor/github.com/free/postgresqlbaseapi/dbx ~/Documents/postgrebase-super-backend/rai-backend/vendor/github.com/free/postgresqlbaseapi/

# Option B: Use go mod vendor
cd rai-backend
go mod edit -replace github.com/free/postgresqlbaseapi/dbx=/path/to/dbx
go mod vendor
```

#### Step 6.2: Update import paths
```bash
# Find all dbx imports
cd rai-backend
grep -r "import.*dbx" --include="*.go" | cut -d: -f1 | sort | uniq

# Update imports (if needed)
# This depends on whether we keep postgrebase's package path or change it
```

#### Step 6.3: Test basic queries
```go
// Create test file: rai-backend/core/db_test.go
package core

import (
	"testing"
)

func TestConnectDB_PostgreSQL(t *testing.T) {
	// Mock test - will replace with real test later
	dsn := "postgresql://localhost:5432/test"
	db, err := connectDB(dsn)
	if err != nil {
		t.Skip("PostgreSQL not available:", err)
	}
	defer db.Close()
	
	// Test ping
	if err := db.DB().Ping(); err != nil {
		t.Fatal("Ping failed:", err)
	}
}
```

**Deliverables**:
- [ ] dbx package vendored
- [ ] Import paths updated (if needed)
- [ ] Basic connection test passes
- [ ] SQL dialect documentation created

### Day 7: Migration Conversion

**Tasks**:
1. Analyze existing migrations in `pb_migrations/`
2. Convert SQLite syntax → PostgreSQL syntax
3. Create migration converter script
4. Test migrations on PostgreSQL

**Migration Conversion Matrix**:

| SQLite Syntax | PostgreSQL Syntax |
|---------------|-------------------|
| `INTEGER PRIMARY KEY AUTOINCREMENT` | `SERIAL PRIMARY KEY` |
| `TEXT` | `TEXT` (same) |
| `datetime('now')` | `NOW()` or `CURRENT_TIMESTAMP` |
| `datetime(created, '+1 day')` | `created + INTERVAL '1 day'` |
| `json_extract(col, '$.key')` | `col->>'key'` |
| `||` (concat) | `||` (same) or `CONCAT()` |

**Implementation Steps**:

#### Step 7.1: Create migration converter script
```bash
# Create: scripts/convert-migrations.sh
#!/bin/bash
# Convert SQLite migrations to PostgreSQL syntax

for file in pb_migrations/*.sql; do
    echo "Converting $file..."
    
    # Replace datetime('now') with NOW()
    sed -i "s/datetime('now')/NOW()/g" "$file"
    
    # Replace INTEGER PRIMARY KEY AUTOINCREMENT
    sed -i "s/INTEGER PRIMARY KEY AUTOINCREMENT/SERIAL PRIMARY KEY/g" "$file"
    
    # Replace TEXT DEFAULT (datetime('now'))
    sed -i "s/TEXT DEFAULT (datetime('now'))/TIMESTAMP DEFAULT NOW()/g" "$file"
    
    echo "Converted: $file"
done
```

#### Step 7.2: Backup and convert migrations
```bash
# Backup original migrations
cp -r pb_migrations pb_migrations.sqlite.backup

# Run converter
chmod +x scripts/convert-migrations.sh
./scripts/convert-migrations.sh

# Review changes
git diff pb_migrations/
```

#### Step 7.3: Test migrations
```bash
# Start PostgreSQL test database
docker run -d --name pg-test \
  -e POSTGRES_PASSWORD=test \
  -e POSTGRES_DB=rai_test \
  -p 5433:5432 \
  postgres:16

# Run migrations
./pb migrate up --dataDsn "postgresql://postgres:test@localhost:5433/rai_test?sslmode=disable"

# Verify tables created
docker exec pg-test psql -U postgres -d rai_test -c "\dt"
```

**Deliverables**:
- [ ] Migration converter script created
- [ ] All migrations converted to PostgreSQL syntax
- [ ] Migrations tested on PostgreSQL
- [ ] Rollback tested

---

## Phase 4: Testing & Validation (Ngày 8-10)

### Day 8: Unit Testing

**Tasks**:
1. Create unit tests for database connection
2. Create unit tests for Redis integration
3. Test query generation
4. Test transaction handling

**Test Files to Create**:

#### File: `rai-backend/core/db_test.go`
```go
package core

import (
	"os"
	"testing"
)

func TestConnectDB_PostgreSQL(t *testing.T) {
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN not set")
	}

	db, err := connectDB(dsn)
	if err != nil {
		t.Fatal("Connection failed:", err)
	}
	defer db.Close()

	if err := db.DB().Ping(); err != nil {
		t.Fatal("Ping failed:", err)
	}
}

func TestConnectDB_MySQL(t *testing.T) {
	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("TEST_MYSQL_DSN not set")
	}

	db, err := connectDB(dsn)
	if err != nil {
		t.Fatal("Connection failed:", err)
	}
	defer db.Close()

	if err := db.DB().Ping(); err != nil {
		t.Fatal("Ping failed:", err)
	}
}
```

#### File: `rai-backend/core/redis_test.go`
```go
package core

import (
	"os"
	"testing"
)

func TestInitRedis_Success(t *testing.T) {
	redisDsn := os.Getenv("TEST_REDIS_DSN")
	if redisDsn == "" {
		t.Skip("TEST_REDIS_DSN not set")
	}

	app := NewBaseApp(BaseAppConfig{
		RedisDsn: redisDsn,
	})

	if err := app.initRedis(); err != nil {
		t.Fatal("Redis init failed:", err)
	}

	// Test publish
	if err := app.Publish("test", map[string]string{"msg": "hello"}); err != nil {
		t.Fatal("Publish failed:", err)
	}
}

func TestInitRedis_Fallback(t *testing.T) {
	// Test graceful degradation when Redis unavailable
	app := NewBaseApp(BaseAppConfig{
		RedisDsn: "redis://invalid:6379/0",
	})

	// Should not fail, just disable Redis
	if err := app.initRedis(); err != nil {
		t.Fatal("Init should not fail:", err)
	}

	// Publish should fallback to local
	if err := app.Publish("test", map[string]string{"msg": "hello"}); err != nil {
		t.Fatal("Local publish failed:", err)
	}
}
```

**Run Tests**:
```bash
# Set test environment variables
export TEST_POSTGRES_DSN="postgresql://postgres:test@localhost:5432/rai_test?sslmode=disable"
export TEST_REDIS_DSN="redis://localhost:6379/0"

# Run tests
cd rai-backend
go test ./core -v
```

**Deliverables**:
- [ ] Database connection tests pass
- [ ] Redis integration tests pass
- [ ] Query generation tests pass
- [ ] Test coverage >= 80%

### Day 9: Integration Testing

**Tasks**:
1. Start full Docker stack (PostgreSQL + Redis + MinIO)
2. Test multi-tenant isolation
3. Test realtime with Redis Pub/Sub
4. Test clustering (2 instances)

**Test Scenarios**:

#### Scenario 1: Multi-Tenant Isolation
```bash
# Start services
docker-compose -f docker-compose.dev.yml up -d

# Create collections for app1
curl -X POST http://localhost:8090/api/collections \
  -H "X-App-ID: app1" \
  -d '{"name":"users_app1","type":"auth"}'

# Create collections for app2
curl -X POST http://localhost:8090/api/collections \
  -H "X-App-ID: app2" \
  -d '{"name":"users_app2","type":"auth"}'

# Create user in app1
curl -X POST http://localhost:8090/api/collections/users_app1/records \
  -d '{"email":"user1@app1.com","password":"Test123!"}'

# Verify app2 cannot access app1's users
curl -X GET http://localhost:8090/api/collections/users_app1/records \
  -H "X-App-ID: app2"
# Expected: 403 Forbidden
```

#### Scenario 2: Redis Pub/Sub Testing
```bash
# Start 2 instances
./pb serve --dataDsn "$POSTGRES_DSN" --redisDsn "$REDIS_DSN" --http localhost:8090 &
./pb serve --dataDsn "$POSTGRES_DSN" --redisDsn "$REDIS_DSN" --http localhost:8091 &

# Connect SSE to instance 1
curl -N http://localhost:8090/api/realtime &

# Create record via instance 2
curl -X POST http://localhost:8091/api/collections/posts/records \
  -d '{"title":"Test Post"}'

# Verify SSE event received on instance 1
# Expected: Realtime event with "create" action
```

**Deliverables**:
- [ ] Multi-tenant isolation verified
- [ ] Redis Pub/Sub working across instances
- [ ] Clustering test passed
- [ ] No race conditions detected

### Day 10: Performance Testing & Documentation

**Tasks**:
1. Run load tests (concurrent writes)
2. Compare performance vs SQLite
3. Document benchmark results
4. Finalize migration guide

**Load Test Script**:
```bash
# Create: scripts/load-test.sh
#!/bin/bash
# Load test with 100 concurrent clients

echo "Running load test..."

ab -n 10000 -c 100 -p post-data.json \
  -T "application/json" \
  http://localhost:8090/api/collections/posts/records

# Expected:
# - PostgreSQL: >= 1000 req/s
# - SQLite: ~100 req/s (write-lock bottleneck)
```

**Benchmark Comparison**:

| Metric | SQLite | PostgreSQL | Improvement |
|--------|--------|------------|-------------|
| **Concurrent Writes** | ~100 req/s | ~1000 req/s | **10x** |
| **Read Latency (p50)** | 5ms | 8ms | 1.6x slower |
| **Read Latency (p99)** | 50ms | 30ms | **1.7x faster** |
| **Connection Pool** | N/A (file lock) | 120 connections | ✅ |
| **Clustering** | ❌ Not supported | ✅ Multi-node | ✅ |

**Deliverables**:
- [ ] Load test results documented
- [ ] Performance comparison vs SQLite
- [ ] Migration guide finalized
- [ ] Production deployment checklist

---

## Rollback Strategy

### If Migration Fails (Ngày bất kỳ)

**Immediate Rollback**:
```bash
cd ~/Documents/postgrebase-super-backend
git reset --hard v0.1.0-pre-migration
docker-compose -f docker-compose.dev.yml down
docker-compose -f docker-compose.dev.yml up -d
```

**Partial Rollback** (keep some changes):
```bash
# Rollback specific files
git checkout v0.1.0-pre-migration -- rai-backend/core/base.go
git checkout v0.1.0-pre-migration -- rai-backend/cmd/serve.go
```

**Data Backup** (before testing):
```bash
# Backup PostgreSQL database
docker exec postgrebase-postgres pg_dump -U postgrebase rai_backend_dev > backup.sql

# Restore if needed
docker exec -i postgrebase-postgres psql -U postgrebase rai_backend_dev < backup.sql
```

---

## Success Criteria

### Phase 1 Complete
- [ ] PocketBase v0.37.5 source analyzed
- [ ] Postgrebase code extracted and documented
- [ ] Dependency plan created

### Phase 2 Complete
- [ ] Database connection layer implemented
- [ ] Redis integration working
- [ ] Command-line flags added
- [ ] Code compiles successfully

### Phase 3 Complete
- [ ] dbx package vendored
- [ ] Migrations converted to PostgreSQL
- [ ] Migrations tested on PostgreSQL

### Phase 4 Complete
- [ ] Unit tests passing (>= 80% coverage)
- [ ] Integration tests passing
- [ ] Performance benchmarks meet targets
- [ ] Documentation complete

### Final Acceptance
- [ ] PostgreSQL + MySQL support working
- [ ] Redis clustering working
- [ ] Multi-tenant isolation verified
- [ ] >= 10x write performance vs SQLite
- [ ] No API breaking changes
- [ ] Production deployment ready

---

## Risk Mitigation

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| **dbx compatibility issues** | 🟡 Medium | 🔴 High | Use postgrebase's fork, test thoroughly |
| **Migration syntax errors** | 🟡 Medium | 🟠 Medium | Automated converter + manual review |
| **Redis connection failures** | 🟢 Low | 🟡 Low | Graceful fallback to in-memory |
| **Performance regression** | 🟢 Low | 🟠 Medium | Benchmark before/after, tune connection pool |
| **Data loss during migration** | 🟢 Low | 🔴 High | Backup before testing, test on dummy data |

---

## Timeline Summary

| Phase | Days | Focus | Key Deliverables |
|-------|------|-------|------------------|
| **Phase 1** | 2 | Preparation & Analysis | Code extraction, dependency analysis |
| **Phase 2** | 3 | Core Implementation | Database layer, Redis, CLI flags |
| **Phase 3** | 2 | Query Builder | dbx vendoring, migrations |
| **Phase 4** | 3 | Testing & Validation | Unit tests, integration tests, benchmarks |
| **TOTAL** | **10 days** | **2 weeks** | **Production-ready PostgreSQL backend** |

---

## Next Action

**Start Date**: 2026-05-03 (Today!)  
**First Task**: Phase 1, Day 1 Morning - Repository Setup

**Command to begin**:
```bash
cd ~/Documents/postgrebase-super-backend
git checkout -b feature/postgres-migration
git add docs/POCKETBASE_V0375_TO_POSTGRES_PLAN.md
git commit -m "Add PostgreSQL migration plan"
```

**Status**: 📋 **Ready to Execute**  
**Last Updated**: 2026-05-03
