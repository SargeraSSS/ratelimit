package ratelimit

import (
	"sync"
	"time"
)

type Limiter struct {
	capacity   float64
	refillRate float64
	mutex      sync.Mutex
	buckets    map[string]*clientEntry
}
type clientEntry struct {
	bucket   *TokenBucket
	lastSeen time.Time
}

func NewLimiter(capacity, refillRate float64) *Limiter {
	return &Limiter{
		capacity:   capacity,
		refillRate: refillRate,
		buckets:    make(map[string]*clientEntry),
	}
}

func (l *Limiter) Allow(clientID string) bool {
	bucket := l.getOrCreateBucket(clientID)
	return bucket.Request(1)
}

func (l *Limiter) getOrCreateBucket(clientID string) *TokenBucket {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if entry, exists := l.buckets[clientID]; exists {
		entry.lastSeen = time.Now()
		return entry.bucket
	}
	newBucket := NewTokenBucket(l.capacity, l.refillRate)
	l.buckets[clientID] = &clientEntry{
		bucket:   newBucket,
		lastSeen: time.Now(),
	}
	return newBucket
}
