package token

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   float64
	refillRate float64
	tokens     float64
	lastRefill time.Time
	mutex      sync.Mutex // Không export (viết thường)
}

func NewTokenBucket(capacity, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		refillRate: refillRate,
		tokens:     capacity,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow(tokens float64) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.lastRefill = now

	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	if tb.tokens >= tokens {
		tb.tokens -= tokens
		return true
	}
	return false
}

// Thêm method mới
func (tb *TokenBucket) GetLastRefill() time.Time {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	return tb.lastRefill
}
