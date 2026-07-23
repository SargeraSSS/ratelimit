package ratelimit

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
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
			if l.Allow("client1") {
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
			if !l.Allow(clientID) {
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
