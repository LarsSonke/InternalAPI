package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	rate       int           // requests per interval
	interval   time.Duration // time window
	buckets    map[string]*bucket
	mu         sync.RWMutex
	cleanupInt time.Duration
}

type bucket struct {
	tokens     int
	lastRefill time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		rate:       rate,
		interval:   interval,
		buckets:    make(map[string]*bucket),
		cleanupInt: interval * 10, // cleanup old buckets every 10 intervals
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow checks if a request should be allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.RLock()
	b, exists := rl.buckets[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		b = &bucket{
			tokens:     rl.rate,
			lastRefill: time.Now(),
		}
		rl.buckets[key] = b
		rl.mu.Unlock()
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	// Refill tokens based on time passed
	now := time.Now()
	elapsed := now.Sub(b.lastRefill)
	if elapsed >= rl.interval {
		b.tokens = rl.rate
		b.lastRefill = now
	}

	// Check if request is allowed
	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}

// cleanup removes stale buckets
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInt)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, b := range rl.buckets {
			b.mu.Lock()
			if now.Sub(b.lastRefill) > rl.cleanupInt {
				delete(rl.buckets, key)
			}
			b.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// RateLimitByIP creates middleware that rate limits by IP address
func RateLimitByIP(rate int, interval time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, interval)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "Too many requests. Please try again later.",
				"retry_after": interval.Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByUser creates middleware that rate limits by authenticated user
func RateLimitByUser(rate int, interval time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, interval)

	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userID, exists := c.Get("userID")
		if !exists {
			// If no user ID, fall back to IP-based limiting
			userID = c.ClientIP()
		}

		key := userID.(string)
		
		if !limiter.Allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "Too many requests. Please try again later.",
				"retry_after": interval.Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// StrictRateLimitByIP creates middleware with stricter limits (e.g., for login)
func StrictRateLimitByIP(rate int, interval time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, interval)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "Too many login attempts. Please try again later.",
				"retry_after": interval.Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
