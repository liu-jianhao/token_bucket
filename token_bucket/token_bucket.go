package token_bucket

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	cap   int64
	avail int64

	timer *time.Ticker
	mutex *sync.Mutex
}

func New(interval time.Duration, cap int64) *TokenBucket {
	if interval < 0 {
		panic(fmt.Sprintf("interval %v < 0", interval))
	}

	if cap < 0 {
		panic(fmt.Sprintf("cap %v < 0", cap))
	}

	tb := &TokenBucket{
		cap:   cap,
		avail: cap,
		timer: time.NewTicker(interval),
		mutex: &sync.Mutex{},
	}

	go tb.daemon()
	return tb
}

func (tb *TokenBucket) Stop() {
	tb.timer.Stop()
}

func (tb *TokenBucket) Capability() int64 {
	return tb.cap
}

func (tb *TokenBucket) Available() int64 {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	return tb.avail
}

func (tb *TokenBucket) TryTake(count int64) bool {
	if count <= 0 || count > tb.cap {
		return false
	}

	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	if count <= tb.avail {
		tb.avail -= count
		return true
	}

	return false
}

func (tb *TokenBucket) daemon() {
	for range tb.timer.C {
		tb.mutex.Lock()

		if tb.avail < tb.cap {
			tb.avail++
		}

		tb.mutex.Unlock()
	}
}
