package token

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokens     float64
	capacity   float64
	refillRate float64
	lastRefill time.Time
	lastUsed   time.Time // Thêm trường mới
	mutex      sync.Mutex
}

func NewTokenBucket(capacity float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity,
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
		lastUsed:   time.Now(), // Khởi tạo
	}
}

// Bổ sung phương thức LastUsed
func (tb *TokenBucket) LastUsed() time.Time {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	return tb.lastUsed
}

func (tb *TokenBucket) Allow(tokensCost float64) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	// Cập nhật thời gian sử dụng cuối
	tb.lastUsed = time.Now()

	// Refill logic
	now := time.Now()
	duration := now.Sub(tb.lastRefill)
	tokensToAdd := tb.refillRate * duration.Seconds()
	tb.tokens = tb.tokens + tokensToAdd
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastRefill = now

	if tb.tokens < tokensCost {
		return false
	}

	tb.tokens -= tokensCost
	return true
}
