# PostgreSQL Migration - Quick Reference Card

**Cheat Sheet cho Multi-Agent Workflow**

---

## 🚀 Khởi động nhanh

### Bắt đầu một agent
```bash
cd /Users/accompany/Documents/postgres-migration-research

# Đọc prompt template cho agent
grep -A 20 "AGENT-01: Repository Setup Specialist" AGENT_WORKFLOW_GUIDE.md

# Khởi động agent với Claude Code
claude "Tôi là AGENT-01: Repository Setup Specialist trong dự án PostgreSQL Migration..."
```

### Kiểm tra dependencies
```bash
# Agent nào ready to start?
grep "🟢 Ready to Start" PROJECT_TRACKING.md -B 3

# Agent nào đang blocked?
grep "🟡 Blocked" PROJECT_TRACKING.md -B 3

# Agent nào đã xong?
grep "✅ Complete" PROJECT_TRACKING.md -B 3
```

---

## 📋 Agent Overview

| Agent | Role | Phase | Dependencies | ETA |
|-------|------|-------|--------------|-----|
| **01** | Repo Setup | 1 | None | 3h |
| **02** | Code Extraction | 1 | 01 | 4h |
| **03** | Dependency Analysis | 1 | 01 | 3h |
| **04** | dbx Package | 1 | 01 | 5h |
| **05** | DB Connection | 2 | 02, 03 | 6h |
| **06** | Redis Integration | 2 | 02, 05 | 6h |
| **07** | CLI Flags | 2 | 02, 05 | 3h |
| **08** | dbx Vendor | 3 | 04, 05 | 8h |
| **09** | Migrations | 3 | 08, 04 | 8h |
| **10** | Unit Tests | 4 | 05, 06, 08 | 8h |
| **11** | Integration Tests | 4 | 10 | 8h |
| **12** | Benchmarks | 4 | 11 | 6h |

---

## 🎯 Agent Prompt Templates (Copy-Paste)

### AGENT-01 Prompt
```
Tôi là AGENT-01: Repository Setup Specialist trong dự án PostgreSQL Migration.

Nhiệm vụ: Backup codebase, clone repos, compare structure

Hãy đọc PROJECT_TRACKING.md, tìm "AGENT-01" và thực hiện:
- TASK-01-A: Backup current codebase
- TASK-01-B: Clone reference repositories
- TASK-01-C: Directory structure comparison

Sau khi xong, tạo completion report theo Agent Communication Protocol.
```

### AGENT-02 Prompt
```
Tôi là AGENT-02: Code Extraction Specialist trong dự án PostgreSQL Migration.

Nhiệm vụ: Extract postgrebase implementation code

Dependencies: Chờ AGENT-01 hoàn thành TASK-01-B

Hãy đọc PROJECT_TRACKING.md, tìm "AGENT-02" và thực hiện:
- TASK-02-A: Extract db_postgresql.go
- TASK-02-B: Extract Redis functions
- TASK-02-C: Extract CLI flags
- TASK-02-D: Create code mapping document

Sau khi xong, notify AGENT-05, 06, 07 có thể bắt đầu.
```

### AGENT-03 Prompt
```
Tôi là AGENT-03: Dependency Analysis Specialist trong dự án PostgreSQL Migration.

Nhiệm vụ: Analyze Go dependencies và conflicts

Dependencies: Chờ AGENT-01 hoàn thành TASK-01-B

Hãy đọc PROJECT_TRACKING.md, tìm "AGENT-03" và thực hiện:
- TASK-03-A: Compare go.mod files
- TASK-03-B: List new dependencies
- TASK-03-C: Conflict resolution plan

Sau khi xong, notify AGENT-05.
```

### AGENT-04 Prompt
```
Tôi là AGENT-04: dbx Package Specialist trong dự án PostgreSQL Migration.

Nhiệm vụ: Analyze dbx package và vendoring strategy

Dependencies: Chờ AGENT-01 hoàn thành TASK-01-B

Hãy đọc PROJECT_TRACKING.md, tìm "AGENT-04" và thực hiện:
- TASK-04-A: Locate dbx package
- TASK-04-B: Compare dbx versions
- TASK-04-C: Vendoring strategy
- TASK-04-D: SQL dialect conversion matrix

Sau khi xong, notify AGENT-08, 09.
```

### AGENT-05 Prompt
```
Tôi là AGENT-05: Database Connection Layer Engineer trong dự án PostgreSQL Migration.

Nhiệm vụ: Implement PostgreSQL/MySQL connection layer

Dependencies: Chờ AGENT-02 (all) và AGENT-03 (TASK-03-B)

Hãy đọc PROJECT_TRACKING.md, tìm "AGENT-05" và thực hiện:
- TASK-05-A: Create core/db_postgresql.go
- TASK-05-B: Update core/base.go
- TASK-05-C: Update go.mod
- TASK-05-D: Compilation test

Sau khi xong, notify AGENT-06, 07, 08.
```

