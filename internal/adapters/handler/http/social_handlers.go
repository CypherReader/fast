package http

import (
\t"fastinghero/internal/core/domain"
\t"net/http"
\t"strconv"

\t"github.com/gin-gonic/gin"
\t"github.com/google/uuid"
)

func (h *Handler) AddFriend(c *gin.Context) {
\tuserIDVal, exists := c.Get("user_id")
\tif !exists {
\t\tc.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
\t\treturn
\t}
\tuserID := userIDVal.(uuid.UUID)

\tvar req struct {
\t\tFriendID string `json:"friend_id"`
\t}
\tif err := c.BindJSON(&req); err != nil {
\t\tc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
\t\treturn
\t}

\tfriendUUID, err := uuid.Parse(req.FriendID)
\tif err != nil {
\t\tc.JSON(http.StatusBadRequest, gin.H{"error": "invalid friend id"})
\t\treturn
\t}

\terr = h.socialService.AddFriend(c.Request.Context(), userID, friendUUID)
\tif err != nil {
\t\tc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
\t\treturn
\t}

\tc.JSON(http.StatusOK, gin.H{"message": "Friend request sent"})
}

func (h *Handler) GetFriends(c *gin.Context) {
\tuserIDVal, exists := c.Get("user_id")
\tif !exists {
\t\tc.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
\t\treturn
\t}
\tuserID := userIDVal.(uuid.UUID)

\tfriends, err := h.socialService.GetFriends(c.Request.Context(), userID)
\tif err != nil {
\t\tc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
\t\treturn
\t}

\tc.JSON(http.StatusOK, friends)
}

func (h *Handler) CreateTribe(c *gin.Context) {
\tuserIDVal, exists := c.Get("user_id")
\tif !exists {
\t\tc.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
\t\treturn
\t}
\tuserID := userIDVal.(uuid.UUID)

\tvar req struct {
\t\tName        string `json:"name"`
\t\tDescription string `json:"description"`
\t\tIsPublic    bool   `json:"is_public"`
\t}
\tif err := c.BindJSON(&req); err != nil {
\t\tc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
\t\treturn
\t}

\ttribe, err := h.socialService.CreateTribe(c.Request.Context(), userID, req.Name, req.Description, req.IsPublic)
\tif err != nil {
\t\tc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
\t\treturn
\t}

\tc.JSON(http.StatusCreated, tribe)
}

func (h *Handler) GetTribe(c *gin.Context) {
\tid := c.Param("id")
\ttribeUUID, err := uuid.Parse(id)
\tif err != nil {
\t\tc.JSON(http.StatusBadRequest, gin.H{"error": "invalid tribe id"})
\t\treturn
\t}

\ttribe, err := h.socialService.GetTribe(c.Request.Context(), tribeUUID)
\tif err != nil {
\t\tc.JSON(http.StatusNotFound, gin.H{"error": "tribe not found"})
\t\treturn
\t}

\tc.JSON(http.StatusOK, tribe)
}

func (h *Handler) ListTribes(c *gin.Context) {
\t// Parse query parameters with defaults
\tlimit := 20
\toffset := 0

\tif limitStr := c.Query("limit"); limitStr != "" {
\t\tif l, err := strconv.Atoi(limitStr); err == nil {
\t\t\tlimit = l
\t\t}
\t}

\tif offsetStr := c.Query("offset"); offsetStr != "" {
\t\tif o, err := strconv.Atoi(offsetStr); err == nil {
\t\t\toffset = o
\t\t}
\t}

\ttribes, err := h.socialService.ListTribes(c.Request.Context(), limit, offset)
\tif err != nil {
\t\tc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
\t\treturn
\t}

\t// Return response in format expected by frontend
\tc.JSON(http.StatusOK, gin.H{
\t\t"tribes": tribes,
\t\t"total":  len(tribes), // TODO: Get actual total from service
\t\t"limit":  limit,
\t\t"offset": offset,
\t})
}

// JoinTribe handles joining a tribe
func (h *Handler) JoinTribe(c *gin.Context) {
\tuserIDVal, exists := c.Get("user_id")
\tif !exists {
\t\tc.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
\t\treturn
\t}
\tuserID := userIDVal.(uuid.UUID)
\ttribeID := c.Param("id")

\ttribeUUID, err := uuid.Parse(tribeID)
\tif err != nil {
\t\tc.JSON(http.StatusBadRequest, gin.H{"error": "invalid tribe id"})
\t\treturn
\t}

\terr = h.socialService.JoinTribe(c.Request.Context(), userID, tribeUUID)
\tif err != nil {
\t\tc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
\t\treturn
\t}

\tc.JSON(http.StatusOK, gin.H{"message": "successfully joined tribe"})
}

