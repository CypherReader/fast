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
	var tribeRepo ports.TribeRepository
	var socialRepo ports.SocialRepository
	var leaderboardRepo ports.LeaderboardRepository
	var gamificationRepo ports.GamificationRepository
	var referralRepo ports.ReferralRepository
	var notificationRepo ports.NotificationRepository

	// Check for DB connection string
	dsn := os.Getenv("DSN")
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

	var db *sql.DB
	var err error
	useMemory := false

	if dsn != "" {
		log.Println("Connecting to PostgreSQL...")
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Failed to open DB: %v. Switching to IN-MEMORY mode.", err)
			useMemory = true
		} else {
			if err := db.Ping(); err != nil {
				log.Printf("Failed to ping DB: %v. Switching to IN-MEMORY mode.", err)
				useMemory = true
			}
		}
	} else {
		log.Println("DSN not set. Switching to IN-MEMORY mode.")
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
		tribeRepo = postgres.NewPostgresTribeRepository(db)
		socialRepo = postgres.NewPostgresSocialRepository(db)
		leaderboardRepo = postgres.NewPostgresLeaderboardRepository(db)
		gamificationRepo = postgres.NewPostgresGamificationRepository(db)
		referralRepo = postgres.NewPostgresReferralRepository(db)
		notificationRepo = postgres.NewPostgresNotificationRepository(db)
	} else {
		log.Println("!!! RUNNING IN IN-MEMORY MODE (DATA WILL BE LOST ON RESTART) !!!")
		userRepo = memory.NewUserRepository()
		fastingRepo = memory.NewFastingRepository()
		ketoRepo = memory.NewKetoRepository()
		tribeRepo = memory.NewTribeRepository()
		socialRepo = memory.NewSocialRepository()
		leaderboardRepo = memory.NewLeaderboardRepository()
		gamificationRepo = memory.NewGamificationRepository()
		referralRepo = memory.NewReferralRepository()
		notificationRepo = memory.NewNotificationRepository()
	}

	// Activity Repo (Memory only for now)
	activityRepo := memory.NewActivityRepository()
	telemetryRepo := memory.NewTelemetryRepository()
	mealRepo := memory.NewMealRepository()
	recipeRepo := memory.NewRecipeRepository()

	// 2. Initialize Services (Core)
	vaultService := services.NewVaultService(userRepo, nil) // Payment gateway injected later if needed for refunds
	referralService := services.NewReferralService(referralRepo, userRepo, vaultService)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	if jwtSecret == "your-secret-key" {
		log.Fatal("JWT_SECRET must be changed from the default value")
	}

	authService := services.NewAuthService(userRepo, referralService, jwtSecret)
	fastingService := services.NewFastingService(fastingRepo, vaultService, userRepo)
	ketoService := services.NewKetoService(ketoRepo, userRepo)
	socialService := services.NewSocialService(socialRepo, userRepo)
	leaderboardService := services.NewLeaderboardService(leaderboardRepo)
	gamificationService := services.NewGamificationService(gamificationRepo, fastingRepo)
	activityService := services.NewActivityService(activityRepo)
	telemetryService := services.NewTelemetryService(telemetryRepo)

	// Initialize LLM Adapter
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		log.Println("Warning: DEEPSEEK_API_KEY not set, Cortex will fail")
	}
	llmAdapter := llm.NewDeepSeekAdapter(apiKey)
	cortexService := services.NewCortexService(llmAdapter, fastingRepo, userRepo)

	mealService := services.NewMealService(mealRepo, cortexService)
	recipeService := services.NewRecipeService(recipeRepo)
	tribeService := services.NewTribeService(tribeRepo, userRepo)

	// Payment Adapter
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	paymentAdapter := payment.NewStripeAdapter(stripeKey, stripeWebhookSecret)

	// Notification Service
	firebaseServiceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")
	if firebaseServiceAccountPath == "" {
		log.Println("Warning: FIREBASE_SERVICE_ACCOUNT_PATH not set, notifications will be disabled")
	}
	// notificationRepo is already initialized above
	notificationService, err := services.NewNotificationService(notificationRepo, firebaseServiceAccountPath)
	if err != nil {
		log.Printf("Warning: Failed to initialize notification service: %v", err)
		notificationService = nil // Continue without notifications
	}

	handler := http.NewHandler(
		authService,
		fastingService,
		ketoService,
		socialService,
		leaderboardService,
		gamificationService,
		cortexService,
		activityService,
		telemetryService,
		mealService,
		recipeService,
		tribeService,
		paymentAdapter,
		referralService,
		notificationService,
		userRepo,
	)

	// 4. Setup Router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
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

	// 5. Start Server
	log.Println("Starting FastingHero Server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
