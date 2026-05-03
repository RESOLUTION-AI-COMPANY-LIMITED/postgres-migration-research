# Integration Test Plan (AGENT-11)

**Agent**: AGENT-11 - Integration Test Engineer  
**Date**: 2026-05-03  
**Project**: PocketBase v0.37.5 → PostgreSQL + Redis Migration

---

## 📋 Overview

End-to-end integration testing plan cho multi-node clustering và API compatibility. Đây là **RESEARCH PROJECT** - documentation ONLY.

### Success Criteria
- ✅ Multi-node clustering tests documented
- ✅ API compatibility tests documented
- ✅ SSE cross-node broadcast tests documented
- ✅ Data migration tests documented

---

## 🎯 Phase 4 - AGENT-11 Tasks

### Dependencies
- 🟡 AGENT-10 (Unit Tests) - Complete

### Estimated Time
8 hours (research documentation)

---

## 🧪 TEST SUITE 1: Multi-Node Clustering

### File: `tests/integration/cluster_test.go`

```go
package integration

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Test 1: Two-Node Redis Pub/Sub Communication
// =============================================================================

func TestMultiNode_RedisPubSub_MessageBroadcast(t *testing.T) {
	// Prerequisites:
	// - PostgreSQL running on localhost:5432
	// - Redis running on localhost:6379
	
	sharedPostgresDB := "postgres://postgres:password@localhost:5432/testdb?sslmode=disable"
	sharedRedis := "redis://localhost:6379/0"
	
	// =========================================================================
	// Setup: Start 2 PocketBase instances
	// =========================================================================
	
	// Node 1
	app1 := NewTestApp(BaseAppConfig{
		DataDir:  "./test_data_node1",
		DataDsn:  sharedPostgresDB,
		RedisDsn: sharedRedis,
	})
	defer app1.Cleanup()
	require.NotNil(t, app1)
	
	// Node 2
	app2 := NewTestApp(BaseAppConfig{
		DataDir:  "./test_data_node2",
		DataDsn:  sharedPostgresDB,
		RedisDsn: sharedRedis,
	})
	defer app2.Cleanup()
	require.NotNil(t, app2)
	
	// =========================================================================
	// Test: Node 1 publishes → Node 2 receives
	// =========================================================================
	
	// Register listener on Node 2
	node2Received := false
	var receivedPayload []byte
	
	app2.OnRealtimeBroadcast().Add(func(e *RealtimeBroadcastEvent) error {
		node2Received = true
		receivedPayload = e.Payload
		t.Logf("Node 2 received: %s", string(e.Payload))
		return nil
	})
	
	// Node 1 publishes message
	testData := map[string]any{
		"action": "create",
		"record": map[string]any{
			"id":   "test123",
			"name": "Test Record",
		},
	}
	
	err := app1.Publish("realtime", testData)
	require.NoError(t, err)
	
	// Wait for message propagation
	time.Sleep(200 * time.Millisecond)
	
	// Verify Node 2 received the message
	assert.True(t, node2Received, "Node 2 should receive message from Node 1")
	assert.Contains(t, string(receivedPayload), "test123")
}

// =============================================================================
// Test 2: Three-Node Load Balancing
// =============================================================================

func TestMultiNode_ThreeNodes_LoadBalancing(t *testing.T) {
	sharedPostgresDB := "postgres://postgres:password@localhost:5432/testdb?sslmode=disable"
	sharedRedis := "redis://localhost:6379/0"
	
	// Start 3 nodes
	nodes := make([]*TestApp, 3)
	for i := 0; i < 3; i++ {
		nodes[i] = NewTestApp(BaseAppConfig{
			DataDir:  fmt.Sprintf("./test_data_node%d", i+1),
			DataDsn:  sharedPostgresDB,
			RedisDsn: sharedRedis,
		})
		defer nodes[i].Cleanup()
	}
	
	// Register listeners on all nodes
	receiveCount := make([]int, 3)
	for i, node := range nodes {
		nodeIdx := i
		node.OnRealtimeBroadcast().Add(func(e *RealtimeBroadcastEvent) error {
			receiveCount[nodeIdx]++
			return nil
		})
	}
	
	// Node 0 publishes 10 messages
	for i := 0; i < 10; i++ {
		err := nodes[0].Publish("realtime", map[string]any{"msg": i})
		require.NoError(t, err)
		time.Sleep(50 * time.Millisecond)
	}
	
	// Verify all nodes received messages
	for i, count := range receiveCount {
		assert.Equal(t, 10, count, fmt.Sprintf("Node %d should receive all 10 messages", i))
	}
}

// =============================================================================
// Test 3: Node Failure Resilience
// =============================================================================

func TestMultiNode_NodeFailure_GracefulDegradation(t *testing.T) {
	sharedPostgresDB := "postgres://postgres:password@localhost:5432/testdb?sslmode=disable"
	sharedRedis := "redis://localhost:6379/0"
	
	// Start 2 nodes
	app1 := NewTestApp(BaseAppConfig{
		DataDir:  "./test_data_node1",
		DataDsn:  sharedPostgresDB,
		RedisDsn: sharedRedis,
	})
	defer app1.Cleanup()
	
	app2 := NewTestApp(BaseAppConfig{
		DataDir:  "./test_data_node2",
		DataDsn:  sharedPostgresDB,
		RedisDsn: sharedRedis,
	})
	
	// Verify communication works
	received := false
	app2.OnRealtimeBroadcast().Add(func(e *RealtimeBroadcastEvent) error {
		received = true
		return nil
	})
	
	err := app1.Publish("realtime", map[string]any{"test": "pre-shutdown"})
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	assert.True(t, received)
	
	// Shutdown Node 2
	app2.Cleanup()
	t.Log("Node 2 shut down")
	
	// Node 1 should still be able to publish (no crash)
	received = false
	err = app1.Publish("realtime", map[string]any{"test": "post-shutdown"})
	assert.NoError(t, err, "Node 1 should continue working after Node 2 shutdown")
}
```

