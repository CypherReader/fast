package http

import (
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"fastinghero/internal/core/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	authService     ports.AuthService
	fastingService  ports.FastingService
	ketoService     ports.KetoService
	socialService   *services.SocialService
	cortexService   ports.CortexService
	activityService ports.ActivityService
}

func NewHandler(auth ports.AuthService, fasting ports.FastingService, keto ports.KetoService, social *services.SocialService, cortex ports.CortexService, activity ports.ActivityService) *Handler {
	return &Handler{
		authService:     auth,
		fastingService:  fasting,
		ketoService:     keto,
		socialService:   social,
		cortexService:   cortex,
		activityService: activity,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}

	// Middleware for protected routes
	protected := api.Group("/")
	protected.Use(AuthMiddleware(h.authService.(*services.AuthService)))

	fasting := protected.Group("/fasting")
	{
		fasting.POST("/start", h.StartFast)
		fasting.POST("/stop", h.StopFast)
		fasting.GET("/current", h.GetCurrentFast)
	}

	keto := protected.Group("/keto")
	{
		keto.POST("/log", h.LogKeto)
	}

	social := protected.Group("/social")
	{
		social.GET("/feed", h.GetFeed)
	}

	cortex := protected.Group("/cortex")
	{
		cortex.POST("/chat", h.Chat)
		cortex.POST("/insight", h.GetInsight)
	}

	activity := protected.Group("/activity")
	{
		activity.POST("/sync", h.SyncActivity)
		activity.GET("/", h.GetActivities)
		activity.GET("/:id", h.GetActivity)
	}
}

func (h *Handler) Chat(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req struct {
		Message string `json:"message"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.cortexService.Chat(c.Request.Context(), userID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

func (h *Handler) GetInsight(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req struct {
		FastingHours float64 `json:"fasting_hours"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insight, err := h.cortexService.GenerateInsight(c.Request.Context(), userID, req.FastingHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"insight": insight})
}

func (h *Handler) GetFeed(c *gin.Context) {
	feed, err := h.socialService.GetFeed(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, feed)
}

func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, refresh, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "refresh_token": refresh})
}

func (h *Handler) StartFast(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req struct {
		PlanType  domain.FastingPlanType `json:"plan_type"`
		GoalHours int                    `json:"goal_hours"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	session, err := h.fastingService.StartFast(c.Request.Context(), userID, req.PlanType, req.GoalHours)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, session)
}

func (h *Handler) StopFast(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)
	session, err := h.fastingService.StopFast(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, session)
}

func (h *Handler) GetCurrentFast(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)
	session, err := h.fastingService.GetCurrentFast(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no active session"})
		return
	}
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no active session"})
		return
	}
	c.JSON(http.StatusOK, session)
}

func (h *Handler) LogKeto(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)
	var req domain.KetoEntry
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.ketoService.LogEntry(c.Request.Context(), userID, req)
	if err != nil {
		if err.Error() == "premium subscription required for hard data inputs" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "logged"})
}

func (h *Handler) SyncActivity(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req domain.Activity
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.activityService.SyncActivity(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "synced"})
}

func (h *Handler) GetActivities(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	activities, err := h.activityService.GetActivities(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activities)
}

func (h *Handler) GetActivity(c *gin.Context) {
	id := c.Param("id")
	activity, err := h.activityService.GetActivity(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "activity not found"})
		return
	}
	c.JSON(http.StatusOK, activity)
}
