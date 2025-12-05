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
	userID := userIDVal.(uuid.UUID)

	var req domain.UserProfileUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.onboardingService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	userID := userIDVal.(uuid.UUID)

	user, err := h.onboardingService.CompleteOnboarding(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
