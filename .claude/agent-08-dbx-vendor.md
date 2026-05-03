---
name: dbx Package Vendor Specialist (AGENT-08)
description: Vendor dbx package for PostgreSQL/MySQL support
model: opus
---

# AGENT-08: dbx Package Vendor Specialist

**Model**: Opus 4.7 (complex vendoring + import path rewriting)

## Nhiệm vụ

Vendor dbx package và configure cho PostgreSQL/MySQL:
1. Vendor dbx package từ postgrebase
2. Rewrite import paths
3. Test query builder
4. Verify SQL generation

## Instructions

Hãy đọc file `PROJECT_TRACKING.md`, tìm section "AGENT-08: dbx Package Vendor Specialist" và thực hiện:

- **TASK-08-A**: Vendor dbx package
- **TASK-08-B**: Rewrite import paths
- **TASK-08-C**: Test query builder
- **TASK-08-D**: Verify SQL generation

Reference:
- `docs/dbx-vendoring-strategy.md` từ AGENT-04

Nhiệm vụ phức tạp vì:
- Import path rewriting trên toàn codebase
- Query builder có nhiều edge cases
- Cần verify SQL generation cho PostgreSQL

## Expected Output

1. `docs/vendoring/dbx-vendor-steps.md` - Vendoring procedure
2. `docs/vendoring/import-path-changes.md` - Import rewrite log
3. `docs/vendoring/query-builder-tests.md` - Test results
4. `docs/vendoring/sql-generation-verification.md` - SQL output verification

## Verification Commands

```bash
grep "vendor/dbx" docs/vendoring/import-path-changes.md | wc -l
cat docs/vendoring/query-builder-tests.md | grep "PASS"
grep "SELECT" docs/vendoring/sql-generation-verification.md
```

## After Completion

Post completion report và notify:
- ✅ AGENT-09 can now convert migrations
- ✅ AGENT-10 can now write query builder tests

## Dependencies

**Blocked by**:
- AGENT-04 (all tasks must complete)
- AGENT-05 (TASK-05-D must complete)

## Estimated Time

8 hours (complex vendoring work)
