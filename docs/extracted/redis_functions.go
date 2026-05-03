package core

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/redis/go-redis/v9"
)

// =============================================================================
// REDIS INTEGRATION FOR POCKETBASE
// =============================================================================
// Purpose: Enable multi-node clustering via Redis Pub/Sub
// Use Case: Realtime subscriptions (SSE) across multiple PocketBase instances
// Design Pattern: Pub/Sub with graceful degradation (fallback to local)
//
// Integration Points:
// 1. core/base.go: Add these functions to BaseApp struct
// 2. cmd/serve.go: Add --redisDsn flag
// 3. Bootstrap(): Call initRedis() during app startup
// =============================================================================

// BaseApp struct additions (add these fields to existing struct):
//
// type BaseApp struct {
//     ... existing fields ...
//
//     // Redis integration
//     redisDsn         string         // Redis connection string (from CLI flag)
//     redisCache       *redis.Client  // Redis client instance
//     redisContext     context.Context // Context for Redis operations
// }

// Config struct additions:
//
// type BaseAppConfig struct {
//     ... existing fields ...
//
//     RedisDsn         string // Redis DSN from --redisDsn flag
// }

// =============================================================================
// FUNCTION 1: initRedis() - Initialize Redis Client
// =============================================================================
// Called during app.Bootstrap() after database initialization
// Graceful Degradation: If Redis fails, app continues without clustering
//
// Error Handling Strategy:
// - Empty DSN: Skip Redis (single-node mode)
// - Invalid DSN: Return error (fail fast)
// - Connection fails: Set redisCache = nil, log warning, continue
//
// Background Goroutine:
// - Subscribes to "realtime" channel
// - Forwards Redis messages to PocketBase realtime system
// - Enables cross-node SSE broadcasting

func (app *BaseApp) initRedis() error {
	// Graceful skip: No Redis DSN provided
	if app.redisDsn == "" {
		return nil // Single-node mode (no clustering)
	}

	// Parse Redis connection URL
	// Format: redis://user:pass@localhost:6379/0
	opt, err := redis.ParseURL(app.redisDsn)
	if err != nil {
		return err // Fail fast on invalid DSN
	}

	// Create Redis client
	app.redisCache = redis.NewClient(opt)

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(app.redisContext, 5*time.Second)
	defer cancel()

	if _, err := app.redisCache.Ping(ctx).Result(); err != nil {
		// Connection failed - disable Redis and continue
		app.redisCache = nil // Graceful degradation
		if app.IsDebug() {
			color.Red("Redis connection failed: %v", err)
			color.Yellow("Continuing in single-node mode (no clustering)")
		}
	} else {
		// Connection successful
		if app.IsDebug() {
			color.Green("Redis connected successfully")
			color.Cyan("Multi-node clustering enabled")
		}

		// Subscribe to realtime channel for cross-node messaging
		// This goroutine runs for the lifetime of the application
		go func() {
			pubsub := app.redisCache.Subscribe(app.redisContext, "realtime")
			defer pubsub.Close()

			ch := pubsub.Channel()

			// Listen for Redis Pub/Sub messages
			for msg := range ch {
				// Forward Redis message to PocketBase realtime system
				event := &RealtimeBroadcastEvent{
					App:     app,
					Channel: msg.Channel,
					Payload: []byte(msg.Payload),
				}

				// Trigger PocketBase realtime hooks
				// This broadcasts to local SSE connections
				if err := app.OnRealtimeBroadcast().Trigger(event); err != nil && app.IsDebug() {
					log.Println("Realtime broadcast error:", err)
				}
			}
		}()
	}

	return nil
}

// =============================================================================
// FUNCTION 2: Publish() - Broadcast Messages to Redis
// =============================================================================
// Purpose: Send realtime updates to all nodes in cluster
//
// Behavior:
// - If Redis available: Publish to Redis channel
// - If Redis unavailable: Fallback to local broadcast only
//
// Flow:
// 1. Marshal data to JSON
// 2. Try Redis publish (if enabled)
// 3. Fallback to local broadcast if Redis disabled
//
// Example Usage:
//   app.Publish("realtime", map[string]any{
//       "action": "create",
//       "record": record,
//   })

func (app *BaseApp) Publish(channel string, data any) error {
	// Marshal data to JSON payload
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Try Redis publish if available
	if app.redisCache != nil {
		// Publish to Redis channel
		// All subscribed nodes will receive this message
		return app.redisCache.Publish(app.redisContext, channel, payload).Err()
	}

	// Fallback: Local broadcast only (single-node mode)
	// This triggers realtime hooks for local SSE connections only
	event := &RealtimeBroadcastEvent{
		App:     app,
		Channel: channel,
		Payload: payload,
	}
	return app.OnRealtimeBroadcast().Trigger(event)
}

// =============================================================================
// FUNCTION 3: RedisCache() - Get Redis Client
// =============================================================================
// Purpose: Expose Redis client for custom caching logic
//
// Usage Example:
//   if redis := app.RedisCache(); redis != nil {
//       redis.Set(ctx, "key", "value", time.Hour)
//   }

func (app *BaseApp) RedisCache() *redis.Client {
	return app.redisCache
}

// =============================================================================
// INTEGRATION CHECKLIST
// =============================================================================
// [ ] 1. Add fields to BaseApp struct (redisDsn, redisCache, redisContext)
// [ ] 2. Add RedisDsn to BaseAppConfig struct
// [ ] 3. Update NewBaseApp() constructor to set redisContext
// [ ] 4. Call app.initRedis() in app.Bootstrap() (after DB init)
// [ ] 5. Add --redisDsn flag to cmd/serve.go
// [ ] 6. Add redis/go-redis/v9 to go.mod
// [ ] 7. Test single-node mode (no Redis)
// [ ] 8. Test multi-node mode (with Redis)

// =============================================================================
// DEPLOYMENT CONSIDERATIONS
// =============================================================================
//
// Single-Node Mode:
//   pocketbase serve --dataDsn="postgres://localhost/db"
//   (No --redisDsn flag = no clustering)
//
// Multi-Node Mode:
//   # Node 1
//   pocketbase serve --dataDsn="postgres://shared-db/db" \
//                    --redisDsn="redis://shared-redis:6379/0"
//
//   # Node 2
//   pocketbase serve --dataDsn="postgres://shared-db/db" \
//                    --redisDsn="redis://shared-redis:6379/0"
//
// Load Balancer:
//   Use sticky sessions or Redis-backed session store

// File Stats:
// - Total Lines: ~65 (original code)
// - With Comments: ~180 (annotated)
// - Complexity: ⭐⭐ Moderate (Pub/Sub pattern)
// - Dependencies: github.com/redis/go-redis/v9