---

## 🧪 TEST SUITE 2: API Compatibility

### File: `tests/integration/api_compatibility_test.go`

```go
package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Test 1: CRUD API Endpoints (PostgreSQL Backend)
// =============================================================================

func TestAPI_CRUD_WithPostgreSQL(t *testing.T) {
	// Setup app with PostgreSQL backend
	app := NewTestApp(BaseAppConfig{
		DataDsn: "postgres://postgres:password@localhost:5432/testdb?sslmode=disable",
	})
	defer app.Cleanup()
	
	// Setup test HTTP server
	server := httptest.NewServer(app.Handler())
	defer server.Close()
	
	// =========================================================================
	// Test: CREATE record
	// =========================================================================
	
	createPayload := map[string]any{
		"name":  "Test User",
		"email": "test@example.com",
		"age":   30,
	}
	
	createBody, _ := json.Marshal(createPayload)
	resp, err := http.Post(server.URL+"/api/collections/users/records", "application/json", bytes.NewReader(createBody))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	var createdRecord map[string]any
	json.NewDecoder(resp.Body).Decode(&createdRecord)
	recordID := createdRecord["id"].(string)
	t.Logf("Created record ID: %s", recordID)
	
	// =========================================================================
	// Test: READ record
	// =========================================================================
	
	resp, err = http.Get(server.URL + "/api/collections/users/records/" + recordID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	var fetchedRecord map[string]any
	json.NewDecoder(resp.Body).Decode(&fetchedRecord)
	assert.Equal(t, "Test User", fetchedRecord["name"])
	
	// =========================================================================
	// Test: UPDATE record
	// =========================================================================
	
	updatePayload := map[string]any{"age": 31}
	updateBody, _ := json.Marshal(updatePayload)
	
	req, _ := http.NewRequest(http.MethodPatch, server.URL+"/api/collections/users/records/"+recordID, bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// =========================================================================
	// Test: DELETE record
	// =========================================================================
	
	req, _ = http.NewRequest(http.MethodDelete, server.URL+"/api/collections/users/records/"+recordID, nil)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// =============================================================================
// Test 2: Authentication Endpoints
// =============================================================================

func TestAPI_Authentication_WithPostgreSQL(t *testing.T) {
	app := NewTestApp(BaseAppConfig{
		DataDsn: "postgres://postgres:password@localhost:5432/testdb?sslmode=disable",
	})
	defer app.Cleanup()
	
	server := httptest.NewServer(app.Handler())
	defer server.Close()
	
	// Test: User signup
	signupPayload := map[string]any{
		"email":           "newuser@example.com",
		"password":        "password123",
		"passwordConfirm": "password123",
	}
	
	signupBody, _ := json.Marshal(signupPayload)
	resp, err := http.Post(server.URL+"/api/collections/users/auth-with-password", "application/json", bytes.NewReader(signupBody))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// Test: User login
	loginPayload := map[string]any{
		"identity": "newuser@example.com",
		"password": "password123",
	}
	
	loginBody, _ := json.Marshal(loginPayload)
	resp, err = http.Post(server.URL+"/api/collections/users/auth-with-password", "application/json", bytes.NewReader(loginBody))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	var authResponse map[string]any
	json.NewDecoder(resp.Body).Decode(&authResponse)
	assert.NotEmpty(t, authResponse["token"])
}
```

