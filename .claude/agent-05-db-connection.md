---
name: Database Connection Layer Engineer (AGENT-05)
description: Implement PostgreSQL/MySQL connection layer
model: sonnet
---

# AGENT-05: Database Connection Layer Engineer

**Model**: Sonnet 4.5 (cần code generation với best practices)

## Nhiệm vụ

Implement database connection layer cho PostgreSQL/MySQL:
1. Create `core/db_postgresql.go`
2. Update `core/base.go` struct
3. Update go.mod với dependencies mới
4. Run compilation test

## Instructions

Hãy đọc file `PROJECT_TRACKING.md`, tìm section "AGENT-05: Database Connection Layer Engineer" và thực hiện:

- **TASK-05-A**: Create core/db_postgresql.go
- **TASK-05-B**: Update core/base.go struct
- **TASK-05-C**: Update go.mod
- **TASK-05-D**: Compilation test

Reference:
- `docs/extracted/db_postgresql.go` từ AGENT-02
- `docs/new-dependencies.txt` từ AGENT-03

## Expected Output

1. `docs/implementation/core_db_postgresql.go` - Implementation guide
2. `docs/implementation/core_base_changes.md` - Struct changes
3. `docs/compilation-test-results.md` - Compilation verification

## Verification Commands

```bash
cat docs/implementation/core_db_postgresql.go | wc -l  # Should be ~29 lines
grep "connectDB" docs/implementation/core_db_postgresql.go
cat docs/compilation-test-results.md | grep "PASS"
```

## After Completion

Post completion report và notify:
- ✅ AGENT-06 can now start Redis integration
- ✅ AGENT-07 can now start CLI flags
- ✅ AGENT-08 can now start dbx vendoring

## Dependencies

**Blocked by**: 
- AGENT-02 (all tasks must complete)
- AGENT-03 (TASK-03-B must complete)

## Estimated Time

6 hours
