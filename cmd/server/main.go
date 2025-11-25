package main

import (
	"database/sql"
	"fastinghero/internal/adapters/handler/http"
	"fastinghero/internal/adapters/repository/mariadb"
	"fastinghero/internal/adapters/repository/memory"
	"fastinghero/internal/adapters/secondary/llm"
	"fastinghero/internal/core/ports"
	"fastinghero/internal/core/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var userRepo ports.UserRepository
	var fastingRepo ports.FastingRepository
	var ketoRepo ports.KetoRepository

	// Check for DB connection string
	dsn := os.Getenv("DSN")
	if dsn != "" {
		log.Println("Connecting to MariaDB...")
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Failed to open DB: %v", err)
		}
		if err := db.Ping(); err != nil {
			log.Fatalf("Failed to ping DB: %v", err)
		}
		userRepo = mariadb.NewUserRepository(db)
		fastingRepo = mariadb.NewFastingRepository(db)
		ketoRepo = mariadb.NewKetoRepository(db)
	} else {
		log.Println("Using In-Memory Repositories")
		userRepo = memory.NewUserRepository()
		fastingRepo = memory.NewFastingRepository()
		ketoRepo = memory.NewKetoRepository()
	}

	// 2. Initialize Services (Core)
	// 2. Initialize Services (Core)
	pricingService := services.NewPricingService()
	authService := services.NewAuthService(userRepo)
	fastingService := services.NewFastingService(fastingRepo, pricingService, userRepo)
	ketoService := services.NewKetoService(ketoRepo, userRepo)
	socialService := services.NewSocialService()

	// Initialize LLM Adapter
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		log.Println("Warning: DEEPSEEK_API_KEY not set, Cortex will fail")
	}
	llmAdapter := llm.NewDeepSeekAdapter(apiKey)
	cortexService := services.NewCortexService(llmAdapter, fastingRepo, userRepo)

	// 3. Initialize Handlers (Adapters)
	handler := http.NewHandler(authService, fastingService, ketoService, socialService, cortexService)

	// 4. Setup Router
	router := gin.Default()

	// CORS Middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	handler.RegisterRoutes(router)

	// 5. Start Server
	log.Println("Starting FastingHero Server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
