# PostgreSQL Migration Project - Agent Task Tracking

**Project**: PocketBase v0.37.5 → PostgreSQL + Redis Migration  
**Timeline**: 10 ngày làm việc (2 tuần)  
**Created**: 2026-05-03  
**Status**: 🟡 In Planning

---

## 📊 Project Overview

### Objective
Migrate PocketBase v0.37.5 to PostgreSQL/MySQL + Redis using postgrebase methodology (~200 lines core changes)

### Success Criteria
- [ ] PostgreSQL connection working
- [ ] Redis integration complete
- [ ] All original APIs compatible
- [ ] Unit tests passing
- [ ] Integration tests passing
- [ ] Performance benchmarks acceptable

---

## 🎯 Agent Task Distribution Matrix

### Legend
- 🟢 **Ready to Start**: All dependencies met
- 🟡 **Blocked**: Waiting for dependencies
- 🔵 **In Progress**: Agent working
- ✅ **Complete**: Verified & merged
- 🔴 **Failed**: Needs retry

---

## Phase 1: Preparation & Analysis (Day 1-2)

### AGENT-01: Repository Setup Specialist
**Focus**: Codebase preparation & environment setup  
**Status**: 🟢 Ready to Start  
**Dependencies**: None  
**Estimated Time**: 3 hours

#### Tasks
- [ ] **TASK-01-A**: Backup current codebase
  - Create backup commit
  - Create git tag `v0.1.0-pre-migration`
  - Verify backup integrity
  - **Deliverable**: Commit hash + tag confirmation

- [ ] **TASK-01-B**: Clone reference repositories
  - Clone PocketBase v0.37.5
  - Clone postgrebase repository
  - Verify versions match
  - **Deliverable**: Repository paths documented

- [ ] **TASK-01-C**: Directory structure comparison
  - Run diff between repositories
  - Document key differences
  - Identify modified files
  - **Deliverable**: `docs/v0375-diff.txt` report

**Output Files**:
- `docs/v0375-diff.txt` - Structure comparison
- `docs/repo-setup-report.md` - Setup verification

**Verification Command**:
```bash
git tag | grep v0.1.0-pre-migration
ls /tmp/pocketbase-v0.37.5/
ls /tmp/postgrebase/
cat docs/v0375-diff.txt | wc -l
```

---

### AGENT-02: Code Extraction Specialist
**Focus**: Extract postgrebase implementation patterns  
**Status**: 🟡 Blocked (wait for AGENT-01)  
**Dependencies**: AGENT-01 (TASK-01-B)  
**Estimated Time**: 4 hours

#### Tasks
- [ ] **TASK-02-A**: Extract database connection layer
  - Copy `core/db_postgresql.go` (29 lines)
  - Annotate with comments
  - Document function signatures
  - **Deliverable**: `docs/extracted/db_postgresql.go`

- [ ] **TASK-02-B**: Extract Redis integration code
  - Extract `initRedis()` function from `core/base.go`
  - Extract `Publish()` function
  - Extract Redis client setup
  - **Deliverable**: `docs/extracted/redis_functions.go` (~150 lines)

- [ ] **TASK-02-C**: Extract command-line flags
  - Extract `--dataDsn` flag definition
  - Extract `--redisDsn` flag definition
  - Document flag behavior
  - **Deliverable**: `docs/extracted/cmd_serve_flags.go`

- [ ] **TASK-02-D**: Create code annotation document
  - Map extracted code to PocketBase v0.37.5 structure
  - Identify insertion points
  - Document dependencies
  - **Deliverable**: `docs/code-mapping.md`

**Output Files**:
- `docs/extracted/db_postgresql.go`
- `docs/extracted/redis_functions.go`
- `docs/extracted/cmd_serve_flags.go`
- `docs/code-mapping.md`

**Verification Command**:
```bash
wc -l docs/extracted/*.go
grep -c "func " docs/extracted/*.go
```

---

### AGENT-03: Dependency Analysis Specialist
**Focus**: Go module dependencies & version conflicts  
**Status**: 🟡 Blocked (wait for AGENT-01)  
**Dependencies**: AGENT-01 (TASK-01-B)  
**Estimated Time**: 3 hours

