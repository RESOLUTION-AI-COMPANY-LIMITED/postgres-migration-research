# ⚠️ IMPLEMENTATION GUIDE - NOT ACTUAL CODE

**This is a GUIDE on how to implement, NOT working code.**

The actual implementation has NOT been done yet. This document provides:
- ✅ Step-by-step instructions
- ✅ Code examples to copy
- ✅ Testing procedures

**To use this guide**: Follow steps in actual PocketBase fork repository.

---

# Redis Integration Implementation Guide (AGENT-06)

**Agent**: AGENT-06 - Redis Integration Engineer  
**Date**: 2026-05-03  
**Project**: PocketBase v0.37.5 → PostgreSQL + Redis Migration

---

## 📋 Overview

Hướng dẫn chi tiết để implement Redis Pub/Sub cho multi-node clustering trong PocketBase. Đây là **RESEARCH PROJECT** - tạo documentation, KHÔNG phải actual implementation.

### Success Criteria
- ✅ Redis client initialization documented
- ✅ `Publish()` method implementation guide complete
- ✅ `Subscribe()` mechanics documented
- ✅ Graceful degradation strategy defined
- ✅ Multi-node clustering test plan provided

---

## 🎯 Phase 2 - AGENT-06 Tasks

### Dependencies
- ✅ AGENT-02 (Code Extraction) - `redis_functions.go` extracted
- 🟡 AGENT-05 (DB Connection Layer) - Needs `base.go` struct changes
- 🟢 Ready to start (can work in parallel with AGENT-05)

### Estimated Time
6 hours (research documentation)

---

## 🧩 Architecture Overview

### Design Pattern: Pub/Sub with Graceful Degradation

```
┌─────────────────────────────────────────────────────────────┐
│                    Load Balancer                            │
└────────────┬────────────────────────────┬───────────────────┘
             │                            │
             ▼                            ▼
    ┌────────────────┐          ┌────────────────┐
    │  PocketBase    │          │  PocketBase    │
    │  Node 1        │          │  Node 2        │
    │  (Port 8090)   │          │  (Port 8091)   │
    └────────┬───────┘          └────────┬───────┘
             │                            │
             │   ┌────────────────┐       │
             └───►  Redis Pub/Sub ◄───────┘
                 │  (Port 6379)   │
                 └────────────────┘
                         │
                         ▼
              ┌──────────────────┐
              │  PostgreSQL DB   │
              │  (Port 5432)     │
              └──────────────────┘

Flow:
1. User connects to Node 1 via SSE (Server-Sent Events)
2. Node 2 creates a new record
3. Node 2 publishes event to Redis "realtime" channel
4. Node 1 subscribes to Redis, receives event
5. Node 1 broadcasts to its SSE connections
6. User receives realtime update
```

### Graceful Degradation Strategy

| Scenario | Behavior | User Impact |
|----------|----------|-------------|
| **No Redis DSN** | Single-node mode (in-memory broadcast) | ✅ Full functionality (single node only) |
| **Redis DSN invalid** | Fail fast, return error | ❌ App won't start (quick feedback) |
| **Redis connection fails** | Disable Redis, log warning, continue | ⚠️ Single-node mode (no cross-node SSE) |
| **Redis crashes mid-run** | Next publish fails silently, fallback to local | ⚠️ Temporary single-node mode |

---

## 📦 STEP 1: Add Redis Fields to BaseApp Struct

### 1.1 Locate BaseApp Struct

**File:** `rai-backend/core/base.go` (existing file, modified by AGENT-05)

**Current Struct (After AGENT-05 changes):**
```go
type BaseApp struct {
	// Existing fields (from AGENT-05)
	dataDir           string
	encryptionEnv     string
	isDebug           bool
	db                *dbx.DB
	dao               *daos.Dao
	store             *store.Store<any]
	settings          *settings.Settings
	cache             *store.Store[any]
	
	// Event hooks
	onBeforeServe     *hook.Hook[*ServeEvent]
	onRealtimeConnect *hook.Hook[*RealtimeConnectEvent]
	onRealtimeBroadcast *hook.Hook[*RealtimeBroadcastEvent] // NEW: Add if not exists
	// ... more hooks ...
	
	// Database DSN (from AGENT-05)
	dataDsn           string
	
	// =========================================================================
	// AGENT-06: Redis Integration Fields
	// =========================================================================
	
	redisDsn          string         // Redis connection string (from --redisDsn flag)
	redisCache        *redis.Client  // Redis client instance (nil if disabled)
	redisContext      context.Context // Context for Redis operations
}
```

