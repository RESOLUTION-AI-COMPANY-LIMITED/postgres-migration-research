# Code Mapping: Postgrebase → PocketBase v0.37.5

**Purpose**: Map extracted postgrebase code to PocketBase v0.37.5 integration points  
**Date**: 2026-05-03  
**Agent**: AGENT-02 - Code Extraction Specialist

---

## 📋 Overview

This document maps the 3 core implementation files from postgrebase to their integration points in PocketBase v0.37.5.

**Total Implementation**: ~200 lines of code across 3 files

---

## 🗺️ File Mapping Matrix

| Postgrebase Source | Extracted Doc | Target Location | Action | Lines |
|-------------------|---------------|-----------------|--------|-------|
| `core/db_postgresql.go` | `docs/extracted/db_postgresql.go` | `core/db_postgresql.go` | **CREATE NEW** | 29 |
| `core/base.go` (Redis) | `docs/extracted/redis_functions.go` | `core/base.go` | **MODIFY** | ~100 |
| `pocketbase.go` (Flags) | `docs/extracted/cmd_serve_flags.go` | `pocketbase.go` | **MODIFY** | ~30 |

---

## 📁 File 1: Database Connection Layer

### Source
- **File**: `/tmp/postgrebase/core/db_postgresql.go`
- **Extracted**: `docs/extracted/db_postgresql.go`
- **Lines**: 29

### Integration Point
**Target**: `core/db_postgresql.go` (CREATE NEW FILE)

**Action**: Copy entire file to PocketBase core directory

**Dependencies to Add**:
```go
import (
    _ "github.com/lib/pq"                       // PostgreSQL driver
    _ "github.com/go-sql-driver/mysql"          // MySQL driver
    "github.com/free/postgresqlbaseapi/dbx"     // Query builder
)
```

**Integration Steps**:
1. Create new file `core/db_postgresql.go`
2. Copy `connectDB()` function from extracted file
3. Update imports to use vendored dbx package
4. Add connection pooling configuration (optional)

**Original Code Location in Postgrebase**:
- Lines 1-29 of `/tmp/postgrebase/core/db_postgresql.go`

**Modification in base.go**:
Replace SQLite `connectDB()` call with new `connectDB(dsn)` implementation.

---

## 📁 File 2: Redis Integration

### Source
- **File**: `/tmp/postgrebase/core/base.go` (lines 1061-1125)
- **Extracted**: `docs/extracted/redis_functions.go`
- **Lines**: ~100

### Integration Points in `core/base.go`

#### 2.1: BaseApp Struct Fields (INSERT)
**Location**: `core/base.go`, inside `type BaseApp struct`

**Add these fields**:
```go
type BaseApp struct {
    // ... existing fields ...

    // Redis integration
    redisDsn         string
    redisCache       *redis.Client
    redisContext     context.Context
}
```

**Line Numbers**: Insert after existing database fields (~line 50)

---

#### 2.2: BaseAppConfig Struct (INSERT)
**Location**: `core/base.go`, inside `type BaseAppConfig struct`

**Add field**:
```go
type BaseAppConfig struct {
    // ... existing fields ...

    RedisDsn         string // Redis DSN for clustering
}
```

**Line Numbers**: Insert near DataDir field (~line 185)

---

#### 2.3: NewBaseApp Constructor (MODIFY)
**Location**: `core/base.go`, `NewBaseApp()` function

**Add initialization**:
```go
app := &BaseApp{
    // ... existing initializations ...

    redisDsn:            config.RedisDsn,
    redisContext:        context.Background(),
}
```

**Line Numbers**: ~line 195-205

---

#### 2.4: Bootstrap() Method (INSERT CALL)
**Location**: `core/base.go`, `Bootstrap()` method

**Add after database initialization**:
```go
func (app *BaseApp) Bootstrap() error {
    // ... existing code ...

    // Initialize databases
    if err := app.initDataDB(); err != nil {
        return err
    }

    // NEW: Initialize Redis (after DB init)
    if err := app.initRedis(); err != nil {
        return err
    }

    // ... rest of bootstrap ...
}
```

**Line Numbers**: Insert call around line 356

---

#### 2.5: New Methods (INSERT)
**Location**: `core/base.go`, end of file or near other utility methods

**Add 3 new methods**:
1. `initRedis()` - Lines 1061-1095 from postgrebase
2. `Publish()` - Lines 1097-1115 from postgrebase
3. `RedisCache()` - Lines 477-480 from postgrebase

**Suggested Line Numbers**: End of file, before closing brace

---

## 📁 File 3: CLI Flags

### Source
- **File**: `/tmp/postgrebase/pocketbase.go` (lines 200-225)
- **Extracted**: `docs/extracted/cmd_serve_flags.go`
- **Lines**: ~30

### Integration Points in `pocketbase.go`

#### 3.1: PocketBase Struct Fields (INSERT)
**Location**: `pocketbase.go`, inside `type PocketBase struct`

**Add fields**:
```go
type PocketBase struct {
    // ... existing fields ...

    dataDataFlag      string // --dataDsn flag
    redisFlag         string // --redisDsn flag
}
```

**Line Numbers**: Near other flag fields

---

