# Migration Conversion Guide (AGENT-09)

**Agent**: AGENT-09 - Migration System Engineer  
**Date**: 2026-05-03  
**Project**: PocketBase v0.37.5 → PostgreSQL + Redis Migration

---

## 📋 Overview

Hướng dẫn convert SQLite migrations sang PostgreSQL syntax. Đây là **RESEARCH PROJECT** - tạo documentation, KHÔNG phải actual implementation.

### Success Criteria
- ✅ SQLite → PostgreSQL conversion rules documented
- ✅ Automated conversion script defined
- ✅ Migration audit process documented
- ✅ Rollback testing strategy defined

---

## 🎯 Phase 3 - AGENT-09 Tasks

### Dependencies
- ✅ AGENT-04 (SQL Dialect Matrix) - Complete
- ✅ AGENT-08 (dbx Vendoring) - Complete

### Estimated Time
8 hours (research documentation)

---

## 🔍 STEP 1: Audit Existing Migrations

### 1.1 Locate Migration Files

**Expected Location:**
```
migrations/
├── 1640988000_initial.sql
├── 1650120000_add_users_table.sql
├── 1660250000_add_collections.sql
└── ... (more migrations)
```

**Command to list migrations:**
```bash
ls -la migrations/*.sql | wc -l
# Expected: ~20-50 migration files
```

### 1.2 Identify SQLite-Specific Syntax

**Common SQLite patterns to convert:**

| SQLite Syntax | PostgreSQL Equivalent | Frequency |
|---------------|----------------------|-----------|
| `INTEGER PRIMARY KEY AUTOINCREMENT` | `SERIAL PRIMARY KEY` | High (every table) |
| `DATETIME('now')` | `NOW()` | High (timestamps) |
| `DATETIME('now', '+1 day')` | `NOW() + INTERVAL '1 day'` | Medium |
| `IFNULL(col, default)` | `COALESCE(col, default)` | Low |
| `` `column_name` `` (backticks) | `"column_name"` (double quotes) | Low |
| `TEXT` | `TEXT` | Same ✅ |
| `INTEGER` | `INTEGER` or `BIGINT` | Same ✅ |
| `REAL` | `DOUBLE PRECISION` | Low |

### 1.3 Audit Script

**File:** `scripts/audit-migrations.sh`

```bash
#!/bin/bash

echo "=== Migration Audit Report ==="
echo "Date: $(date)"
echo ""

MIGRATIONS_DIR="migrations"
TOTAL_FILES=$(ls -1 $MIGRATIONS_DIR/*.sql 2>/dev/null | wc -l)

echo "Total migration files: $TOTAL_FILES"
echo ""

echo "--- SQLite-Specific Syntax Counts ---"
echo "AUTOINCREMENT:     $(grep -c "AUTOINCREMENT" $MIGRATIONS_DIR/*.sql 2>/dev/null || echo 0)"
echo "DATETIME('now'):   $(grep -c "DATETIME('now')" $MIGRATIONS_DIR/*.sql 2>/dev/null || echo 0)"
echo "IFNULL:            $(grep -c "IFNULL" $MIGRATIONS_DIR/*.sql 2>/dev/null || echo 0)"
echo "Backticks:         $(grep -c '`' $MIGRATIONS_DIR/*.sql 2>/dev/null || echo 0)"
echo ""

