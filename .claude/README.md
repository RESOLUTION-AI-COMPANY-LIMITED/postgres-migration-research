# .claude/ - Agent Definitions

Thư mục này chứa **12 agent definitions** cho PostgreSQL Migration project, mỗi agent được cấu hình với model phù hợp cho nhiệm vụ của nó.

---

## 🤖 Model Selection Strategy

### Haiku 4.5 (Fast & Cost-Effective)
**Dùng cho**: Tác vụ đơn giản, không cần reasoning phức tạp

- **AGENT-01**: Repository Setup (backup, clone, diff)
- **AGENT-07**: CLI Flags (simple flag additions)

**Ưu điểm**: Nhanh, tiết kiệm chi phí, đủ cho structural tasks

---

### Sonnet 4.5 (Balanced)
**Dùng cho**: Code generation, pattern analysis, test writing

- **AGENT-02**: Code Extraction (pattern recognition)
- **AGENT-03**: Dependency Analysis (version comparison)
- **AGENT-05**: Database Connection Layer (code implementation)
- **AGENT-06**: Redis Integration (Pub/Sub patterns)
- **AGENT-09**: Migration System (SQL conversion)
- **AGENT-10**: Unit Tests (test generation)
- **AGENT-11**: Integration Tests (E2E scenarios)

**Ưu điểm**: Balance giữa speed và capability, tốt cho standard development tasks

---

### Opus 4.7 (Deep Analysis)
**Dùng cho**: Phức tạp nhất - deep analysis, optimization

- **AGENT-04**: dbx Package Analysis (query builder internals)
- **AGENT-08**: dbx Vendoring (complex import rewriting)
- **AGENT-12**: Performance Benchmarks (optimization recommendations)

**Ưu điểm**: Strongest reasoning, cần thiết cho complex architectural decisions

---

## 📋 Agent Overview

| Agent | Model | Complexity | Focus Area | Est. Time |
|-------|-------|------------|------------|-----------|
| 01 | Haiku | ⭐ Simple | Repo Setup | 3h |
| 02 | Sonnet | ⭐⭐ Moderate | Code Extraction | 4h |
| 03 | Sonnet | ⭐⭐ Moderate | Dependencies | 3h |
| 04 | Opus | ⭐⭐⭐ Complex | dbx Package | 5h |
| 05 | Sonnet | ⭐⭐ Moderate | DB Connection | 6h |
| 06 | Sonnet | ⭐⭐ Moderate | Redis | 6h |
| 07 | Haiku | ⭐ Simple | CLI Flags | 3h |
| 08 | Opus | ⭐⭐⭐ Complex | dbx Vendor | 8h |
| 09 | Sonnet | ⭐⭐ Moderate | Migrations | 8h |
| 10 | Sonnet | ⭐⭐ Moderate | Unit Tests | 8h |
| 11 | Sonnet | ⭐⭐ Moderate | Integration Tests | 8h |
| 12 | Opus | ⭐⭐⭐ Complex | Benchmarks | 6h |

**Total**: 68 hours (8.5 working days)

---

## 🚀 Cách sử dụng Agent Definitions

### Option 1: Manual Invocation (với model selection)

```bash
# Start agent với model được định nghĩa
claude --agent .claude/agent-01-repo-setup.md

# Hoặc với model override
claude --agent .claude/agent-04-dbx-package.md --model opus
```

### Option 2: Copy Prompt Template

```bash
# Đọc agent definition
cat .claude/agent-01-repo-setup.md

# Copy prompt và instructions
# Khởi động Claude Code với prompt đó
```

### Option 3: Reference from QUICK_REFERENCE.md

```bash
# QUICK_REFERENCE.md có sẵn prompts cho tất cả agents
grep -A 20 "AGENT-01 Prompt" ../QUICK_REFERENCE.md
```

---

## 📁 File Structure

