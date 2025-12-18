package http

import (
	"fastinghero/internal/core/ports"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService ports.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]
		user, err := authService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Set user in context
		c.Set("user_id", user.ID)
		c.Set("user", user)
		c.Next()
	}
}

// OptionalAuthMiddleware extracts user info from token if present, but doesn't require auth
// This is useful for public routes that may want to personalize content for logged-in users
func OptionalAuthMiddleware(authService ports.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		user, err := authService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			// Token invalid, but don't fail - just continue without user context
			c.Next()
			return
		}

		// Set user in context
		c.Set("user_id", user.ID)
		c.Set("user", user)
		c.Next()
	}
}
