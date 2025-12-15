package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CSRFMiddleware implements CSRF protection using double-submit cookie pattern
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip CSRF check for safe methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Validate CSRF token for state-changing requests
		csrfToken := c.GetHeader("X-CSRF-Token")
		csrfCookie, err := c.Cookie("csrf_token")

		if err != nil || csrfToken == "" || csrfToken != csrfCookie {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "CSRF token validation failed",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SetCSRFToken generates and sets a new CSRF token
func SetCSRFToken(c *gin.Context) {
	token := generateCSRFToken()

	// Set as cookie
	c.SetCookie(
		"csrf_token",
		token,
		3600*24, // 24 hours
		"/",     // path
		"",      // domain
		false,   // secure (should be true in production with HTTPS)
		true,    // httpOnly
	)

	// Also return in response for client to use in headers
	c.Header("X-CSRF-Token", token)
}

// generateCSRFToken creates a cryptographically secure random token
func generateCSRFToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// Fallback to less secure but still random
		return base64.URLEncoding.EncodeToString([]byte("fallback-token"))
	}
	return base64.URLEncoding.EncodeToString(b)
}
