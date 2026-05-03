# PostgreSQL Migration Strategy - PocketBase v0.37.5

**Document Version**: 1.0  
**Created**: 2026-05-03  
**Target**: Migrate PocketBase v0.37.5 to PostgreSQL + Redis  
**Reference**: [zhenruyan/postgrebase](https://github.com/zhenruyan/postgrebase) methodology

---

## Executive Summary

**Goal**: Học cách postgrebase đã chuyển đổi PocketBase sang PostgreSQL + Redis, sau đó áp dụng cho PocketBase v0.37.5 (latest stable).

**Approach**: Phân tích postgrebase source code → Xác định các thay đổi quan trọng → Áp dụng lên v0.37.5

**Timeline**: 2-3 tuần (phân tích + implement + test)

---

## 1. Postgrebase Architecture Analysis

### 1.1 Core Changes (từ README)

**Database Layer**:
- ✅ Thay thế SQLite → PostgreSQL/MySQL
- ✅ DSN-based configuration: `--dataDsn` flag
- ✅ Dual engine support: PostgreSQL (default) + MySQL

**Caching Layer**:
- ✅ Redis cache với `--redisDsn` flag
- ✅ Fallback to in-memory cache nếu không có Redis
- ✅ SSE Realtime sử dụng Redis Pub/Sub cho clustering

**Clustering**:
- ✅ Multi-node deployment support
- ✅ Load balancer compatible
- ✅ No local database file dependency

**API Compatibility**:
- ✅ 100% tương thích PocketBase APIs
- ✅ Admin UI không thay đổi
- ✅ Business logic giữ nguyên

### 1.2 Key Architectural Decisions

| Aspect | SQLite (Original) | PostgreSQL/MySQL (Postgrebase) |
|--------|-------------------|--------------------------------|
| **Concurrency** | Write-lock bottleneck | True multi-connection concurrency |
| **Clustering** | Single-node only | Multi-node with Redis sync |
| **Storage** | Local file | Network database |
| **Realtime** | Local memory | Redis Pub/Sub |
| **Cache** | In-memory only | Redis + in-memory fallback |

---

## 2. Code Analysis Plan

### Phase 1: Repository Structure Comparison (Week 1, Day 1-2)

**Tasks**:
1. Clone postgrebase repository
2. Clone PocketBase v0.37.5
3. Compare directory structures
4. Identify modified files

**Commands**:
```bash
# Clone repositories
cd /tmp
git clone https://github.com/zhenruyan/postgrebase.git
git clone --branch v0.37.5 https://github.com/pocketbase/pocketbase.git pocketbase-v0.37.5

# Compare structures
diff -qr postgrebase/ pocketbase-v0.37.5/ | tee /tmp/structure-diff.txt

# Find Go files with database logic
cd postgrebase
find . -name "*.go" | xargs grep -l "sqlite\|database\|sql.DB" | tee /tmp/db-files.txt
```

**Expected Findings**:
- Modified files in: `core/`, `daos/`, `models/`, `tools/`
- New files: Redis integration, PostgreSQL drivers
- Configuration changes: DSN flags, cache settings

### Phase 2: Database Layer Analysis (Week 1, Day 3-5)

**Focus Areas**:
1. **Database Abstraction Layer** (`core/db.go`):
   - Connection pool management
   - SQL dialect handling (PostgreSQL vs SQLite syntax)
   - Transaction management
   
2. **Schema Management** (`daos/schema.go`):
   - Table creation (SQLite → PostgreSQL DDL)
   - Index creation (syntax differences)
   - Migration system changes

3. **Query Builder** (`daos/query.go`):
   - SQL query generation
   - Placeholder syntax (`?` vs `$1`)
   - Date/time functions (SQLite vs PostgreSQL)

**Analysis Method**:
```bash
# Extract database-related code
cd /tmp/postgrebase

# Find DSN handling
rg "dataDsn|database.*dsn" -A 5 -B 5

# Find PostgreSQL-specific code
rg "postgres|pq\.|lib/pq" -A 3

# Find MySQL-specific code
rg "mysql|go-sql-driver" -A 3

# Find Redis integration
rg "redis|redisDsn" -A 5
```

**Expected Changes**:
- SQL placeholders: `?` → `$1, $2, $3` (PostgreSQL)
- Date functions: `datetime()` → `NOW()`, `CURRENT_TIMESTAMP`
- Auto-increment: `AUTOINCREMENT` → `SERIAL`, `BIGSERIAL`
- JSON handling: SQLite JSON → PostgreSQL JSONB
- Full-text search: SQLite FTS5 → PostgreSQL `tsvector`

### Phase 3: Cache Layer Analysis (Week 2, Day 1-2)

**Focus Areas**:
1. **Redis Integration**:
   - Connection management
   - Cache key naming
   - TTL configuration
   - Pub/Sub for realtime

2. **Fallback Mechanism**:
   - In-memory cache when Redis unavailable
   - Cache interface abstraction

**Analysis Method**:
```bash
# Find cache-related code
cd /tmp/postgrebase
rg "cache|Cache" --type go -A 3

# Find Redis Pub/Sub
rg "Subscribe|Publish" --type go -A 5

# Find realtime/SSE integration
rg "realtime|SSE" --type go -A 5
```

### Phase 4: API Layer Analysis (Week 2, Day 3-4)

**Focus Areas**:
1. **Connection Pool Configuration**:
   - `SetMaxOpenConns()`, `SetMaxIdleConns()`
   - Connection lifetime management

2. **Transaction Handling**:
   - Multi-statement transactions
   - Savepoints (PostgreSQL)
   - Rollback behavior

3. **Performance Optimizations**:
   - Prepared statements
   - Batch inserts
   - Query caching

### Phase 5: Migration System Analysis (Week 2, Day 5)

**Focus Areas**:
1. **Schema Migrations**:
   - Migration file format changes
   - SQL dialect conversion
   - Data migration scripts

2. **Backward Compatibility**:
   - Existing PocketBase migrations
   - Data import from SQLite

---

## 3. Key Files to Analyze

### 3.1 Core Files (Expected Changes)

| File | Change Type | Priority |
|------|-------------|----------|
| `core/db.go` | 🔴 Major | P0 - Database connection |
| `core/app.go` | 🟡 Medium | P1 - App initialization |
| `daos/base.go` | 🔴 Major | P0 - DAO foundation |
| `daos/query.go` | 🔴 Major | P0 - Query builder |
| `models/schema.go` | 🟡 Medium | P1 - Schema definitions |
| `tools/migrate.go` | 🟡 Medium | P1 - Migration runner |

### 3.2 New Files (Expected Additions)

| File | Purpose | Priority |
|------|---------|----------|
| `core/postgres.go` | PostgreSQL driver | P0 |
| `core/mysql.go` | MySQL driver | P1 |
| `core/redis.go` | Redis cache | P0 |
| `core/cache.go` | Cache abstraction | P0 |
| `daos/postgres_query.go` | PostgreSQL query builder | P0 |

### 3.3 Modified Command-Line Flags

```go
// Expected new flags in cmd/serve.go
app.RootCmd.PersistentFlags().StringVar(
    &app.DataDsn(),
    "dataDsn",
    "",
    "Database DSN (postgresql:// or mysql://)",
)

app.RootCmd.PersistentFlags().StringVar(
    &app.RedisDsn(),
    "redisDsn",
    "",
    "Redis DSN (optional, falls back to in-memory cache)",
)
```

---

## 4. SQL Dialect Conversion Matrix

### 4.1 DDL (Data Definition Language)

| Feature | SQLite | PostgreSQL | MySQL |
|---------|--------|------------|-------|
| **Auto-increment** | `INTEGER PRIMARY KEY AUTOINCREMENT` | `SERIAL PRIMARY KEY` or `BIGSERIAL` | `INT AUTO_INCREMENT PRIMARY KEY` |
| **Boolean** | `INTEGER` (0/1) | `BOOLEAN` | `BOOLEAN` or `TINYINT(1)` |
| **JSON** | `TEXT` with JSON functions | `JSONB` | `JSON` |
| **Datetime** | `TEXT` (ISO8601) | `TIMESTAMP` or `TIMESTAMPTZ` | `DATETIME` or `TIMESTAMP` |
| **Text** | `TEXT` | `TEXT` or `VARCHAR` | `TEXT` or `VARCHAR` |

### 4.2 DML (Data Manipulation Language)

| Feature | SQLite | PostgreSQL | MySQL |
|---------|--------|------------|-------|
| **Placeholders** | `?` | `$1, $2, $3` | `?` |
| **Current time** | `datetime('now')` | `NOW()` or `CURRENT_TIMESTAMP` | `NOW()` |
| **String concat** | `||` | `||` or `CONCAT()` | `CONCAT()` |
| **Limit offset** | `LIMIT ? OFFSET ?` | `LIMIT $1 OFFSET $2` | `LIMIT ? OFFSET ?` |
| **UPSERT** | `INSERT OR REPLACE` | `INSERT ... ON CONFLICT` | `INSERT ... ON DUPLICATE KEY UPDATE` |

### 4.3 Full-Text Search

| Feature | SQLite | PostgreSQL | MySQL |
|---------|--------|------------|-------|
| **FTS** | `CREATE VIRTUAL TABLE ... USING fts5` | `CREATE INDEX ... USING gin(to_tsvector(...))` | `FULLTEXT INDEX` |
| **Search** | `MATCH` | `@@` operator with `to_tsquery()` | `MATCH ... AGAINST` |
| **Ranking** | `rank` | `ts_rank()` | `relevance` |

### 4.4 JSON Operations

| Feature | SQLite | PostgreSQL |
|---------|--------|------------|
| **Extract** | `json_extract(col, '$.path')` | `col->>'path'` or `col#>>'{path}'` |
| **Set** | `json_set(col, '$.path', value)` | `jsonb_set(col, '{path}', value)` |
| **Array** | `json_each()` | `jsonb_array_elements()` |

---

## 5. Implementation Roadmap

### Week 1: Analysis & Preparation

**Day 1-2**: Repository comparison
- Clone postgrebase + PocketBase v0.37.5
- Generate file-by-file diff
- Identify all modified Go files
- Document architectural changes

**Day 3-4**: Database layer deep dive
- Analyze `core/db.go` changes
- Map SQL dialect conversions
- Study connection pool configuration
- Document query builder changes

**Day 5**: Cache layer analysis
- Analyze Redis integration
- Study Pub/Sub for realtime
- Document fallback mechanism

### Week 2: Core Implementation

**Day 1-2**: Database abstraction layer
- Implement PostgreSQL driver
- Implement MySQL driver (optional)
- Create database factory pattern
- Add DSN parsing logic

**Day 3-4**: Query builder conversion
- Convert SQLite queries → PostgreSQL
- Handle placeholder differences (`?` → `$1`)
- Implement dialect-specific functions
- Add transaction support

**Day 5**: Redis cache integration
- Implement Redis client
- Add cache abstraction layer
- Create fallback mechanism
- Integrate with realtime/SSE

### Week 3: Testing & Validation

**Day 1-2**: Unit tests
- Test database connections
- Test query generation
- Test cache operations
- Test transactions

**Day 3-4**: Integration tests
- Test multi-tenant isolation
- Test realtime with Redis
- Test cache fallback
- Test clustering (multi-node)

**Day 5**: Performance testing
- Load testing (concurrent writes)
- Cache hit rate analysis
- Query performance comparison
- Realtime latency testing

---

## 6. Risk Assessment & Mitigation

### 6.1 High-Risk Areas

| Risk | Severity | Mitigation |
|------|----------|------------|
| **SQL syntax incompatibilities** | 🔴 High | Comprehensive test suite with real PostgreSQL |
| **Data migration failures** | 🔴 High | Backup strategy + rollback plan |
| **Performance regression** | 🟡 Medium | Benchmark before/after |
| **Cache synchronization issues** | 🟡 Medium | Test multi-node scenarios |
| **Transaction deadlocks** | 🟡 Medium | Proper isolation levels |

### 6.2 Testing Strategy

**Unit Tests**:
- Test each SQL dialect conversion
- Test cache operations (Redis + fallback)
- Test connection pool behavior

**Integration Tests**:
- Test multi-tenant data isolation
- Test realtime with Redis Pub/Sub
- Test clustering scenarios

**Performance Tests**:
- Concurrent write performance (vs SQLite)
- Cache hit rate analysis
- Query latency comparison

**Regression Tests**:
- Run all existing PocketBase tests
- Verify API compatibility
- Verify Admin UI functionality

---

## 7. Deliverables

### 7.1 Code Deliverables

1. **Modified Core** (`rai-backend/core/`):
   - `db_postgres.go` - PostgreSQL driver
   - `db_mysql.go` - MySQL driver (optional)
   - `redis.go` - Redis cache
   - `cache.go` - Cache abstraction

2. **Modified DAOs** (`rai-backend/daos/`):
   - `query_postgres.go` - PostgreSQL query builder
   - `migrate_postgres.go` - PostgreSQL migrations

3. **Configuration**:
   - Updated `cmd/serve.go` with DSN flags
   - Docker Compose with PostgreSQL + Redis

### 7.2 Documentation Deliverables

1. **Technical Documentation**:
   - SQL Dialect Conversion Guide
   - Cache Architecture Guide
   - Migration Guide (SQLite → PostgreSQL)

2. **Operations Documentation**:
   - PostgreSQL Setup Guide
   - Redis Configuration Guide
   - Clustering Guide

3. **Developer Documentation**:
   - Contributing Guide for PostgreSQL changes
   - Testing Guide

---

## 8. Next Steps (Immediate Actions)

### Step 1: Clone Repositories (Today)

```bash
cd /tmp

# Clone postgrebase (reference implementation)
git clone https://github.com/zhenruyan/postgrebase.git
cd postgrebase
git log --oneline --graph --all | head -20

# Clone PocketBase v0.37.5 (target version)
cd /tmp
git clone --branch v0.37.5 https://github.com/pocketbase/pocketbase.git pocketbase-v0.37.5
cd pocketbase-v0.37.5
ls -la
```

### Step 2: Generate File Diff (Today)

```bash
# Compare directory structures
cd /tmp
diff -qr postgrebase/ pocketbase-v0.37.5/ | grep -E "\.go$" > ~/Documents/postgrebase-super-backend/docs/files-diff.txt

# Count modified files
wc -l ~/Documents/postgrebase-super-backend/docs/files-diff.txt
```

### Step 3: Analyze Core Changes (Tomorrow)

```bash
# Find database-related changes
cd /tmp/postgrebase
rg "postgres|pq\.|database" --type go -l | sort | uniq > /tmp/postgres-files.txt

# Find Redis changes
rg "redis|cache" --type go -l | sort | uniq > /tmp/redis-files.txt

# Compare with PocketBase v0.37.5
cd /tmp/pocketbase-v0.37.5
rg "sqlite|database" --type go -l | sort | uniq > /tmp/sqlite-files.txt
```

### Step 4: Create Detailed Analysis Document (Next Week)

Document này sẽ chứa:
- Line-by-line comparison của các file quan trọng
- SQL query conversion examples
- Redis integration patterns
- Performance benchmarks

---

## 9. Success Criteria

### 9.1 Functional Requirements

- ✅ **Database**: PostgreSQL + MySQL support with DSN configuration
- ✅ **Cache**: Redis cache with in-memory fallback
- ✅ **Realtime**: Redis Pub/Sub for multi-node SSE
- ✅ **Clustering**: Multi-instance deployment support
- ✅ **API**: 100% backward compatible with PocketBase v0.37.5 APIs
- ✅ **Admin UI**: No changes (works as-is)

### 9.2 Non-Functional Requirements

- ✅ **Performance**: >= 10x concurrent write throughput vs SQLite
- ✅ **Latency**: <= 50ms p99 for read queries
- ✅ **Cache Hit Rate**: >= 80% for hot data
- ✅ **Availability**: Multi-node clustering with zero downtime
- ✅ **Scalability**: Horizontal scaling via load balancer

### 9.3 Quality Requirements

- ✅ **Test Coverage**: >= 80% for modified code
- ✅ **Documentation**: Complete migration guide + API docs
- ✅ **Security**: No regressions from PocketBase v0.37.5
- ✅ **Compatibility**: Pass all PocketBase v0.37.5 tests

---

## 10. Timeline Summary

| Week | Focus | Deliverables |
|------|-------|--------------|
| **Week 1** | Analysis & Planning | Repository diff, architectural analysis, SQL conversion matrix |
| **Week 2** | Core Implementation | Database layer, query builder, Redis cache |
| **Week 3** | Testing & Documentation | Unit tests, integration tests, migration guide |

**Total Effort**: 3 weeks (15 working days)

**Team Size**: 1 developer + 1 Claude Code assistant

**Risk Level**: 🟡 Medium (manageable with proper testing)

---

## 11. Resources

### 11.1 Reference Repositories

- **Postgrebase**: https://github.com/zhenruyan/postgrebase
- **PocketBase v0.37.5**: https://github.com/pocketbase/pocketbase/tree/v0.37.5
- **Your Fork**: `/Users/accompany/Documents/postgrebase-super-backend/rai-backend`

### 11.2 Key Documentation

- **PostgreSQL Docs**: https://www.postgresql.org/docs/16/
- **Redis Docs**: https://redis.io/docs/
- **PocketBase Docs**: https://pocketbase.io/docs/
- **Go Database/SQL**: https://pkg.go.dev/database/sql

### 11.3 Related Files in This Repo

- `POCKETBASE_V0.37.5_UPGRADE_ANALYSIS.md` - Upgrade analysis
- `POSTGREBASE_UPGRADE_PLAN.md` - Main implementation plan
- `BACKEND_FEATURE_MATRIX.md` - Feature comparison

---

**Document Status**: 📋 **Ready for Implementation**  
**Next Action**: Clone repositories and generate file diff  
**Owner**: Development Team  
**Last Updated**: 2026-05-03
