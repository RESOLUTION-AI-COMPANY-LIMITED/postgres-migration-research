# Multi-Agent Workflow Guide

**Project**: PostgreSQL Migration  
**Approach**: Parallel agent-based task execution  
**Created**: 2026-05-03

---

## 🎯 Workflow Overview

Dự án này được chia thành **12 agents** làm việc song song và tuần tự theo dependency graph. Mỗi agent có một lĩnh vực chuyên môn và một tập hợp các tasks cụ thể.

---

## 🚀 Quick Start: Khởi động các Agent

### Bước 1: Chuẩn bị môi trường

```bash
# Di chuyển vào thư mục dự án
cd /Users/accompany/Documents/postgres-migration-research

# Đọc project tracking
cat PROJECT_TRACKING.md

# Kiểm tra dependencies của agent bạn sẽ khởi động
grep "AGENT-01" PROJECT_TRACKING.md -A 20
```

### Bước 2: Khởi động agent với Claude Code

Mỗi agent có một prompt template riêng. Sử dụng các prompt dưới đây để khởi động agent:

---

## 📋 Agent Prompt Templates

### AGENT-01: Repository Setup Specialist
```
Tôi là AGENT-01: Repository Setup Specialist trong dự án PostgreSQL Migration.

Nhiệm vụ của tôi:
- Backup current codebase
- Clone reference repositories (PocketBase v0.37.5 + postgrebase)
- Directory structure comparison

Hãy đọc file PROJECT_TRACKING.md, tìm section "AGENT-01: Repository Setup Specialist" và thực hiện tất cả tasks trong đó (TASK-01-A, TASK-01-B, TASK-01-C).

Sau khi hoàn thành, tạo completion report theo Agent Communication Protocol format.
```

---

### AGENT-02: Code Extraction Specialist
```
Tôi là AGENT-02: Code Extraction Specialist trong dự án PostgreSQL Migration.

Nhiệm vụ của tôi:
- Extract core/db_postgresql.go từ postgrebase
- Extract Redis integration code
- Extract CLI flags
- Create code mapping document

Dependencies: Chờ AGENT-01 hoàn thành TASK-01-B.

Hãy đọc file PROJECT_TRACKING.md, tìm section "AGENT-02: Code Extraction Specialist" và thực hiện tất cả tasks.

Sau khi hoàn thành, tạo completion report và notify AGENT-05, AGENT-06, AGENT-07 có thể bắt đầu.
```

---

### AGENT-03: Dependency Analysis Specialist
```
Tôi là AGENT-03: Dependency Analysis Specialist trong dự án PostgreSQL Migration.

Nhiệm vụ của tôi:
- Compare go.mod files (postgrebase vs PocketBase v0.37.5 vs rai-backend)
- List new dependencies (lib/pq, go-sql-driver, go-redis)
- Create conflict resolution plan

Dependencies: Chờ AGENT-01 hoàn thành TASK-01-B.

Hãy đọc file PROJECT_TRACKING.md, tìm section "AGENT-03: Dependency Analysis Specialist" và thực hiện tất cả tasks.

Sau khi hoàn thành, tạo completion report và notify AGENT-05 có thể update go.mod.
```

---

### AGENT-04: dbx Package Specialist
```
Tôi là AGENT-04: dbx Package Specialist trong dự án PostgreSQL Migration.

Nhiệm vụ của tôi:
- Locate dbx package in postgrebase
- Compare dbx versions (original vs fork)
- Create vendoring strategy
- Create SQL dialect conversion matrix

Dependencies: Chờ AGENT-01 hoàn thành TASK-01-B.

Hãy đọc file PROJECT_TRACKING.md, tìm section "AGENT-04: dbx Package Specialist" và thực hiện tất cả tasks.

Sau khi hoàn thành, tạo completion report và notify AGENT-08, AGENT-09 có thể bắt đầu.
```

---

### AGENT-05: Database Connection Layer Engineer
```
Tôi là AGENT-05: Database Connection Layer Engineer trong dự án PostgreSQL Migration.

Nhiệm vụ của tôi:
- Create core/db_postgresql.go
- Update core/base.go struct
- Update go.mod with new dependencies
- Run compilation test

Dependencies: Chờ AGENT-02 (all tasks) và AGENT-03 (TASK-03-B).

Hãy đọc file PROJECT_TRACKING.md, tìm section "AGENT-05: Database Connection Layer Engineer" và thực hiện tất cả tasks.

Sau khi hoàn thành, tạo completion report và notify AGENT-06, AGENT-07, AGENT-08 có thể bắt đầu.
```

---

