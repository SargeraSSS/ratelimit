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
		if allowed, _ := b.Request(1); !allowed {
			t.Fatalf("request %d: want accepted, got denied", i)
		}
	}
	if allowed, _ := b.Request(1); allowed {
		t.Error("6th request: want denied (bucket empty), got accepted")
	}
}

func TestTokenBucket_ExhaustsCapacity_Table(t *testing.T) {
	tests := []struct {
		name     string
		capacity float64
	}{
		{"capacity of 0", 0},
		{"capacity of 1", 1},
		{"capacity of 5", 5},
		{"capacity of 6", 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewTokenBucket(tt.capacity, 1)
			for i := 0; i < int(tt.capacity); i++ {
				if allowed, _ := b.Request(1); !allowed {
					t.Fatalf("request %d: want accepted, got denied", i)
				}
			}
			if allowed, _ := b.Request(1); allowed {
				t.Error("request after exhaustion: want denied, got accepted")
			}
		})
	}
}

func TestTokenBucket_NeverExceedsCapacity(t *testing.T) {
	capacity := 3.0
	rate := 100.0

	b := NewTokenBucket(capacity, rate)

	time.Sleep(200 * time.Millisecond)

	accepted := 0
	attempts := int(capacity) + 2

	for i := 0; i < attempts; i++ {
		if allowed, _ := b.Request(1); allowed {
			accepted++
		}
	}

	if accepted != int(capacity) {
		t.Errorf("accepted %d requests after long sleep, want exactly %d (capacity)", accepted, int(capacity))
	}
}

func TestTokenBucket_RefillsOverTime(t *testing.T) {
	b := NewTokenBucket(1, 100)

	if allowed, _ := b.Request(1); !allowed {
		t.Fatal("first request must pass")
	}
	if allowed, _ := b.Request(1); allowed {
		t.Fatal("second immediate request must fail")
	}

	time.Sleep(50 * time.Millisecond)

	if allowed, _ := b.Request(1); !allowed {
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
			if allowed, _ := b.Request(1); allowed {
				accepted.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := accepted.Load(); got != 100 {
		t.Errorf("want exactly 100 accepted, got %d", got)
	}
}