#### Tasks
- [ ] **TASK-03-A**: Compare go.mod files
  - Compare postgrebase vs PocketBase v0.37.5
  - Compare with current rai-backend
  - Identify version conflicts
  - **Deliverable**: `docs/dependency-matrix.md`

- [ ] **TASK-03-B**: List new dependencies
  - Document `lib/pq` version (PostgreSQL driver)
  - Document `go-sql-driver/mysql` version (MySQL driver)
  - Document `redis/go-redis/v9` version (Redis client)
  - **Deliverable**: `docs/new-dependencies.md`

- [ ] **TASK-03-C**: Conflict resolution plan
  - Check for breaking changes
  - Plan gradual upgrade path
  - Document rollback strategy
  - **Deliverable**: `docs/dependency-upgrade-plan.md`

**Output Files**:
- `docs/dependency-matrix.md`
- `docs/new-dependencies.md`
- `docs/dependency-upgrade-plan.md`

**Verification Command**:
```bash
grep "lib/pq\|go-sql-driver\|redis" docs/new-dependencies.md
```

---

### AGENT-04: dbx Package Specialist
**Focus**: Query builder (dbx) analysis & vendoring strategy  
**Status**: 🟡 Blocked (wait for AGENT-01)  
**Dependencies**: AGENT-01 (TASK-01-B)  
**Estimated Time**: 5 hours

#### Tasks
- [ ] **TASK-04-A**: Locate dbx package in postgrebase
  - Find vendor directory
  - Identify custom fork modifications
  - Document import paths
  - **Deliverable**: `docs/dbx-location.md`

- [ ] **TASK-04-B**: Compare dbx versions
  - Original dbx (github.com/go-ozzo/ozzo-dbx)
  - Postgrebase fork (github.com/free/postgresqlbaseapi/dbx)
  - Document differences
  - **Deliverable**: `docs/dbx-diff-report.md`

- [ ] **TASK-04-C**: Vendoring strategy
  - Decision: Full vendor vs minimal changes
  - Import path rewrite strategy
  - Build script modifications
  - **Deliverable**: `docs/dbx-vendoring-strategy.md`

- [ ] **TASK-04-D**: SQL dialect conversion matrix
  - SQLite → PostgreSQL syntax mapping
  - Placeholder conversion (`?` → `$1`)
  - Date function mapping
  - **Deliverable**: `docs/sql-dialect-matrix.md`

**Output Files**:
- `docs/dbx-location.md`
- `docs/dbx-diff-report.md`
- `docs/dbx-vendoring-strategy.md`
- `docs/sql-dialect-matrix.md`

**Verification Command**:
```bash
grep "postgresqlbaseapi/dbx" /tmp/postgrebase -r --include="*.go" | wc -l
```

---

## Phase 2: Core Implementation (Day 3-5)

### AGENT-05: Database Connection Layer Engineer
**Focus**: Implement PostgreSQL/MySQL connection abstraction  
**Status**: 🟡 Blocked (wait for Phase 1)  
**Dependencies**: AGENT-02 (all tasks), AGENT-03 (TASK-03-B)  
**Estimated Time**: 6 hours

#### Tasks
- [ ] **TASK-05-A**: Create `core/db_postgresql.go`
  - Implement `connectDB(dsn string)` function
  - Add driver detection logic (postgres/mysql)
  - Handle DSN prefix parsing
  - **Deliverable**: `rai-backend/core/db_postgresql.go`

- [ ] **TASK-05-B**: Update `core/base.go` struct
  - Add DSN fields
  - Add Redis client field
  - Update Bootstrap() method signature
  - **Deliverable**: Modified `rai-backend/core/base.go`

- [ ] **TASK-05-C**: Update go.mod
  - Add `lib/pq` dependency
  - Add `go-sql-driver/mysql` dependency
  - Run `go mod tidy`
  - **Deliverable**: Modified `go.mod` + `go.sum`

