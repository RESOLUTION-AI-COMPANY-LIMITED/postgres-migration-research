package main

// =============================================================================
// CLI FLAGS FOR POSTGRESQL/MYSQL + REDIS SUPPORT
// =============================================================================
// Location: pocketbase.go (eagerParseFlags function)
// Purpose: Add --dataDsn and --redisDsn flags to PocketBase CLI
//
// Integration: Add these flags to RootCmd.PersistentFlags() before Execute()
// =============================================================================

import (
	"os"

	"github.com/spf13/cobra"
)

// PocketBase struct additions (add these fields):
//
// type PocketBase struct {
//     ... existing fields ...
//
//     dataDataFlag      string // --dataDsn flag value
//     redisFlag         string // --redisDsn flag value
// }

// Config struct additions:
//
// type Config struct {
//     ... existing fields ...
//
//     DefaultDataDsn    string // Default dataDsn (e.g., from env var)
//     RedisDsn          string // Redis DSN (optional, for clustering)
// }

// eagerParseFlags parses the global app flags before calling pb.RootCmd.Execute()
// so we can have all PocketBase flags ready for use on initialization.
//
// IMPORTANT: This function is called BEFORE app.Bootstrap()
// Flags are parsed early so they can be used during initialization
func (pb *PocketBase) eagerParseFlags(config *Config) error {
	// -------------------------------------------------------------------------
	// Existing Flags (keep as-is)
	// -------------------------------------------------------------------------

	// --dir: Data directory (SQLite mode, now optional)
	pb.RootCmd.PersistentFlags().StringVar(
		&pb.dataDirFlag,
		"dir",
		config.DefaultDataDir,
		"the PocketBase data directory",
	)

	// -------------------------------------------------------------------------
	// NEW FLAG 1: --dataDsn (PostgreSQL/MySQL DSN)
	// -------------------------------------------------------------------------
	// Purpose: Specify PostgreSQL or MySQL database connection string
	// Default: Empty (falls back to SQLite if not provided)
	// Priority: --dataDsn > DB_DSN env var > SQLite (./pb_data)
	//
	// Examples:
	//   --dataDsn="postgres://user:pass@localhost:5432/db?sslmode=disable"
	//   --dataDsn="postgresql://user:pass@localhost:5432/db"
	//   --dataDsn="mysql://user:pass@tcp(localhost:3306)/db"
	//
	// Environment Variable Alternative:
	//   export DB_DSN="postgres://..."
	//   pocketbase serve  (no --dataDsn needed)

	pb.RootCmd.PersistentFlags().StringVar(
		&pb.dataDataFlag, // Store flag value in this field
		"dataDsn",        // Flag name: --dataDsn
		config.DefaultDataDsn, // Default value (usually empty or from env)
		// Help text (shown in --help)
		"store data postgresql/mysql dsn (e.g. postgres://user:pass@127.0.0.1:5432/db?sslmode=disable OR mysql://user:pass@tcp(127.0.0.1:3306)/db)",
	)

	// -------------------------------------------------------------------------
	// NEW FLAG 2: --redisDsn (Redis Pub/Sub for clustering)
	// -------------------------------------------------------------------------
	// Purpose: Enable multi-node clustering via Redis Pub/Sub
	// Default: Empty (single-node mode, no clustering)
	// Optional: Can run without Redis (graceful degradation)
	//
	// Examples:
	//   --redisDsn="redis://localhost:6379/0"
	//   --redisDsn="redis://:password@localhost:6379/0"
	//   --redisDsn="redis://user:pass@redis-cluster:6379/0"
	//
	// Use Case:
	//   Multi-node PocketBase behind load balancer
	//   Realtime subscriptions (SSE) work across all nodes
	//
	// Environment Variable Alternative:
	//   export REDIS_DSN="redis://..."
	//   pocketbase serve --dataDsn="postgres://..."  (no --redisDsn needed)

	pb.RootCmd.PersistentFlags().StringVar(
		&pb.redisFlag,  // Store flag value in this field
		"redisDsn",     // Flag name: --redisDsn
		config.RedisDsn, // Default value (usually empty)
		// Help text (shown in --help)
		"Cache data Redis dsn(default  redis://<user>:<pass>@localhost:6379/<db>  redis://localhost:6379/0)",
	)

	// -------------------------------------------------------------------------
	// Existing Flags (keep as-is)
	// -------------------------------------------------------------------------

	// --encryptionEnv: Encryption key environment variable
	pb.RootCmd.PersistentFlags().StringVar(
		&pb.encryptionEnvFlag,
		"encryptionEnv",
		config.DefaultEncryptionEnv,
		"the env variable whose value of 32 characters will be used \nas encryption key for the app settings (default none)",
	)

	// --debug: Enable debug logging
	pb.RootCmd.PersistentFlags().BoolVar(
		&pb.debugFlag,
		"debug",
		config.DefaultDebug,
		"enable debug mode, aka. showing more detailed logs",
	)

	// Parse flags early (before app.Bootstrap())
	return pb.RootCmd.ParseFlags(os.Args[1:])
}

