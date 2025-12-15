package http

import (
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OnboardingHandler struct {
	onboardingService ports.OnboardingService
}

func NewOnboardingHandler(onboardingService ports.OnboardingService) *OnboardingHandler {
	return &OnboardingHandler{
		onboardingService: onboardingService,
	}
}

func (h *OnboardingHandler) UpdateProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Secure type assertion with check to prevent panic
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		// Log the incident for security monitoring
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var req domain.UserProfileUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		// Don't expose internal error details to client
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	user, err := h.onboardingService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		// Don't expose internal errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *OnboardingHandler) CompleteOnboarding(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Secure type assertion with check to prevent panic
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	user, err := h.onboardingService.CompleteOnboarding(c.Request.Context(), userID)
	if err != nil {
		// Don't expose internal errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete onboarding"})
		return
	}

	c.JSON(http.StatusOK, user)
}
