package main

import (
	"context"
	"database/sql"
	"fastinghero/internal/adapters/handler/http"
	"fastinghero/internal/adapters/middleware"
	"fastinghero/internal/adapters/payment"
	"fastinghero/internal/adapters/repository/memory"
	"fastinghero/internal/adapters/repository/postgres"
	"fastinghero/internal/adapters/secondary/llm"
	"fastinghero/internal/core/ports"
	"fastinghero/internal/core/services"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"golang.org/x/time/rate"

	"fastinghero/pkg/logger"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize logger
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logger.Init(logLevel)
	logger.Info().Msg("Starting FastingHero application...")

	var userRepo ports.UserRepository
	var fastingRepo ports.FastingRepository
	var ketoRepo ports.KetoRepository
	var leaderboardRepo ports.LeaderboardRepository
	var gamificationRepo ports.GamificationRepository
	var referralRepo ports.ReferralRepository
	var notificationRepo ports.NotificationRepository
	var subscriptionRepo ports.SubscriptionRepository
	var vaultRepo ports.VaultRepository
	var socialRepo ports.SocialRepository
	var progressRepo ports.ProgressRepository
	var tribeRepo ports.TribeRepository

	// Check for DB connection string
	// Priority: DATABASE_URL (Cloud Run) > DSN > individual env vars
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("DSN")
	}
	if dsn == "" {
		// Default to Postgres DSN if not set, or use a specific env var
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		dbname := os.Getenv("POSTGRES_DB")
		if user != "" && password != "" && dbname != "" {
			dsn = "postgres://" + user + ":" + password + "@" + host + ":5432/" + dbname + "?sslmode=disable"
		}
	}

	// Log connection attempt (without sensitive data)
	if dsn != "" {
		if strings.Contains(dsn, "/cloudsql/") {
			logger.Info().Msg("Database connection: Cloud SQL Unix socket detected")
		} else {
			logger.Info().Msg("Database connection: TCP connection detected")
		}
	}

	var db *sql.DB
	var err error
	useMemory := false

	if dsn != "" {
		logger.Info().Msg("Connecting to PostgreSQL...")
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to open DB connection. Switching to IN-MEMORY mode.")
			useMemory = true
		} else {
			logger.Info().Msg("Database connection opened, testing connectivity...")
			if err := db.Ping(); err != nil {
				logger.Error().Err(err).Msg("Failed to ping DB. Switching to IN-MEMORY mode.")
				useMemory = true
			} else {
				logger.Info().Msg("Database connectivity verified successfully")
			}
		}
	} else {
		logger.Info().Msg("No database configuration found. Using IN-MEMORY mode.")
		useMemory = true
	}

	if !useMemory {
		// Run Migrations
		log.Println("Running database migrations...")
		if err := postgres.RunMigrations(db, "migrations"); err != nil {
			log.Printf("Failed to run migrations: %v. Switching to IN-MEMORY mode.", err)
			useMemory = true
		}
	}

	if !useMemory {
		userRepo = postgres.NewPostgresUserRepository(db)
		fastingRepo = postgres.NewPostgresFastingRepository(db)
		ketoRepo = postgres.NewPostgresKetoRepository(db)
		leaderboardRepo = postgres.NewPostgresLeaderboardRepository(db)
		gamificationRepo = postgres.NewPostgresGamificationRepository(db)
		referralRepo = postgres.NewPostgresReferralRepository(db)
		notificationRepo = postgres.NewPostgresNotificationRepository(db)
		subscriptionRepo = postgres.NewPostgresSubscriptionRepository(db)
		vaultRepo = postgres.NewPostgresVaultRepository(db)
		socialRepo = postgres.NewPostgresSocialRepository(db)
		progressRepo = postgres.NewPostgresProgressRepository(db)
		tribeRepo = postgres.NewPostgresTribeRepository(db)
	} else {
		log.Println("!!! RUNNING IN IN-MEMORY MODE (DATA WILL BE LOST ON RESTART) !!!")
		userRepo = memory.NewUserRepository()
		fastingRepo = memory.NewFastingRepository()
		ketoRepo = memory.NewKetoRepository()
		leaderboardRepo = memory.NewLeaderboardRepository()
		gamificationRepo = memory.NewGamificationRepository()
		referralRepo = memory.NewReferralRepository()
		notificationRepo = memory.NewNotificationRepository()
		subscriptionRepo = memory.NewSubscriptionRepository()
		vaultRepo = memory.NewVaultRepository()
		socialRepo = memory.NewSocialRepository()
		progressRepo = memory.NewProgressRepository()
		tribeRepo = nil // No in-memory implementation for tribes yet
	}

	// Activity Repo (Memory only for now)
	activityRepo := memory.NewActivityRepository()
	telemetryRepo := memory.NewTelemetryRepository()
	mealRepo := memory.NewMealRepository()
	recipeRepo := memory.NewRecipeRepository()

	// Payment Adapter
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	paymentAdapter := payment.NewStripeAdapter(stripeKey, stripeWebhookSecret)

	// 2. Initialize Services (Core)
	vaultService := services.NewVaultService(userRepo, vaultRepo, paymentAdapter)
	referralService := services.NewReferralService(referralRepo, userRepo, vaultService)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	// Enforce minimum entropy requirements
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters long")
	}

	// Check for common weak patterns
	weakPatterns := []string{
		"your-secret-key", "secret", "password", "123456",
		"admin", "test", "demo", "changeme",
	}
	jwtSecretLower := strings.ToLower(jwtSecret)
	for _, weak := range weakPatterns {
		if strings.Contains(jwtSecretLower, weak) {
			log.Fatal("JWT_SECRET contains weak patterns. Use a cryptographically random string.")
		}
	}

	authService := services.NewAuthService(userRepo, referralService, jwtSecret)
	fastingService := services.NewFastingService(fastingRepo, vaultService, userRepo)
	ketoService := services.NewKetoService(ketoRepo, userRepo)
	leaderboardService := services.NewLeaderboardService(leaderboardRepo)
	gamificationService := services.NewGamificationService(gamificationRepo, fastingRepo)
	activityService := services.NewActivityService(activityRepo)
	telemetryService := services.NewTelemetryService(telemetryRepo)
	socialService := services.NewSocialService(socialRepo)
	progressService := services.NewProgressService(progressRepo)
	tribeService := services.NewTribeService(tribeRepo)

	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		log.Println("Warning: DEEPSEEK_API_KEY not set, Cortex will fail")
	} else {
		// Validate API key format (DeepSeek keys typically start with "sk-")
		if !strings.HasPrefix(apiKey, "sk-") || len(apiKey) < 20 {
			log.Println("Warning: DEEPSEEK_API_KEY appears invalid. Expected format: sk-...")
			log.Println("Cortex functionality may not work properly")
		}

		// Check for common placeholder values
		placeholders := []string{"your-api-key", "test-key", "demo-key", "sk-test"}
		apiKeyLower := strings.ToLower(apiKey)
		for _, placeholder := range placeholders {
			if strings.Contains(apiKeyLower, placeholder) {
				log.Println("Warning: DEEPSEEK_API_KEY contains placeholder value. Use a real API key.")
				log.Println("Cortex functionality may not work properly")
				break
			}
		}
	}
	llmAdapter := llm.NewDeepSeekAdapter(apiKey)
	cortexService := services.NewCortexService(llmAdapter, fastingRepo, userRepo)

	mealService := services.NewMealService(mealRepo, cortexService)
	recipeService := services.NewRecipeService(recipeRepo)

	stripeService := services.NewStripeService(paymentAdapter, subscriptionRepo, userRepo)

	// Notification Service
	firebaseServiceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")
	if firebaseServiceAccountPath == "" {
		log.Println("Warning: FIREBASE_SERVICE_ACCOUNT_PATH not set, notifications will be disabled")
	}
	// notificationRepo is already initialized above
	var notificationService ports.NotificationService
	realNotificationService, err := services.NewNotificationService(notificationRepo, firebaseServiceAccountPath)
	if err != nil {
		log.Printf("Warning: Failed to initialize notification service: %v. Using NoOp service.", err)
		notificationService = services.NewNoOpNotificationService()
	} else {
		notificationService = realNotificationService
	}

	handler := http.NewHandler(
		authService,
		fastingService,
		ketoService,
		leaderboardService,
		gamificationService,
		cortexService,
		activityService,
		telemetryService,
		mealService,
		recipeService,
		stripeService,
		referralService,
		notificationService,
		socialService,
		progressService,
		userRepo,
	)

	// Initialize OAuth Service
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	if googleRedirectURL == "" {
		// Default redirect URL based on environment
		env := os.Getenv("ENVIRONMENT")
		if env == "uat" {
			googleRedirectURL = "https://fastinghero-uat-537397575496.us-central1.run.app/api/v1/auth/google/callback"
		} else {
			googleRedirectURL = "http://localhost:8080/api/v1/auth/google/callback"
		}
	}

	oauthService := services.NewOAuthService(
		userRepo,
		jwtSecret,
		googleClientID,
		googleClientSecret,
		googleRedirectURL,
	)

	// Initialize OAuth handler and set it in main handler
	oauthHandler := http.NewOAuthHandler(oauthService)
	handler.SetOAuthHandler(oauthHandler)

	// Initialize Tribe handler
	tribeHandler := http.NewTribeHandler(tribeService)

	// 4. Setup Router
	router := gin.Default()

	// Configure CORS based on environment
	allowedOrigins := []string{"http://localhost:5173"}

	if prodOrigins := os.Getenv("ALLOWED_ORIGINS"); prodOrigins != "" {
		// Support multiple origins: "https://app.example.com,https://admin.example.com"
		allowedOrigins = strings.Split(prodOrigins, ",")
		for i, origin := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(origin)
		}
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Security Headers
	router.Use(middleware.SecurityHeaders())

	// Rate Limiting
	// 100 requests per minute for anonymous users
	anonymousLimiter := middleware.NewRateLimiter(rate.Limit(100.0/60.0), 10)
	anonymousLimiter.CleanupOldVisitors()
	router.Use(anonymousLimiter.Middleware())

	// Request Logging
	router.Use(middleware.RequestLogger())

	// 5. Setup Cron Jobs
	cronScheduler := cron.New()
	_, err = cronScheduler.AddFunc("@daily", func() {
		log.Println("Running daily vault earnings calculation...")
		ctx := context.Background()
		if err := vaultService.ProcessDailyEarnings(ctx); err != nil {
			log.Printf("Error processing daily earnings: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}
	cronScheduler.Start()
	log.Println("Cron scheduler started")

	handler.RegisterRoutes(router)

	// Register Tribe routes
	api := router.Group("/api/v1")
	authMiddleware := http.AuthMiddleware(authService)
	http.RegisterTribesRoutes(api, tribeHandler, authMiddleware)

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting FastingHero Server on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