---

## 🧪 TEST SUITE 3: SSE Cross-Node Broadcast

### File: `tests/integration/sse_broadcast_test.go`

```go
package integration

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Test: SSE Event Broadcasting Across Nodes
// =============================================================================

func TestSSE_CrossNodeBroadcast(t *testing.T) {
	sharedPostgresDB := "postgres://postgres:password@localhost:5432/testdb?sslmode=disable"
	sharedRedis := "redis://localhost:6379/0"
	
	// =========================================================================
	// Setup: Start 2 nodes
	// =========================================================================
	
	app1 := NewTestApp(BaseAppConfig{
		DataDir:  "./test_data_node1",
		DataDsn:  sharedPostgresDB,
		RedisDsn: sharedRedis,
	})
	defer app1.Cleanup()
	
	app2 := NewTestApp(BaseAppConfig{
		DataDir:  "./test_data_node2",
		DataDsn:  sharedPostgresDB,
		RedisDsn: sharedRedis,
	})
	defer app2.Cleanup()
	
	// =========================================================================
	// Setup: Connect SSE client to Node 1
	// =========================================================================
	
	server1 := httptest.NewServer(app1.Handler())
	defer server1.Close()
	
	sseURL := server1.URL + "/api/realtime"
	req, _ := http.NewRequest("GET", sseURL, nil)
	req.Header.Set("Accept", "text/event-stream")
	
	client := &http.Client{Timeout: 0} // No timeout for SSE
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))
	
	// =========================================================================
	// Test: Node 2 creates record → Node 1's SSE client receives event
	// =========================================================================
	
	// Start SSE reader goroutine
	eventReceived := make(chan string, 1)
	go func() {
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data:") {
				eventReceived <- line
				return
			}
		}
	}()
	
	// Node 2 creates record (triggers broadcast)
	server2 := httptest.NewServer(app2.Handler())
	defer server2.Close()
	
	createPayload := `{"name": "Test Record", "value": "test"}`
	resp2, err := http.Post(server2.URL+"/api/collections/items/records", "application/json", strings.NewReader(createPayload))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
	
	// Wait for SSE event (timeout after 3 seconds)
	select {
	case event := <-eventReceived:
		t.Logf("SSE event received: %s", event)
		assert.Contains(t, event, "create")
		assert.Contains(t, event, "Test Record")
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout: SSE event not received from Node 2")
	}
}
```

---

## 🧪 TEST SUITE 4: Data Migration

### File: `tests/integration/data_migration_test.go`

