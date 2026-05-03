# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

---

## ⚠️ CRITICAL: Project Type

**This is a RESEARCH & DOCUMENTATION repository, NOT a working codebase.**

- ❌ NO actual Go code to run/test
- ❌ NO `go.mod` to build
- ❌ NO binary to execute
- ✅ Implementation GUIDES only
- ✅ Code TEMPLATES only (extracted from postgrebase)
- ✅ Test PLANS only (not actual tests)

**Purpose**: Planning and preparation for future implementation

---

## 📋 Tổng quan dự án

Dự án này là **nghiên cứu và document hóa** migration path từ PocketBase v0.37.5 sang PostgreSQL + Redis, dựa trên phân tích [zhenruyan/postgrebase](https://github.com/zhenruyan/postgrebase).

**Triết lý chính**: Minimal Changes (~200 dòng code) để đạt được PostgreSQL/MySQL support.

**Current Status**: Research complete, implementation NOT started.

---

## 🎯 Multi-Agent Workflow

Dự án này được thiết kế để thực thi với **12 agents chuyên môn** làm việc song song và tuần tự. Mỗi agent có một lĩnh vực chuyên môn và một tập hợp các tasks cụ thể.

### Khi nào dùng Multi-Agent?

**LUÔN** sử dụng multi-agent workflow khi:
- User yêu cầu thực thi implementation
- Cần phân chia công việc phức tạp thành các subtasks
- Có nhiều tasks độc lập có thể chạy song song

### Cách khởi động Agent

```bash
# Xem agent nào sẵn sàng
grep "🟢 Ready to Start" PROJECT_TRACKING.md -B 3

# Copy prompt template từ QUICK_REFERENCE.md
grep -A 20 "AGENT-01 Prompt" QUICK_REFERENCE.md

# Khởi động agent (ví dụ AGENT-01)
# Sử dụng prompt từ QUICK_REFERENCE.md
```

### Agent Dependencies Graph

```
Phase 1: AGENT-01 → (02, 03, 04 parallel)
Phase 2: (02+03) → AGENT-05 → (06, 07 parallel)
Phase 3: (04+05) → AGENT-08 → AGENT-09
Phase 4: (05+06+08) → AGENT-10 → AGENT-11 → AGENT-12
```

---

## 📁 Cấu trúc thư mục

```
postgres-migration-research/
├── README.md                              # Tổng quan dự án (ĐỌC ĐẦU TIÊN)
├── START_HERE.md                          # Quick start guide (2 phút)
├── PROJECT_TRACKING.md                    # ⭐ Main task tracking
├── AGENT_WORKFLOW_GUIDE.md                # Agent prompt templates
├── DASHBOARD.md                           # Visual progress tracking
├── QUICK_REFERENCE.md                     # Cheat sheet (BOOKMARK)
├── MULTI_AGENT_SUMMARY.md                 # System overview
│
├── POSTGRES_MIGRATION_STRATEGY.md         # Chiến lược tổng quan
├── POSTGREBASE_CODE_ANALYSIS.md           # Line-by-line code analysis
├── POCKETBASE_V0375_TO_POSTGRES_PLAN.md  # 10-day implementation plan
│
└── docs/                                  # Agent deliverables
    ├── extracted/                         # Code snippets từ postgrebase
    ├── *.md                               # Documentation từ agents
    └── *.txt                              # Reports & diffs
```

### File ưu tiên đọc

1. **START_HERE.md** (2 phút) - Quick orientation
2. **README.md** (5 phút) - Project overview
3. **QUICK_REFERENCE.md** (bookmark!) - Daily cheat sheet
4. **PROJECT_TRACKING.md** - Khi cần biết task chi tiết
5. **AGENT_WORKFLOW_GUIDE.md** - Khi khởi động agent

---

## 🚀 Lệnh thường dùng

### Tracking Progress

```bash
# Xem agent nào ready to start
grep "🟢 Ready to Start" PROJECT_TRACKING.md -B 3

# Xem agent nào đang blocked
grep "🟡 Blocked" PROJECT_TRACKING.md -B 3

# Đếm tasks đã hoàn thành
grep -c "\[x\]" PROJECT_TRACKING.md

# Tính % hoàn thành
echo "scale=2; $(grep -c '\[x\]' PROJECT_TRACKING.md) * 100 / $(grep -c '\[ \]\|\[x\]' PROJECT_TRACKING.md)" | bc
```

### Agent Status

```bash
# View all agent status
grep -E "AGENT-[0-9]{2}.*Status" PROJECT_TRACKING.md

# Count completed agents
grep -c "Status.*✅ Complete" PROJECT_TRACKING.md

# Find next agent to run
grep "🟢 Ready to Start" PROJECT_TRACKING.md -B 3 | grep "AGENT"
```

### Verification Commands

```bash
# Verify backups (AGENT-01)
git tag | grep v0.1.0-pre-migration

# Check cloned repos (AGENT-01)
ls /tmp/pocketbase-v0.37.5/
ls /tmp/postgrebase/

# Count extracted code lines (AGENT-02)
wc -l docs/extracted/*.go
```

---

## 🏗️ Kiến trúc Migration

### Core Changes (chỉ ~200 dòng!)

| Component | Lines | Complexity | Files |
|-----------|-------|------------|-------|
| Database connection | 29 | ⭐ Very Simple | `core/db_postgresql.go` |
| Redis integration | ~100 | ⭐⭐ Simple | `core/base.go` |
| Config/flags | ~10 | ⭐ Trivial | `cmd/serve.go` |
| Query builder | (vendor) | ⭐⭐⭐ Complex | `vendor/dbx/` |

### Key Design Patterns

1. **Dependency Injection**: DSN strings passed to config
2. **Factory Pattern**: `connectDB(dsn)` auto-detects driver (postgres/mysql)
3. **Graceful Degradation**: Redis fails → fallback to in-memory
4. **Connection Pooling**: Separate concurrent/nonconcurrent pools
5. **Pub/Sub Pattern**: Redis channels for multi-node sync

### SQL Dialect Conversion

```go
// SQLite → PostgreSQL
"AUTOINCREMENT" → "SERIAL"
"DATETIME('now')" → "NOW()"
"||" (concat) → "||" (same)
"`backticks`" → "\"quotes\""
```

---

## 🔧 Development Workflow

### Bắt đầu một Agent Task

1. **Đọc task details**:
   ```bash
   grep "AGENT-XX" PROJECT_TRACKING.md -A 50
   ```

2. **Copy prompt template**:
   ```bash
   grep -A 20 "AGENT-XX Prompt" QUICK_REFERENCE.md
   ```

3. **Khởi động agent** với prompt đã copy

4. **Thực hiện tasks** theo checklist trong PROJECT_TRACKING.md

5. **Update progress**:
   - Đổi `[ ]` → `[x]` cho mỗi task hoàn thành
   - Update agent status: 🟢 → 🔵 → ✅

6. **Post completion report** theo format:
   ```markdown
   ## AGENT-XX Completion Report
   **Agent ID**: AGENT-XX
   **Status**: ✅ Complete
   **Completion Time**: 2026-05-03 14:30
   
   ### Deliverables
   - [file1.md](path/to/file1.md)
   - [file2.go](path/to/file2.go)
   
   ### Verification
   [Command outputs proving success]
   
   ### Blockers Removed
   - ✅ AGENT-YY can now start
   
   ### Next Agents Ready
   - AGENT-YY, AGENT-ZZ (can start in parallel)
   ```

### Commit Convention

```bash
git commit -m "AGENT-XX/TASK-XX-Y: Brief description"

# Ví dụ
git commit -m "AGENT-01/TASK-01-A: Create backup tag v0.1.0-pre-migration"
git commit -m "AGENT-05/TASK-05-A: Implement core/db_postgresql.go"
```

---

## 📚 Tài liệu kỹ thuật chính

### 1. POSTGRES_MIGRATION_STRATEGY.md
- Architecture overview (SQLite → PostgreSQL)
- 3-week implementation roadmap
- Risk assessment & mitigation
- **Đọc khi**: Cần hiểu chiến lược tổng quan

### 2. POSTGREBASE_CODE_ANALYSIS.md
- Line-by-line code analysis
- Database connection layer (29 dòng)
- Redis integration (~100 dòng)
- Code examples với annotations
- **Đọc khi**: Cần implement cụ thể

### 3. POCKETBASE_V0375_TO_POSTGRES_PLAN.md
- Day-by-day implementation plan (10 ngày)
- Step-by-step instructions
- Test scenarios & validation
- **Đọc khi**: Đang implement

---

## 🎯 Timeline & Phases

| Phase | Days | Agents | Focus | Status |
|-------|------|--------|-------|--------|
| **Research** | ✅ | N/A | Analysis & planning | ✅ Done |
| **Phase 1** | 2 | 01-04 | Preparation | ⏳ Next |
| **Phase 2** | 3 | 05-07 | Core implementation | ⏳ Pending |
| **Phase 3** | 2 | 08-09 | Query builder | ⏳ Pending |
| **Phase 4** | 3 | 10-12 | Testing | ⏳ Pending |
| **TOTAL** | **10** | **12** | **2 weeks** | 📅 ETA: 2026-05-17 |

---

## 🔍 Troubleshooting

### Agent bị blocked không rõ lý do
```bash
# Tìm dependencies của agent
grep "AGENT-XX" PROJECT_TRACKING.md -A 5 | grep "Dependencies"

# Kiểm tra agent dependency đã xong chưa
grep "AGENT-YY" PROJECT_TRACKING.md -A 50 | grep -E "Status|\\[x\\]"
```

### Task thất bại
Post failure report trong PROJECT_TRACKING.md:
```markdown
## AGENT-XX Failure Report

**Agent ID**: AGENT-XX
**Failed Task**: TASK-XX-Y
**Failure Time**: YYYY-MM-DD HH:MM

### Error Description
[Mô tả lỗi chi tiết]

### Root Cause
[Nguyên nhân]

### Resolution
- [ ] Retry with different approach
- [ ] Request help from another agent
- [ ] Escalate to team lead

### Impact
- Blocks: AGENT-ZZ (TASK-ZZ-A)
- ETA delay: +X hours
```

### Không tìm thấy file reference
```bash
# Kiểm tra deliverables section
grep "AGENT-XX" PROJECT_TRACKING.md -A 100 | grep "Deliverable"

# Tìm file trong docs/
find docs/ -name "*agent*" -o -name "*extracted*"
```

---

## 💡 Best Practices

### 1. Luôn đọc context trước khi bắt đầu
- Đọc PROJECT_TRACKING.md cho task details
- Đọc reference docs (POSTGREBASE_CODE_ANALYSIS.md) cho implementation guidance
- Check dependencies đã hoàn thành chưa

### 2. Atomic commits
```bash
# One commit per task
git commit -m "AGENT-05/TASK-05-A: Create core/db_postgresql.go"
```

### 3. Verify before handoff
```bash
# Run verification commands trước khi post completion report
wc -l docs/extracted/*.go
grep -c "func " docs/extracted/*.go
git diff --stat
```

### 4. Document everything
Mỗi agent nên tạo deliverable files:
```
docs/
├── agent-01-repo-setup.md
├── agent-02-code-extraction.md
├── extracted/
│   ├── db_postgresql.go
│   ├── redis_functions.go
│   └── ...
```

### 5. Update tracking immediately
- Đổi checkbox `[ ]` → `[x]` ngay sau khi hoàn thành task
- Update agent status: 🟢 → 🔵 → ✅
- Update DASHBOARD.md với progress mới

---

## 🔗 External References

- **Postgrebase**: https://github.com/zhenruyan/postgrebase
- **PocketBase v0.37.5**: https://github.com/pocketbase/pocketbase/releases/tag/v0.37.5
- **PostgreSQL Driver**: https://github.com/lib/pq
- **Redis Client**: https://github.com/redis/go-redis
- **dbx Query Builder**: https://github.com/go-ozzo/ozzo-dbx

---

## 🚨 Important Notes

### Không phải Go project để implement!
Đây là **research project** với documentation only. Không có actual Go codebase để edit trong repo này.

### Purpose của dự án
- Nghiên cứu cách postgrebase implement PostgreSQL support
- Document implementation patterns
- Create detailed migration plan
- Chuẩn bị cho implementation ở repo khác (postgrebase-super-backend)

### Khi nào implement thực sự?
Implementation sẽ diễn ra ở repo khác:
```bash
cd ~/Documents/postgrebase-super-backend
# Follow plan từ POCKETBASE_V0375_TO_POSTGRES_PLAN.md
```

---

## 📞 Getting Help

### Quick Reference
```bash
# Bookmark this!
cat QUICK_REFERENCE.md
```

### Full Agent Guide
```bash
# Detailed instructions cho tất cả agents
cat AGENT_WORKFLOW_GUIDE.md
```

### Visual Dashboard
```bash
# Visual progress tracking
cat DASHBOARD.md
```

---

**Last Updated**: 2026-05-03  
**Project Status**: 📋 Research Complete - Ready for Implementation  
**Next Step**: Start AGENT-01 (Repository Setup Specialist)
