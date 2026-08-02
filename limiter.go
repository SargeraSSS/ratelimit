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
	stopChan   chan struct{}
}

type clientEntry struct {
	bucket   *TokenBucket
	lastSeen time.Time
}

const cleanupThreshold = 5 * time.Minute

func NewLimiter(capacity, refillRate float64) *Limiter {
	l := &Limiter{
		capacity:   capacity,
		refillRate: refillRate,
		buckets:    make(map[string]*clientEntry),
		stopChan:   make(chan struct{}),
	}
	go l.cleanupLoop()
	return l
}

func (l *Limiter) Allow(clientID string) (bool, time.Duration) {
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

func (l *Limiter) cleanup() {
	for clientID, entry := range l.buckets {
		if time.Since(entry.lastSeen) > cleanupThreshold {
			delete(l.buckets, clientID)
		}
	}
}

func (l *Limiter) cleanupLoop() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.mutex.Lock()
			l.cleanup()
			l.mutex.Unlock()
		case <-l.stopChan:
			return
		}
	}
}

func (l *Limiter) Close() {
	close(l.stopChan)
}
