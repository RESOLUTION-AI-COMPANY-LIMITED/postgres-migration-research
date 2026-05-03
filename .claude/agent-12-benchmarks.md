---
name: Performance Benchmark Engineer (AGENT-12)
description: Performance testing and optimization recommendations
model: opus
---

# AGENT-12: Performance Benchmark Engineer

**Model**: Opus 4.7 (complex performance analysis + optimization recommendations)

## Nhiệm vụ

Benchmark performance và provide optimization recommendations:
1. Write database query benchmarks
2. Write Redis cache benchmarks
3. Run benchmark suite and compare with baseline
4. Provide performance optimization recommendations

## Instructions

Hãy đọc file `PROJECT_TRACKING.md`, tìm section "AGENT-12: Performance Benchmark Engineer" và thực hiện:

- **TASK-12-A**: Database query benchmarks
- **TASK-12-B**: Redis cache benchmarks
- **TASK-12-C**: Run benchmark suite
- **TASK-12-D**: Performance optimization

Benchmark scenarios:
- Database: SELECT/INSERT/UPDATE/DELETE operations
- Redis: Cache hit/miss, Pub/Sub latency
- Concurrent operations (100+ simultaneous requests)
- Compare với SQLite baseline

## Expected Output

1. `docs/benchmarks/db_query_bench_test.go` - Database benchmarks
2. `docs/benchmarks/redis_cache_bench_test.go` - Redis benchmarks
3. `docs/benchmarks/benchmark-results.md` - Results với graphs
4. `docs/benchmarks/optimization-recommendations.md` - Performance tuning guide

## Verification Commands

```bash
grep "func Benchmark" docs/benchmarks/*.go | wc -l
cat docs/benchmarks/benchmark-results.md | grep "ops/sec"
cat docs/benchmarks/optimization-recommendations.md | grep "##" | wc -l
```

## After Completion

Post **FINAL** completion report!

This is the last agent. Project complete! 🎉

Deliverables summary:
- All 12 agents completed
- All tasks verified
- Complete migration documentation
- Ready for implementation

## Dependencies

**Blocked by**:
- AGENT-11 (all tasks must complete)

## Estimated Time

6 hours

---

## 🎯 Success Criteria

- [ ] Benchmarks show acceptable performance (within 20% of SQLite baseline)
- [ ] Redis caching improves multi-node performance
- [ ] Optimization recommendations documented
- [ ] Final project report generated
