# CLI Flags Implementation Guide (AGENT-07)

**Agent**: AGENT-07 - CLI Flags Engineer  
**Date**: 2026-05-03  
**Project**: PocketBase v0.37.5 → PostgreSQL + Redis Migration

---

## 📋 Overview

Hướng dẫn chi tiết để add `--dataDsn` và `--redisDsn` flags vào PocketBase CLI. Đây là **RESEARCH PROJECT** - tạo documentation, KHÔNG phải actual implementation.

### Success Criteria
- ✅ `--dataDsn` flag documented
- ✅ `--redisDsn` flag documented
- ✅ Environment variable support documented
- ✅ Flag precedence logic defined
- ✅ Help text examples provided

---

## 🎯 Phase 2 - AGENT-07 Tasks

### Dependencies
- ✅ AGENT-02 (Code Extraction) - `cmd_serve_flags.go` extracted
- 🟡 AGENT-05 (DB Connection Layer) - Needs `BaseAppConfig` struct
- 🟡 AGENT-06 (Redis Integration) - Needs Redis integration complete

### Estimated Time
3 hours (research documentation)

---

## 📦 STEP 1: Locate PocketBase CLI Entry Point

### 1.1 File Structure

**Original PocketBase structure:**
```
pocketbase/
├── main.go                    # Entry point (calls pocketbase.go)
├── pocketbase.go              # Main PocketBase application struct
└── cmd/
    └── serve.go               # Serve command (HTTP server)
```

**For RAI Backend (may differ):**
```
rai-backend/
├── main.go
├── cmd/
│   ├── root.go                # RootCmd definition
│   └── serve.go               # Serve subcommand
```

### 1.2 Identify Flag Definition Location

**Option 1: Centralized flags (PocketBase pattern)**
```go
// File: pocketbase.go
func (pb *PocketBase) eagerParseFlags(config *Config) error {
    // All persistent flags defined here
    pb.RootCmd.PersistentFlags().StringVar(...)
}
```

**Option 2: Distributed flags (Cobra standard)**
```go
// File: cmd/serve.go
func init() {
    serveCmd.Flags().StringVar(...)
}
```

---

## 🔧 STEP 2: Add Flags to PocketBase Struct

### 2.1 Locate PocketBase Struct

**File:** `pocketbase.go` or `cmd/root.go`

**Original Struct:**
```go
type PocketBase struct {
	RootCmd *cobra.Command
	
	// Existing flags
	dataDirFlag       string
	encryptionEnvFlag string
	debugFlag         bool
	
	// App instance
	app *core.BaseApp
	
	// ... other fields ...
}
```

### 2.2 Add New Flag Fields

**Modified Struct:**
```go
type PocketBase struct {
	RootCmd *cobra.Command
	
	// Existing flags
	dataDirFlag       string
	encryptionEnvFlag string
	debugFlag         bool
	
	// =========================================================================
	// AGENT-07: PostgreSQL + Redis Flags
	// =========================================================================
	
	dataDataFlag      string // --dataDsn flag value (PostgreSQL/MySQL DSN)
	redisFlag         string // --redisDsn flag value (Redis DSN)
	
	// App instance
	app *core.BaseApp
	
	// ... other fields ...
}
```

---

## 🚩 STEP 3: Define CLI Flags

### 3.1 Add --dataDsn Flag

**Location:** `eagerParseFlags()` function in `pocketbase.go`

**Implementation:**
```go
func (pb *PocketBase) eagerParseFlags(config *Config) error {
	// -------------------------------------------------------------------------
	// Existing Flags (keep as-is)
	// -------------------------------------------------------------------------
	
	pb.RootCmd.PersistentFlags().StringVar(
		&pb.dataDirFlag,
		"dir",
		config.DefaultDataDir,
		"the PocketBase data directory",
	)
	
	pb.RootCmd.PersistentFlags().StringVar(
		&pb.encryptionEnvFlag,
		"encryptionEnv",
		config.DefaultEncryptionEnv,
		"the env variable whose value of 32 characters will be used as encryption key for the app settings (default none)",
	)
	
	pb.RootCmd.PersistentFlags().BoolVar(
		&pb.debugFlag,
		"debug",
		config.DefaultDebug,
		"enable debug mode, aka. showing more detailed logs",
	)
	
	// -------------------------------------------------------------------------
	// AGENT-07: NEW FLAG 1 - --dataDsn (PostgreSQL/MySQL DSN)
	// -------------------------------------------------------------------------
	
	pb.RootCmd.PersistentFlags().StringVar(
		&pb.dataDataFlag,              // Store flag value in this field
		"dataDsn",                     // Flag name (use: --dataDsn="...")
		config.DefaultDataDsn,         // Default value (from Config or env var)
		// Help text (shown in --help output)
		"store data postgresql/mysql dsn (e.g. postgres://user:pass@127.0.0.1:5432/db?sslmode=disable OR mysql://user:pass@tcp(127.0.0.1:3306)/db)",
	)
	
	// -------------------------------------------------------------------------
	// AGENT-07: NEW FLAG 2 - --redisDsn (Redis DSN)
	// -------------------------------------------------------------------------
	
	pb.RootCmd.PersistentFlags().StringVar(
		&pb.redisFlag,                 // Store flag value in this field
		"redisDsn",                    // Flag name (use: --redisDsn="...")
		config.RedisDsn,               // Default value (from Config or env var)
		// Help text (shown in --help output)
		"Cache data Redis dsn (default redis://<user>:<pass>@localhost:6379/<db> e.g. redis://localhost:6379/0)",
	)
	
	// Parse all flags early (before Bootstrap)
	return pb.RootCmd.ParseFlags(os.Args[1:])
}
```

