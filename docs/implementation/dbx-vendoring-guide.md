# dbx Package Vendoring Guide (AGENT-08)

**Agent**: AGENT-08 - dbx Package Vendor Specialist  
**Date**: 2026-05-03  
**Project**: PocketBase v0.37.5 → PostgreSQL + Redis Migration

---

## 📋 Overview

Hướng dẫn vendor custom dbx fork từ postgrebase repository. Đây là **RESEARCH PROJECT** - tạo documentation, KHÔNG phải actual implementation.

### Success Criteria
- ✅ dbx vendoring strategy documented
- ✅ Import path rewrite steps defined
- ✅ Build configuration updated
- ✅ SQL placeholder conversion verified

---

## 🎯 Phase 3 - AGENT-08 Tasks

### Dependencies
- ✅ AGENT-04 (dbx Analysis) - Complete (dbx location identified)
- 🟡 AGENT-05 (DB Connection) - Needs connectDB() implementation

### Estimated Time
8 hours (research documentation)

---

## 🔍 STEP 1: Analyze dbx Fork

### 1.1 Original vs Forked dbx

| Aspect | Original (go-ozzo/ozzo-dbx) | Forked (postgresqlbaseapi/dbx) |
|--------|----------------------------|--------------------------------|
| **Import Path** | `github.com/go-ozzo/ozzo-dbx` | `github.com/free/postgresqlbaseapi/dbx` |
| **Placeholder** | `?` (SQLite/MySQL style) | `$1, $2, $3` (PostgreSQL style) |
| **Drivers** | SQLite, MySQL, MSSQL | PostgreSQL, MySQL |
| **Custom Changes** | None | Placeholder conversion, query builder mods |

### 1.2 Key Modifications in Fork

**File: `postgresqlbaseapi/dbx/builder.go`**
```go
// Modified placeholder conversion
func (q *Query) Build() (string, []interface{}) {
	sql := q.sql
	
	// Original: Uses ? for all databases
	// Modified: Converts ? to $1, $2, $3 for PostgreSQL
	if q.db.DriverName() == "postgres" {
		sql = convertPlaceholders(sql)
	}
	
	return sql, q.params
}

func convertPlaceholders(sql string) string {
	// ? → $1, $2, $3, ...
	count := 0
	return placeholderRegex.ReplaceAllStringFunc(sql, func(match string) string {
		count++
		return fmt.Sprintf("$%d", count)
	})
}
```

**Example Conversion:**
```sql
-- Original (SQLite style)
SELECT * FROM users WHERE email = ? AND age > ?

-- Converted (PostgreSQL style)
SELECT * FROM users WHERE email = $1 AND age > $2
```

---

## 📦 STEP 2: Vendoring Strategy (2 Options)

### Option 1: go.mod replace directive (RECOMMENDED)

**Advantages:**
- ✅ Simple (1 line change)
- ✅ Maintains original import paths
- ✅ Easy to update

**Disadvantages:**
- ⚠️ Requires fork hosted on GitHub/GitLab

**Implementation:**

**1. Add replace directive to `go.mod`:**
```go
module github.com/pocketbase/pocketbase

go 1.25.0

require (
	github.com/go-ozzo/ozzo-dbx v1.5.0  // Keep original import
	// ... other dependencies ...
)

// =========================================================================
// AGENT-08: Replace original dbx with custom fork
// =========================================================================
replace github.com/go-ozzo/ozzo-dbx => github.com/free/postgresqlbaseapi/dbx v1.0.0
```

**2. Update imports (NO CHANGE NEEDED):**
```go
// Imports remain unchanged
import "github.com/go-ozzo/ozzo-dbx"
```

**3. Run go mod tidy:**
```bash
go mod tidy
go mod download
```

---

### Option 2: Manual vendor directory (ALTERNATIVE)

**Advantages:**
- ✅ Works offline (no external dependencies)
- ✅ Full control over vendored code

**Disadvantages:**
- ⚠️ More complex setup
- ⚠️ Harder to update
- ⚠️ Requires import path rewrite

**Implementation:**

**1. Copy dbx fork to vendor directory:**
```bash
# Clone postgrebase repo
git clone https://github.com/free/postgresqlbaseapi /tmp/postgrebase

# Copy dbx package to vendor
mkdir -p vendor/github.com/free/postgresqlbaseapi/
cp -r /tmp/postgrebase/dbx vendor/github.com/free/postgresqlbaseapi/

# Verify vendor structure
tree vendor/github.com/free/postgresqlbaseapi/dbx
```

**Expected structure:**
```
vendor/
└── github.com/
    └── free/
        └── postgresqlbaseapi/
            └── dbx/
                ├── builder.go
                ├── query.go
                ├── db.go
                └── ... (all dbx files)
```