### AGENT-06: Redis Integration Engineer
```
Tôi là AGENT-06: Redis Integration Engineer trong dự án PostgreSQL Migration.

Nhiệm vụ của tôi:
- Add Redis client initialization
- Implement Publish() method
- Implement Subscribe() for SSE
- Update go.mod for Redis

Dependencies: Chờ AGENT-02 (TASK-02-B) và AGENT-05 (TASK-05-B).

Hãy đọc file PROJECT_TRACKING.md, tìm section "AGENT-06: Redis Integration Engineer" và thực hiện tất cả tasks.

Sau khi hoàn thành, tạo completion report và notify AGENT-10 có thể viết Redis tests.
```

---

### AGENT-07: CLI Flags Engineer
```
Tôi là AGENT-07: CLI Flags Engineer trong dự án PostgreSQL Migration.

Nhiệm vụ của tôi:
- Add --dataDsn and --redisDsn flags to cmd/serve.go
- Pass DSN to app.Bootstrap()
- Add environment variable support
- Update help text

Dependencies: Chờ AGENT-02 (TASK-02-C) và AGENT-05 (TASK-05-B).

Hãy đọc file PROJECT_TRACKING.md, tìm section "AGENT-07: CLI Flags Engineer" và thực hiện tất cả tasks.

Sau khi hoàn thành, tạo completion report.
```

---

### AGENT-08: dbx Package Vendor Specialist
```
Tôi là AGENT-08: dbx Package Vendor Specialist trong dự án PostgreSQL Migration.

Nhiệm vụ của tôi:
- Vendor dbx package
- Rewrite import paths
- Test query builder
- Verify SQL generation

Dependencies: Chờ AGENT-04 (all tasks) và AGENT-05 (TASK-05-D).

Hãy đọc file PROJECT_TRACKING.md, tìm section "AGENT-08: dbx Package Vendor Specialist" và thực hiện tất cả tasks.

Sau khi hoàn thành, tạo completion report và notify AGENT-09, AGENT-10 có thể bắt đầu.
```

---

### AGENT-09: Migration System Engineer
```
Tôi là AGENT-09: Migration System Engineer trong dự án PostgreSQL Migration.

Nhiệm vụ của tôi:
- Analyze existing SQLite migrations
- Create SQL conversion script
- Convert migrations to PostgreSQL
- Test migration rollback

Dependencies: Chờ AGENT-08 (all tasks) và AGENT-04 (TASK-04-D).

Hãy đọc file PROJECT_TRACKING.md, tìm section "AGENT-09: Migration System Engineer" và thực hiện tất cả tasks.

Sau khi hoàn thành, tạo completion report và notify AGENT-11 có thể test migrations.
```

---

### AGENT-10: Unit Test Engineer
```
Tôi là AGENT-10: Unit Test Engineer trong dự án PostgreSQL Migration.

Nhiệm vụ của tôi:
- Write database connection tests
- Write Redis integration tests
- Write query builder tests
- Run all unit tests and collect coverage

Dependencies: Chờ AGENT-05, AGENT-06, AGENT-08 hoàn thành.

Hãy đọc file PROJECT_TRACKING.md, tìm section "AGENT-10: Unit Test Engineer" và thực hiện tất cả tasks.

Sau khi hoàn thành, tạo completion report và notify AGENT-11 có thể bắt đầu integration tests.
```

---

### AGENT-11: Integration Test Engineer
```
Tôi là AGENT-11: Integration Test Engineer trong dự án PostgreSQL Migration.

Nhiệm vụ của tôi:
- Setup test infrastructure (Docker Compose)
- Write API compatibility tests
- Write multi-node clustering test
- Test data migration from SQLite

Dependencies: Chờ AGENT-10 hoàn thành all tasks.

Hãy đọc file PROJECT_TRACKING.md, tìm section "AGENT-11: Integration Test Engineer" và thực hiện tất cả tasks.

Sau khi hoàn thành, tạo completion report và notify AGENT-12 có thể bắt đầu benchmarks.
```

---

### AGENT-12: Performance Benchmark Engineer
```
Tôi là AGENT-12: Performance Benchmark Engineer trong dự án PostgreSQL Migration.

Nhiệm vụ của tôi:
- Write database query benchmarks
- Write Redis cache benchmarks
- Run benchmark suite and compare with baseline
- Provide performance optimization recommendations

Dependencies: Chờ AGENT-11 hoàn thành all tasks.

Hãy đọc file PROJECT_TRACKING.md, tìm section "AGENT-12: Performance Benchmark Engineer" và thực hiện tất cả tasks.

Sau khi hoàn thành, tạo completion report. Đây là agent cuối cùng!
```

---

## 🔄 Dependency Graph

