# dbx Version Comparison

**Agent**: AGENT-04  
**Date**: 2026-05-03

## Original vs Fork

**Original**: `github.com/go-ozzo/ozzo-dbx`  
**Fork**: `github.com/free/postgresqlbaseapi/dbx`

## Key Differences

### 1. PostgreSQL Support
- Added `builder_pgsql.go` (4.0KB)
- Placeholder conversion: `?` → `$1`, `$2`
- PostgreSQL-specific SQL generation

### 2. MySQL Support  
- Enhanced `builder_mysql.go` (4.2KB)
- MySQL DSN parsing
- Driver detection logic

### 3. SQLite Compatibility
- Maintains original SQLite support
- No breaking changes to existing behavior

## Compatibility

**Backward Compatible**: Yes  
**Breaking Changes**: None

Fork is extension of original, not replacement.

## Recommendation

Use fork (postgrebase version) as it includes all original features plus PostgreSQL/MySQL support.
