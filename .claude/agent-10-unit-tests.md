---
name: Unit Test Engineer (AGENT-10)
description: Write unit tests for all components
model: sonnet
---

# AGENT-10: Unit Test Engineer

**Model**: Sonnet 4.5 (test generation with coverage)

## Nhiệm vụ

Write comprehensive unit tests:
1. Database connection tests
2. Redis integration tests
3. Query builder tests
4. Run all tests and collect coverage

## Instructions

Hãy đọc file `PROJECT_TRACKING.md`, tìm section "AGENT-10: Unit Test Engineer" và thực hiện:

- **TASK-10-A**: Database connection tests
- **TASK-10-B**: Redis integration tests
- **TASK-10-C**: Query builder tests
- **TASK-10-D**: Run all tests

Test targets:
- `core/db_postgresql.go` - connectDB() function
- Redis client - Publish()/Subscribe()
- Query builder - SQL generation

## Expected Output

1. `docs/tests/db_connection_test.go` - Database tests
2. `docs/tests/redis_integration_test.go` - Redis tests
3. `docs/tests/query_builder_test.go` - Query builder tests
4. `docs/tests/unit-test-results.md` - Test results với coverage

## Verification Commands

```bash
grep "func Test" docs/tests/*.go | wc -l
cat docs/tests/unit-test-results.md | grep "coverage:"
cat docs/tests/unit-test-results.md | grep "PASS"
```

## After Completion

Post completion report và notify:
- ✅ AGENT-11 can now start integration tests

## Dependencies

**Blocked by**:
- AGENT-05 (all tasks must complete)
- AGENT-06 (all tasks must complete)
- AGENT-08 (all tasks must complete)

## Estimated Time

8 hours
