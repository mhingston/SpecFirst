# RFC: Distributed Rate Limiting with Redis

## Summary
We propose replacing our in-memory rate limiter with a distributed Redis-backed solution to support our new multi-region architecture.

## Design
We will use the sliding window algorithm implementation in `go-redis/redis_rate`.
Keys will be namespaced by `{region}:{user_id}`.

## Trade-offs
- **Latency**: Adding a network hop to Redis adds ~2ms.
- **Consistency**: Redis allows us to enforce global limits, which in-memory cannot do.
- **Cost**: Redis instance costs are negligible compared to current API spend.

## Migration
We will roll out to 1% of traffic, then 10%, then 100%. Fallback to in-memory if Redis is down.