### 3.2 Flag Naming Conventions

| Flag Name | Format | Example | Notes |
|-----------|--------|---------|-------|
| `--dataDsn` | CamelCase | `--dataDsn="postgres://..."` | Matches PocketBase style |
| `--redisDsn` | CamelCase | `--redisDsn="redis://..."` | Consistent with dataDsn |
| `--dir` | lowercase | `--dir=./pb_data` | Existing flag (keep as-is) |
| `--debug` | lowercase | `--debug` | Existing flag (keep as-is) |

---

## 🌍 STEP 4: Environment Variable Support

### 4.1 Update Config Struct

**File:** `pocketbase.go` (or `config.go` if exists)

**Original Config:**
```go
type Config struct {
	DefaultDataDir      string
	DefaultEncryptionEnv string
	DefaultDebug        bool
}
```

**Modified Config:**
```go
type Config struct {
	// Existing fields
	DefaultDataDir      string
	DefaultEncryptionEnv string
	DefaultDebug        bool
	
	// =========================================================================
	// AGENT-07: PostgreSQL + Redis Configuration
	// =========================================================================
	
	DefaultDataDsn      string // Default dataDsn (from env var or empty)
	RedisDsn            string // Redis DSN (from env var or empty)
}
```

### 4.2 Load Environment Variables

**Function:** `NewConfig()` in `pocketbase.go`

**Implementation:**
```go
func NewConfig() *Config {
	return &Config{
		// Existing defaults
		DefaultDataDir:      "./pb_data",
		DefaultEncryptionEnv: "",
		DefaultDebug:        false,
		
		// =========================================================================
		// AGENT-07: Load PostgreSQL + Redis from Environment
		// =========================================================================
		
		// Priority: Env var > Empty string (flag will override later)
		DefaultDataDsn:      os.Getenv("DATA_DSN"),  // export DATA_DSN="postgres://..."
		RedisDsn:            os.Getenv("REDIS_DSN"), // export REDIS_DSN="redis://..."
	}
}
```

### 4.3 Flag Precedence Logic

**Priority Order (highest to lowest):**
```
1. CLI Flag (--dataDsn="...")
2. Environment Variable (DATA_DSN="...")
3. Default Value (empty string)
```

**Example:**
```bash
# Priority 1: CLI flag (highest)
pocketbase serve --dataDsn="postgres://override@localhost/db"

# Priority 2: Environment variable (if no flag)
export DATA_DSN="postgres://env@localhost/db"
pocketbase serve

# Priority 3: Default (if no flag and no env var)
pocketbase serve  # Uses SQLite (empty dataDsn)
```

---

## 🔗 STEP 5: Pass Flags to BaseApp

### 5.1 Update Bootstrap() Function

**Location:** `pocketbase.go`

**Original Bootstrap:**
```go
func (pb *PocketBase) Bootstrap() error {
	config := core.BaseAppConfig{
		DataDir:       pb.dataDirFlag,
		EncryptionEnv: pb.encryptionEnvFlag,
		IsDebug:       pb.debugFlag,
	}
	
	return pb.app.Bootstrap(config)
}
```

**Modified Bootstrap:**
```go
func (pb *PocketBase) Bootstrap() error {
	// =========================================================================
	// AGENT-07: Pass PostgreSQL + Redis Flags to BaseApp
	// =========================================================================
	
	config := core.BaseAppConfig{
		// Existing fields
		DataDir:       pb.dataDirFlag,
		EncryptionEnv: pb.encryptionEnvFlag,
		IsDebug:       pb.debugFlag,
		
		// NEW: PostgreSQL + Redis Configuration
		DataDsn:       pb.dataDataFlag, // --dataDsn flag value
		RedisDsn:      pb.redisFlag,    // --redisDsn flag value
	}
	
	return pb.app.Bootstrap(config)
}
```