**2. Rewrite ALL import paths:**
```bash
# Find all files importing dbx
grep -r "github.com/go-ozzo/ozzo-dbx" . --include="*.go"

# Replace with vendored path
find . -name "*.go" -exec sed -i '' 's|github.com/go-ozzo/ozzo-dbx|github.com/free/postgresqlbaseapi/dbx|g' {} +
```

**3. Update go.mod:**
```go
require (
	// Remove original dbx
	// github.com/go-ozzo/ozzo-dbx v1.5.0  // REMOVE THIS
	
	// Add vendored dbx path
	github.com/free/postgresqlbaseapi/dbx v1.0.0
)
```

**4. Run go mod vendor:**
```bash
go mod vendor
go mod tidy
```

---

## 🔧 STEP 3: Import Path Rewrite (for Option 2)

### 3.1 Files Requiring Changes

**Expected files with dbx imports:**
```
core/db_postgresql.go
core/base.go
daos/*.go
models/*.go
```

### 3.2 Automated Rewrite Script

**File:** `scripts/rewrite-dbx-imports.sh`

```bash
#!/bin/bash

echo "Rewriting dbx import paths..."

OLD_IMPORT="github.com/go-ozzo/ozzo-dbx"
NEW_IMPORT="github.com/free/postgresqlbaseapi/dbx"

# Find all Go files with old import
FILES=$(grep -r "$OLD_IMPORT" . --include="*.go" -l)

if [ -z "$FILES" ]; then
  echo "No files found with old import path"
  exit 0
fi

# Replace imports
for FILE in $FILES; do
  echo "Rewriting: $FILE"
  sed -i '' "s|$OLD_IMPORT|$NEW_IMPORT|g" "$FILE"
done

echo "✅ Import rewrite complete"
echo "Files modified: $(echo "$FILES" | wc -l)"
```

**Run script:**
```bash
chmod +x scripts/rewrite-dbx-imports.sh
./scripts/rewrite-dbx-imports.sh
```

### 3.3 Verify Import Changes

**Command:**
```bash
# Verify old imports are gone
grep -r "github.com/go-ozzo/ozzo-dbx" . --include="*.go"
# Expected: (no output)

# Verify new imports exist
grep -r "github.com/free/postgresqlbaseapi/dbx" . --include="*.go"
# Expected: List of files using new import
```

---

## 🧪 STEP 4: Test SQL Generation

### 4.1 Unit Test: Placeholder Conversion

**File:** `core/dbx_test.go`

```go
package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/free/postgresqlbaseapi/dbx"
)

func TestPlaceholderConversion_PostgreSQL(t *testing.T) {
	// Setup: Connect to test PostgreSQL
	db, err := dbx.Open("postgres", "postgres://localhost/testdb?sslmode=disable")
	assert.NoError(t, err)
	defer db.Close()
	
	// Build query with ? placeholders
	query := db.NewQuery("SELECT * FROM users WHERE email = ? AND age > ?")
	query.Bind("test@example.com", 18)
	
	// Get generated SQL
	sql, params := query.Build()
	
	// Verify placeholders converted to $1, $2
	assert.Contains(t, sql, "$1")
	assert.Contains(t, sql, "$2")
	assert.NotContains(t, sql, "?")
	
	// Verify params bound correctly
	assert.Equal(t, 2, len(params))
	assert.Equal(t, "test@example.com", params[0])
	assert.Equal(t, 18, params[1])
}

func TestPlaceholderConversion_MySQL(t *testing.T) {
	// Setup: Connect to test MySQL
	db, err := dbx.Open("mysql", "root:pass@tcp(localhost:3306)/testdb")
	assert.NoError(t, err)
	defer db.Close()
	
	// Build query with ? placeholders
	query := db.NewQuery("SELECT * FROM users WHERE email = ? AND age > ?")
	query.Bind("test@example.com", 18)
	
	// Get generated SQL
	sql, params := query.Build()
	
	// Verify placeholders remain as ?
	assert.Contains(t, sql, "?")
	assert.NotContains(t, sql, "$1")
	
	// Verify params bound correctly
	assert.Equal(t, 2, len(params))
}
```

### 4.2 Integration Test: Query Execution

**File:** `tests/integration/dbx_queries_test.go`

```go
package integration

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDBX_CRUD_PostgreSQL(t *testing.T) {
	// Setup test database
	app := NewTestApp(BaseAppConfig{
		DataDsn: "postgres://postgres:password@localhost:5432/testdb?sslmode=disable",
	})
	defer app.Cleanup()
	
	// Test INSERT
	_, err := app.db.Insert("users", dbx.Params{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
	}).Execute()
	assert.NoError(t, err)
	
	// Test SELECT with WHERE
	var user struct {
		ID    int    `db:"id"`
		Name  string `db:"name"`
		Email string `db:"email"`
		Age   int    `db:"age"`
	}
	
	err = app.db.Select("*").
		From("users").
		Where(dbx.HashExp{"email": "john@example.com"}).
		One(&user)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, 30, user.Age)
	
	// Test UPDATE
	_, err = app.db.Update("users", dbx.Params{
		"age": 31,
	}, dbx.HashExp{"email": "john@example.com"}).Execute()
	assert.NoError(t, err)
	
	// Test DELETE
	_, err = app.db.Delete("users", dbx.HashExp{"email": "john@example.com"}).Execute()
	assert.NoError(t, err)
}
```

