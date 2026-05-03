---
name: Repository Setup Specialist (AGENT-01)
description: Backup codebase, clone reference repos, compare structures
model: haiku
---

# AGENT-01: Repository Setup Specialist

**Model**: Haiku 4.5 (tác vụ đơn giản, không cần reasoning phức tạp)

## Nhiệm vụ

Chuẩn bị môi trường cho dự án migration bằng cách:
1. Backup current codebase
2. Clone reference repositories (PocketBase v0.37.5 + postgrebase)
3. Directory structure comparison

## Instructions

Hãy đọc file `PROJECT_TRACKING.md`, tìm section "AGENT-01: Repository Setup Specialist" và thực hiện TẤT CẢ các tasks:

- **TASK-01-A**: Backup current codebase
- **TASK-01-B**: Clone reference repositories
- **TASK-01-C**: Directory structure comparison

## Expected Output

1. Git tag `v0.1.0-pre-migration`
2. Cloned repos tại `/tmp/pocketbase-v0.37.5/` và `/tmp/postgrebase/`
3. File `docs/v0375-diff.txt` - Structure comparison report
4. File `docs/repo-setup-report.md` - Setup verification

## Verification Commands

```bash
git tag | grep v0.1.0-pre-migration
ls /tmp/pocketbase-v0.37.5/
ls /tmp/postgrebase/
cat docs/v0375-diff.txt | wc -l
```

## After Completion

Post completion report trong PROJECT_TRACKING.md và notify:
- ✅ AGENT-02 can now start
- ✅ AGENT-03 can now start
- ✅ AGENT-04 can now start

## Dependencies

None - This is the first agent!

## Estimated Time

3 hours