echo "--- Files Requiring Manual Review ---"
grep -l "DATETIME.*+" $MIGRATIONS_DIR/*.sql 2>/dev/null
grep -l "strftime" $MIGRATIONS_DIR/*.sql 2>/dev/null
grep -l "PRAGMA" $MIGRATIONS_DIR/*.sql 2>/dev/null
```

**Expected Output:**
```
=== Migration Audit Report ===
Date: 2026-05-03

Total migration files: 25

--- SQLite-Specific Syntax Counts ---
AUTOINCREMENT:     12
DATETIME('now'):   8
IFNULL:            2
Backticks:         0

--- Files Requiring Manual Review ---
migrations/1650120000_add_users_table.sql
migrations/1660250000_add_expiry_logic.sql
```

---

## 🔄 STEP 2: Conversion Rules

### 2.1 Rule 1: AUTO INCREMENT

**SQLite:**
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT
);
```

**PostgreSQL:**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT
);
```

**Regex Pattern:**
```bash
# Automated conversion
s/INTEGER PRIMARY KEY AUTOINCREMENT/SERIAL PRIMARY KEY/g
```

### 2.2 Rule 2: DATETIME('now')

**SQLite:**
```sql
INSERT INTO logs (created) VALUES (DATETIME('now'));

UPDATE users SET updated = DATETIME('now') WHERE id = 1;

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    created TIMESTAMP DEFAULT DATETIME('now')
);
```

**PostgreSQL:**
```sql
INSERT INTO logs (created) VALUES (NOW());

UPDATE users SET updated = NOW() WHERE id = 1;

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    created TIMESTAMP DEFAULT NOW()
);
```

**Regex Pattern:**
```bash
s/DATETIME\('now'\)/NOW()/g
```

### 2.3 Rule 3: DATE/TIME Arithmetic

**SQLite:**
```sql
-- Add 1 day
SELECT DATETIME('now', '+1 day');

-- Subtract 7 days
SELECT DATETIME('now', '-7 days');

-- Add 1 hour
SELECT DATETIME('now', '+1 hour');
```

**PostgreSQL:**
```sql
-- Add 1 day
SELECT NOW() + INTERVAL '1 day';

-- Subtract 7 days
SELECT NOW() - INTERVAL '7 days';

-- Add 1 hour
SELECT NOW() + INTERVAL '1 hour';
```

**Regex Pattern (Complex - Manual Review Recommended):**
```bash
# Pattern: DATETIME('now', '+N unit') → NOW() + INTERVAL 'N unit'
s/DATETIME\('now', '\+([0-9]+) (day|hour|minute)s?'\)/NOW() + INTERVAL '\1 \2'/g

# Pattern: DATETIME('now', '-N unit') → NOW() - INTERVAL 'N unit'
s/DATETIME\('now', '-([0-9]+) (day|hour|minute)s?'\)/NOW() - INTERVAL '\1 \2'/g
```

### 2.4 Rule 4: IFNULL → COALESCE

**SQLite:**
```sql
SELECT IFNULL(email, 'no-email@example.com') FROM users;
```

**PostgreSQL:**
```sql
SELECT COALESCE(email, 'no-email@example.com') FROM users;
```

**Regex Pattern:**
```bash
s/IFNULL/COALESCE/g
```

### 2.5 Rule 5: Backticks → Double Quotes

**SQLite:**
```sql
SELECT `user_id`, `created_at` FROM `orders`;
```

**PostgreSQL:**
```sql
SELECT "user_id", "created_at" FROM "orders";
```

**Regex Pattern:**
```bash
s/`/"/g
```

---

## 🤖 STEP 3: Automated Conversion Script

### 3.1 Conversion Script Implementation

**File:** `scripts/convert-migrations.sh`

```bash
#!/bin/bash

# PostgreSQL Migration Converter
# Converts SQLite migrations to PostgreSQL syntax

set -e  # Exit on error

MIGRATIONS_DIR="migrations"
OUTPUT_DIR="migrations_postgres"
DRY_RUN=false

# Parse arguments
if [ "$1" == "--dry-run" ]; then
  DRY_RUN=true
  echo "DRY RUN MODE: No files will be modified"
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo "=== Converting SQLite Migrations to PostgreSQL ==="
echo "Source: $MIGRATIONS_DIR"
echo "Output: $OUTPUT_DIR"
echo ""

CONVERTED_COUNT=0
ERROR_COUNT=0

for FILE in $MIGRATIONS_DIR/*.sql; do
  BASENAME=$(basename "$FILE")
  OUTPUT_FILE="$OUTPUT_DIR/$BASENAME"
  
  echo "Converting: $BASENAME"
  
  if [ "$DRY_RUN" == true ]; then
    # Dry run: Show changes, don't write
    cat "$FILE" | \
      sed 's/INTEGER PRIMARY KEY AUTOINCREMENT/SERIAL PRIMARY KEY/g' | \
      sed "s/DATETIME('now')/NOW()/g" | \
      sed 's/IFNULL/COALESCE/g' | \
      sed "s/DATETIME('now', '\+\([0-9]\+\) \(day\|hour\|minute\)s\?')/NOW() + INTERVAL '\1 \2'/g" | \
      sed "s/DATETIME('now', '-\([0-9]\+\) \(day\|hour\|minute\)s\?')/NOW() - INTERVAL '\1 \2'/g" | \
      sed 's/`/"/g'
    
    echo "  ✅ Dry run completed"
  else
    # Real conversion: Write to output directory
    cat "$FILE" | \
      sed 's/INTEGER PRIMARY KEY AUTOINCREMENT/SERIAL PRIMARY KEY/g' | \
      sed "s/DATETIME('now')/NOW()/g" | \
      sed 's/IFNULL/COALESCE/g' | \
      sed "s/DATETIME('now', '\+\([0-9]\+\) \(day\|hour\|minute\)s\?')/NOW() + INTERVAL '\1 \2'/g" | \
      sed "s/DATETIME('now', '-\([0-9]\+\) \(day\|hour\|minute\)s\?')/NOW() - INTERVAL '\1 \2'/g" | \
      sed 's/`/"/g' \
      > "$OUTPUT_FILE"
    
    if [ $? -eq 0 ]; then
      echo "  ✅ Converted successfully"
      CONVERTED_COUNT=$((CONVERTED_COUNT + 1))
    else
      echo "  ❌ Conversion failed"
      ERROR_COUNT=$((ERROR_COUNT + 1))
    fi
  fi
  
  echo ""
done

echo "=== Conversion Complete ==="
echo "Converted: $CONVERTED_COUNT files"
echo "Errors:    $ERROR_COUNT files"
echo ""

if [ "$DRY_RUN" == false ]; then
  echo "Converted migrations saved to: $OUTPUT_DIR"
  echo "Next steps:"
  echo "  1. Review converted migrations manually"
  echo "  2. Test migrations on PostgreSQL test database"
  echo "  3. Replace original migrations if all tests pass"
fi
```

**Usage:**
```bash
# Dry run (preview changes)
./scripts/convert-migrations.sh --dry-run

# Real conversion
./scripts/convert-migrations.sh
```

---

## 🧪 STEP 4: Test Migrations on PostgreSQL

### 4.1 Setup Test Database

**Docker Compose for Testing:**
```yaml
# docker-compose.migration-test.yml
version: '3.8'

services:
  postgres-test:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: migration_test
    ports:
      - "5432:5432"
    volumes:
      - ./migrations_postgres:/migrations
```

**Start test database:**
```bash
docker-compose -f docker-compose.migration-test.yml up -d
```

### 4.2 Apply Migrations

**Manual Application (for testing):**
```bash
# Apply each migration sequentially
for FILE in migrations_postgres/*.sql; do
  echo "Applying: $FILE"
  psql -U postgres -d migration_test -h localhost -f "$FILE"
  
  if [ $? -eq 0 ]; then
    echo "✅ Success"
  else
    echo "❌ Failed"
    exit 1
  fi
done
```

**Expected Output:**
```
Applying: migrations_postgres/1640988000_initial.sql
CREATE TABLE
CREATE TABLE
✅ Success

Applying: migrations_postgres/1650120000_add_users_table.sql
CREATE TABLE
CREATE INDEX
✅ Success

...
```

### 4.3 Migration Testing Script

**File:** `scripts/test-migrations.sh`

```bash
#!/bin/bash

set -e

DB_HOST="localhost"
DB_USER="postgres"
DB_PASS="password"
DB_NAME="migration_test"
MIGRATIONS_DIR="migrations_postgres"

export PGPASSWORD=$DB_PASS

echo "=== Testing PostgreSQL Migrations ==="
echo ""

# Drop and recreate test database
echo "Recreating test database..."
psql -U $DB_USER -h $DB_HOST -c "DROP DATABASE IF EXISTS $DB_NAME"
psql -U $DB_USER -h $DB_HOST -c "CREATE DATABASE $DB_NAME"
echo "✅ Database ready"
echo ""

# Apply migrations
MIGRATION_COUNT=0
for FILE in $MIGRATIONS_DIR/*.sql; do
  BASENAME=$(basename "$FILE")
  echo "Applying: $BASENAME"
  
  psql -U $DB_USER -h $DB_HOST -d $DB_NAME -f "$FILE" > /dev/null 2>&1
  
  if [ $? -eq 0 ]; then
    echo "  ✅ Success"
    MIGRATION_COUNT=$((MIGRATION_COUNT + 1))
  else
    echo "  ❌ Failed"
    psql -U $DB_USER -h $DB_HOST -d $DB_NAME -f "$FILE"  # Show error
    exit 1
  fi
done

echo ""
echo "=== Migration Test Complete ==="
echo "Applied: $MIGRATION_COUNT migrations"
echo ""

# Verify tables created
echo "Verifying database schema..."
TABLE_COUNT=$(psql -U $DB_USER -h $DB_HOST -d $DB_NAME -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'")
echo "Tables created: $TABLE_COUNT"
```

---

## 🔙 STEP 5: Test Migration Rollback

### 5.1 Create Down Migrations

**Convention:**
```
migrations/
├── 1640988000_initial_up.sql     # CREATE TABLE
├── 1640988000_initial_down.sql   # DROP TABLE
├── 1650120000_add_users_up.sql
├── 1650120000_add_users_down.sql
```

**Example Down Migration:**
```sql
-- File: migrations_postgres/1650120000_add_users_down.sql
DROP TABLE IF EXISTS users CASCADE;
DROP INDEX IF EXISTS idx_users_email;
```

### 5.2 Rollback Testing Script

**File:** `scripts/test-rollback.sh`

```bash
#!/bin/bash

set -e

DB_HOST="localhost"
DB_USER="postgres"
DB_PASS="password"
DB_NAME="rollback_test"
MIGRATIONS_DIR="migrations_postgres"

export PGPASSWORD=$DB_PASS

echo "=== Testing Migration Rollback ==="
echo ""

# Create test database
psql -U $DB_USER -h $DB_HOST -c "DROP DATABASE IF EXISTS $DB_NAME"
psql -U $DB_USER -h $DB_HOST -c "CREATE DATABASE $DB_NAME"

# Apply all UP migrations
echo "Applying UP migrations..."
for FILE in $MIGRATIONS_DIR/*_up.sql; do
  psql -U $DB_USER -h $DB_HOST -d $DB_NAME -f "$FILE" > /dev/null
done
echo "✅ All UP migrations applied"

# Count tables before rollback
BEFORE_COUNT=$(psql -U $DB_USER -h $DB_HOST -d $DB_NAME -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'")
echo "Tables before rollback: $BEFORE_COUNT"

# Apply all DOWN migrations (reverse order)
echo ""
echo "Applying DOWN migrations..."
for FILE in $(ls -r $MIGRATIONS_DIR/*_down.sql); do
  psql -U $DB_USER -h $DB_HOST -d $DB_NAME -f "$FILE" > /dev/null
done
echo "✅ All DOWN migrations applied"

# Count tables after rollback
AFTER_COUNT=$(psql -U $DB_USER -h $DB_HOST -d $DB_NAME -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'")
echo "Tables after rollback: $AFTER_COUNT"

if [ "$AFTER_COUNT" -eq 0 ]; then
  echo "✅ Rollback successful (all tables dropped)"
else
  echo "⚠️ Rollback incomplete ($AFTER_COUNT tables remain)"
fi
```

---

## 📁 Deliverables Summary

### Scripts Created (Documentation)
1. ✅ `scripts/audit-migrations.sh` - Identify SQLite-specific syntax
2. ✅ `scripts/convert-migrations.sh` - Automated conversion
3. ✅ `scripts/test-migrations.sh` - Test converted migrations
4. ✅ `scripts/test-rollback.sh` - Test rollback logic

### Migrations Converted (Documentation)
1. `migrations_postgres/` directory with converted migrations
2. Manual review report for complex conversions

---

## ✅ Verification Checklist

**Before marking AGENT-09 complete:**

- [ ] All migrations audited (audit script run)
- [ ] Conversion script created and tested (dry-run mode)
- [ ] All migrations converted (real conversion run)
- [ ] Converted migrations tested on PostgreSQL test database
- [ ] No SQL syntax errors
- [ ] All tables created successfully
- [ ] Rollback migrations created (if not exists)
- [ ] Rollback tested (DOWN migrations work)
- [ ] Migration conversion report documented

---

## 🔗 Dependencies & Handoffs

### Upstream Dependencies
- ✅ **AGENT-04** - SQL dialect matrix complete
- ✅ **AGENT-08** - dbx vendoring complete

### Downstream Dependencies
- 🟡 **AGENT-11** - Integration testing (needs converted migrations)

---

**AGENT-09 Status:** 🟢 Ready for Review  
**Document Version:** 1.0  
**Last Updated:** 2026-05-03
