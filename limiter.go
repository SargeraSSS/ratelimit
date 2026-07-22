package ratelimit

import (
	"sync"
	"time"
)

type Limiter struct {
	capacity   float64
	refillRate float64
	mutex      sync.Mutex
	buckets    map[string]*TokenBucket
}

func NewLimiter(capacity, refillRate float64) *Limiter {
	return &Limiter{
		capacity:       max_capacity,
		max_capacity:   max_capacity,
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
}

func (l *Limiter) Allow(clientID string) bool {
	if bucket, ok := l.buckets[clientID]; ok {
		return bool
	}
	newBucket := NewTokenBucket(l.capacity, l.refillRate)
	l.buckets[clientID] = newBucket
	return newBucket
}
