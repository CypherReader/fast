package http

import (
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"fastinghero/internal/core/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	authService      ports.AuthService
	fastingService   ports.FastingService
	ketoService      ports.KetoService
	socialService    *services.SocialService
	cortexService    ports.CortexService
	activityService  ports.ActivityService
	telemetryService ports.TelemetryService
	mealService      ports.MealService
	recipeService    ports.RecipeService
}

func NewHandler(authService ports.AuthService, fastingService ports.FastingService, ketoService ports.KetoService, socialService *services.SocialService, cortexService ports.CortexService, activityService ports.ActivityService, telemetryService ports.TelemetryService, mealService ports.MealService, recipeService ports.RecipeService) *Handler {
	return &Handler{
		authService:      authService,
		fastingService:   fastingService,
		ketoService:      ketoService,
		socialService:    socialService,
		cortexService:    cortexService,
		activityService:  activityService,
		telemetryService: telemetryService,
		mealService:      mealService,
		recipeService:    recipeService,
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

	telemetry := protected.Group("/telemetry")
	{
		telemetry.POST("/connect", h.ConnectDevice)
		telemetry.POST("/sync", h.SyncTelemetry)
		telemetry.GET("/status", h.GetTelemetryStatus)
		telemetry.POST("/manual", h.LogManualTelemetry)
		telemetry.GET("/metric", h.GetLatestMetric)
		telemetry.GET("/weekly", h.GetWeeklyStats)
	}

	meals := protected.Group("/meals")
	{
		meals.POST("/", h.LogMeal)
		meals.GET("/", h.GetMeals)
	}

	recipes := protected.Group("/recipes")
	{
		recipes.GET("/", h.GetRecipes)
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
		StartTime *time.Time             `json:"start_time"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	session, err := h.fastingService.StartFast(c.Request.Context(), userID, req.PlanType, req.GoalHours, req.StartTime)
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

func (h *Handler) ConnectDevice(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req struct {
		Source domain.TelemetrySource `json:"source"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := h.telemetryService.ConnectDevice(c.Request.Context(), userID, req.Source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, conn)
}

func (h *Handler) SyncTelemetry(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req struct {
		Source domain.TelemetrySource `json:"source"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.telemetryService.SyncData(c.Request.Context(), userID, req.Source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "synced"})
}

func (h *Handler) GetTelemetryStatus(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	status, err := h.telemetryService.GetDeviceStatus(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

func (h *Handler) LogManualTelemetry(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req struct {
		Type  domain.MetricType `json:"type"`
		Value float64           `json:"value"`
		Unit  string            `json:"unit"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.telemetryService.LogManualData(c.Request.Context(), userID, req.Type, req.Value, req.Unit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, data)
}

func (h *Handler) GetLatestMetric(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	metricType := c.Query("type")
	if metricType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "metric type required"})
		return
	}

	data, err := h.telemetryService.GetLatestMetric(c.Request.Context(), userID, domain.MetricType(metricType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no data found"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetWeeklyStats(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	metricType := c.Query("type")
	if metricType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "metric type required"})
		return
	}

	stats, err := h.telemetryService.GetWeeklyStats(c.Request.Context(), userID, domain.MetricType(metricType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) LogMeal(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req struct {
		Image       string `json:"image"`
		Description string `json:"description"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meal, err := h.mealService.LogMeal(c.Request.Context(), userID, req.Image, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, meal)
}

func (h *Handler) GetMeals(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	meals, err := h.mealService.GetMeals(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, meals)
}

func (h *Handler) GetRecipes(c *gin.Context) {
	diet := c.Query("diet")
	recipes, err := h.recipeService.GetRecipes(c.Request.Context(), domain.DietType(diet))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipes)
}