### 1.2 Add Redis Import Statements

**At top of `core/base.go`:**
```go
import (
	"context"          // For Redis context
	"encoding/json"    // For message marshaling
	"log"              // For error logging
	"time"             // For connection timeout
	
	// Existing imports...
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/store"
	
	// NEW: Redis client
	"github.com/redis/go-redis/v9"
	
	// NEW: Colored logging (optional)
	"github.com/fatih/color"
)
```

### 1.3 Update BaseAppConfig Struct

**Already modified by AGENT-05, verify this exists:**
```go
type BaseAppConfig struct {
	DataDir       string
	EncryptionEnv string
	IsDebug       bool
	DataDsn       string // From AGENT-05
	
	// =========================================================================
	// AGENT-06: Redis Configuration
	// =========================================================================
	
	RedisDsn      string // Redis DSN (optional, e.g., "redis://localhost:6379/0")
}
```

### 1.4 Update NewBaseApp() Constructor

**Modified Constructor (combine with AGENT-05 changes):**
```go
func NewBaseApp(config BaseAppConfig) *BaseApp {
	app := &BaseApp{
		// Existing fields
		dataDir:       config.DataDir,
		encryptionEnv: config.EncryptionEnv,
		isDebug:       config.IsDebug,
		store:         store.New[any](nil),
		cache:         store.New[any](nil),
		
		// Database DSN (from AGENT-05)
		dataDsn:       config.DataDsn,
		
		// =========================================================================
		// AGENT-06: Redis Configuration
		// =========================================================================
		
		redisDsn:      config.RedisDsn,      // Redis DSN from config
		redisContext:  context.Background(), // Initialize Redis context
		// redisCache will be set in initRedis() during Bootstrap()
	}
	
	// Initialize hooks
	app.initHooks()
	
	return app
}
```

---

## 🔌 STEP 2: Implement `initRedis()` Method

### 2.1 Complete Implementation

**Add to `core/base.go`:**

