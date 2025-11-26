# FastingHero - Detailed Technical Configuration Guide

## Table of Contents

1. [System Architecture](#system-architecture)
2. [Backend Configuration](#backend-configuration)
3. [Frontend Configuration](#frontend-configuration)
4. [Service Layer Details](#service-layer-details)
5. [API Reference](#api-reference)
6. [Production Deployment](#production-deployment)

---

## System Architecture

### Technology Stack

- **Backend**: Go 1.21+
- **Frontend**: React 18 + TypeScript + Vite
- **Database**: In-memory (MariaDB support available)
- **AI**: DeepSeek LLM API
- **Authentication**: JWT
- **HTTP Framework**: Gin (Go)
- **UI Framework**: Tailwind CSS + shadcn/ui

### Directory Structure

```
fastinghero/
├── cmd/server/              # Application entry point
├── internal/
│   ├── core/
│   │   ├── domain/         # Business entities
│   │   ├── ports/          # Interface definitions
│   │   └── services/       # Business logic
│   └── adapters/
│       ├── handler/http/   # HTTP handlers
│       ├── repository/     # Data persistence
│       └── secondary/llm/  # External integrations
├── frontend/
│   ├── src/
│   │   ├── pages/         # Main views
│   │   ├── components/    # Reusable components
│   │   └── api/           # API client
└── docs/
```

---

## Backend Configuration

### 1. Main Server (`cmd/server/main.go`)

**Port Configuration**:

```go
// Line 95: Hardcoded port
router.Run(":8080")

// Production: Use environment variable
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
router.Run(":" + port)
```

**Database Selection**:

```go
// Lines 24-42: Database initialization
dsn := os.Getenv("DSN")
if dsn != "" {
    // MariaDB mode
    db, err := sql.Open("mysql", dsn)
    userRepo = mariadb.NewUserRepository(db)
} else {
    // In-memory mode (default)
    userRepo = memory.NewUserRepository()
}
```

**Environment Variables**:

```bash
# Optional - Database connection
DSN="user:password@tcp(localhost:3306)/fastinghero?parseTime=true"

# Required - DeepSeek API
DEEPSEEK_API_KEY="sk-xxxxxxxxxxxxx"

# Recommended - Server port
PORT="8080"
```

**CORS Configuration** (Lines 77-89):

```go
router.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
    c.Writer.Header().Set("Access-Control-Allow-Headers", 
        "Content-Type, Authorization, ...")
    c.Writer.Header().Set("Access-Control-Allow-Methods", 
        "POST, OPTIONS, GET, PUT")
})
```

**Production CORS**:

```go
// Replace "*" with specific origin
allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
if allowedOrigin == "" {
    allowedOrigin = "https://yourdomain.com"
}
c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
```

### 2. Service Initialization

**Dependency Injection Order** (Lines 50-68):

```go
// 1. Repositories
userRepo, fastingRepo, ketoRepo, activityRepo, 
telemetryRepo, mealRepo, recipeRepo

// 2. Core Services
pricingService := services.NewPricingService()
authService := services.NewAuthService(userRepo)
fastingService := services.NewFastingService(
    fastingRepo, pricingService, userRepo)

// 3. External Adapters
llmAdapter := llm.NewDeepSeekAdapter(apiKey)
cortexService := services.NewCortexService(
    llmAdapter, fastingRepo, userRepo)

// 4. Dependent Services
mealService := services.NewMealService(mealRepo, cortexService)

// 5. HTTP Handler
handler := http.NewHandler(authService, fastingService, ...)
```

### 3. Authentication Service

**JWT Configuration** (`internal/core/services/auth_service.go`):

```go
// CRITICAL: Hardcoded secret (Line 13)
var jwtSecret = []byte("your-secret-key")

// Production: Use environment variable
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
if len(jwtSecret) == 0 {
    log.Fatal("JWT_SECRET environment variable required")
}
```

**Token Expiration**:

```go
// Line 41: 24 hour expiration
ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour))

// Configurable:
expHours := 24
if env := os.Getenv("JWT_EXPIRY_HOURS"); env != "" {
    expHours, _ = strconv.Atoi(env)
}
ExpiresAt: jwt.NewNumericDate(
    time.Now().Add(time.Duration(expHours) * time.Hour))
```

**Password Storage**:

```go
// CRITICAL ISSUE: Plain text passwords (Line 24)
user := &domain.User{
    Email:    email,
    Password: password, // INSECURE!
}

// Required Fix: Use bcrypt
import "golang.org/x/crypto/bcrypt"

hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password), bcrypt.DefaultCost)
user.Password = string(hashedPassword)

// Login verification:
err := bcrypt.CompareHashAndPassword(
    []byte(user.Password), []byte(password))
```

### 4. Pricing Service

**Constants** (`internal/core/services/pricing_service.go`):

```go
const (
    MonthlyCharge = 30.0  // Base subscription
    BaseFee       = 10.0  // Minimum fee
    VaultDeposit  = 20.0  // Initial deposit
    DailyMax      = 2.0   // Max earnings per day
)
```

**Discipline Calculation** (Lines 50-64):

```go
func UpdateDisciplineIndex(user *domain.User, 
    completedFast bool, verifiedKetosis bool) {
    
    if completedFast {
        user.DisciplineIndex += 1
    }
    if verifiedKetosis {
        user.DisciplineIndex += 2
    }
    
    // Clamping
    if user.DisciplineIndex > 100 {
        user.DisciplineIndex = 100
    }
    if user.DisciplineIndex < 0 {
        user.DisciplineIndex = 0
    }
}
```

**Customization**:

```go
type PricingConfig struct {
    MonthlyCharge float64
    BaseFee       float64
    VaultDeposit  float64
    DailyMax      float64
    FastBonus     float64
    KetoBonus     float64
}

// Load from config file or env
config := loadPricingConfig()
```

### 5. Fasting Service

**Plan Types** (`internal/core/domain/fasting.go`):

```go
type FastingPlanType string

const (
    Plan168  FastingPlanType = "16:8"   // 16h fast, 8h eat
    Plan186  FastingPlanType = "18:6"   // 18h fast, 6h eat
    Plan204  FastingPlanType = "20:4"   // 20h fast, 4h eat
    PlanOMAD FastingPlanType = "OMAD"   // One meal a day
    Plan5_2  FastingPlanType = "5:2"    // 5 days normal, 2 days 500cal
)
```

**Goal Calculation** (`internal/core/services/fasting_service.go`):

```go
func (s *FastingService) StopFast(ctx context.Context, 
    userID uuid.UUID) (*domain.FastingSession, error) {
    
    // Calculate duration
    duration := time.Since(session.StartTime)
    goalMet := duration.Hours() >= float64(session.GoalHours)
    
    // Update discipline
    if goalMet {
        user.DisciplineIndex += 5  // Success bonus
    } else {
        user.DisciplineIndex -= 10 // Failure penalty
    }
}
```

### 6. DeepSeek Integration

**API Configuration** (`internal/adapters/secondary/llm/deepseek.go`):

```go
const (
    deepseekURL = "https://api.deepseek.com/v1/chat/completions"
    model       = "deepseek-chat"
    maxTokens   = 500
    temperature = 0.7
)
```

**Request Structure**:

```go
type chatRequest struct {
    Model       string    `json:"model"`
    Messages    []message `json:"messages"`
    MaxTokens   int       `json:"max_tokens"`
    Temperature float64   `json:"temperature"`
}

type message struct {
    Role    string `json:"role"`    // "system" or "user"
    Content string `json:"content"`
}
```

**Multimodal Support** (Lines 86-140):

```go
// Vision API structure (not yet supported by deepseek-chat)
type visionMessage struct {
    Role    string        `json:"role"`
    Content []contentPart `json:"content"`
}

type contentPart struct {
    Type     string    `json:"type"` // "text" or "image_url"
    Text     string    `json:"text,omitempty"`
    ImageURL *imageURL `json:"image_url,omitempty"`
}

type imageURL struct {
    URL string `json:"url"` // "data:image/jpeg;base64,..."
}
```

**Error Handling**:

```go
// Fallback for vision failures (Line 130)
if err != nil {
    return "Analysis: Mock response - DeepSeek V2 is text-only. " +
           "Image analysis simulated based on description.", nil
}
```

### 7. Meal Analysis

**Cortex Service Configuration** (`internal/core/services/cortex_service.go`):

```go
func (s *CortexService) AnalyzeMeal(ctx context.Context, 
    imageBase64, description string) (string, bool, bool, error) {
    
    prompt := fmt.Sprintf(`Analyze this meal based on: "%s".
    1. Is this real food or fake? (Authenticity)
    2. Estimate carbs. Keto-friendly (under 10g net carbs)?
    
    Output format:
    Analysis: [description and carb estimate]
    Authenticity: [Verified/Suspicious]
    Keto-Friendly: [Yes/No]
    `, description)
    
    response, err := s.llm.AnalyzeImage(ctx, imageBase64, prompt)
    
    // Parse response
    lowerResp := strings.ToLower(response)
    isKeto := !strings.Contains(lowerResp, "keto-friendly: no")
    isAuthentic := !strings.Contains(lowerResp, "suspicious")
    
    return response, isAuthentic, isKeto, nil
}
```

**Meal Logging Rewards** (`internal/core/services/meal_service.go`):

```go
// After successful meal log
const mealReward = 0.50 // $0.50 per meal
const dailyMealLimit = 3

// In handler:
if mealsLoggedToday < dailyMealLimit {
    pricingService.AddDailyEarnings(ctx, user, mealReward)
}
```

### 8. Recipe Service

**Seeded Data** (`internal/adapters/repository/memory/memory_repos.go`):

```go
func NewRecipeRepository() *RecipeRepository {
    return &RecipeRepository{
        recipes: []domain.Recipe{
            {
                ID:          "1",
                Title:       "Avocado & Egg Breakfast Bowl",
                Description: "A simple, high-fat breakfast",
                Ingredients: []string{"2 Eggs", "1 Avocado", ...},
                Instructions: []string{"Boil eggs", "Slice avocado", ...},
                Diet:        domain.DietVegetarian,
                IsSimple:    true,
                Calories:    450,
                Carbs:       4,
                Image:       "https://images.unsplash.com/...",
            },
            // ... 4 more recipes
        },
    }
}
```

**Filtering Logic** (`internal/core/services/recipe_service.go`):

```go
func (s *RecipeService) GetRecipes(ctx context.Context, 
    diet domain.DietType) ([]domain.Recipe, error) {
    
    allRecipes, _ := s.repo.FindAll(ctx)
    
    if diet == "" || diet == domain.DietNormal {
        return allRecipes, nil // Return all
    }
    
    var filtered []domain.Recipe
    for _, r := range allRecipes {
        if diet == domain.DietVegetarian {
            // Include vegan + vegetarian
            if r.Diet == domain.DietVegetarian || 
               r.Diet == domain.DietVegan {
                filtered = append(filtered, r)
            }
        } else if diet == domain.DietVegan {
            // Vegan only
            if r.Diet == domain.DietVegan {
                filtered = append(filtered, r)
            }
        }
    }
    return filtered, nil
}
```

---

## Frontend Configuration

### 1. API Client (`frontend/src/api/client.ts`)

**Base URL Configuration**:

```typescript
const API_URL = 'http://localhost:8080/api/v1';

// Production:
const API_URL = import.meta.env.VITE_API_URL || 
                'https://api.yourdomain.com/api/v1';
```

**Token Interceptor** (Lines 11-17):

```typescript
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});
```

**Auto-Logout on 401** (Lines 20-29):

```typescript
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            localStorage.removeItem('token');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);
```

### 2. Vite Configuration (`frontend/vite.config.ts`)

```typescript
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
  }
})
```

### 3. Environment Variables

**Development** (`.env.development`):

```bash
VITE_API_URL=http://localhost:8080/api/v1
VITE_ENABLE_MOCK=false
```

**Production** (`.env.production`):

```bash
VITE_API_URL=https://api.fastinghero.com/api/v1
VITE_ENABLE_ANALYTICS=true
```

### 4. Component Configuration

**Progress Page** (`frontend/src/pages/Progress.tsx`):

```typescript
// Hydration goal
const waterGoal = 8; // glasses per day

// Meal logging limits
const dailyMealLimit = 3;
const mealReward = 0.50;

// Weight data mock
const weightData = [
  { date: "Mon", weight: 180 },
  { date: "Tue", weight: 179.5 },
  // ...
];
```

**Dashboard** (`frontend/src/pages/Dashboard.tsx`):

```typescript
// Fasting plans
const fastingPlans = [
  { type: "16:8", hours: 16, description: "Beginner friendly" },
  { type: "18:6", hours: 18, description: "Intermediate" },
  { type: "20:4", hours: 20, description: "Advanced" },
  { type: "OMAD", hours: 23, description: "Expert" },
];
```

**Community** (`frontend/src/pages/Community.tsx`):

```typescript
// Diet filters
const dietFilters = ["all", "vegan", "vegetarian", "normal"];

// Feed refresh interval
const FEED_REFRESH_MS = 30000; // 30 seconds

// Leaderboard size
const LEADERBOARD_SIZE = 8;
```

---

## API Reference

### Authentication Endpoints

**POST /api/v1/auth/register**

```json
Request:
{
  "email": "user@example.com",
  "password": "securepassword"
}

Response: 201 Created
{
  "id": "uuid",
  "email": "user@example.com",
  "subscription_tier": "free",
  "discipline_index": 0,
  "current_price": 50,
  "vault_deposit": 0,
  "earned_refund": 0,
  "created_at": "2025-11-26T12:00:00Z"
}
```

**POST /api/v1/auth/login**

```json
Request:
{
  "email": "user@example.com",
  "password": "securepassword"
}

Response: 200 OK
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "refresh_token_here"
}
```

### Fasting Endpoints

**POST /api/v1/fasting/start**

```json
Request:
{
  "plan_type": "16:8",
  "goal_hours": 16,
  "start_time": "2025-11-26T20:00:00Z" // optional
}

Response: 201 Created
{
  "id": "uuid",
  "user_id": "uuid",
  "plan_type": "16:8",
  "start_time": "2025-11-26T20:00:00Z",
  "goal_hours": 16,
  "is_active": true
}
```

**POST /api/v1/fasting/stop**

```json
Response: 200 OK
{
  "id": "uuid",
  "end_time": "2025-11-27T12:00:00Z",
  "duration_hours": 16,
  "goal_met": true,
  "is_active": false
}
```

### Meal Endpoints

**POST /api/v1/meals/**

```json
Request:
{
  "image": "data:image/jpeg;base64,/9j/4AAQSkZJRg...",
  "description": "Grilled chicken salad"
}

Response: 201 Created
{
  "id": "uuid",
  "user_id": "uuid",
  "image": "data:image/jpeg;base64,...",
  "description": "Grilled chicken salad",
  "logged_at": "2025-11-26T12:30:00Z",
  "calories": 0,
  "analysis": "Analysis: Grilled chicken with leafy greens...",
  "is_keto": true,
  "is_authentic": true
}
```

**GET /api/v1/meals/**

```json
Response: 200 OK
[
  {
    "id": "uuid",
    "image": "data:image/jpeg;base64,...",
    "description": "Grilled chicken salad",
    "logged_at": "2025-11-26T12:30:00Z",
    "analysis": "...",
    "is_keto": true,
    "is_authentic": true
  }
]
```

### Recipe Endpoints

**GET /api/v1/recipes/?diet=vegan**

```json
Response: 200 OK
[
  {
    "id": "4",
    "title": "Vegan Keto Tofu Stir-fry",
    "description": "Crispy tofu with low-carb veggies",
    "ingredients": ["Firm Tofu", "Broccoli", ...],
    "instructions": ["Press tofu", "Fry until golden", ...],
    "diet": "vegan",
    "is_simple": false,
    "calories": 320,
    "carbs": 9,
    "image": "https://images.unsplash.com/..."
  }
]
```

### Telemetry Endpoints

**POST /api/v1/telemetry/manual**

```json
Request:
{
  "type": "weight",
  "value": 175.5,
  "unit": "lbs"
}

Response: 201 Created
{
  "id": "uuid",
  "type": "weight",
  "value": 175.5,
  "unit": "lbs",
  "timestamp": "2025-11-26T08:00:00Z"
}
```

**GET /api/v1/telemetry/metric?type=weight**

```json
Response: 200 OK
{
  "value": 175.5,
  "unit": "lbs",
  "timestamp": "2025-11-26T08:00:00Z"
}
```

### Cortex Endpoints

**POST /api/v1/cortex/chat**

```json
Request:
{
  "message": "Should I break my fast?"
}

Response: 200 OK
{
  "response": "No. You're at 14 hours. Push through..."
}
```

---

## Production Deployment

### Environment Checklist

```bash
# Required
export JWT_SECRET="random-256-bit-secret"
export DEEPSEEK_API_KEY="sk-xxxxx"
export DSN="user:pass@tcp(db:3306)/fastinghero?parseTime=true"

# Recommended
export PORT="8080"
export ALLOWED_ORIGIN="https://yourdomain.com"
export LOG_LEVEL="info"
export ENABLE_METRICS="true"

# Optional
export JWT_EXPIRY_HOURS="24"
export MAX_UPLOAD_SIZE="10485760"  # 10MB
export RATE_LIMIT_PER_MINUTE="60"
```

### Docker Deployment

**Dockerfile** (Backend):

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o fastinghero cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/fastinghero .
EXPOSE 8080
CMD ["./fastinghero"]
```

**Dockerfile** (Frontend):

```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

**docker-compose.yml**:

```yaml
version: '3.8'
services:
  db:
    image: mariadb:10.11
    environment:
      MYSQL_ROOT_PASSWORD: rootpass
      MYSQL_DATABASE: fastinghero
      MYSQL_USER: fhuser
      MYSQL_PASSWORD: fhpass
    volumes:
      - db_data:/var/lib/mysql
    
  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      DSN: "fhuser:fhpass@tcp(db:3306)/fastinghero?parseTime=true"
      JWT_SECRET: "${JWT_SECRET}"
      DEEPSEEK_API_KEY: "${DEEPSEEK_API_KEY}"
    depends_on:
      - db
    
  frontend:
    build: ./frontend
    ports:
      - "80:80"
    depends_on:
      - backend

volumes:
  db_data:
```

### Database Migrations

**Initial Schema** (`migrations/001_init.sql`):

```sql
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    subscription_tier VARCHAR(50) DEFAULT 'free',
    discipline_index FLOAT DEFAULT 0,
    current_price FLOAT DEFAULT 50,
    vault_deposit FLOAT DEFAULT 0,
    earned_refund FLOAT DEFAULT 0,
    signed_contract BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE fasting_sessions (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    plan_type VARCHAR(50) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    goal_hours INT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE meals (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    image LONGTEXT,
    description TEXT,
    logged_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    calories INT DEFAULT 0,
    analysis TEXT,
    is_keto BOOLEAN DEFAULT TRUE,
    is_authentic BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Monitoring & Logging

**Structured Logging**:

```go
import "github.com/rs/zerolog/log"

log.Info().
    Str("user_id", userID.String()).
    Str("action", "start_fast").
    Int("goal_hours", goalHours).
    Msg("Fasting session started")
```

**Health Check Endpoint**:

```go
router.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "healthy",
        "timestamp": time.Now(),
        "version": "1.0.0",
    })
})
```

---

## Security Hardening

### Critical Fixes Required

1. **Password Hashing** (IMMEDIATE):

```go
import "golang.org/x/crypto/bcrypt"

// Registration
hashedPassword, _ := bcrypt.GenerateFromPassword(
    []byte(password), bcrypt.DefaultCost)

// Login
err := bcrypt.CompareHashAndPassword(
    []byte(user.Password), []byte(password))
```

2. **JWT Secret from Environment**:

```go
jwtSecret := []byte(os.Getenv("JWT_SECRET"))
if len(jwtSecret) == 0 {
    log.Fatal("JWT_SECRET required")
}
```

3. **Rate Limiting**:

```go
import "github.com/ulule/limiter/v3"

rate := limiter.Rate{
    Period: 1 * time.Minute,
    Limit:  60,
}
middleware := limiter.NewMiddleware(limiter.New(store, rate))
router.Use(middleware)
```

4. **Input Validation**:

```go
import "github.com/go-playground/validator/v10"

type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

validate := validator.New()
if err := validate.Struct(req); err != nil {
    return errors.New("invalid input")
}
```

---

This configuration guide provides detailed setup instructions for every component. For production deployment, prioritize the security hardening section and ensure all environment variables are properly configured.
