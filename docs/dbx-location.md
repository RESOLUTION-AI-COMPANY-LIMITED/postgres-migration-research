# dbx Package Location Analysis

**Agent**: AGENT-04 - dbx Package Specialist  
**Date**: 2026-05-03

## Package Location

**Path**: `/tmp/postgrebase/dbx/`

**Total Go Files**: 29 files

## Structure

```
dbx/
├── builder.go                 (16.7KB - Core query builder)
├── builder_pgsql.go          (4.0KB - PostgreSQL dialect)
├── builder_mysql.go          (4.2KB - MySQL dialect)  
├── builder_sqlite.go         (4.3KB - SQLite dialect)
├── db.go                     (10.0KB - Database wrapper)
└── [test files]
```

## Import Path

**Postgrebase uses**: `github.com/free/postgresqlbaseapi/dbx`

## Key Modifications

1. PostgreSQL support (`builder_pgsql.go`)
2. MySQL support (`builder_mysql.go`)
3. Placeholder conversion (? → $1, $2, ...)

## Vendoring Strategy

**Recommended**: Vendor entire dbx directory into project
