# SQL Dialect Conversion Matrix

**Agent**: AGENT-04  
**Date**: 2026-05-03

## SQLite → PostgreSQL Conversion

| SQLite | PostgreSQL | Notes |
|--------|-----------|-------|
| `AUTOINCREMENT` | `SERIAL` or `BIGSERIAL` | Primary key auto-increment |
| `INTEGER PRIMARY KEY` | `SERIAL PRIMARY KEY` | Auto-incrementing ID |
| `DATETIME('now')` | `NOW()` or `CURRENT_TIMESTAMP` | Current timestamp |
| `DATE('now')` | `CURRENT_DATE` | Current date only |
| `?` | `$1`, `$2`, `$3` | Placeholder conversion |
| `` `backticks` `` | `"double quotes"` | Identifier quoting |
| `||` | `||` | String concatenation (same) |
| `SUBSTR()` | `SUBSTRING()` | String functions |
| `IFNULL()` | `COALESCE()` | NULL handling |

## Placeholder Conversion

**SQLite**:
```sql
SELECT * FROM users WHERE email = ? AND age > ?
```

**PostgreSQL**:
```sql
SELECT * FROM users WHERE email = $1 AND age > $2
```

**Handled by**: dbx query builder automatically

## Date Functions

**SQLite**:
```sql
INSERT INTO logs (created) VALUES (DATETIME('now'))
```

**PostgreSQL**:
```sql
INSERT INTO logs (created) VALUES (NOW())
```

## Auto-Increment

**SQLite**:
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT
)
```

**PostgreSQL**:
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT
)
```

## dbx Handles Automatically

- Placeholder conversion (? → $1, $2)
- Quote style conversion
- Basic function mapping

## Manual Migration Required

- AUTOINCREMENT → SERIAL (in migration files)
- DATETIME('now') → NOW() (in migration files)
- Complex date math functions
