# ratelimit

Token bucket rate limiter in Go.

- Lazy refill (no background goroutines)
- Thread-safe, race-tested
- Weighted requests support

## Test
    go test -v -race ./...