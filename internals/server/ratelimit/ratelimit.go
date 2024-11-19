package ratelimit

import (
    "sync"
    "time"
)

type IPRateLimiter struct {
    requests map[string]*requestInfo
    mu       sync.RWMutex
}

type requestInfo struct {
    count    int
    firstReq time.Time
}

const (
    maxRequests = 5            // Maximum requests allowed
    timeWindow  = time.Hour    // Time window for rate limiting
)

var Limiter = NewIPRateLimiter()

func NewIPRateLimiter() *IPRateLimiter {
    return &IPRateLimiter{
        requests: make(map[string]*requestInfo),
    }
}

func (rl *IPRateLimiter) IsAllowed(ip string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    now := time.Now()
    
    // Clean old entries
    if info, exists := rl.requests[ip]; exists {
        if now.Sub(info.firstReq) > timeWindow {
            delete(rl.requests, ip)
        }
    }

    // Check if IP exists in map
    if info, exists := rl.requests[ip]; exists {
        if now.Sub(info.firstReq) <= timeWindow {
            if info.count >= maxRequests {
                return false
            }
            info.count++
            return true
        }
    }

    // First request from this IP
    rl.requests[ip] = &requestInfo{
        count:    1,
        firstReq: now,
    }
    return true
}

func IsAllowed(ip string) bool {
	return Limiter.IsAllowed(ip)
}