---

## 🐛 STEP 5: Common Issues & Fixes

### Issue 1: Import Cycle Detected

**Error:**
```
package github.com/free/postgresqlbaseapi/dbx
  imports github.com/pocketbase/pocketbase/core
  imports github.com/free/postgresqlbaseapi/dbx
  import cycle not allowed
```

**Root Cause:** Circular dependency between dbx and core package

**Solution:** Ensure dbx package is standalone (no imports from app packages)

```go
// BAD: dbx importing app code
package dbx
import "github.com/pocketbase/pocketbase/core" // ❌ Import cycle

// GOOD: dbx is standalone
package dbx
import "database/sql" // ✅ Only standard library
```

### Issue 2: Placeholder Conversion Not Working

**Symptom:**
```
pq: syntax error at or near "?"
```

**Root Cause:** dbx fork not using PostgreSQL placeholder conversion

**Verify Fix:**
```go
// Check db.DriverName() is correctly detected
log.Printf("Driver: %s", db.DriverName())
// Expected output: Driver: postgres

// Check placeholder conversion function exists
func convertPlaceholders(sql string) string {
	// Implementation should exist in builder.go
}
```

### Issue 3: Vendor Directory Not Recognized

**Error:**
```
cannot find package "github.com/free/postgresqlbaseapi/dbx" in any of:
  ...
```

**Solution 1: Enable vendoring**
```bash
go mod vendor
export GOFLAGS=-mod=vendor
go build
```

**Solution 2: Use replace directive instead (Option 1)**

---

## 📊 STEP 6: Performance Verification

### 6.1 Benchmark Placeholder Conversion

**File:** `benchmarks/dbx_bench_test.go`

```go
package benchmarks

import (
	"testing"
	"github.com/free/postgresqlbaseapi/dbx"
)

func BenchmarkPlaceholderConversion(b *testing.B) {
	db, _ := dbx.Open("postgres", "postgres://localhost/bench?sslmode=disable")
	defer db.Close()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := db.NewQuery("SELECT * FROM users WHERE email = ? AND age > ? AND status = ?")
		query.Bind("test@example.com", 18, "active")
		_, _ = query.Build()
	}
}

// Expected: < 1000 ns/op (very fast)
```

**Run benchmark:**
```bash
go test -bench=BenchmarkPlaceholderConversion -benchmem
```

**Expected Output:**
```
BenchmarkPlaceholderConversion-8   2000000   800 ns/op   128 B/op   3 allocs/op
```

---

## 📁 Deliverables Summary

### Files Modified (Documentation)
1. ✅ `go.mod` - Add replace directive (Option 1) OR vendor dbx (Option 2)
2. ✅ `core/db_postgresql.go` - Import path updated (if Option 2)
3. ✅ `core/base.go` - Import path updated (if Option 2)
4. ✅ All `*.go` files with dbx imports - Updated (if Option 2)

### Testing Files (Documentation)
1. `core/dbx_test.go` - Placeholder conversion tests
2. `tests/integration/dbx_queries_test.go` - CRUD tests
3. `benchmarks/dbx_bench_test.go` - Performance tests

### Scripts (Documentation)
1. `scripts/rewrite-dbx-imports.sh` - Automated import rewrite

---

## ✅ Verification Checklist

**Before marking AGENT-08 complete:**

- [ ] Vendoring strategy chosen (Option 1 or 2)
- [ ] `go.mod` updated with replace directive (Option 1) OR vendor directory created (Option 2)
- [ ] Import paths rewritten (if Option 2)
- [ ] Placeholder conversion tested (PostgreSQL: `?` → `$1, $2, $3`)
- [ ] CRUD queries tested (INSERT, SELECT, UPDATE, DELETE)
- [ ] Build successful (`go build ./...`)
- [ ] No import cycles detected
- [ ] Performance benchmarks acceptable

---

## 🔗 Dependencies & Handoffs

### Upstream Dependencies
- ✅ **AGENT-04** - dbx location analysis complete
- ✅ **AGENT-05** - `connectDB()` implementation needs dbx

### Downstream Dependencies
- 🟡 **AGENT-09** - Migration conversion (needs dbx placeholder conversion)
- 🟡 **AGENT-10** - Unit testing (needs dbx vendored)

---

**AGENT-08 Status:** 🟢 Ready for Review  
**Document Version:** 1.0  
**Last Updated:** 2026-05-03
