# ratelimit

Token bucket rate limiter in Go.

- Lazy refill on request
- Per-client limiter with periodic cleanup of idle buckets
- Thread-safe, race-tested
- Weighted requests support
- Gin middleware (limits by client IP)

## Usage

```go
limiter := ratelimit.NewLimiter(100, 10) // capacity 100 tokens, refill 10 tokens/sec
defer limiter.Close()

router := gin.New()
router.Use(ratelimit.RateLimitMiddleware(limiter))
```

## Test
    go test -v -race ./...
