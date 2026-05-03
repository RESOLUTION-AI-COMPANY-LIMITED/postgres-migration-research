# dbx Vendoring Strategy

**Agent**: AGENT-04  
**Date**: 2026-05-03

## Decision: Full Vendor

**Strategy**: Copy entire dbx package to project vendor directory

## Implementation Steps

1. Create vendor directory structure
2. Copy `/tmp/postgrebase/dbx/` → `vendor/dbx/`
3. Rewrite all imports from `github.com/free/postgresqlbaseapi/dbx` to local vendor path
4. Update go.mod with replace directive

## Import Path Rewrite

```go
// Before
import "github.com/free/postgresqlbaseapi/dbx"

// After  
import "your-project/vendor/dbx"
```

## Build Configuration

Add to go.mod:
```go
replace github.com/free/postgresqlbaseapi/dbx => ./vendor/dbx
```

## Files to Modify

Search and replace in:
- core/db_postgresql.go
- Any files using dbx.DB, dbx.Query, etc.

## Verification

```bash
grep -r "postgresqlbaseapi/dbx" . --include="*.go"
# Should return 0 results after rewrite
```