```
Phase 1 (Parallel):
┌─────────────┐
│  AGENT-01   │ (No dependencies)
│  Repo Setup │
└──────┬──────┘
       │
       ├───────────────┬───────────────┬──────────────┐
       ▼               ▼               ▼              ▼
┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐
│  AGENT-02  │  │  AGENT-03  │  │  AGENT-04  │  │            │
│   Code     │  │    Deps    │  │    dbx     │  │            │
│ Extraction │  │  Analysis  │  │  Package   │  │            │
└──────┬─────┘  └──────┬─────┘  └──────┬─────┘  │            │
       │               │               │         │            │
       └───────────────┴───────┬───────┘         │            │
                               ▼                 │            │

Phase 2 (Sequential then Parallel):
                        ┌─────────────┐
                        │  AGENT-05   │ (Wait: AGENT-02 + AGENT-03)
                        │     DB      │
                        │ Connection  │
                        └──────┬──────┘
                               │
                ┌──────────────┼──────────────┐
                ▼              ▼              ▼
         ┌────────────┐  ┌────────────┐  ┌────────────┐
         │  AGENT-06  │  │  AGENT-07  │  │            │
         │   Redis    │  │ CLI Flags  │  │            │
         │Integration │  │            │  │            │
         └────────────┘  └────────────┘  │            │

Phase 3 (Sequential):
                        ┌─────────────┐
                        │  AGENT-08   │ (Wait: AGENT-04 + AGENT-05)
                        │     dbx     │
                        │   Vendor    │
                        └──────┬──────┘
                               ▼
                        ┌─────────────┐
                        │  AGENT-09   │ (Wait: AGENT-08)
                        │ Migrations  │
                        └──────┬──────┘
                               │

Phase 4 (Sequential):
                               │
                               ▼
                        ┌─────────────┐
                        │  AGENT-10   │ (Wait: AGENT-05+06+08)
                        │ Unit Tests  │
                        └──────┬──────┘
                               ▼
                        ┌─────────────┐
                        │  AGENT-11   │ (Wait: AGENT-10)
                        │Integration  │
                        │   Tests     │
                        └──────┬──────┘
                               ▼
                        ┌─────────────┐
                        │  AGENT-12   │ (Wait: AGENT-11)
                        │ Benchmarks  │
                        └─────────────┘
```

---

## 📊 Progress Tracking

### Cách update progress

Sau mỗi task completion, update `PROJECT_TRACKING.md`:

1. Tìm task tương ứng (e.g., `TASK-01-A`)
2. Đổi `[ ]` thành `[x]`
3. Update status emoji (🟢 → 🔵 → ✅)
4. Update completion percentage

**Ví dụ**:
```markdown
## Phase 1: Preparation & Analysis (Day 1-2)

### AGENT-01: Repository Setup Specialist
**Status**: 🔵 In Progress  <!-- Was: 🟢 Ready to Start -->
**Dependencies**: None  
**Estimated Time**: 3 hours

#### Tasks
- [x] **TASK-01-A**: Backup current codebase  <!-- Was: [ ] -->
  - ✅ Create backup commit: abc1234
  - ✅ Create git tag: v0.1.0-pre-migration
  - **Deliverable**: Commit hash: abc1234

- [ ] **TASK-01-B**: Clone reference repositories
  - In progress...
```

---

## 🔔 Notification System

### Cách notify agent tiếp theo

Khi agent hoàn thành, post một comment trong PROJECT_TRACKING.md:

```markdown
---

## AGENT-01 Completion Report
**Agent ID**: AGENT-01  
**Status**: ✅ Complete  
**Completion Time**: 2026-05-03 14:30  

### Deliverables
- [docs/v0375-diff.txt](docs/v0375-diff.txt) - 127 lines
- [docs/repo-setup-report.md](docs/repo-setup-report.md)

### Verification
```
$ git tag | grep v0.1.0-pre-migration
v0.1.0-pre-migration

$ ls /tmp/pocketbase-v0.37.5/
cmd/  core/  daos/  migrations/  ...

$ wc -l docs/v0375-diff.txt
127 docs/v0375-diff.txt
```

### Blockers Removed
- ✅ AGENT-02 can now start (TASK-01-B completed)
- ✅ AGENT-03 can now start (TASK-01-B completed)
- ✅ AGENT-04 can now start (TASK-01-B completed)

### Next Agents Ready
- AGENT-02, AGENT-03, AGENT-04 (can start in parallel)

---
```

---

## 🛠️ Useful Commands

