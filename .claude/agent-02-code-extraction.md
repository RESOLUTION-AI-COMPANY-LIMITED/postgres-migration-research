---
name: Code Extraction Specialist (AGENT-02)
description: Extract core implementation code from postgrebase
model: sonnet
---

# AGENT-02: Code Extraction Specialist

**Model**: Sonnet 4.5 (cần hiểu code structure, annotate patterns)

## Nhiệm vụ

Extract và document implementation patterns từ postgrebase:
1. Database connection layer (`core/db_postgresql.go` - 29 lines)
2. Redis integration functions (~100 lines)
3. CLI flags
4. Create code mapping document

## Instructions

Hãy đọc file `PROJECT_TRACKING.md`, tìm section "AGENT-02: Code Extraction Specialist" và thực hiện:

- **TASK-02-A**: Extract database connection layer
- **TASK-02-B**: Extract Redis integration code
- **TASK-02-C**: Extract command-line flags
- **TASK-02-D**: Create code mapping document

Khi extract code, hãy:
- Giữ nguyên code với annotations
- Document function signatures
- Explain design patterns
- Cross-reference với POSTGREBASE_CODE_ANALYSIS.md

## Expected Output

1. `docs/extracted/db_postgresql.go` - Database connection layer
2. `docs/extracted/redis_functions.go` - Redis integration (~150 lines)
3. `docs/extracted/cmd_serve_flags.go` - CLI flags
4. `docs/code-mapping.md` - Implementation mapping guide

## Verification Commands

```bash
wc -l docs/extracted/*.go
grep -c "func " docs/extracted/*.go
cat docs/code-mapping.md | grep "##" | wc -l
```

## After Completion

Post completion report và notify:
- ✅ AGENT-05 can now start (needs TASK-02-A completed)
- ✅ AGENT-06 can now start (needs TASK-02-B completed)
- ✅ AGENT-07 can now start (needs TASK-02-C completed)

## Dependencies

**Blocked by**: AGENT-01 (TASK-01-B must complete first)

## Estimated Time

4 hours
