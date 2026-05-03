---
name: Migration System Engineer (AGENT-09)
description: Convert SQLite migrations to PostgreSQL
model: sonnet
---

# AGENT-09: Migration System Engineer

**Model**: Sonnet 4.5 (SQL conversion patterns)

## Nhiệm vụ

Convert existing SQLite migrations sang PostgreSQL:
1. Analyze existing SQLite migrations
2. Create SQL conversion script
3. Convert migrations to PostgreSQL
4. Test migration rollback

## Instructions

Hãy đọc file `PROJECT_TRACKING.md`, tìm section "AGENT-09: Migration System Engineer" và thực hiện:

- **TASK-09-A**: Analyze existing migrations
- **TASK-09-B**: Create conversion script
- **TASK-09-C**: Convert migrations
- **TASK-09-D**: Test migration rollback

Reference:
- `docs/sql-dialect-conversion-matrix.md` từ AGENT-04

Key conversions:
```sql
AUTOINCREMENT → SERIAL
DATETIME('now') → NOW()
`backticks` → "quotes"
```

## Expected Output

1. `docs/migrations/sqlite-migration-analysis.md` - Existing migrations inventory
2. `docs/migrations/conversion-script.sh` - Automated conversion script
3. `docs/migrations/postgresql-migrations/` - Converted migration files
4. `docs/migrations/rollback-test-results.md` - Rollback verification

## Verification Commands

```bash
cat docs/migrations/sqlite-migration-analysis.md | grep "CREATE TABLE" | wc -l
bash docs/migrations/conversion-script.sh --dry-run
ls -la docs/migrations/postgresql-migrations/ | wc -l
cat docs/migrations/rollback-test-results.md | grep "SUCCESS"
```

## After Completion

Post completion report và notify:
- ✅ AGENT-11 can now test data migration

## Dependencies

**Blocked by**:
- AGENT-08 (all tasks must complete)
- AGENT-04 (TASK-04-D must complete)

## Estimated Time

8 hours
