package ratelimit

import (
	"math"
	"sync"
	"time"
)

type TokenBucket struct {
	mu             sync.Mutex
	tokens         float64
	maxTokens      float64
	refillRate     float64
	lastRefillTime time.Time
}

func NewTokenBucket(maxTokens float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		tokens:         maxTokens,
		maxTokens:      maxTokens,
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	duration := now.Sub(tb.lastRefillTime)
	tokensToAdd := tb.refillRate * duration.Seconds()
	tb.tokens = math.Min(tb.tokens+tokensToAdd, tb.maxTokens)
	tb.lastRefillTime = now
}

func (tb *TokenBucket) Request(weight float64) (bool, time.Duration) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	if weight <= tb.tokens {
		tb.tokens -= weight
		return true, 0
	}

	return false, time.Duration((weight - tb.tokens) / tb.refillRate * float64(time.Second))
}
