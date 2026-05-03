---
name: Dependency Analysis Specialist (AGENT-03)
description: Analyze Go dependencies and conflicts
model: sonnet
---

# AGENT-03: Dependency Analysis Specialist

**Model**: Sonnet 4.5 (cần phân tích dependencies phức tạp)

## Nhiệm vụ

Analyze Go dependencies và create conflict resolution plan:
1. Compare go.mod files (postgrebase vs PocketBase v0.37.5 vs rai-backend)
2. List new dependencies (lib/pq, go-sql-driver, go-redis)
3. Create conflict resolution plan

## Instructions

Hãy đọc file `PROJECT_TRACKING.md`, tìm section "AGENT-03: Dependency Analysis Specialist" và thực hiện:

- **TASK-03-A**: Compare go.mod files
- **TASK-03-B**: List new dependencies
- **TASK-03-C**: Conflict resolution plan

## Expected Output

1. `docs/dependency-comparison.md` - go.mod comparison matrix
2. `docs/new-dependencies.txt` - List of new packages
3. `docs/conflict-resolution-plan.md` - Dependency conflict resolution

## Verification Commands

```bash
cat docs/dependency-comparison.md | grep "lib/pq"
cat docs/new-dependencies.txt | wc -l
cat docs/conflict-resolution-plan.md | grep "##" | wc -l
```

## After Completion

Post completion report và notify:
- ✅ AGENT-05 can now update go.mod (needs TASK-03-B completed)

## Dependencies

**Blocked by**: AGENT-01 (TASK-01-B must complete first)

## Estimated Time

3 hours
