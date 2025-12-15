package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter manages rate limiting for different clients
type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter with the specified rate and burst
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

// GetLimiter returns a rate limiter for the given key (typically IP address)
func (rl *RateLimiter) GetLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.visitors[key]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[key] = limiter
	}

	return limiter
}

// CleanupOldVisitors removes inactive visitors periodically
func (rl *RateLimiter) CleanupOldVisitors() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			// In a production app, you'd track last access time
			// For simplicity, we're just clearing the map periodically
			if len(rl.visitors) > 10000 {
				rl.visitors = make(map[string]*rate.Limiter)
			}
			rl.mu.Unlock()
		}
	}()
}

// Middleware returns a Gin middleware handler for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use IP address as the key
		key := c.ClientIP()

		// Check if user is authenticated - could apply different limits
		userID, authenticated := c.Get("user_id")
		if authenticated {
			key = userID.(string) // Use user ID for authenticated users
		}

		limiter := rl.GetLimiter(key)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthRateLimiter creates a strict rate limiter for authentication endpoints
// Allows only 5 requests per minute to prevent brute-force attacks
func AuthRateLimiter() gin.HandlerFunc {
	limiter := NewRateLimiter(rate.Limit(5.0/60.0), 5)
	limiter.CleanupOldVisitors()
	return limiter.Middleware()
}
