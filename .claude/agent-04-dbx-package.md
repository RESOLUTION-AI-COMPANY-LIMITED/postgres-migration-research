---
name: dbx Package Specialist (AGENT-04)
description: Analyze dbx package and SQL dialect conversion
model: opus
---

# AGENT-04: dbx Package Specialist

**Model**: Opus 4.7 (phức tạp nhất - cần deep analysis của query builder)

## Nhiệm vụ

Analyze dbx package và create vendoring strategy:
1. Locate dbx package in postgrebase
2. Compare dbx versions (original vs fork)
3. Create vendoring strategy
4. Create SQL dialect conversion matrix (SQLite → PostgreSQL)

## Instructions

Hãy đọc file `PROJECT_TRACKING.md`, tìm section "AGENT-04: dbx Package Specialist" và thực hiện:

- **TASK-04-A**: Locate dbx package
- **TASK-04-B**: Compare dbx versions
- **TASK-04-C**: Vendoring strategy
- **TASK-04-D**: SQL dialect conversion matrix

Nhiệm vụ này phức tạp vì:
- dbx là modified fork
- Cần hiểu query builder internals
- SQL dialect conversion có nhiều edge cases

## Expected Output

1. `docs/dbx-package-location.md` - Package structure analysis
2. `docs/dbx-version-comparison.md` - Fork vs original diff
3. `docs/dbx-vendoring-strategy.md` - How to vendor the package
4. `docs/sql-dialect-conversion-matrix.md` - Comprehensive conversion rules

## Verification Commands

```bash
cat docs/dbx-package-location.md | grep "vendor/"
cat docs/dbx-version-comparison.md | grep "diff"
cat docs/sql-dialect-conversion-matrix.md | grep "AUTOINCREMENT"
```

## After Completion

Post completion report và notify:
- ✅ AGENT-08 can now vendor dbx package
- ✅ AGENT-09 can use SQL conversion matrix

## Dependencies

**Blocked by**: AGENT-01 (TASK-01-B must complete first)

## Estimated Time

5 hours (most complex in Phase 1)
