---
name: Integration Test Engineer (AGENT-11)
description: E2E integration testing with PostgreSQL + Redis
model: sonnet
---

# AGENT-11: Integration Test Engineer

**Model**: Sonnet 4.5 (E2E test scenarios)

## Nhiệm vụ

Write end-to-end integration tests:
1. Setup test infrastructure (Docker Compose)
2. Write API compatibility tests
3. Write multi-node clustering test
4. Test data migration from SQLite

## Instructions

Hãy đọc file `PROJECT_TRACKING.md`, tìm section "AGENT-11: Integration Test Engineer" và thực hiện:

- **TASK-11-A**: Setup test infrastructure
- **TASK-11-B**: API compatibility tests
- **TASK-11-C**: Multi-node clustering test
- **TASK-11-D**: Data migration test

Test infrastructure:
- Docker Compose với PostgreSQL + Redis
- Multi-node setup (2+ PocketBase instances)
- Test database với sample data

## Expected Output

1. `docs/tests/docker-compose.test.yml` - Test infrastructure
2. `docs/tests/api_compatibility_test.go` - API tests
3. `docs/tests/clustering_test.go` - Multi-node tests
4. `docs/tests/data_migration_test.go` - Migration tests
5. `docs/tests/integration-test-results.md` - E2E test results

## Verification Commands

```bash
docker-compose -f docs/tests/docker-compose.test.yml config
grep "func TestAPI" docs/tests/api_compatibility_test.go
grep "func TestClustering" docs/tests/clustering_test.go
cat docs/tests/integration-test-results.md | grep "PASS"
```

## After Completion

Post completion report và notify:
- ✅ AGENT-12 can now start performance benchmarks

## Dependencies

**Blocked by**:
- AGENT-10 (all tasks must complete)

## Estimated Time

8 hours