// =============================================================================
// USAGE IN Bootstrap()
// =============================================================================
// After flags are parsed, pass them to app.Bootstrap():
//
// func (pb *PocketBase) Execute() error {
//     config := NewConfig()
//
//     // Parse flags early
//     if err := pb.eagerParseFlags(config); err != nil {
//         return err
//     }
//
//     // Pass flags to Bootstrap()
//     if !pb.skipBootstrap() {
//         if err := pb.Bootstrap(); err != nil {
//             return err
//         }
//     }
//
//     return pb.RootCmd.Execute()
// }

// =============================================================================
// PASSING FLAGS TO BaseApp
// =============================================================================
// In Bootstrap(), create BaseAppConfig with flag values:
//
// func (pb *PocketBase) Bootstrap() error {
//     config := core.BaseAppConfig{
//         DataDir:      pb.dataDirFlag,      // Existing
//         DataDsn:      pb.dataDataFlag,     // NEW: --dataDsn value
//         RedisDsn:     pb.redisFlag,        // NEW: --redisDsn value
//         EncryptionEnv: pb.encryptionEnvFlag,
//         IsDebug:      pb.debugFlag,
//     }
//
//     return pb.app.Bootstrap(config)
// }

// =============================================================================
// COMMAND LINE EXAMPLES
// =============================================================================
//
// SQLite Mode (Original Behavior):
//   pocketbase serve
//   pocketbase serve --dir=./data
//
// PostgreSQL Mode (Single Node):
//   pocketbase serve --dataDsn="postgres://user:pass@localhost:5432/db"
//
// PostgreSQL + Redis (Multi-Node Clustering):
//   # Node 1
//   pocketbase serve --dataDsn="postgres://shared-db:5432/db" \
//                    --redisDsn="redis://shared-redis:6379/0"
//
//   # Node 2 (same flags, different server)
//   pocketbase serve --dataDsn="postgres://shared-db:5432/db" \
//                    --redisDsn="redis://shared-redis:6379/0"
//
// MySQL Mode:
//   pocketbase serve --dataDsn="mysql://user:pass@tcp(localhost:3306)/db"
//
// Environment Variables (Alternative):
//   export DB_DSN="postgres://..."
//   export REDIS_DSN="redis://..."
//   pocketbase serve  (flags not needed)

// =============================================================================
// INTEGRATION CHECKLIST
// =============================================================================
// [ ] 1. Add dataDataFlag and redisFlag fields to PocketBase struct
// [ ] 2. Add DefaultDataDsn and RedisDsn to Config struct
// [ ] 3. Update eagerParseFlags() with new flag definitions
// [ ] 4. Pass flag values to BaseAppConfig in Bootstrap()
// [ ] 5. Update app.Bootstrap() to accept DataDsn and RedisDsn
// [ ] 6. Test with --help to verify flag descriptions
// [ ] 7. Test with invalid DSN to verify error handling

// File Stats:
// - Total Lines: ~30 (flag definitions only)
// - With Comments: ~180 (annotated)
// - Complexity: ⭐ Trivial (just flag definitions)
// - Dependencies: github.com/spf13/cobra