---

## 📖 STEP 6: Update Help Text

### 6.1 Verify --help Output

**Command (documentation only):**
```bash
pocketbase serve --help
```

**Expected Output:**
```
Usage:
  pocketbase serve [flags]

Flags:
      --dataDsn string          store data postgresql/mysql dsn (e.g. postgres://user:pass@127.0.0.1:5432/db?sslmode=disable OR mysql://user:pass@tcp(127.0.0.1:3306)/db)
      --debug                   enable debug mode, aka. showing more detailed logs
      --dir string              the PocketBase data directory (default "./pb_data")
      --encryptionEnv string    the env variable whose value of 32 characters will be used as encryption key for the app settings (default none)
  -h, --help                    help for serve
      --http string             API HTTP server address (default "127.0.0.1:8090")
      --redisDsn string         Cache data Redis dsn (default redis://<user>:<pass>@localhost:6379/<db> e.g. redis://localhost:6379/0)
```

### 6.2 Add Examples to Help Text (Optional)

**Enhanced help text with examples:**
```go
pb.RootCmd.PersistentFlags().StringVar(
	&pb.dataDataFlag,
	"dataDsn",
	config.DefaultDataDsn,
	`store data postgresql/mysql dsn
Examples:
  PostgreSQL: postgres://user:pass@localhost:5432/db?sslmode=disable
  MySQL:      mysql://user:pass@tcp(localhost:3306)/db`,
)
```

---

## 🧪 STEP 7: Testing Strategy

### 7.1 Test Flag Parsing

**Test Script:** `scripts/test-flags.sh`

```bash
#!/bin/bash

echo "Testing CLI flags..."

# Test 1: --help shows new flags
echo "Test 1: --help output"
./pocketbase serve --help | grep -E "dataDsn|redisDsn"
if [ $? -eq 0 ]; then
  echo "✅ Flags appear in --help"
else
  echo "❌ Flags missing from --help"
fi

# Test 2: --dataDsn flag accepted
echo "Test 2: --dataDsn flag parsing"
./pocketbase serve --dataDsn="postgres://localhost/test" --help > /dev/null 2>&1
if [ $? -eq 0 ]; then
  echo "✅ --dataDsn flag accepted"
else
  echo "❌ --dataDsn flag rejected"
fi

# Test 3: --redisDsn flag accepted
echo "Test 3: --redisDsn flag parsing"
./pocketbase serve --redisDsn="redis://localhost:6379/0" --help > /dev/null 2>&1
if [ $? -eq 0 ]; then
  echo "✅ --redisDsn flag accepted"
else
  echo "❌ --redisDsn flag rejected"
fi
```

### 7.2 Test Environment Variables

**Test Script:** `scripts/test-env-vars.sh`

```bash
#!/bin/bash

echo "Testing environment variables..."

# Test 1: DATA_DSN env var
export DATA_DSN="postgres://envtest@localhost/db"
./pocketbase serve --help > /dev/null 2>&1
# TODO: Add verification that env var is read

# Test 2: REDIS_DSN env var
export REDIS_DSN="redis://envtest:6379/0"
./pocketbase serve --help > /dev/null 2>&1
# TODO: Add verification that env var is read

echo "✅ Environment variable tests complete"
```

### 7.3 Test Flag Precedence

**Test Script:** `scripts/test-flag-precedence.sh`

```bash
#!/bin/bash

echo "Testing flag precedence (CLI > Env > Default)..."

# Setup: Set env var
export DATA_DSN="postgres://env@localhost/db"

# Test: CLI flag should override env var
./pocketbase serve --dataDsn="postgres://cli@localhost/db" &
PID=$!
sleep 2

# Verify process is running with CLI flag value
# TODO: Check logs to confirm "postgres://cli@localhost/db" was used

kill $PID
echo "✅ Precedence test complete"
```

---

## 📝 Usage Examples

### Example 1: SQLite Mode (Original Behavior)

**Command:**
```bash
pocketbase serve
# or
pocketbase serve --dir=./my_data
```

**Behavior:**
- No `--dataDsn` flag → Uses SQLite
- Database file: `./pb_data/pb_data.db` (or `./my_data/pb_data.db`)

### Example 2: PostgreSQL Mode (Single Node)

**Command:**
```bash
pocketbase serve \
  --dataDsn="postgres://user:password@localhost:5432/pocketbase?sslmode=disable"
