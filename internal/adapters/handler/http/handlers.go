package http

import (
	"context"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"fastinghero/internal/core/services"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	authService         ports.AuthService
	fastingService      ports.FastingService
	ketoService         ports.KetoService
	leaderboardService  ports.LeaderboardService
	gamificationService ports.GamificationService
	cortexService       ports.CortexService
	activityService     ports.ActivityService
	telemetryService    ports.TelemetryService
	mealService         ports.MealService
	recipeService       ports.RecipeService
	paymentHandler      *PaymentHandler
	onboardingHandler   *OnboardingHandler
	oauthHandler        *OAuthHandler
	notificationService ports.NotificationService
	socialService       ports.SocialService
	progressService     ports.ProgressService
}

func NewHandler(
	authService ports.AuthService,
	fastingService ports.FastingService,
	ketoService ports.KetoService,
	leaderboardService ports.LeaderboardService,
	gamificationService ports.GamificationService,
	cortexService ports.CortexService,
	activityService ports.ActivityService,
	telemetryService ports.TelemetryService,
	mealService ports.MealService,
	recipeService ports.RecipeService,
	paymentService ports.PaymentService,
	referralService ports.ReferralService,
	notificationService ports.NotificationService,
	socialService ports.SocialService,
	progressService ports.ProgressService,
	userRepo ports.UserRepository,
) *Handler {
	// Note: OAuth handler will be initialized in main.go with OAuth credentials
	return &Handler{
		authService:         authService,
		fastingService:      fastingService,
		ketoService:         ketoService,
		leaderboardService:  leaderboardService,
		gamificationService: gamificationService,
		cortexService:       cortexService,
		activityService:     activityService,
		telemetryService:    telemetryService,
		mealService:         mealService,
		recipeService:       recipeService,
		paymentHandler:      NewPaymentHandler(paymentService, referralService, userRepo),
		onboardingHandler:   NewOnboardingHandler(services.NewOnboardingService(userRepo)),
		oauthHandler:        nil, // Will be set in main.go
		notificationService: notificationService,
		socialService:       socialService,
		progressService:     progressService,
	}
}

// SetOAuthHandler sets the OAuth handler (called from main.go after handler construction)
func (h *Handler) SetOAuthHandler(oauthHandler *OAuthHandler) {
	h.oauthHandler = oauthHandler
}

func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		Name         string `json:"name"`
		ReferralCode string `json:"referral_code"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password, req.Name, req.ReferralCode)
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
	}
	token, refresh, user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "refresh_token": refresh, "user": user})
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)

		// Google OAuth routes
		if h.oauthHandler != nil {
			auth.GET("/google", h.oauthHandler.HandleGoogleLogin)
			auth.GET("/google/callback", h.oauthHandler.HandleGoogleCallback)
		}
	}

	onboarding := api.Group("/onboarding")
	onboarding.Use(AuthMiddleware(h.authService))
	{
		onboarding.PUT("/profile", h.onboardingHandler.UpdateProfile)
		onboarding.POST("/complete", h.onboardingHandler.CompleteOnboarding)
	}

	// Middleware for protected routes
	protected := api.Group("/")
	protected.Use(AuthMiddleware(h.authService))

	// User routes
	user := protected.Group("/user")
	{
		user.GET("/profile", h.GetUserProfile)
		user.GET("/stats", h.GetUserStats)
	}

	fasting := protected.Group("/fasting")
	{
		fasting.POST("/start", h.StartFast)
		fasting.POST("/stop", h.StopFast)
		fasting.GET("/current", h.GetCurrentFast)
		fasting.GET("/insight", h.GetFastingInsight)
	}

	keto := protected.Group("/keto")
	{
		keto.POST("/log", h.LogKeto)
	}

	cortex := protected.Group("/cortex")
	{
		cortex.POST("/chat", h.Chat)
		cortex.POST("/insight", h.GetInsight)
		cortex.POST("/craving-help", h.GetCravingHelp)
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

	progress := protected.Group("/progress")
	{
		progress.POST("/weight", h.LogWeight)
		progress.GET("/weight", h.GetWeightHistory)
		progress.POST("/hydration", h.LogHydration)
		progress.GET("/hydration/daily", h.GetDailyHydration)
	}

	// Social Routes
	social := protected.Group("/social")
	{
		social.POST("/friends/add", h.AddFriend)
		social.GET("/friends", h.GetFriends)
		social.POST("/tribes", h.CreateTribe)
		social.GET("/tribes", h.ListTribes)
		social.GET("/tribes/:id", h.GetTribe)
		social.POST("/challenges", h.CreateChallenge)
		social.GET("/challenges", h.GetChallenges)
		social.GET("/feed", h.GetFeed)
	}

	leaderboardGroup := protected.Group("/leaderboard")
	{
		leaderboardGroup.GET("/", h.GetLeaderboard)
	}

	gamificationGroup := protected.Group("/gamification")
	{
		gamificationGroup.GET("/profile", h.GetGamificationProfile)
	}

	// Payment Routes
	payment := api.Group("/payments")
	{
		payment.POST("/deposit", h.paymentHandler.HandleDeposit)
		payment.POST("/webhook", h.paymentHandler.HandleWebhook)
	}

	// Notification Routes
	notifications := protected.Group("/notifications")
	{
		notifications.POST("/register-token", h.RegisterFCMToken)
		notifications.POST("/unregister-token", h.UnregisterFCMToken)
		notifications.GET("/history", h.GetNotificationHistory)
	}

	// Serve frontend static files
	router.Static("/assets", "./frontend/dist/assets")
	router.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")

	// SPA fallback - serve index.html for all non-API routes
	router.NoRoute(func(c *gin.Context) {
		// Don't serve index.html for API routes
		if !strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.File("./frontend/dist/index.html")
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
		}
	})
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
		c.JSON(http.StatusOK, gin.H{"insight": "Stay hydrated and keep going! (AI insights unavailable)"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"insight": insight})
}

func (h *Handler) GetCravingHelp(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req struct {
		CravingDescription string `json:"craving_description"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cravingHelp, err := h.cortexService.GetCravingHelp(c.Request.Context(), userID, req.CravingDescription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate craving help"})
		return
	}

	c.JSON(http.StatusOK, cravingHelp)
}