- [ ] **TASK-05-D**: Compilation test
  - Run `go build`
  - Fix import errors
  - Document compilation issues
  - **Deliverable**: `docs/compilation-report.md`

**Output Files**:
- `rai-backend/core/db_postgresql.go` (new)
- `rai-backend/core/base.go` (modified)
- `go.mod` (modified)
- `docs/compilation-report.md`

**Verification Command**:
```bash
go build -v ./rai-backend/... 2>&1 | tee docs/compilation-report.md
grep "db_postgresql.go" rai-backend/core/
```

---

### AGENT-06: Redis Integration Engineer
**Focus**: Implement Redis client & Pub/Sub for realtime  
**Status**: 🟡 Blocked (wait for AGENT-05)  
**Dependencies**: AGENT-02 (TASK-02-B), AGENT-05 (TASK-05-B)  
**Estimated Time**: 6 hours

#### Tasks
- [ ] **TASK-06-A**: Add Redis client initialization
  - Implement `initRedis()` in `core/base.go`
  - Add connection pooling
  - Add graceful fallback to in-memory
  - **Deliverable**: Redis init logic in `core/base.go`

- [ ] **TASK-06-B**: Implement Publish() method
  - Add `app.Publish(channel, message)` method
  - Integrate with existing SSE logic
  - Add error handling
  - **Deliverable**: Publish method in `core/base.go`

- [ ] **TASK-06-C**: Implement Subscribe() for SSE
  - Update realtime SSE handlers
  - Replace in-memory broadcast with Redis Pub/Sub
  - Test multi-node message delivery
  - **Deliverable**: Modified SSE handlers

- [ ] **TASK-06-D**: Update go.mod for Redis
  - Add `github.com/redis/go-redis/v9`
  - Run `go mod tidy`
  - **Deliverable**: Updated `go.mod`

**Output Files**:
- `rai-backend/core/base.go` (Redis methods added)
- `rai-backend/apis/realtime.go` (modified, if exists)
- `go.mod` (modified)
- `docs/redis-integration-report.md`

**Verification Command**:
```bash
grep "redis.NewClient\|Publish\|Subscribe" rai-backend/core/base.go
go mod graph | grep redis
```

---

### AGENT-07: CLI Flags Engineer
**Focus**: Add --dataDsn and --redisDsn command-line flags  
**Status**: 🟡 Blocked (wait for AGENT-05)  
**Dependencies**: AGENT-02 (TASK-02-C), AGENT-05 (TASK-05-B)  
**Estimated Time**: 3 hours

#### Tasks
- [ ] **TASK-07-A**: Add flags to `cmd/serve.go`
  - Add `--dataDsn` flag definition
  - Add `--redisDsn` flag definition
  - Update flag parsing logic
  - **Deliverable**: Modified `cmd/serve.go`

- [ ] **TASK-07-B**: Pass DSN to app.Bootstrap()
  - Update Bootstrap() call signature
  - Pass dataDsn parameter
  - Pass redisDsn parameter
  - **Deliverable**: Modified `cmd/serve.go`

- [ ] **TASK-07-C**: Add environment variable support
  - Support `DATA_DSN` env var
  - Support `REDIS_DSN` env var
  - Document precedence (flag > env var > default)
  - **Deliverable**: `docs/config-precedence.md`

- [ ] **TASK-07-D**: Update help text
  - Add flag descriptions
  - Add example DSN formats
  - **Deliverable**: Updated `--help` output

**Output Files**:
- `rai-backend/cmd/serve.go` (modified)
- `docs/config-precedence.md`

**Verification Command**:
```bash
go run rai-backend/main.go serve --help | grep -E "dataDsn|redisDsn"
```

---

## Phase 3: Query Builder & Migrations (Day 6-7)

### AGENT-08: dbx Package Vendor Specialist
**Focus**: Vendor custom dbx package for PostgreSQL/MySQL  
**Status**: 🟡 Blocked (wait for AGENT-04)  
**Dependencies**: AGENT-04 (all tasks), AGENT-05 (TASK-05-D)  
**Estimated Time**: 8 hours