```
.claude/
├── README.md                          # This file
├── agent-01-repo-setup.md             # Haiku - Repo setup
├── agent-02-code-extraction.md        # Sonnet - Code extraction
├── agent-03-dependency-analysis.md    # Sonnet - Dependency analysis
├── agent-04-dbx-package.md            # Opus - dbx deep analysis
├── agent-05-db-connection.md          # Sonnet - DB connection
├── agent-06-redis.md                  # Sonnet - Redis integration
├── agent-07-cli-flags.md              # Haiku - CLI flags
├── agent-08-dbx-vendor.md             # Opus - dbx vendoring
├── agent-09-migrations.md             # Sonnet - SQL migrations
├── agent-10-unit-tests.md             # Sonnet - Unit tests
├── agent-11-integration-tests.md      # Sonnet - Integration tests
└── agent-12-benchmarks.md             # Opus - Performance analysis
```

---

## 🎯 Model Selection Rationale

### Khi nào dùng Haiku?
- ✅ File operations (copy, move, backup)
- ✅ Simple string manipulation
- ✅ Directory structure comparisons
- ✅ Basic CLI modifications
- ❌ NOT for: Deep analysis, code generation, optimization

### Khi nào dùng Sonnet?
- ✅ Code implementation (standard patterns)
- ✅ Test generation
- ✅ Documentation writing
- ✅ Pattern recognition & extraction
- ✅ SQL conversion (with reference guide)
- ❌ NOT for: Query builder internals, complex optimizations

### Khi nào dùng Opus?
- ✅ Deep architectural analysis
- ✅ Query builder internals
- ✅ Complex import path rewriting
- ✅ Performance optimization recommendations
- ✅ Complex dependency graphs
- ⚠️ Use sparingly: High cost, save for truly complex tasks

---

## 💰 Cost Optimization Strategy

### Parallel Execution (Save Time, Not Cost)

**Phase 1** - Run in parallel:
- AGENT-02 (Sonnet) + AGENT-03 (Sonnet) + AGENT-04 (Opus)
- Total: ~5 hours parallel vs 12 hours sequential

**Phase 2** - Run in parallel:
- AGENT-06 (Sonnet) + AGENT-07 (Haiku)
- Total: ~6 hours parallel vs 9 hours sequential

**Phase 4** - Sequential (dependencies):
- AGENT-10 → AGENT-11 → AGENT-12 (must run in order)

### Model Selection Saves Cost

| Scenario | All Opus | Optimized | Savings |
|----------|----------|-----------|---------|
| **Phase 1** | 4 × Opus | 1 Opus + 2 Sonnet + 1 Haiku | ~60% |
| **Phase 2** | 3 × Opus | 2 Sonnet + 1 Haiku | ~70% |
| **Phase 3** | 2 × Opus | 1 Opus + 1 Sonnet | ~40% |
| **Phase 4** | 3 × Opus | 1 Opus + 2 Sonnet | ~50% |
| **TOTAL** | ~$XXX | ~$XX | **~50-60% savings** |

*(Exact pricing depends on Claude API rates)*

---

## 🔍 Verification

Sau khi mỗi agent hoàn thành, verify:

```bash
# Check agent completed tasks
grep "AGENT-XX" ../PROJECT_TRACKING.md -A 50 | grep "\[x\]"

# Check deliverables exist
ls -la ../docs/agent-XX-*

# Verify against success criteria
cat .claude/agent-XX-*.md | grep "Verification Commands" -A 10
```

---

## 📚 References

- **Main Project Docs**: `../CLAUDE.md`
- **Task Tracking**: `../PROJECT_TRACKING.md`
- **Quick Reference**: `../QUICK_REFERENCE.md`
- **Agent Workflow Guide**: `../AGENT_WORKFLOW_GUIDE.md`

---

## 🎓 Best Practices

### 1. Always Check Dependencies

```bash
# Before starting AGENT-XX, verify dependencies
grep "AGENT-XX" ../PROJECT_TRACKING.md -A 5 | grep "Dependencies"
grep "AGENT-YY" ../PROJECT_TRACKING.md | grep "Status.*✅"
```

### 2. Use Correct Model

Don't override model unless necessary:
- Haiku → Opus: Only if task proves too complex
- Opus → Sonnet: Never downgrade complex tasks
- Sonnet → Haiku: Only for trivial subtasks

### 3. Document Model Performance

If agent struggles or excels, note it:
```markdown
## Model Performance Note
**Agent**: AGENT-XX
**Model**: Sonnet
**Observation**: Task took 8h (estimated 6h), consider Opus for similar tasks
```

---

**Created**: 2026-05-03  
**Version**: 1.0  
**Last Updated**: 2026-05-03
