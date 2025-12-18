package http

import (
	"net/http"
	"strconv"

	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TribeHandler handles HTTP requests for tribes feature
type TribeHandler struct {
	tribeService ports.TribeService
}

// NewTribeHandler creates a new tribe handler
func NewTribeHandler(tribeService ports.TribeService) *TribeHandler {
	return &TribeHandler{
		tribeService: tribeService,
	}
}

// CreateTribe handles POST /api/v1/tribes
func (h *TribeHandler) CreateTribe(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID).String()

	var req domain.CreateTribeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tribe, err := h.tribeService.CreateTribe(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tribe)
}

// ListTribes handles GET /api/v1/tribes
func (h *TribeHandler) ListTribes(c *gin.Context) {
	// Get current user ID if authenticated (optional for browsing)
	var currentUserID *string
	if userIDVal, exists := c.Get("user_id"); exists {
		userIDStr := userIDVal.(uuid.UUID).String()
		currentUserID = &userIDStr
	}

	// Parse query parameters
	query := domain.ListTribesQuery{
		Search:          c.Query("search"),
		FastingSchedule: c.Query("fasting_schedule"),
		PrimaryGoal:     c.Query("primary_goal"),
		Privacy:         c.Query("privacy"),
		SortBy:          c.DefaultQuery("sort_by", "newest"),
		Limit:           20,
		Offset:          0,
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			query.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			query.Offset = offset
		}
	}

	tribes, total, err := h.tribeService.ListTribes(c.Request.Context(), query, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tribes": tribes,
		"total":  total,
		"limit":  query.Limit,
		"offset": query.Offset,
	})
}

// GetTribe handles GET /api/v1/tribes/:id
func (h *TribeHandler) GetTribe(c *gin.Context) {
	tribeID := c.Param("id")

	// Get current user ID if authenticated (optional)
	var currentUserID *string
	if userIDVal, exists := c.Get("user_id"); exists {
		userIDStr := userIDVal.(uuid.UUID).String()
		currentUserID = &userIDStr
	}

	tribe, err := h.tribeService.GetTribe(c.Request.Context(), tribeID, currentUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tribe not found"})
		return
	}

	c.JSON(http.StatusOK, tribe)
}

// JoinTribe handles POST /api/v1/tribes/:id/join
func (h *TribeHandler) JoinTribe(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID).String()
	tribeID := c.Param("id")

	err := h.tribeService.JoinTribe(c.Request.Context(), tribeID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully joined tribe"})
}

// LeaveTribe handles POST /api/v1/tribes/:id/leave
func (h *TribeHandler) LeaveTribe(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID).String()
	tribeID := c.Param("id")

	err := h.tribeService.LeaveTribe(c.Request.Context(), tribeID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully left tribe"})
}

// GetTribeMembers handles GET /api/v1/tribes/:id/members
func (h *TribeHandler) GetTribeMembers(c *gin.Context) {
	tribeID := c.Param("id")

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

	members, err := h.tribeService.GetTribeMembers(c.Request.Context(), tribeID, limit, offset)
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

// GetMyTribes handles GET /api/v1/users/me/tribes
func (h *TribeHandler) GetMyTribes(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID).String()

	tribes, err := h.tribeService.GetMyTribes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tribes": tribes})
}

// GetTribeStats handles GET /api/v1/tribes/:id/stats
func (h *TribeHandler) GetTribeStats(c *gin.Context) {
	tribeID := c.Param("id")

	stats, err := h.tribeService.GetTribeStats(c.Request.Context(), tribeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// RegisterTribesRoutes registers all tribe-related routes
func RegisterTribesRoutes(router *gin.RouterGroup, handler *TribeHandler, authMiddleware gin.HandlerFunc, optionalAuthMiddleware gin.HandlerFunc) {
	// Public routes (browsing tribes) - use optional auth to detect logged-in users
	publicTribes := router.Group("/tribes")
	publicTribes.Use(optionalAuthMiddleware)
	{
		publicTribes.GET("", handler.ListTribes)
		publicTribes.GET("/:id", handler.GetTribe)
	}

	// Protected routes - require authentication
	tribes := router.Group("/tribes")
	tribes.Use(authMiddleware)
	{
		tribes.POST("", handler.CreateTribe)
		tribes.POST("/:id/join", handler.JoinTribe)
		tribes.POST("/:id/leave", handler.LeaveTribe)
		tribes.GET("/:id/members", handler.GetTribeMembers)
		tribes.GET("/:id/stats", handler.GetTribeStats)
	}

	// User's tribes
	users := router.Group("/users")
	users.Use(authMiddleware)
	{
		users.GET("/me/tribes", handler.GetMyTribes)
	}
}
