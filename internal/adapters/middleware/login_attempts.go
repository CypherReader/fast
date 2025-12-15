package middleware

import (
	"sync"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

// AttemptInfo stores information about login attempts for an email
type AttemptInfo struct {
	Count       int
	LastAttempt time.Time
	LockedUntil time.Time
}

// LoginAttemptTracker tracks failed login attempts and implements account lockout
type LoginAttemptTracker struct {
	mu       sync.RWMutex
	attempts map[string]*AttemptInfo
}

// NewLoginAttemptTracker creates a new login attempt tracker
func NewLoginAttemptTracker() *LoginAttemptTracker {
	tracker := &LoginAttemptTracker{
		attempts: make(map[string]*AttemptInfo),
	}

	// Start cleanup goroutine
	go tracker.cleanupOldAttempts()

	return tracker
}

// RecordFailedAttempt records a failed login attempt and returns false if account is locked
func (t *LoginAttemptTracker) RecordFailedAttempt(email string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	info, exists := t.attempts[email]
	if !exists {
		info = &AttemptInfo{Count: 0}
		t.attempts[email] = info
	}

	// Check if account is locked
	if time.Now().Before(info.LockedUntil) {
		return false // Account is locked
	}

	info.Count++
	info.LastAttempt = time.Now()

	// Lock account after 5 failed attempts for 15 minutes
	if info.Count >= 5 {
		info.LockedUntil = time.Now().Add(15 * time.Minute)
		return false
	}

	return true // Attempt allowed
}

// RecordSuccessfulLogin resets the attempt counter for an email
func (t *LoginAttemptTracker) RecordSuccessfulLogin(email string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.attempts, email)
}

// IsLocked checks if an account is currently locked
func (t *LoginAttemptTracker) IsLocked(email string) (bool, time.Time) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	info, exists := t.attempts[email]
	if !exists {
		return false, time.Time{}
	}

	if time.Now().Before(info.LockedUntil) {
		return true, info.LockedUntil
	}

	return false, time.Time{}
}

// cleanupOldAttempts removes old attempt records every hour
func (t *LoginAttemptTracker) cleanupOldAttempts() {
	for {
		time.Sleep(time.Hour)

		t.mu.Lock()
		for email, info := range t.attempts {
			// Remove attempts older than 24 hours that aren't locked
			if time.Since(info.LastAttempt) > 24*time.Hour && time.Now().After(info.LockedUntil) {
				delete(t.attempts, email)
			}
		}
		t.mu.Unlock()
	}
}

// LoginAttemptMiddleware creates a middleware that checks for account lockout
func LoginAttemptMiddleware(tracker *LoginAttemptTracker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This middleware should only be applied to login endpoints
		// The actual checking happens in the login handler
		// We store the tracker in context for the handler to use
		c.Set("login_tracker", tracker)
		c.Next()
	}
}

// CheckAccountLocked is a helper to check if an account is locked and return appropriate response
func CheckAccountLocked(c *gin.Context, email string) bool {
	trackerVal, exists := c.Get("login_tracker")
	if !exists {
		return false
	}

	tracker, ok := trackerVal.(*LoginAttemptTracker)
	if !ok {
		return false
	}

	locked, until := tracker.IsLocked(email)
	if locked {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":        "Account temporarily locked due to too many failed login attempts",
			"locked_until": until.Format(time.RFC3339),
		})
		return true
	}

	return false
}
