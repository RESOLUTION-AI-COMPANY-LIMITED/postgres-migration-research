package core

import (
	"strings"

	// PostgreSQL driver - imported for side effects (registers "postgres" driver)
	_ "github.com/lib/pq"

	// MySQL driver - imported for side effects (registers "mysql" driver)
	_ "github.com/go-sql-driver/mysql"

	// Modified dbx query builder with PostgreSQL/MySQL support
	// NOTE: This is a custom fork, not the original go-ozzo/ozzo-dbx
	"github.com/free/postgresqlbaseapi/dbx"
)

// connectDB creates a database connection from a DSN string
//
// Design Pattern: Factory Pattern with Auto-Detection
// - Detects driver type from DSN prefix
// - Supports: postgres://, postgresql://, mysql://
// - Defaults to postgres if no prefix
//
// Parameters:
//   dsn - Data Source Name (connection string)
//         Examples:
//         - "postgres://user:pass@localhost/dbname"
//         - "postgresql://user:pass@localhost/dbname"
//         - "mysql://user:pass@localhost/dbname"
//
// Returns:
//   *dbx.DB - Database connection with query builder
//   error   - Connection error if any
//
// Usage in PocketBase:
//   This replaces the original SQLite-only connection logic in core/base.go
//   Original: db, err := connectDB(dataDir)
//   Modified: db, err := connectDB(dataDsn) // DSN passed from CLI flag
func connectDB(dsn string) (*dbx.DB, error) {
	// Default driver is PostgreSQL
	driver := "postgres"

	// Parse driver from DSN prefix
	// NOTE: This allows users to specify driver explicitly via DSN
	if strings.HasPrefix(dsn, "mysql://") {
		driver = "mysql"
		// Remove mysql:// prefix for driver compatibility
		dsn = strings.TrimPrefix(dsn, "mysql://")
	} else if strings.HasPrefix(dsn, "postgres://") {
		driver = "postgres"
		// Keep postgres:// prefix - lib/pq expects it
	} else if strings.HasPrefix(dsn, "postgresql://") {
		driver = "postgres"
		// postgresql:// is alias for postgres://
	}
	// If no prefix, assumes PostgreSQL with raw DSN

	// Open database connection using dbx query builder
	// dbx.Open() wraps database/sql Open() and adds query builder
	db, err := dbx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	// TODO: Connection pooling configuration
	// db.DB().SetMaxOpenConns(100)
	// db.DB().SetMaxIdleConns(10)
	// db.DB().SetConnMaxLifetime(time.Hour)

	return db, nil
}

// Integration Points in PocketBase:
// 1. core/base.go: Replace connectDB() call with this implementation
// 2. cmd/serve.go: Add --dataDsn flag to pass DSN string
// 3. go.mod: Add lib/pq and go-sql-driver/mysql dependencies
// 4. Vendor dbx package: Copy github.com/free/postgresqlbaseapi/dbx

// File Stats:
// - Total Lines: 29 (original)
// - With Comments: 73 (annotated)
// - Complexity: ⭐ Very Simple
// - Dependencies: lib/pq, go-sql-driver/mysql, dbx

// Testing:
// func TestConnectDB_Postgres(t *testing.T) {
//     db, err := connectDB("postgres://localhost/test")
//     assert.NoError(t, err)
//     assert.NotNil(t, db)
// }