### AGENT-06 Prompt
```
Tôi là AGENT-06: Redis Integration Engineer trong dự án PostgreSQL Migration.

Nhiệm vụ: Implement Redis client & Pub/Sub

Dependencies: Chờ AGENT-02 (TASK-02-B) và AGENT-05 (TASK-05-B)

Hãy đọc PROJECT_TRACKING.md, tìm "AGENT-06" và thực hiện:
- TASK-06-A: Redis client initialization
- TASK-06-B: Publish() method
- TASK-06-C: Subscribe() for SSE
- TASK-06-D: Update go.mod

Sau khi xong, notify AGENT-10.
```

### AGENT-07 Prompt
```
Tôi là AGENT-07: CLI Flags Engineer trong dự án PostgreSQL Migration.

Nhiệm vụ: Add --dataDsn và --redisDsn flags

Dependencies: Chờ AGENT-02 (TASK-02-C) và AGENT-05 (TASK-05-B)

Hãy đọc PROJECT_TRACKING.md, tìm "AGENT-07" và thực hiện:
- TASK-07-A: Add flags to cmd/serve.go
- TASK-07-B: Pass DSN to Bootstrap()
- TASK-07-C: Environment variable support
- TASK-07-D: Update help text
```

### AGENT-08 Prompt
```
Tôi là AGENT-08: dbx Package Vendor Specialist trong dự án PostgreSQL Migration.

Nhiệm vụ: Vendor dbx package cho PostgreSQL/MySQL

Dependencies: Chờ AGENT-04 (all) và AGENT-05 (TASK-05-D)

Hãy đọc PROJECT_TRACKING.md, tìm "AGENT-08" và thực hiện:
- TASK-08-A: Vendor dbx package
- TASK-08-B: Rewrite import paths
- TASK-08-C: Test query builder
- TASK-08-D: Verify SQL generation

Sau khi xong, notify AGENT-09, 10.
```

### AGENT-09 Prompt
```
Tôi là AGENT-09: Migration System Engineer trong dự án PostgreSQL Migration.

Nhiệm vụ: Convert SQLite migrations sang PostgreSQL

Dependencies: Chờ AGENT-08 (all) và AGENT-04 (TASK-04-D)

Hãy đọc PROJECT_TRACKING.md, tìm "AGENT-09" và thực hiện:
- TASK-09-A: Analyze existing migrations
- TASK-09-B: Create conversion script
- TASK-09-C: Convert migrations
- TASK-09-D: Test migration rollback

Sau khi xong, notify AGENT-11.
```

### AGENT-10 Prompt
```
Tôi là AGENT-10: Unit Test Engineer trong dự án PostgreSQL Migration.

Nhiệm vụ: Write unit tests cho all components

Dependencies: Chờ AGENT-05, 06, 08

Hãy đọc PROJECT_TRACKING.md, tìm "AGENT-10" và thực hiện:
- TASK-10-A: Database connection tests
- TASK-10-B: Redis integration tests
- TASK-10-C: Query builder tests
- TASK-10-D: Run all tests

Sau khi xong, notify AGENT-11.
```

### AGENT-11 Prompt
```
Tôi là AGENT-11: Integration Test Engineer trong dự án PostgreSQL Migration.

Nhiệm vụ: E2E integration testing với PostgreSQL + Redis

Dependencies: Chờ AGENT-10

Hãy đọc PROJECT_TRACKING.md, tìm "AGENT-11" và thực hiện:
- TASK-11-A: Setup test infrastructure
- TASK-11-B: API compatibility tests
- TASK-11-C: Multi-node clustering test
- TASK-11-D: Data migration test

Sau khi xong, notify AGENT-12.
```

### AGENT-12 Prompt
```
Tôi là AGENT-12: Performance Benchmark Engineer trong dự án PostgreSQL Migration.

Nhiệm vụ: Performance testing & optimization

Dependencies: Chờ AGENT-11

Hãy đọc PROJECT_TRACKING.md, tìm "AGENT-12" và thực hiện:
- TASK-12-A: Database query benchmarks
- TASK-12-B: Redis cache benchmarks
- TASK-12-C: Run benchmark suite
- TASK-12-D: Performance optimization

Đây là agent cuối cùng!
```

---

## 🔍 Useful Commands

### Check Agent Status
```bash
# View all agent status
grep -E "AGENT-[0-9]{2}.*Status" PROJECT_TRACKING.md

# Count completed agents
grep -c "Status.*✅ Complete" PROJECT_TRACKING.md

# Count blocked agents
grep -c "Status.*🟡 Blocked" PROJECT_TRACKING.md
```