// LeaveTribe handles leaving a tribe
func (h *Handler) LeaveTribe(c *gin.Context) {
\tuserIDVal, exists := c.Get("user_id")
\tif !exists {
\t\tc.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
\t\treturn
\t}
\tuserID := userIDVal.(uuid.UUID)
\ttribeID := c.Param("id")

\ttribeUUID, err := uuid.Parse(tribeID)
\tif err != nil {
\t\tc.JSON(http.StatusBadRequest, gin.H{"error": "invalid tribe id"})
\t\treturn
\t}

\terr = h.socialService.LeaveTribe(c.Request.Context(), userID, tribeUUID)
\tif err != nil {
\t\tc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
\t\treturn
\t}

\tc.JSON(http.StatusOK, gin.H{"message": "successfully left tribe"})
}

// GetTribeMembers handles getting members of a tribe
func (h *Handler) GetTribeMembers(c *gin.Context) {
\ttribeID := c.Param("id")
\ttribeUUID, err := uuid.Parse(tribeID)
\tif err != nil {
\t\tc.JSON(http.StatusBadRequest, gin.H{"error": "invalid tribe id"})
\t\treturn
\t}

\tlimit := 20
\toffset := 0

\tif limitStr := c.Query("limit"); limitStr != "" {
\t\tif l, err := strconv.Atoi(limitStr); err == nil {
\t\t\tlimit = l
\t\t}
\t}

\tif offsetStr := c.Query("offset"); offsetStr != "" {
\t\tif o, err := strconv.Atoi(offsetStr); err == nil {
\t\t\toffset = o
\t\t}
\t}

\tmembers, err := h.socialService.GetTribeMembers(c.Request.Context(), tribeUUID, limit, offset)
\tif err != nil {
\t\tc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
\t\treturn
\t}

\tc.JSON(http.StatusOK, gin.H{
\t\t"members": members,
\t\t"limit":   limit,
\t\t"offset":  offset,
\t})
}

// GetTribeStats handles getting statistics for a tribe
func (h *Handler) GetTribeStats(c *gin.Context) {
\ttribeID := c.Param("id")
\ttribeUUID, err := uuid.Parse(tribeID)
\tif err != nil {
\t\tc.JSON(http.StatusBadRequest, gin.H{"error": "invalid tribe id"})
\t\treturn
\t}

\tstats, err := h.socialService.GetTribeStats(c.Request.Context(), tribeUUID)
\tif err != nil {
\t\tc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
\t\treturn
\t}

\tc.JSON(http.StatusOK, stats)
}

// GetMyTribes handles getting the current user's tribes
func (h *Handler) GetMyTribes(c *gin.Context) {
\tuserIDVal, exists := c.Get("user_id")
\tif !exists {
\t\tc.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
\t\treturn
\t}
\tuserID := userIDVal.(uuid.UUID)

\ttribes, err := h.socialService.GetMyTribes(c.Request.Context(), userID)
\tif err != nil {
\t\tc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
\t\treturn
\t}

\tc.JSON(http.StatusOK, gin.H{"tribes": tribes})
}

func (h *Handler) CreateChallenge(c *gin.Context) {
\t// userIDVal, exists := c.Get("user_id")
\t// if !exists {
\t// \tc.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
\t// \treturn
\t// }
\t// userID := userIDVal.(uuid.UUID)

\tvar req struct {
\t\tName          string               `json:"name"`
\t\tChallengeType domain.ChallengeType `json:"challenge_type"`
\t\tGoal          int                  `json:"goal"`
\t\tStartDate     string               `json:"start_date"` // ISO string
\t\tEndDate       string               `json:"end_date"`   // ISO string
\t}
\tif err := c.BindJSON(&req); err != nil {
\t\tc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
\t\treturn
\t}

\t// Parse dates
\t// ... (simplified for now, assuming ISO string parsing in service or here)
\t// Using time.Parse(time.RFC3339, ...)

\t// For now, let's skip date parsing complexity and assume service handles it or just pass zero time if not critical for this step
\t// But service expects time.Time.

\t// TODO: Implement date parsing
\tc.JSON(http.StatusNotImplemented, gin.H{"error": "Challenge creation not fully implemented"})
}

func (h *Handler) GetChallenges(c *gin.Context) {
\tuserIDVal, exists := c.Get("user_id")
\tif !exists {
\t\tc.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
\t\treturn
\t}
\tuserID := userIDVal.(uuid.UUID)

\tchallenges, err := h.socialService.GetChallenges(c.Request.Context(), userID)
\tif err != nil {
\t\tc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
\t\treturn
\t}

\tc.JSON(http.StatusOK, challenges)
}

func (h *Handler) GetFeed(c *gin.Context) {
\tuserIDVal, exists := c.Get("user_id")
\tif !exists {
\t\tc.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
\t\treturn
\t}
\tuserID := userIDVal.(uuid.UUID)

\t// Pagination params (simplified)
\tlimit := 20
\toffset := 0

\tfeed, err := h.socialService.GetFeed(c.Request.Context(), userID, limit, offset)
\tif err != nil {
\t\tc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
\t\treturn
\t}

\tc.JSON(http.StatusOK, feed)
}
