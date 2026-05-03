---
name: CLI Flags Engineer (AGENT-07)
description: Add --dataDsn and --redisDsn flags
model: haiku
---

# AGENT-07: CLI Flags Engineer

**Model**: Haiku 4.5 (simple CLI modification)

## Nhiệm vụ

Add command-line flags cho database và Redis DSN:
1. Add `--dataDsn` và `--redisDsn` flags to cmd/serve.go
2. Pass DSN to `app.Bootstrap()`
3. Add environment variable support
4. Update help text

## Instructions

Hãy đọc file `PROJECT_TRACKING.md`, tìm section "AGENT-07: CLI Flags Engineer" và thực hiện:

- **TASK-07-A**: Add flags to cmd/serve.go
- **TASK-07-B**: Pass DSN to Bootstrap()
- **TASK-07-C**: Environment variable support
- **TASK-07-D**: Update help text

Reference:
- `docs/extracted/cmd_serve_flags.go` từ AGENT-02

## Expected Output

1. `docs/implementation/cmd_serve_flags.go` - Flag implementation
2. `docs/cli-usage-guide.md` - CLI usage examples

## Verification Commands

```bash
grep "dataDsn" docs/implementation/cmd_serve_flags.go
grep "redisDsn" docs/implementation/cmd_serve_flags.go
grep "Environment" docs/cli-usage-guide.md
```

## After Completion

Post completion report (no downstream dependencies)

## Dependencies

**Blocked by**:
- AGENT-02 (TASK-02-C must complete)
- AGENT-05 (TASK-05-B must complete)

## Estimated Time

3 hours