#### Tasks
- [ ] **TASK-08-A**: Vendor dbx package
  - Copy postgrebase dbx fork to `vendor/`
  - Or: Use `go mod vendor` with replace directive
  - Update import paths in all files
  - **Deliverable**: `vendor/github.com/.../dbx/` (or go.mod replace)

- [ ] **TASK-08-B**: Rewrite import paths
  - Find all `import "dbx"` statements
  - Rewrite to `postgresqlbaseapi/dbx`
  - Verify no circular imports
  - **Deliverable**: Updated import statements across codebase

- [ ] **TASK-08-C**: Test query builder
  - Write unit tests for dbx layer
  - Test PostgreSQL placeholder conversion
  - Test MySQL placeholder conversion
  - **Deliverable**: `core/db_test.go`

- [ ] **TASK-08-D**: Verify SQL generation
  - Test SELECT queries
  - Test INSERT queries
  - Test UPDATE/DELETE queries
  - **Deliverable**: `docs/sql-generation-tests.md`

**Output Files**:
- `vendor/` directory (or go.mod replace directive)
- `core/db_test.go`
- `docs/sql-generation-tests.md`

**Verification Command**:
```bash
go list -m all | grep dbx
go test ./core/... -v -run TestDB
```

---

### AGENT-09: Migration System Engineer
**Focus**: Convert SQLite migrations to PostgreSQL  
**Status**: 🟡 Blocked (wait for AGENT-08)  
**Dependencies**: AGENT-08 (all tasks), AGENT-04 (TASK-04-D)  
**Estimated Time**: 8 hours

#### Tasks
- [ ] **TASK-09-A**: Analyze existing migrations
  - List all migration files in `migrations/`
  - Identify SQLite-specific syntax
  - Document required conversions
  - **Deliverable**: `docs/migration-audit.md`

- [ ] **TASK-09-B**: Create conversion script
  - Write script to convert SQL syntax
  - Handle AUTOINCREMENT → SERIAL
  - Handle datetime() → NOW()
  - Handle JSON → JSONB
  - **Deliverable**: `scripts/convert-migrations.sh`

- [ ] **TASK-09-C**: Convert migrations
  - Run conversion script
  - Manual review of each migration
  - Test on PostgreSQL test database
  - **Deliverable**: Converted migration files

- [ ] **TASK-09-D**: Test migration rollback
  - Test up migrations
  - Test down migrations
  - Verify data integrity
  - **Deliverable**: `docs/migration-test-report.md`

**Output Files**:
- `docs/migration-audit.md`
- `scripts/convert-migrations.sh`
- Converted `migrations/*.sql` files
- `docs/migration-test-report.md`

**Verification Command**:
```bash
./scripts/convert-migrations.sh --dry-run | tee docs/migration-audit.md
psql -U postgres -d test_db < migrations/001_init.sql
```

---

## Phase 4: Testing & Validation (Day 8-10)

### AGENT-10: Unit Test Engineer
**Focus**: Write & execute unit tests for all new components  
**Status**: 🟡 Blocked (wait for Phase 2 & 3)  
**Dependencies**: AGENT-05, AGENT-06, AGENT-08  
**Estimated Time**: 8 hours

#### Tasks
- [ ] **TASK-10-A**: Database connection tests
  - Test `connectDB()` with PostgreSQL DSN
  - Test `connectDB()` with MySQL DSN
  - Test invalid DSN handling
  - **Deliverable**: `core/db_postgresql_test.go`

- [ ] **TASK-10-B**: Redis integration tests
  - Test `initRedis()` with valid DSN
  - Test fallback to in-memory when Redis unavailable
  - Test Publish/Subscribe mechanics
  - **Deliverable**: `core/redis_test.go`

- [ ] **TASK-10-C**: Query builder tests
  - Test SQL placeholder conversion
  - Test date function conversion
  - Test transaction handling
  - **Deliverable**: `core/query_test.go`

- [ ] **TASK-10-D**: Run all unit tests
  - `go test ./...` execution
  - Collect coverage report
  - Fix failing tests
  - **Deliverable**: `docs/unit-test-report.md`

