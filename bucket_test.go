package ratelimit

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTokenBucket_ExhaustsCapacity(t *testing.T) {
	b := NewTokenBucket(5, 1)

	for i := 0; i < 5; i++ {
		if !b.Request(1) {
			t.Fatalf("request %d: want accepted, got denied", i)
		}
	}
	if b.Request(1) {
		t.Error("6th request: want denied (bucket empty), got accepted")
	}
}
func TestTokenBucket_RefillsOverTime(t *testing.T) {
	b := NewTokenBucket(1, 100)

	if !b.Request(1) {
		t.Fatal("first request must pass")
	}
	if b.Request(1) {
		t.Fatal("second immediate request must fail")
	}

	time.Sleep(50 * time.Millisecond)

	if !b.Request(1) {
		t.Error("after refill window, want accepted")
	}
}

func TestTokenBucket_ConcurrentExactCount(t *testing.T) {
	b := NewTokenBucket(100, 0)

	var wg sync.WaitGroup
	var accepted atomic.Int64

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if b.Request(1) {
				accepted.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := accepted.Load(); got != 100 {
		t.Errorf("want exactly 100 accepted, got %d", got)
	}
}