func (h *Handler) LogWeight(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req struct {
		Weight float64 `json:"weight"`
		Unit   string  `json:"unit"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log, err := h.progressService.LogWeight(c.Request.Context(), userID, req.Weight, req.Unit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, log)
}

func (h *Handler) GetWeightHistory(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	logs, err := h.progressService.GetWeightHistory(c.Request.Context(), userID, 30) // Default 30 days
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func (h *Handler) LogHydration(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req struct {
		Amount float64 `json:"amount"`
		Unit   string  `json:"unit"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log, err := h.progressService.LogHydration(c.Request.Context(), userID, req.Amount, req.Unit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, log)
}

func (h *Handler) GetDailyHydration(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	log, err := h.progressService.GetDailyHydration(c.Request.Context(), userID)
	if err != nil {
		// Return empty or zero if not found, or error
		c.JSON(http.StatusOK, gin.H{"glasses_count": 0})
		return
	}
	if log == nil {
		c.JSON(http.StatusOK, gin.H{"glasses_count": 0})
		return
	}
	c.JSON(http.StatusOK, log)
}

func (h *Handler) GetLeaderboard(c *gin.Context) {
	ctx := c.Request.Context()
	// Check if tribe_id query param is present
	tribeIDStr := c.Query("tribe_id")
	if tribeIDStr != "" {
		// Tribe Leaderboard logic here if needed
	}

	leaderboard, err := h.leaderboardService.GetGlobalLeaderboard(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}
	c.JSON(http.StatusOK, leaderboard)
}

func (h *Handler) GetGamificationProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	streak, badges, err := h.gamificationService.GetUserGamificationProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"streak": streak,
		"badges": badges,
	})
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

	// Trigger Gamification Updates (Async or Sync)
	// Ideally async, but sync for now for simplicity
	go func() {
		ctx := context.Background() // Use background context for async
		h.gamificationService.UpdateStreak(ctx, userID)
		h.gamificationService.CheckAndAwardBadges(ctx, userID, "fast_completed", session)
	}()

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
		Name        string `json:"name"`
		Calories    int    `json:"calories"`
		MealType    string `json:"meal_type"`
		Image       string `json:"image"`
		Description string `json:"description"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meal, err := h.mealService.LogMeal(c.Request.Context(), userID, req.Name, req.Calories, req.MealType, req.Image, req.Description)
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

func (h *Handler) RegisterFCMToken(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req struct {
		Token      string `json:"token"`
		DeviceType string `json:"device_type"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.DeviceType == "" {
		req.DeviceType = "web"
	}

	err := h.notificationService.RegisterFCMToken(c.Request.Context(), userID, req.Token, req.DeviceType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token registered successfully"})
}

func (h *Handler) UnregisterFCMToken(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req struct {
		Token string `json:"token"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.notificationService.UnregisterFCMToken(c.Request.Context(), userID, req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unregister token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token unregistered successfully"})
}

func (h *Handler) GetNotificationHistory(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	notifications, err := h.notificationService.GetHistory(c.Request.Context(), userID, 50) // Default limit 50
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

func (h *Handler) GetUserProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	user, err := h.authService.(*services.AuthService).GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUserStats(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	user, err := h.authService.(*services.AuthService).GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	history, _ := h.fastingService.GetFastingHistory(c.Request.Context(), userID)
	streak, _, _ := h.gamificationService.GetUserGamificationProfile(c.Request.Context(), userID)

	fastsCompleted := 0
	totalHours := 0.0
	for _, f := range history {
		if f.Status == domain.StatusCompleted {
			fastsCompleted++
			if f.EndTime != nil {
				totalHours += f.EndTime.Sub(f.StartTime).Hours()
			}
		}
	}

	currentStreak := 0
	longestStreak := 0
	if streak != nil {
		currentStreak = streak.CurrentStreak
		longestStreak = streak.LongestStreak
	}

	c.JSON(http.StatusOK, gin.H{
		"fasts_completed":     fastsCompleted,
		"total_fasting_hours": totalHours,
		"current_streak":      currentStreak,
		"longest_streak":      longestStreak,
		"vault_balance":       user.EarnedRefund,
		"vault_total":         user.VaultDeposit,
	})
}

// GetFastingInsight returns AI-generated insights about the user's current fast
func (h *Handler) GetFastingInsight(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	// Get hours from query parameter
	hoursStr := c.Query("hours")
	if hoursStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hours parameter required"})
		return
	}

	var hours float64
	if _, err := fmt.Sscanf(hoursStr, "%f", &hours); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hours parameter"})
		return
	}

	// Call cortex service for insight
	insight, err := h.cortexService.(*services.CortexService).GetFastingMilestoneInsight(c.Request.Context(), userID, hours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate insight"})
		return
	}

	c.JSON(http.StatusOK, insight)
}
