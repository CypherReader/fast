package http

import (
	"fastinghero/internal/core/domain"
	"net/http"
	"strconv"

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

func (h *Handler) ListTribes(c *gin.Context) {
	// Parse query parameters with defaults
	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	tribes, err := h.socialService.ListTribes(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response in format expected by frontend
	c.JSON(http.StatusOK, gin.H{
		"tribes": tribes,
		"total":  len(tribes), // TODO: Get actual total from service
		"limit":  limit,
		"offset": offset,
	})
}

// JoinTribe handles joining a tribe
func (h *Handler) JoinTribe(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)
	tribeID := c.Param("id")

	tribeUUID, err := uuid.Parse(tribeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tribe id"})
		return
	}

	err = h.socialService.JoinTribe(c.Request.Context(), userID, tribeUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully joined tribe"})
}

// LeaveTribe handles leaving a tribe
func (h *Handler) LeaveTribe(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)
	tribeID := c.Param("id")

	tribeUUID, err := uuid.Parse(tribeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tribe id"})
		return
	}

	err = h.socialService.LeaveTribe(c.Request.Context(), userID, tribeUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully left tribe"})
}

// GetTribeMembers handles getting members of a tribe
func (h *Handler) GetTribeMembers(c *gin.Context) {
	tribeID := c.Param("id")
	tribeUUID, err := uuid.Parse(tribeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tribe id"})
		return
	}

	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	members, err := h.socialService.GetTribeMembers(c.Request.Context(), tribeUUID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"members": members,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetTribeStats handles getting statistics for a tribe
func (h *Handler) GetTribeStats(c *gin.Context) {
	tribeID := c.Param("id")
	tribeUUID, err := uuid.Parse(tribeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tribe id"})
		return
	}

	stats, err := h.socialService.GetTribeStats(c.Request.Context(), tribeUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetMyTribes handles getting the current user's tribes
func (h *Handler) GetMyTribes(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	tribes, err := h.socialService.GetMyTribes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tribes": tribes})
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
