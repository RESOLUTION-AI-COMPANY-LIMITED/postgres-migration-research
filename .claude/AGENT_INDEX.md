# Agent Index - Quick Reference

**12 specialized agents cho PostgreSQL Migration project**

---

## 🚦 Agent Status Legend

- 🟢 **Ready to Start** - All dependencies met
- 🟡 **Blocked** - Waiting for dependencies
- 🔵 **In Progress** - Agent working
- ✅ **Complete** - Verified & merged

---

## 📋 Complete Agent List

### Phase 1: Preparation & Analysis (Day 1-2)

#### AGENT-01: Repository Setup Specialist
- **File**: [agent-01-repo-setup.md](agent-01-repo-setup.md)
- **Model**: Haiku 4.5
- **Time**: 3h
- **Dependencies**: None
- **Tasks**: Backup, clone repos, compare structure

#### AGENT-02: Code Extraction Specialist
- **File**: [agent-02-code-extraction.md](agent-02-code-extraction.md)
- **Model**: Sonnet 4.5
- **Time**: 4h
- **Dependencies**: AGENT-01
- **Tasks**: Extract db_postgresql.go, Redis code, CLI flags

#### AGENT-03: Dependency Analysis Specialist
- **File**: [agent-03-dependency-analysis.md](agent-03-dependency-analysis.md)
- **Model**: Sonnet 4.5
- **Time**: 3h
- **Dependencies**: AGENT-01
- **Tasks**: Compare go.mod, list new deps, conflict resolution

#### AGENT-04: dbx Package Specialist
- **File**: [agent-04-dbx-package.md](agent-04-dbx-package.md)
- **Model**: Opus 4.7 ⭐
- **Time**: 5h
- **Dependencies**: AGENT-01
- **Tasks**: Locate dbx, compare versions, SQL conversion matrix

---

### Phase 2: Core Implementation (Day 2-3)

#### AGENT-05: Database Connection Layer Engineer
- **File**: [agent-05-db-connection.md](agent-05-db-connection.md)
- **Model**: Sonnet 4.5
- **Time**: 6h
- **Dependencies**: AGENT-02, AGENT-03
- **Tasks**: Create db_postgresql.go, update base.go, compilation test

#### AGENT-06: Redis Integration Engineer
- **File**: [agent-06-redis.md](agent-06-redis.md)
- **Model**: Sonnet 4.5
- **Time**: 6h
- **Dependencies**: AGENT-02, AGENT-05
- **Tasks**: Redis client init, Publish(), Subscribe()

#### AGENT-07: CLI Flags Engineer
- **File**: [agent-07-cli-flags.md](agent-07-cli-flags.md)
- **Model**: Haiku 4.5
- **Time**: 3h
- **Dependencies**: AGENT-02, AGENT-05
- **Tasks**: Add --dataDsn, --redisDsn flags

---

### Phase 3: Query Builder (Day 3-5)

#### AGENT-08: dbx Package Vendor Specialist
- **File**: [agent-08-dbx-vendor.md](agent-08-dbx-vendor.md)
- **Model**: Opus 4.7 ⭐
- **Time**: 8h
- **Dependencies**: AGENT-04, AGENT-05
- **Tasks**: Vendor dbx, rewrite imports, test query builder

#### AGENT-09: Migration System Engineer
- **File**: [agent-09-migrations.md](agent-09-migrations.md)
- **Model**: Sonnet 4.5
- **Time**: 8h
- **Dependencies**: AGENT-08, AGENT-04
- **Tasks**: Analyze migrations, create conversion script, convert to PostgreSQL

---

### Phase 4: Testing & Validation (Day 5-7)

#### AGENT-10: Unit Test Engineer
- **File**: [agent-10-unit-tests.md](agent-10-unit-tests.md)
- **Model**: Sonnet 4.5
- **Time**: 8h
- **Dependencies**: AGENT-05, AGENT-06, AGENT-08
- **Tasks**: DB tests, Redis tests, query builder tests

#### AGENT-11: Integration Test Engineer
- **File**: [agent-11-integration-tests.md](agent-11-integration-tests.md)
- **Model**: Sonnet 4.5
- **Time**: 8h
- **Dependencies**: AGENT-10
- **Tasks**: Docker setup, API tests, clustering tests, migration tests

#### AGENT-12: Performance Benchmark Engineer
- **File**: [agent-12-benchmarks.md](agent-12-benchmarks.md)
- **Model**: Opus 4.7 ⭐
- **Time**: 6h
- **Dependencies**: AGENT-11
- **Tasks**: DB benchmarks, Redis benchmarks, optimization recommendations

---

## 🎯 Quick Selection

### By Model

**Haiku** (fast, simple):
- AGENT-01, AGENT-07

**Sonnet** (balanced):
- AGENT-02, 03, 05, 06, 09, 10, 11

**Opus** (complex):
- AGENT-04, 08, 12

### By Phase

**Phase 1**: AGENT-01, 02, 03, 04  
**Phase 2**: AGENT-05, 06, 07  
**Phase 3**: AGENT-08, 09  
**Phase 4**: AGENT-10, 11, 12

### By Complexity

**Simple** (⭐):
- AGENT-01, 07

**Moderate** (⭐⭐):
- AGENT-02, 03, 05, 06, 09, 10, 11

**Complex** (⭐⭐⭐):
- AGENT-04, 08, 12

---

## 🔄 Execution Order

### Critical Path (Sequential)
```
01 → 02 → 05 → 08 → 09 → 10 → 11 → 12
```

### Parallel Opportunities
```
After 01: 02 + 03 + 04 (parallel)
After 05: 06 + 07 (parallel)
```

---

## 📊 Statistics

- **Total Agents**: 12
- **Total Hours**: 68h (sequential) / ~51h (parallel)
- **Total Tasks**: ~48 tasks
- **Deliverable Files**: ~50+ documentation files

### Model Distribution
- Haiku: 2 agents (9%)
- Sonnet: 7 agents (63%)
- Opus: 3 agents (28%)

---

## 🚀 Quick Start Commands

```bash
# View agent details
cat .claude/agent-01-repo-setup.md

# Launch agent
claude --agent .claude/agent-01-repo-setup.md

# Check dependencies
grep "AGENT-01" ../PROJECT_TRACKING.md -A 5 | grep "Dependencies"

# Verify completion
cat .claude/agent-01-repo-setup.md | grep "Verification Commands" -A 10
```

---

**Last Updated**: 2026-05-03  
**Version**: 1.0