```go
package integration

import (
	"database/sql"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
)

// =============================================================================
// Test: Migrate SQLite Data to PostgreSQL
// =============================================================================

func TestDataMigration_SQLiteToPostgreSQL(t *testing.T) {
	// =========================================================================
	// Step 1: Setup SQLite source database with test data
	// =========================================================================
	
	sqliteDB, err := sql.Open("sqlite3", "./test_data/source.db")
	require.NoError(t, err)
	defer sqliteDB.Close()
	
	// Create table and insert test data
	_, err = sqliteDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT,
			created DATETIME DEFAULT (DATETIME('now'))
		)
	`)
	require.NoError(t, err)
	
	_, err = sqliteDB.Exec(`INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com')`)
	require.NoError(t, err)
	
	// Count SQLite records
	var sqliteCount int
	sqliteDB.QueryRow("SELECT COUNT(*) FROM users").Scan(&sqliteCount)
	t.Logf("SQLite record count: %d", sqliteCount)
	
	// =========================================================================
	// Step 2: Setup PostgreSQL target database
	// =========================================================================
	
	postgresDB, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/migration_test?sslmode=disable")
	require.NoError(t, err)
	defer postgresDB.Close()
	
	// Create table (PostgreSQL syntax)
	_, err = postgresDB.Exec(`
		DROP TABLE IF EXISTS users CASCADE;
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT,
			created TIMESTAMP DEFAULT NOW()
		)
	`)
	require.NoError(t, err)
	
	// =========================================================================
	// Step 3: Migrate data (INSERT INTO ... SELECT FROM ...)
	// =========================================================================
	
	// Fetch data from SQLite
	rows, err := sqliteDB.Query("SELECT name, email FROM users")
	require.NoError(t, err)
	defer rows.Close()
	
	// Insert into PostgreSQL
	insertCount := 0
	for rows.Next() {
		var name, email string
		rows.Scan(&name, &email)
		
		_, err = postgresDB.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
		require.NoError(t, err)
		insertCount++
	}
	
	t.Logf("Migrated %d records to PostgreSQL", insertCount)
	
	// =========================================================================
	// Step 4: Verify data integrity
	// =========================================================================
	
	var postgresCount int
	postgresDB.QueryRow("SELECT COUNT(*) FROM users").Scan(&postgresCount)
	
	assert.Equal(t, sqliteCount, postgresCount, "Record count should match")
	assert.Equal(t, sqliteCount, insertCount, "All records should be migrated")
	
	// Verify data content
	var name, email string
	postgresDB.QueryRow("SELECT name, email FROM users WHERE email = $1", "john@example.com").Scan(&name, &email)
	assert.Equal(t, "John Doe", name)
	assert.Equal(t, "john@example.com", email)
}
```

---

## 🐳 Docker Compose Test Environment

### File: `docker-compose.integration-test.yml`

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: testdb
    ports:
      - "5432:5432"
    volumes:
      - postgres_test_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes
    volumes:
      - redis_test_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  pocketbase-node1:
    build: .
    command: serve --dataDsn="postgres://postgres:password@postgres:5432/testdb?sslmode=disable" --redisDsn="redis://redis:6379/0" --http="0.0.0.0:8090"
    ports:
      - "8090:8090"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy

  pocketbase-node2:
    build: .
    command: serve --dataDsn="postgres://postgres:password@postgres:5432/testdb?sslmode=disable" --redisDsn="redis://redis:6379/0" --http="0.0.0.0:8090"
    ports:
      - "8091:8090"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy

volumes:
  postgres_test_data:
  redis_test_data:
```

**Start environment:**
```bash
docker-compose -f docker-compose.integration-test.yml up -d
```

**Run integration tests:**
```bash
go test ./tests/integration/... -v -timeout=10m
```

---

## ✅ Verification Checklist

**Before marking AGENT-11 complete:**

- [ ] Multi-node clustering tests documented (3+ scenarios)
- [ ] API compatibility tests documented (CRUD + Auth)
- [ ] SSE cross-node broadcast tests documented
- [ ] Data migration tests documented
- [ ] Docker Compose test environment documented
- [ ] Test execution commands documented
- [ ] All integration tests passing

---

**AGENT-11 Status:** 🟢 Ready for Review  
**Document Version:** 1.0  
**Last Updated:** 2026-05-03