```go
// initRedis initializes Redis client for multi-node clustering
//
// Design Pattern: Pub/Sub with Graceful Degradation
// - If no Redis DSN: Skip (single-node mode)
// - If invalid DSN: Fail fast (return error)
// - If connection fails: Disable Redis, log warning, continue
//
// Called During: app.Bootstrap() (after database initialization)
//
// Background Goroutine:
// - Subscribes to "realtime" channel
// - Forwards Redis messages to PocketBase realtime system
// - Enables cross-node SSE broadcasting
func (app *BaseApp) initRedis() error {
	// =========================================================================
	// STEP 1: Graceful Skip (No Redis DSN)
	// =========================================================================
	
	if app.redisDsn == "" {
		if app.IsDebug() {
			color.Yellow("ℹ No Redis DSN provided")
			color.Yellow("Running in single-node mode (no clustering)")
		}
		return nil // Not an error, just skip Redis
	}
	
	// =========================================================================
	// STEP 2: Parse Redis DSN
	// =========================================================================
	
	// Supported DSN formats:
	// - redis://localhost:6379/0
	// - redis://:password@localhost:6379/0
	// - redis://user:pass@localhost:6379/0
	opt, err := redis.ParseURL(app.redisDsn)
	if err != nil {
		// Fail fast: Invalid DSN format
		return fmt.Errorf("invalid Redis DSN: %w", err)
	}
	
	// =========================================================================
	// STEP 3: Create Redis Client
	// =========================================================================
	
	app.redisCache = redis.NewClient(opt)
	
	// =========================================================================
	// STEP 4: Test Connection (with timeout)
	// =========================================================================
	
	ctx, cancel := context.WithTimeout(app.redisContext, 5*time.Second)
	defer cancel()
	
	if _, err := app.redisCache.Ping(ctx).Result(); err != nil {
		// Connection failed - Graceful degradation
		app.redisCache = nil // Disable Redis
		
		if app.IsDebug() {
			color.Red("⚠ Redis connection failed: %v", err)
			color.Yellow("Continuing in single-node mode (no clustering)")
			color.Yellow("Multi-node SSE will NOT work")
		}
		
		// NOT returning error - app continues without Redis
		return nil
	}
	
	// =========================================================================
	// STEP 5: Connection Successful - Start Subscriber
	// =========================================================================
	
	if app.IsDebug() {
		color.Green("✅ Redis connected successfully: %s", app.redisDsn)
		color.Cyan("Multi-node clustering enabled")
	}
	
	// =========================================================================
	// STEP 6: Subscribe to "realtime" Channel (Background Goroutine)
	// =========================================================================
	
	// This goroutine runs for the lifetime of the application
	// It listens for messages published by other nodes
	go func() {
		// Create Redis Pub/Sub subscription
		pubsub := app.redisCache.Subscribe(app.redisContext, "realtime")
		defer pubsub.Close()
		
		// Get message channel
		ch := pubsub.Channel()
		
		if app.IsDebug() {
			color.Cyan("📡 Redis Pub/Sub subscriber started (channel: realtime)")
		}
		
		// Listen for messages indefinitely
		for msg := range ch {
			// Debug log (optional)
			if app.IsDebug() {
				log.Printf("Redis message received: channel=%s payload=%s", msg.Channel, msg.Payload)
			}
			
			// Create realtime broadcast event
			event := &RealtimeBroadcastEvent{
				App:     app,
				Channel: msg.Channel,
				Payload: []byte(msg.Payload),
			}
			
			// Trigger PocketBase realtime hooks
			// This broadcasts to local SSE connections on THIS node
			if err := app.OnRealtimeBroadcast().Trigger(event); err != nil {
				if app.IsDebug() {
					log.Printf("Realtime broadcast error: %v", err)
				}
			}
		}
		
		// Channel closed - Redis connection lost
		if app.IsDebug() {
			color.Yellow("⚠ Redis Pub/Sub subscriber stopped")
		}
	}()
	
	return nil
}
```

### 2.2 Integration Point in Bootstrap()

**Add to `core/base.go` `Bootstrap()` method:**

```go
func (app *BaseApp) Bootstrap() error {
	// ... existing database initialization (from AGENT-05) ...
	
	// =========================================================================
	// AGENT-06: Initialize Redis (after database)
	// =========================================================================
	
	if err := app.initRedis(); err != nil {
		// Redis initialization failed with error (invalid DSN)
		return fmt.Errorf("Redis initialization failed: %w", err)
	}
	
	// Note: If Redis connection fails (not invalid DSN), initRedis() returns nil
	// and app continues in single-node mode
	
	// ... rest of existing Bootstrap() logic ...
	
	return nil
}
```

---

## 📤 STEP 3: Implement `Publish()` Method

### 3.1 Complete Implementation

**Add to `core/base.go`:**