#### 3.2: Config Struct (INSERT)
**Location**: `pocketbase.go`, inside `type Config struct`

**Add fields**:
```go
type Config struct {
    // ... existing fields ...

    DefaultDataDsn    string // Default PostgreSQL/MySQL DSN
    RedisDsn          string // Redis DSN
}
```

---

#### 3.3: eagerParseFlags() Method (INSERT)
**Location**: `pocketbase.go`, inside `eagerParseFlags()` method

**Add flag definitions**:
```go
func (pb *PocketBase) eagerParseFlags(config *Config) error {
    // ... existing flags ...

    // NEW: --dataDsn flag
    pb.RootCmd.PersistentFlags().StringVar(
        &pb.dataDataFlag,
        "dataDsn",
        config.DefaultDataDsn,
        "store data postgresql/mysql dsn (e.g. postgres://...)",
    )

    // NEW: --redisDsn flag
    pb.RootCmd.PersistentFlags().StringVar(
        &pb.redisFlag,
        "redisDsn",
        config.RedisDsn,
        "Cache data Redis dsn (default redis://localhost:6379/0)",
    )

    // ... rest of flags ...
}
```

**Line Numbers**: Insert around line 210-225

---

#### 3.4: Bootstrap() Call (MODIFY)
**Location**: `pocketbase.go`, where `BaseAppConfig` is created

**Pass flag values**:
```go
config := core.BaseAppConfig{
    DataDir:      pb.dataDirFlag,      // Existing
    DataDsn:      pb.dataDataFlag,     // NEW
    RedisDsn:     pb.redisFlag,        // NEW
    EncryptionEnv: pb.encryptionEnvFlag,
    IsDebug:      pb.debugFlag,
}
```

---

## 🔗 Dependency Chain

### Import Order
1. Create `core/db_postgresql.go` first (no dependencies)
2. Modify `core/base.go` (depends on db_postgresql.go)
3. Modify `pocketbase.go` (depends on core/base.go changes)

### Build Order
```
go.mod (add dependencies)
  ↓
core/db_postgresql.go (create)
  ↓
core/base.go (modify - add Redis)
  ↓
pocketbase.go (modify - add flags)
  ↓
go build
```

---

## 📦 go.mod Changes

Add these dependencies:
```go
require (
    github.com/lib/pq v1.10.9                    // PostgreSQL driver
    github.com/go-sql-driver/mysql v1.7.1         // MySQL driver
    github.com/redis/go-redis/v9 v9.0.5           // Redis client
    github.com/free/postgresqlbaseapi/dbx v1.0.0  // Query builder (vendored)
)
```

**Action**: Run `go mod tidy` after adding imports

---

## ✅ Verification Checklist

After integration:

### Compilation
- [ ] `go build` passes without errors
- [ ] No missing import errors
- [ ] All type checks pass

### Functionality
- [ ] `pocketbase --help` shows --dataDsn and --redisDsn flags
- [ ] `pocketbase serve --dataDsn="postgres://..."` starts successfully
- [ ] Redis gracefully disabled when --redisDsn not provided
- [ ] Database connection works with PostgreSQL
- [ ] Database connection works with MySQL

### Regression
- [ ] SQLite mode still works (when --dataDsn not provided)
- [ ] Existing APIs unchanged
- [ ] No breaking changes to existing behavior

---

## 🎯 Implementation Order (For AGENT-05)

**AGENT-05 will follow this sequence**:

1. **TASK-05-A**: Create `core/db_postgresql.go`
   - Copy from `docs/extracted/db_postgresql.go`
   - Update import paths for vendored dbx

2. **TASK-05-B**: Modify `core/base.go`
   - Add Redis fields to BaseApp struct
   - Add Redis config to BaseAppConfig
   - Add initRedis() call in Bootstrap()
   - Add 3 new methods (initRedis, Publish, RedisCache)

3. **TASK-05-C**: Modify `pocketbase.go`
   - Add flag fields to PocketBase struct
   - Add config fields to Config struct
   - Add flag definitions in eagerParseFlags()
   - Pass flags to Bootstrap()

4. **TASK-05-D**: Update `go.mod`
   - Add lib/pq dependency
   - Add go-sql-driver/mysql dependency
   - Add redis/go-redis/v9 dependency
   - Run go mod tidy

---

## 📝 Notes for Implementation

### Line Number Accuracy
Line numbers provided are approximate. Use search patterns instead:
- Search for struct definitions: `type BaseApp struct`
- Search for method signatures: `func (app *BaseApp) Bootstrap()`
- Insert code in logical groupings (all Redis together)

### Import Paths
After vendoring dbx package (AGENT-08), update import:
```go
// Change this:
"github.com/free/postgresqlbaseapi/dbx"

// To this:
"your-project/vendor/dbx"
```

### Environment Variables
Support environment variable alternatives:
```go
dataDsn := os.Getenv("DB_DSN")
if dataDsn == "" {
    dataDsn = flagValue
}
```

---

**Created by**: AGENT-02 - Code Extraction Specialist  
**Date**: 2026-05-03  
**References**:
- `docs/extracted/db_postgresql.go`
- `docs/extracted/redis_functions.go`
- `docs/extracted/cmd_serve_flags.go`
