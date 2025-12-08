package rest

import (
	"net/http"
	"sync"
	"time"

	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	mu         sync.RWMutex
	buckets    map[string]*TokenBucket
	rate       int           // tokens per window
	window     time.Duration // time window
	maxBurst   int           // maximum burst size
	cleanupTTL time.Duration // how long to keep inactive buckets
}

// TokenBucket represents a token bucket for a single key
type TokenBucket struct {
	tokens       int
	lastRefill   time.Time
	lastAccessed time.Time
	mu           sync.Mutex
}

// NewRateLimiter creates a new rate limiter
// rate: number of requests allowed per window
// window: time window duration
// maxBurst: maximum burst size (extra tokens available)
func NewRateLimiter(rate int, window time.Duration, maxBurst int) *RateLimiter {
	rl := &RateLimiter{
		buckets:    make(map[string]*TokenBucket),
		rate:       rate,
		window:     window,
		maxBurst:   maxBurst,
		cleanupTTL: 10 * time.Minute,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow checks if a request should be allowed for the given key
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	bucket, exists := rl.buckets[key]
	if !exists {
		bucket = &TokenBucket{
			tokens:       rl.rate + rl.maxBurst,
			lastRefill:   time.Now(),
			lastAccessed: time.Now(),
		}
		rl.buckets[key] = bucket
	}
	rl.mu.Unlock()

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill)
	tokensToAdd := int(elapsed / rl.window * time.Duration(rl.rate))

	if tokensToAdd > 0 {
		bucket.tokens += tokensToAdd
		if bucket.tokens > rl.rate+rl.maxBurst {
			bucket.tokens = rl.rate + rl.maxBurst
		}
		bucket.lastRefill = now
	}

	bucket.lastAccessed = now

	// Check if we have tokens available
	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// Reset removes the bucket for a given key
func (rl *RateLimiter) Reset(key string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.buckets, key)
}

// cleanup removes stale buckets
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, bucket := range rl.buckets {
			bucket.mu.Lock()
			if now.Sub(bucket.lastAccessed) > rl.cleanupTTL {
				delete(rl.buckets, key)
			}
			bucket.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// Global rate limiters for different scopes
var (
	// API rate limiter: 100 requests per minute per user
	APIRateLimiter = NewRateLimiter(100, time.Minute, 20)

	// RCON command rate limiter: 30 commands per minute per user
	RconRateLimiter = NewRateLimiter(30, time.Minute, 10)

	// Login rate limiter: 5 attempts per minute per IP
	LoginRateLimiter = NewRateLimiter(5, time.Minute, 2)
)

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(limiter *RateLimiter, keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFunc(c)

		if !limiter.Allow(key) {
			c.Set("error", "Rate limit exceeded. Please try again later.")
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}

// RateLimitByUser rate limits by user ID
func RateLimitByUser(limiter *RateLimiter) gin.HandlerFunc {
	return RateLimitMiddleware(limiter, func(c *gin.Context) string {
		if userVal, exists := c.Get("user"); exists {
			if user, ok := userVal.(*models.User); ok {
				return "user:" + string(rune(user.ID))
			}
		}
		// Fallback to IP if no user
		return "ip:" + c.ClientIP()
	})
}

// RateLimitByIP rate limits by IP address
func RateLimitByIP(limiter *RateLimiter) gin.HandlerFunc {
	return RateLimitMiddleware(limiter, func(c *gin.Context) string {
		return "ip:" + c.ClientIP()
	})
}

// RateLimitByUserOrIP rate limits by user ID or IP (whichever is available)
func RateLimitByUserOrIP(limiter *RateLimiter) gin.HandlerFunc {
	return RateLimitMiddleware(limiter, func(c *gin.Context) string {
		if userVal, exists := c.Get("user"); exists {
			if user, ok := userVal.(*models.User); ok {
				return "user:" + string(rune(user.ID))
			}
		}
		return "ip:" + c.ClientIP()
	})
}