```go
// Publish broadcasts a message to all nodes via Redis Pub/Sub
//
// Behavior:
// - If Redis enabled: Publish to Redis channel (all nodes receive)
// - If Redis disabled: Fallback to local broadcast (only this node receives)
//
// Parameters:
//   channel - Redis channel name (e.g., "realtime")
//   data    - Any serializable data (will be marshaled to JSON)
//
// Returns:
//   error - Publish error (or nil on success)
//
// Usage Example:
//   app.Publish("realtime", map[string]any{
//       "action": "create",
//       "record": recordData,
//   })
func (app *BaseApp) Publish(channel string, data any) error {
	// =========================================================================
	// STEP 1: Marshal Data to JSON
	// =========================================================================
	
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal publish data: %w", err)
	}
	
	// =========================================================================
	// STEP 2: Try Redis Publish (if available)
	// =========================================================================
	
	if app.redisCache != nil {
		// Redis enabled - Publish to Redis channel
		// All subscribed nodes (including this one) will receive the message
		
		if app.IsDebug() {
			log.Printf("Publishing to Redis: channel=%s payload=%s", channel, string(payload))
		}
		
		err := app.redisCache.Publish(app.redisContext, channel, payload).Err()
		if err != nil {
			// Redis publish failed - Log error but don't fail
			if app.IsDebug() {
				color.Yellow("⚠ Redis publish failed: %v", err)
				color.Yellow("Falling back to local broadcast")
			}
			
			// Fallback to local broadcast (continue execution)
			goto LocalBroadcast
		}
		
		// Redis publish successful
		return nil
	}
	
	// =========================================================================
	// STEP 3: Fallback to Local Broadcast (single-node mode)
	// =========================================================================
	
LocalBroadcast:
	// Redis not available - Trigger local realtime hooks only
	// This broadcasts to SSE connections on THIS node only
	
	if app.IsDebug() {
		log.Printf("Local broadcast: channel=%s payload=%s", channel, string(payload))
	}
	
	event := &RealtimeBroadcastEvent{
		App:     app,
		Channel: channel,
		Payload: payload,
	}
	
	return app.OnRealtimeBroadcast().Trigger(event)
}
```

### 3.2 Usage in Realtime Handlers

**Example: Trigger realtime event after record creation**

**File:** `rai-backend/apis/record_crud.go` (or similar)

```go
func (api *RecordApi) create(c echo.Context) error {
	// ... existing record creation logic ...
	
	// After record created successfully
	record := // ... created record ...
	
	// =========================================================================
	// NEW: Broadcast realtime event to all nodes
	// =========================================================================
	
	api.app.Publish("realtime", map[string]any{
		"action":     "create",
		"collection": record.Collection().Name,
		"record":     record,
	})
	
	return c.JSON(200, record)
}
```

---

## 📥 STEP 4: Implement `Subscribe()` for SSE

### 4.1 Subscriber Already Running

**Important:** The subscriber is already implemented in `initRedis()` (Step 2)

The background goroutine automatically:
1. Subscribes to "realtime" channel
2. Receives messages from Redis
3. Forwards to `OnRealtimeBroadcast()` hook
4. Triggers local SSE connections

### 4.2 Verify Realtime Hook Exists

**Check `core/base.go` for this hook:**

```go
type BaseApp struct {
	// ... existing fields ...
	
	// Realtime hooks
	onRealtimeConnect   *hook.Hook[*RealtimeConnectEvent]
	onRealtimeDisconnect *hook.Hook[*RealtimeDisconnectEvent]
	onRealtimeBroadcast *hook.Hook[*RealtimeBroadcastEvent] // REQUIRED for Redis
}
```

**If `onRealtimeBroadcast` doesn't exist, add it:**

```go
func (app *BaseApp) initHooks() {
	// ... existing hooks ...
	
	// NEW: Realtime broadcast hook (for Redis Pub/Sub)
	app.onRealtimeBroadcast = &hook.Hook[*RealtimeBroadcastEvent]{}
}

// OnRealtimeBroadcast hook is triggered when a realtime message is broadcast
func (app *BaseApp) OnRealtimeBroadcast() *hook.Hook[*RealtimeBroadcastEvent] {
	return app.onRealtimeBroadcast
}
```

### 4.3 Define RealtimeBroadcastEvent

**Add to `core/base.go` or `core/realtime_events.go`:**

```go
// RealtimeBroadcastEvent is triggered when a realtime message is broadcast
// Used for forwarding Redis Pub/Sub messages to local SSE connections
type RealtimeBroadcastEvent struct {
	App     *BaseApp
	Channel string // Redis channel name
	Payload []byte // JSON payload
}
```

### 4.4 SSE Handler Integration

**Example SSE handler:** `apis/realtime.go`

```go
func (api *RealtimeApi) subscribeHandler(c echo.Context) error {
	// ... existing SSE setup ...
	
	// Register hook listener for this SSE connection
	api.app.OnRealtimeBroadcast().Add(func(e *RealtimeBroadcastEvent) error {
		// Parse payload
		var data map[string]any
		if err := json.Unmarshal(e.Payload, &data); err != nil {
			return err
		}
		
		// Filter by collection (if needed)
		// ... collection filtering logic ...
		
		// Send to SSE client
		fmt.Fprintf(c.Response(), "data: %s\n\n", e.Payload)
		c.Response().Flush()
		
		return nil
	})
	
	// ... keep connection alive ...
}
```