### Track Progress
```bash
# Completed tasks
grep -c "\[x\]" PROJECT_TRACKING.md

# Total tasks
grep -c "\[ \]\|\[x\]" PROJECT_TRACKING.md

# Calculate percentage
echo "scale=2; $(grep -c '\[x\]' PROJECT_TRACKING.md) * 100 / $(grep -c '\[ \]\|\[x\]' PROJECT_TRACKING.md)" | bc
```

### Find Next Agent
```bash
# Which agent can start now?
grep "🟢 Ready to Start" PROJECT_TRACKING.md -B 3 | grep "AGENT"

# Which agents are waiting?
grep "🟡 Blocked" PROJECT_TRACKING.md -B 5 | grep -E "AGENT|Dependencies"
```

---

## 📊 Progress Update Template

### After Task Completion
```markdown
## AGENT-XX Task Completion

**Task ID**: TASK-XX-Y  
**Status**: ✅ Complete  
**Time**: 2026-05-03 14:30  

### Deliverable
- [File link](path/to/file.md)

### Verification
[Command output showing success]

### Next
Handoff to: AGENT-YY
```

---

## 🔔 Notification Checklist

After agent completes:

- [ ] Update PROJECT_TRACKING.md (task checkboxes)
- [ ] Change agent status (🟢 → 🔵 → ✅)
- [ ] Post completion report
- [ ] Update DASHBOARD.md (progress bars)
- [ ] Notify dependent agents
- [ ] Create deliverable files
- [ ] Run verification commands
- [ ] Commit changes

---

## 🚨 Troubleshooting

### Agent không biết làm gì
→ Đọc lại prompt template trong AGENT_WORKFLOW_GUIDE.md

### Agent bị blocked
→ Kiểm tra dependencies: `grep "AGENT-XX" PROJECT_TRACKING.md -A 5 | grep Dependencies`

### Task thất bại
→ Post failure report trong PROJECT_TRACKING.md theo format

### Không tìm thấy file
→ Kiểm tra deliverables section trong PROJECT_TRACKING.md

---

## 📁 File Structure Quick Ref

```
postgres-migration-research/
├── README.md                              # Tổng quan dự án
├── PROJECT_TRACKING.md                    # ⭐ Main tracking doc
├── AGENT_WORKFLOW_GUIDE.md                # Agent instructions
├── DASHBOARD.md                           # Visual progress
├── QUICK_REFERENCE.md                     # This file
│
├── POSTGRES_MIGRATION_STRATEGY.md         # Strategy reference
├── POSTGREBASE_CODE_ANALYSIS.md           # Code reference
├── POCKETBASE_V0375_TO_POSTGRES_PLAN.md  # Implementation plan
│
└── docs/                                  # Agent deliverables
    ├── v0375-diff.txt                     # AGENT-01
    ├── extracted/                         # AGENT-02
    │   ├── db_postgresql.go
    │   ├── redis_functions.go
    │   └── cmd_serve_flags.go
    ├── dependency-matrix.md               # AGENT-03
    ├── dbx-vendoring-strategy.md          # AGENT-04
    └── ...
```

---

## 🎯 Daily Workflow

### Morning Standup
1. Check DASHBOARD.md for current status
2. Identify which agents can start today
3. Assign agents to team members
4. Set daily goals

### During Work
1. Follow agent prompt template
2. Read PROJECT_TRACKING.md for detailed tasks
3. Update task checkboxes as you complete
4. Post completion reports

### End of Day
1. Update DASHBOARD.md with progress
2. Notify blocked agents if dependencies met
3. Commit all changes
4. Review tomorrow's plan

---

## 🔗 Quick Links

| Document | Purpose |
|----------|---------|
| [DASHBOARD.md](DASHBOARD.md) | Visual progress tracking |
| [PROJECT_TRACKING.md](PROJECT_TRACKING.md) | Detailed task tracking |
| [AGENT_WORKFLOW_GUIDE.md](AGENT_WORKFLOW_GUIDE.md) | Agent instructions |
| [POSTGRES_MIGRATION_STRATEGY.md](POSTGRES_MIGRATION_STRATEGY.md) | Strategy overview |
| [POSTGREBASE_CODE_ANALYSIS.md](POSTGREBASE_CODE_ANALYSIS.md) | Code reference |

---

## 💡 Pro Tips

### For Speed
- Run Phase 1 agents (02, 03, 04) in parallel after AGENT-01
- Use multiple Claude Code sessions simultaneously
- Pre-read reference docs before starting agent

### For Quality
- Always run verification commands before completion report
- Double-check deliverable files exist
- Review generated code for correctness

### For Coordination
- Post completion reports immediately after finishing
- Update DASHBOARD.md daily for team visibility
- Use standardized format for easy parsing

---

**Version**: 1.0  
**Created**: 2026-05-03  
**Print This**: For quick desk reference!
