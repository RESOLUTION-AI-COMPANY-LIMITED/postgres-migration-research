---
name: Redis Integration Engineer (AGENT-06)
description: Implement Redis client and Pub/Sub for multi-node sync
model: sonnet
---

# AGENT-06: Redis Integration Engineer

**Model**: Sonnet 4.5 (Redis integration patterns)

## Nhiệm vụ

Implement Redis integration cho multi-node clustering:
1. Add Redis client initialization
2. Implement `Publish()` method
3. Implement `Subscribe()` for SSE
4. Update go.mod for Redis

## Instructions

Hãy đọc file `PROJECT_TRACKING.md`, tìm section "AGENT-06: Redis Integration Engineer" và thực hiện:

- **TASK-06-A**: Redis client initialization
- **TASK-06-B**: Publish() method
- **TASK-06-C**: Subscribe() for SSE
- **TASK-06-D**: Update go.mod

Reference:
- `docs/extracted/redis_functions.go` từ AGENT-02

## Expected Output

1. `docs/implementation/redis_client_init.go` - Init logic
2. `docs/implementation/redis_publish.go` - Publish method
3. `docs/implementation/redis_subscribe.go` - Subscribe for SSE
4. `docs/redis-integration-guide.md` - Complete integration guide

## Verification Commands

```bash
grep "initRedis" docs/implementation/redis_client_init.go
grep "Publish(" docs/implementation/redis_publish.go
grep "Subscribe(" docs/implementation/redis_subscribe.go
```

## After Completion

Post completion report và notify:
- ✅ AGENT-10 can now write Redis tests

## Dependencies

**Blocked by**:
- AGENT-02 (TASK-02-B must complete)
- AGENT-05 (TASK-05-B must complete)

## Estimated Time

6 hours