---

## 🧪 STEP 5: Update `go.mod`

### 5.1 Add Redis Dependency

**Command (documentation only):**
```bash
go get github.com/redis/go-redis/v9@v9.3.0
```

**Resulting `go.mod` changes:**
```go
require (
	// Existing dependencies (from AGENT-05)...
	github.com/lib/pq v1.10.9
	github.com/go-sql-driver/mysql v1.7.1
	
	// NEW: Redis client
	github.com/redis/go-redis/v9 v9.3.0
)
```

### 5.2 Optional: Add Colored Logging

**Command (documentation only):**
```bash
go get github.com/fatih/color@latest
```

**Resulting `go.mod` changes:**
```go
require (
	// ... existing dependencies ...
	
	// Optional: Colored terminal output
	github.com/fatih/color v1.16.0
)
```

---

## 🧪 STEP 6: Testing Strategy

### 6.1 Unit Test: Redis Client Initialization

**File:** `core/redis_test.go`

```go
package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInitRedis_NoDSN(t *testing.T) {
	// Test: No Redis DSN provided (single-node mode)
	app := NewBaseApp(BaseAppConfig{
		DataDir:  "./test_data",
		RedisDsn: "", // Empty DSN
	})
	
	err := app.initRedis()
	
	assert.NoError(t, err)
	assert.Nil(t, app.redisCache) // Redis should be nil
}

func TestInitRedis_InvalidDSN(t *testing.T) {
	// Test: Invalid Redis DSN (should fail fast)
	app := NewBaseApp(BaseAppConfig{
		DataDir:  "./test_data",
		RedisDsn: "invalid://format",
	})
	
	err := app.initRedis()
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid Redis DSN")
}

func TestInitRedis_ValidDSN_ConnectionFails(t *testing.T) {
	// Test: Valid DSN but Redis not running (graceful degradation)
	app := NewBaseApp(BaseAppConfig{
		DataDir:  "./test_data",
		RedisDsn: "redis://localhost:9999/0", // Wrong port
	})
	
	err := app.initRedis()
	
	// Should NOT error (graceful degradation)
	assert.NoError(t, err)
	assert.Nil(t, app.redisCache) // Redis should be disabled
}

func TestInitRedis_Success(t *testing.T) {
	// Test: Valid DSN with running Redis
	// Prerequisite: Redis must be running on localhost:6379
	
	app := NewBaseApp(BaseAppConfig{
		DataDir:  "./test_data",
		RedisDsn: "redis://localhost:6379/0",
	})
	
	err := app.initRedis()
	
	assert.NoError(t, err)
	assert.NotNil(t, app.redisCache) // Redis should be connected
	
	// Test PING
	pong, err := app.redisCache.Ping(app.redisContext).Result()
	assert.NoError(t, err)
	assert.Equal(t, "PONG", pong)
}
```

### 6.2 Integration Test: Multi-Node Pub/Sub

**File:** `tests/integration/redis_pubsub_test.go`

```go
package integration

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestRedisPubSub_MultiNode(t *testing.T) {
	// Prerequisite: Redis must be running
	redisDsn := "redis://localhost:6379/0"
	
	// =========================================================================
	// Setup: Start 2 PocketBase instances
	// =========================================================================
	
	// Node 1
	app1 := NewTestApp(BaseAppConfig{
		DataDir:  "./test_data_node1",
		RedisDsn: redisDsn,
	})
	defer app1.Cleanup()
	
	// Node 2
	app2 := NewTestApp(BaseAppConfig{
		DataDir:  "./test_data_node2",
		RedisDsn: redisDsn,
	})
	defer app2.Cleanup()
	
	// =========================================================================
	// Test: Node 1 publishes, Node 2 receives
	// =========================================================================
	
	// Register listener on Node 2
	received := false
	app2.OnRealtimeBroadcast().Add(func(e *RealtimeBroadcastEvent) error {
		received = true
		t.Logf("Node 2 received: %s", string(e.Payload))
		return nil
	})
	
	// Node 1 publishes message
	err := app1.Publish("realtime", map[string]any{
		"action": "test",
		"data":   "hello from node 1",
	})
	assert.NoError(t, err)
	
	// Wait for message propagation
	time.Sleep(100 * time.Millisecond)
	
	// Verify Node 2 received the message
	assert.True(t, received, "Node 2 should have received message from Node 1")
}
```