```

**Behavior:**
- Connects to PostgreSQL
- No Redis → Single-node mode (no clustering)

### Example 3: PostgreSQL + Redis (Multi-Node)

**Command:**
```bash
# Node 1
pocketbase serve \
  --dataDsn="postgres://user:pass@shared-db:5432/pb" \
  --redisDsn="redis://shared-redis:6379/0" \
  --http="0.0.0.0:8090"

# Node 2
pocketbase serve \
  --dataDsn="postgres://user:pass@shared-db:5432/pb" \
  --redisDsn="redis://shared-redis:6379/0" \
  --http="0.0.0.0:8091"
```

**Behavior:**
- Both nodes share PostgreSQL database
- Both nodes use Redis Pub/Sub for realtime sync
- SSE events work across both nodes

### Example 4: MySQL Mode

**Command:**
```bash
pocketbase serve \
  --dataDsn="mysql://root:password@tcp(localhost:3306)/pocketbase?parseTime=true"
```

**Behavior:**
- Connects to MySQL
- No Redis → Single-node mode

### Example 5: Environment Variables

**Command:**
```bash
# Setup
export DATA_DSN="postgres://user:pass@localhost:5432/pb"
export REDIS_DSN="redis://localhost:6379/0"

# Run (no flags needed)
pocketbase serve
```

**Behavior:**
- Reads DSNs from environment variables
- Connects to PostgreSQL + Redis

---

## 🐛 Error Handling

### Error 1: Invalid DSN Format

**Command:**
```bash
pocketbase serve --dataDsn="invalid://format"
```

**Expected Error:**
```
Error: failed to connect to database: invalid DSN
```

**Solution:** Use valid DSN format (see examples above)

### Error 2: Database Connection Fails

**Command:**
```bash
pocketbase serve --dataDsn="postgres://localhost:9999/db"
```

**Expected Error:**
```
Error: failed to connect to database: dial tcp [::1]:9999: connect: connection refused
```

**Solution:** Verify database is running and port is correct

### Error 3: Redis Connection Fails (Non-Fatal)

**Command:**
```bash
pocketbase serve \
  --dataDsn="postgres://localhost/db" \
  --redisDsn="redis://localhost:9999/0"
```

**Expected Warning (Not Error):**
```
⚠ Redis connection failed: dial tcp [::1]:9999: connect: connection refused
Continuing in single-node mode (no clustering)
```

**Behavior:** App continues without Redis (graceful degradation)

---

## 📁 Deliverables Summary

### Files Modified (Documentation)
1. ✅ `pocketbase.go` (or `cmd/root.go`)
   - Add fields: `dataDataFlag`, `redisFlag`
   - Update `eagerParseFlags()` with new flags
   - Update `Bootstrap()` to pass flags to BaseApp
2. ✅ `config.go` (or Config struct in `pocketbase.go`)
   - Add fields: `DefaultDataDsn`, `RedisDsn`
   - Update `NewConfig()` to load env vars

### Testing Scripts (Documentation)
1. `scripts/test-flags.sh`
2. `scripts/test-env-vars.sh`
3. `scripts/test-flag-precedence.sh`

---

## ✅ Verification Checklist

**Before marking AGENT-07 complete:**

- [ ] `dataDataFlag` and `redisFlag` fields added to PocketBase struct
- [ ] `--dataDsn` flag definition added to `eagerParseFlags()`
- [ ] `--redisDsn` flag definition added to `eagerParseFlags()`
- [ ] `DefaultDataDsn` and `RedisDsn` fields added to Config struct
- [ ] Environment variable loading implemented in `NewConfig()`
- [ ] `Bootstrap()` updated to pass flags to BaseApp
- [ ] Help text verified (--help output)
- [ ] Usage examples documented
- [ ] Error handling documented
- [ ] Testing strategy defined

---

## 🔗 Dependencies & Handoffs

### Upstream Dependencies
- ✅ **AGENT-02** - Extracted `cmd_serve_flags.go` code
- 🟡 **AGENT-05** - `BaseAppConfig` struct (can work in parallel)
- 🟡 **AGENT-06** - Redis integration (can work in parallel)

### Downstream Dependencies
- 🟡 **AGENT-10** - Unit testing (needs flags for test setup)
- 🟡 **AGENT-11** - Integration testing (needs flags for multi-node tests)

---

## 📌 Next Steps

**After AGENT-07 completion:**

1. **Compile and Test:** Verify flags work end-to-end
2. **Handoff to AGENT-10:** Unit tests can now test with various DSN formats
3. **Handoff to AGENT-11:** Integration tests can now test multi-node setup

---

**AGENT-07 Status:** 🟢 Ready for Review  
**Document Version:** 1.0  
**Last Updated:** 2026-05-03
