package main

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu         sync.Mutex
	timestamps map[string][]time.Time
	rateLimit  int
}

func NewRateLimiter(rateLimit int) *RateLimiter {
	return &RateLimiter{
		timestamps: make(map[string][]time.Time),
		rateLimit:  rateLimit,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	times, exists := rl.timestamps[ip]
	if !exists {
		rl.timestamps[ip] = []time.Time{now}
		return true
	}

	if len(times) < rl.rateLimit {
		rl.timestamps[ip] = append(times, now)
		return true
	}

	earliest := times[0]
	if now.Sub(earliest) < time.Second {
		return false
	}

	rl.timestamps[ip] = append(times[1:], now)
	return true
}
