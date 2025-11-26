package http

import (
	"fastinghero/internal/core/ports"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TribeHandler struct {
	service ports.TribeService
}

func NewTribeHandler(service ports.TribeService) *TribeHandler {
	return &TribeHandler{service: service}
}

func (h *TribeHandler) CreateTribe(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	tribe, err := h.service.CreateTribe(c.Request.Context(), req.Name, req.Description, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tribe)
}

func (h *TribeHandler) ListTribes(c *gin.Context) {
	tribes, err := h.service.ListTribes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tribes)
}

func (h *TribeHandler) JoinTribe(c *gin.Context) {
	tribeIDStr := c.Param("id")
	tribeID, err := uuid.Parse(tribeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tribe id"})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.service.JoinTribe(c.Request.Context(), tribeID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "joined tribe successfully"})
}

func (h *TribeHandler) LeaveTribe(c *gin.Context) {
	tribeIDStr := c.Param("id")
	tribeID, err := uuid.Parse(tribeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tribe id"})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.service.LeaveTribe(c.Request.Context(), tribeID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "left tribe successfully"})
}