### Kiểm tra dependencies đã hoàn thành
```bash
# Xem agent nào đang chờ
grep "🟡 Blocked" PROJECT_TRACKING.md -B 5

# Xem agent nào ready to start
grep "🟢 Ready to Start" PROJECT_TRACKING.md -B 5

# Xem agent nào đang làm việc
grep "🔵 In Progress" PROJECT_TRACKING.md -B 5

# Xem agent nào đã hoàn thành
grep "✅ Complete" PROJECT_TRACKING.md -B 5
```

### Track overall progress
```bash
# Count completed tasks
grep -c "\\[x\\]" PROJECT_TRACKING.md

# Count total tasks
grep -c "\\[ \\]\\|\\[x\\]" PROJECT_TRACKING.md

# Calculate percentage
echo "scale=2; $(grep -c '\\[x\\]' PROJECT_TRACKING.md) * 100 / $(grep -c '\\[ \\]\\|\\[x\\]' PROJECT_TRACKING.md)" | bc
```

---

## 🚨 Troubleshooting

### Agent bị block không rõ lý do
```bash
# Tìm dependencies của agent
grep "AGENT-XX" PROJECT_TRACKING.md -A 5 | grep "Dependencies"

# Kiểm tra agent dependency đã xong chưa
grep "AGENT-YY" PROJECT_TRACKING.md -A 50 | grep -E "Status|\\[x\\]"
```

### Task thất bại
```markdown
## AGENT-XX Failure Report

**Agent ID**: AGENT-XX  
**Failed Task**: TASK-XX-Y  
**Failure Time**: 2026-05-03 15:45  

### Error Description
[Mô tả lỗi chi tiết]

### Root Cause
[Nguyên nhân]

### Resolution Attempted
- [ ] Retry with different approach
- [ ] Request help from another agent
- [ ] Escalate to team lead

### Impact
- Blocks: AGENT-ZZ (TASK-ZZ-A)
- ETA delay: +X hours
```

---

## 📝 Best Practices

### 1. Atomic Commits
Mỗi task nên tạo một commit riêng:
```bash
git commit -m "AGENT-05/TASK-05-A: Create core/db_postgresql.go"
```

### 2. Document Everything
Mỗi agent nên tạo các file documentation theo deliverables:
```
docs/
├── agent-01-repo-setup.md
├── agent-02-code-extraction.md
├── agent-03-dependency-analysis.md
└── ...
```

### 3. Verify Before Handoff
Trước khi notify agent tiếp theo, chạy verification commands:
```bash
# Ví dụ cho AGENT-02
wc -l docs/extracted/*.go
grep -c "func " docs/extracted/*.go
```

### 4. Keep Communication Protocol Consistent
Luôn dùng format chuẩn trong completion reports để dễ parse.

---

## 🎓 Example Workflow

### Scenario: Khởi động Phase 1

**Step 1**: User khởi động AGENT-01
```bash
# Terminal 1
claude "Tôi là AGENT-01: Repository Setup Specialist..."
```

**Step 2**: AGENT-01 hoàn thành, update PROJECT_TRACKING.md
```markdown
### AGENT-01: Repository Setup Specialist
**Status**: ✅ Complete
...
```

**Step 3**: AGENT-01 post completion report
```markdown
## AGENT-01 Completion Report
...
### Blockers Removed
- ✅ AGENT-02 can now start
- ✅ AGENT-03 can now start
- ✅ AGENT-04 can now start
```

**Step 4**: User khởi động AGENT-02, 03, 04 song song
```bash
# Terminal 1
claude "Tôi là AGENT-02: Code Extraction Specialist..."

# Terminal 2
claude "Tôi là AGENT-03: Dependency Analysis Specialist..."

# Terminal 3
claude "Tôi là AGENT-04: dbx Package Specialist..."
```

**Step 5**: Chờ Phase 1 hoàn thành, chuyển sang Phase 2
```bash
# Check if Phase 1 done
grep "Phase 1" PROJECT_TRACKING.md -A 200 | grep "✅ Complete" | wc -l
# Should output: 4 (AGENT-01, 02, 03, 04)

# Start Phase 2
claude "Tôi là AGENT-05: Database Connection Layer Engineer..."
```

---

## 📚 Additional Resources

- [PROJECT_TRACKING.md](PROJECT_TRACKING.md) - Main tracking document
- [POSTGRES_MIGRATION_STRATEGY.md](POSTGRES_MIGRATION_STRATEGY.md) - Overall strategy
- [POSTGREBASE_CODE_ANALYSIS.md](POSTGREBASE_CODE_ANALYSIS.md) - Code analysis reference
- [POCKETBASE_V0375_TO_POSTGRES_PLAN.md](POCKETBASE_V0375_TO_POSTGRES_PLAN.md) - Detailed implementation plan

---

**Last Updated**: 2026-05-03  
**Version**: 1.0
