package ratelimit

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLimiter_ConcurrentAccess(t *testing.T) {
	l := NewLimiter(100, 10)
	defer l.Close()

	var wg sync.WaitGroup
	var accepted atomic.Int64
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if allowed, _ := l.Allow("client1"); allowed {
				accepted.Add(1)
			}
		}()
	}
	var wg2 sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg2.Add(1)
		go func(id int) {
			defer wg2.Done()
			clientID := fmt.Sprintf("client-%d", id)
			if allowed, _ := l.Allow(clientID); !allowed {
				t.Errorf("unique client %s: want accepted, got denied", clientID)
			}
		}(i)
	}
	wg2.Wait()
	wg.Wait()

	if got := accepted.Load(); got != 100 {
		t.Errorf("client1: want exactly 100 accepted, got %d", got)
	}
}

func TestLimiter_CleanupRemovesStaleClients(t *testing.T) {
	l := NewLimiter(10, 1)
	defer l.Close()

	l.Allow("active")
	l.Allow("stale")

	l.mutex.Lock()
	l.buckets["stale"].lastSeen = time.Now().Add(-2 * cleanupThreshold)
	l.cleanup()
	_, staleExists := l.buckets["stale"]
	_, activeExists := l.buckets["active"]
	l.mutex.Unlock()

	if staleExists {
		t.Error("stale client: want removed, got kept")
	}
	if !activeExists {
		t.Error("active client: want kept, got removed")
	}
}