**Output Files**:
- `core/db_postgresql_test.go`
- `core/redis_test.go`
- `core/query_test.go`
- `docs/unit-test-report.md`

**Verification Command**:
```bash
go test ./core/... -v -cover | tee docs/unit-test-report.md
go test ./core/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o docs/coverage.html
```

---

### AGENT-11: Integration Test Engineer
**Focus**: End-to-end integration testing with PostgreSQL + Redis  
**Status**: 🟡 Blocked (wait for AGENT-10)  
**Dependencies**: AGENT-10 (all tasks)  
**Estimated Time**: 8 hours

#### Tasks
- [ ] **TASK-11-A**: Setup test infrastructure
  - Docker Compose for PostgreSQL + Redis
  - Test database initialization script
  - Test data seeding
  - **Deliverable**: `docker-compose.test.yml`

- [ ] **TASK-11-B**: API compatibility tests
  - Test all CRUD endpoints
  - Test authentication endpoints
  - Test realtime SSE endpoints
  - **Deliverable**: `tests/integration/api_test.go`

- [ ] **TASK-11-C**: Multi-node clustering test
  - Start 2 instances with shared PostgreSQL
  - Test Redis Pub/Sub message broadcast
  - Verify SSE events across nodes
  - **Deliverable**: `tests/integration/cluster_test.go`

- [ ] **TASK-11-D**: Data migration test
  - Import SQLite data to PostgreSQL
  - Verify data integrity
  - Run API tests on migrated data
  - **Deliverable**: `docs/data-migration-report.md`

**Output Files**:
- `docker-compose.test.yml`
- `tests/integration/api_test.go`
- `tests/integration/cluster_test.go`
- `docs/data-migration-report.md`

**Verification Command**:
```bash
docker-compose -f docker-compose.test.yml up -d
go test ./tests/integration/... -v | tee docs/integration-test-report.md
docker-compose -f docker-compose.test.yml down
```

---

### AGENT-12: Performance Benchmark Engineer
**Focus**: Performance testing & optimization  
**Status**: 🟡 Blocked (wait for AGENT-11)  
**Dependencies**: AGENT-11 (all tasks)  
**Estimated Time**: 6 hours

#### Tasks
- [ ] **TASK-12-A**: Database query benchmarks
  - Benchmark SELECT queries (PostgreSQL vs SQLite)
  - Benchmark INSERT/UPDATE queries
  - Benchmark JOIN queries
  - **Deliverable**: `benchmarks/db_bench_test.go`

- [ ] **TASK-12-B**: Redis cache benchmarks
  - Benchmark cache hit/miss performance
  - Benchmark Pub/Sub latency
  - Compare vs in-memory performance
  - **Deliverable**: `benchmarks/redis_bench_test.go`

- [ ] **TASK-12-C**: Run benchmark suite
  - `go test -bench=.` execution
  - Collect performance metrics
  - Compare with baseline (SQLite)
  - **Deliverable**: `docs/performance-report.md`

- [ ] **TASK-12-D**: Performance optimization
  - Identify bottlenecks
  - Optimize slow queries
  - Tune connection pool settings
  - **Deliverable**: `docs/optimization-recommendations.md`

**Output Files**:
- `benchmarks/db_bench_test.go`
- `benchmarks/redis_bench_test.go`
- `docs/performance-report.md`
- `docs/optimization-recommendations.md`

**Verification Command**:
```bash
go test ./benchmarks/... -bench=. -benchmem | tee docs/performance-report.md
```

---

## 📦 Deliverables Summary

### Documentation
- [x] Repository setup report
- [ ] Code extraction & mapping
- [ ] Dependency analysis matrix
- [ ] dbx vendoring strategy
- [ ] SQL dialect conversion matrix
- [ ] Compilation reports
- [ ] Test reports (unit + integration)
- [ ] Performance benchmarks
- [ ] Migration guide

### Code Artifacts
- [ ] `core/db_postgresql.go` (new)
- [ ] `core/base.go` (modified - Redis support)
- [ ] `cmd/serve.go` (modified - CLI flags)
- [ ] `vendor/dbx/` or `go.mod` replace directive
- [ ] Converted migration files
- [ ] Test suites (unit + integration + benchmarks)