### 6.3 Manual Test: Docker Compose Setup

**File:** `docker-compose.redis-test.yml`

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: pocketbase_test
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data

  pocketbase-node1:
    build: .
    command: serve --dataDsn="postgres://postgres:password@postgres:5432/pocketbase_test?sslmode=disable" --redisDsn="redis://redis:6379/0" --http="0.0.0.0:8090"
    ports:
      - "8090:8090"
    depends_on:
      - postgres
      - redis

  pocketbase-node2:
    build: .
    command: serve --dataDsn="postgres://postgres:password@postgres:5432/pocketbase_test?sslmode=disable" --redisDsn="redis://redis:6379/0" --http="0.0.0.0:8090"
    ports:
      - "8091:8090"
    depends_on:
      - postgres
      - redis

volumes:
  postgres_data:
  redis_data:
```

**Start test environment:**
```bash
docker-compose -f docker-compose.redis-test.yml up -d
```

**Test multi-node SSE:**
```bash
# Terminal 1: Connect SSE to Node 1
curl -N http://localhost:8090/api/realtime

# Terminal 2: Connect SSE to Node 2
curl -N http://localhost:8091/api/realtime

# Terminal 3: Create record on Node 1
curl -X POST http://localhost:8090/api/collections/users/records \
  -H "Content-Type: application/json" \
  -d '{"name": "Test User"}'

# Expected: BOTH Terminal 1 and Terminal 2 should receive realtime event
```

---

## 🐛 Error Handling & Edge Cases

### Edge Case 1: Redis Goes Down Mid-Run

**Scenario:** Redis server crashes while app is running

**Current Behavior:**
- Next `Publish()` call fails
- Falls back to local broadcast automatically
- Subscriber goroutine exits when channel closes

**Improvement (Optional):**
```go
// Add reconnection logic in initRedis()
go func() {
	for {
		pubsub := app.redisCache.Subscribe(app.redisContext, "realtime")
		
		for msg := range pubsub.Channel() {
			// ... handle message ...
		}
		
		// Channel closed - try reconnect
		if app.IsDebug() {
			color.Yellow("Redis disconnected, retrying in 5s...")
		}
		time.Sleep(5 * time.Second)
	}
}()
```

### Edge Case 2: Message Too Large

**Scenario:** Published data exceeds Redis max message size (512MB default)

**Solution:**
```go
func (app *BaseApp) Publish(channel string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	// Check message size (1MB limit recommended)
	const maxSize = 1 * 1024 * 1024 // 1MB
	if len(payload) > maxSize {
		return fmt.Errorf("publish payload too large: %d bytes (max %d)", len(payload), maxSize)
	}
	
	// ... rest of publish logic ...
}
```

### Edge Case 3: Duplicate Events (Echo Problem)

**Scenario:** Node 1 publishes to Redis, receives its own message back

**Solution 1: Add node ID to skip own messages**
```go
type RealtimeMessage struct {
	NodeID  string         `json:"node_id"`
	Channel string         `json:"channel"`
	Data    map[string]any `json:"data"`
}

// In Publish()
msg := RealtimeMessage{
	NodeID:  app.nodeID, // Add nodeID field to BaseApp
	Channel: channel,
	Data:    data.(map[string]any),
}

// In subscriber
for msg := range ch {
	var realtimeMsg RealtimeMessage
	json.Unmarshal([]byte(msg.Payload), &realtimeMsg)
	
	// Skip own messages
	if realtimeMsg.NodeID == app.nodeID {
		continue
	}
	
	// ... trigger hook ...
}
```

**Solution 2: Use separate local broadcast method**
```go
// PublishLocal: Broadcast to local SSE only (no Redis)
func (app *BaseApp) PublishLocal(channel string, data any) error {
	// ... trigger local hook only, skip Redis ...
}

