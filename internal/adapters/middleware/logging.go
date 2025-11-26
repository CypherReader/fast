package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// RequestLogger logs HTTP requests with structured logging
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get user ID if authenticated
		var userID string
		if uid, exists := c.Get("user_id"); exists {
			userID = uid.(uuid.UUID).String()
		}

		// Build log entry
		logEvent := log.Info().
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("ip", c.ClientIP()).
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Int("body_size", c.Writer.Size())

		if raw != "" {
			logEvent = logEvent.Str("query", raw)
		}

		if userID != "" {
			logEvent = logEvent.Str("user_id", userID)
		}

		// Log errors if any
		if len(c.Errors) > 0 {
			errStrings := make([]string, len(c.Errors))
			for i, err := range c.Errors {
				errStrings[i] = err.Error()
			}
			logEvent = logEvent.Strs("errors", errStrings)
		}

		// Determine log level based on status code
		status := c.Writer.Status()
		if status >= 500 {
			logEvent.Msg("Server error")
		} else if status >= 400 {
			logEvent.Msg("Client error")
		} else {
			logEvent.Msg("Request completed")
		}
	}
}