### Infrastructure
- [ ] `docker-compose.test.yml` (test environment)
- [ ] `scripts/convert-migrations.sh` (migration converter)
- [ ] CI/CD pipeline updates (optional)

---

## 🔗 Agent Communication Protocol

### Handoff Format
When an agent completes a task, post in this format:

```markdown
## AGENT-XX Task Completion Report

**Agent ID**: AGENT-XX  
**Task ID**: TASK-XX-Y  
**Status**: ✅ Complete / 🔴 Failed  
**Completion Time**: YYYY-MM-DD HH:MM  

### Deliverables
- [Link to output file 1]
- [Link to output file 2]

### Verification
[Paste command output showing success]

### Blockers Removed
- AGENT-YY can now start (dependency met)

### Issues Found
- [List any problems discovered]

### Next Agent
Handoff to: AGENT-YY (TASK-YY-Z)
```

### Blocker Notification Format
When an agent is blocked:

```markdown
## AGENT-XX Blocker Report

**Agent ID**: AGENT-XX  
**Blocked On**: TASK-YY-Z (AGENT-YY)  
**Blocking Since**: YYYY-MM-DD HH:MM  

### Action Required
Waiting for: [Specific deliverable needed]

### Workaround Attempted
[What was tried, if any]

### ETA Impact
Current delay: X hours
```

---

## 📈 Progress Tracking

### Overall Progress
```
Phase 1: [████░░░░░░] 0/4 agents complete (0%)
Phase 2: [░░░░░░░░░░] 0/3 agents complete (0%)
Phase 3: [░░░░░░░░░░] 0/2 agents complete (0%)
Phase 4: [░░░░░░░░░░] 0/3 agents complete (0%)

Total: [░░░░░░░░░░] 0/12 agents complete (0%)
```

### Timeline Status
```
Day 1-2:  [ ] Phase 1 (Preparation)
Day 3-5:  [ ] Phase 2 (Core Implementation)
Day 6-7:  [ ] Phase 3 (Query Builder & Migrations)
Day 8-10: [ ] Phase 4 (Testing & Validation)
```

---

## 🚨 Risk Register

| Risk ID | Description | Probability | Impact | Mitigation | Owner |
|---------|-------------|-------------|--------|------------|-------|
| R-01 | dbx package incompatibility | Medium | High | Full vendoring with tests | AGENT-08 |
| R-02 | SQL migration conversion errors | High | High | Manual review + testing | AGENT-09 |
| R-03 | Redis connection failures | Low | Medium | Graceful fallback to in-memory | AGENT-06 |
| R-04 | Performance regression vs SQLite | Medium | Medium | Benchmark early, optimize | AGENT-12 |
| R-05 | Breaking API changes | Low | High | Integration test coverage | AGENT-11 |

---

## 📞 Contact & Escalation

### Agent Assignment
| Agent ID | Assigned To | Slack/Contact | Status |
|----------|-------------|---------------|--------|
| AGENT-01 | TBD | @username | 🟢 Available |
| AGENT-02 | TBD | @username | 🟢 Available |
| AGENT-03 | TBD | @username | 🟢 Available |
| AGENT-04 | TBD | @username | 🟢 Available |
| AGENT-05 | TBD | @username | 🟢 Available |
| AGENT-06 | TBD | @username | 🟢 Available |
| AGENT-07 | TBD | @username | 🟢 Available |
| AGENT-08 | TBD | @username | 🟢 Available |
| AGENT-09 | TBD | @username | 🟢 Available |
| AGENT-10 | TBD | @username | 🟢 Available |
| AGENT-11 | TBD | @username | 🟢 Available |
| AGENT-12 | TBD | @username | 🟢 Available |

### Escalation Path
1. **Level 1**: Agent self-resolves (blocker < 2 hours)
2. **Level 2**: Team lead review (blocker 2-4 hours)
3. **Level 3**: Project manager escalation (blocker > 4 hours)

---

**Last Updated**: 2026-05-03  
**Next Review**: Daily standup @ 09:00