// Publish: Broadcast to Redis (all nodes, including this one)
func (app *BaseApp) Publish(channel string, data any) error {
	// ... publish to Redis only, subscriber will trigger local hook ...
}
```

---

## 📊 Performance Considerations

### Redis Connection Pool Tuning

**Recommended settings:**
```go
func (app *BaseApp) initRedis() error {
	// ... parse DSN ...
	
	opt.PoolSize = 100          // Max simultaneous connections
	opt.MinIdleConns = 10       // Keep idle connections ready
	opt.DialTimeout = 5 * time.Second
	opt.ReadTimeout = 3 * time.Second
	opt.WriteTimeout = 3 * time.Second
	
	app.redisCache = redis.NewClient(opt)
	
	// ... rest of init ...
}
```

### Pub/Sub Latency Benchmarks

**Expected latency (local network):**
- Publish → Subscribe: **1-5ms**
- Publish → SSE client: **10-50ms**

**Load testing script:**
```bash
# Test Redis Pub/Sub performance
for i in {1..1000}; do
  curl -X POST http://localhost:8090/api/collections/users/records \
    -H "Content-Type: application/json" \
    -d '{"name": "Test '$i'"}'
done

# Monitor Redis stats
redis-cli INFO stats | grep pubsub
```

---

## 📁 Deliverables Summary

### Files Created (Documentation)
1. ✅ `docs/implementation/redis-integration-guide.md` (this file)
2. 📄 `core/base.go` (Redis methods added)
3. 📄 `core/redis_test.go` (unit tests)
4. 📄 `tests/integration/redis_pubsub_test.go` (integration tests)
5. 📄 `docker-compose.redis-test.yml` (test environment)

### Code Changes Required (Summary)
1. **Modified File:** `core/base.go`
   - Add fields: `redisDsn`, `redisCache`, `redisContext`
   - Add method: `initRedis()` (~100 lines)
   - Add method: `Publish()` (~50 lines)
   - Add method: `RedisCache()` (getter)
   - Update `Bootstrap()` to call `initRedis()`
2. **Modified File:** `go.mod`
   - Add `github.com/redis/go-redis/v9 v9.3.0`

---

## ✅ Verification Checklist

**Before marking AGENT-06 complete:**

- [ ] Redis fields added to `BaseApp` struct
- [ ] `initRedis()` implementation documented
- [ ] `Publish()` implementation documented
- [ ] Subscriber goroutine logic documented
- [ ] `OnRealtimeBroadcast()` hook integration documented
- [ ] Unit tests documented
- [ ] Integration tests documented
- [ ] Docker Compose test environment documented
- [ ] Error handling & edge cases covered
- [ ] Performance tuning recommendations provided

---

## 🔗 Dependencies & Handoffs

### Upstream Dependencies
- ✅ **AGENT-02** - Extracted `redis_functions.go` code
- 🟡 **AGENT-05** - `BaseApp` struct changes (can work in parallel)

### Downstream Dependencies
- 🟡 **AGENT-07** - CLI flags (needs `RedisDsn` field in config)
- 🟡 **AGENT-10** - Unit testing (needs `initRedis()` and `Publish()` methods)
- 🟡 **AGENT-11** - Integration testing (needs multi-node setup)

---

## 📌 Next Steps

**After AGENT-06 completion:**

1. **Handoff to AGENT-07:** CLI Flags Guide
   - Add `--redisDsn` flag
   - Map to `BaseAppConfig.RedisDsn`

2. **Handoff to AGENT-10:** Unit Testing
   - Test `initRedis()` with various DSN formats
   - Test `Publish()` with/without Redis

3. **Handoff to AGENT-11:** Integration Testing
   - Multi-node clustering test
   - SSE cross-node broadcast test

---

**AGENT-06 Status:** 🟢 Ready for Review  
**Document Version:** 1.0  
**Last Updated:** 2026-05-03
