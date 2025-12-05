package http

import (
	"fastinghero/internal/core/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) AddFriend(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req struct {
		FriendID string `json:"friend_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	friendUUID, err := uuid.Parse(req.FriendID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid friend id"})
		return
	}

	err = h.socialService.AddFriend(c.Request.Context(), userID, friendUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request sent"})
}

func (h *Handler) GetFriends(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	friends, err := h.socialService.GetFriends(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, friends)
}

func (h *Handler) CreateTribe(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsPublic    bool   `json:"is_public"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tribe, err := h.socialService.CreateTribe(c.Request.Context(), userID, req.Name, req.Description, req.IsPublic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tribe)
}

func (h *Handler) GetTribe(c *gin.Context) {
	id := c.Param("id")
	tribeUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tribe id"})
		return
	}

	tribe, err := h.socialService.GetTribe(c.Request.Context(), tribeUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tribe not found"})
		return
	}

	c.JSON(http.StatusOK, tribe)
}

func (h *Handler) CreateChallenge(c *gin.Context) {
	// userIDVal, exists := c.Get("user_id")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }
	// userID := userIDVal.(uuid.UUID)

	var req struct {
		Name          string               `json:"name"`
		ChallengeType domain.ChallengeType `json:"challenge_type"`
		Goal          int                  `json:"goal"`
		StartDate     string               `json:"start_date"` // ISO string
		EndDate       string               `json:"end_date"`   // ISO string
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dates
	// ... (simplified for now, assuming ISO string parsing in service or here)
	// Using time.Parse(time.RFC3339, ...)

	// For now, let's skip date parsing complexity and assume service handles it or just pass zero time if not critical for this step
	// But service expects time.Time.

	// TODO: Implement date parsing
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Challenge creation not fully implemented"})
}

func (h *Handler) GetChallenges(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	challenges, err := h.socialService.GetChallenges(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, challenges)
}

func (h *Handler) ListTribes(c *gin.Context) {
	// Pagination params (simplified)
	limit := 20
	offset := 0

	tribes, err := h.socialService.ListTribes(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tribes)
}

func (h *Handler) GetFeed(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	// Pagination params (simplified)
	limit := 20
	offset := 0

	feed, err := h.socialService.GetFeed(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feed)
}
