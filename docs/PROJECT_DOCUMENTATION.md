# FastingHero - Complete Project Documentation

**Version**: 1.0.0  
**Last Updated**: November 26, 2025  
**Architecture**: Hexagonal/Clean Architecture  
**Tech Stack**: Go + React/TypeScript + DeepSeek AI

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [System Architecture](#system-architecture)
3. [Domain Models](#domain-models)
4. [Backend Services](#backend-services)
5. [API Endpoints](#api-endpoints)
6. [Frontend Components](#frontend-components)
7. [Data Flow](#data-flow)
8. [Configuration](#configuration)
9. [Deployment](#deployment)
10. [Production Readiness](#production-readiness)

---

## Project Overview

### What is FastingHero?

FastingHero is a gamified fasting and keto tracking application that uses behavioral economics to help users maintain their health goals. The core concept is a "Commitment Vault" where users deposit money upfront and earn it back through consistent adherence to fasting and dietary commitments.

### Key Features

1. **Fasting Timer & Tracking**
   - Multiple fasting plans (16:8, 18:6, OMAD, 24h, 36h, Extended)
   - Real-time timer with goal tracking
   - Discipline Index (0-100 scale)
   - Success/failure tracking

2. **Commitment Vault System**
   - Monthly deposit ($20 default)
   - Earn back through adherence
   - Dynamic pricing based on discipline
   - Daily earning caps ($2/day)

3. **AI Coaching (Cortex)**
   - DeepSeek-powered conversational AI
   - Context-aware responses
   - Discipline-based tone adjustment
   - Biological insights during fasting

4. **Progress Tracking**
   - Weight logging and trends
   - Hydration tracking (8 glasses/day)
   - Ketosis level estimation
   - Food logging with AI analysis

5. **Community Features**
   - Social feed with user activities
   - Weekly leaderboard
   - Keto recipe database (vegan/vegetarian/normal)
   - Knowledge hub

6. **Tribe System**
   - Group fasting challenges
   - Shared accountability
   - Collective discipline tracking

### Business Model

- **Free Tier**: Basic fasting tracking
- **Premium Tier**: Advanced analytics, blood ketone logging
- **Elite Tier**: Personalized coaching, extended fasting protocols

---

## System Architecture

### Hexagonal Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    External World                        │
│  (HTTP Clients, Databases, External APIs)               │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                   Adapters (Ports)                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   HTTP       │  │  Repository  │  │   DeepSeek   │  │
│  │  Handlers    │  │   (Memory/   │  │     LLM      │  │
│  │              │  │   MariaDB)   │  │   Adapter    │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                   Core Business Logic                    │
│  ┌──────────────────────────────────────────────────┐   │
│  │              Services Layer                       │   │
│  │  - AuthService      - FastingService             │   │
│  │  - PricingService   - CortexService              │   │
│  │  - MealService      - RecipeService              │   │
│  │  - TelemetryService - TribeService               │   │
│  └──────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────┐   │
│  │              Domain Models                        │   │
│  │  - User             - FastingSession             │   │
│  │  - Meal             - Recipe                     │   │
│  │  - TelemetryData    - Tribe                      │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

### Technology Stack

**Backend**:

- **Language**: Go 1.21+
- **HTTP Framework**: Gin
- **Database**: In-memory (MariaDB support available)
- **Authentication**: JWT
- **AI Integration**: DeepSeek API

**Frontend**:

- **Framework**: React 18
- **Language**: TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **UI Components**: shadcn/ui
- **State Management**: React Hooks

**DevOps**:

- **Version Control**: Git
- **Containerization**: Docker (ready)
- **CI/CD**: Not yet configured

---

## Domain Models

### 1. User (`internal/core/domain/user.go`)

```go
type User struct {
    ID               uuid.UUID        // Unique identifier
    Email            string           // User email (unique)
    PasswordHash     string           // Hashed password (currently plain text!)
    SubscriptionTier SubscriptionTier // free, premium, elite
    DisciplineIndex  float64          // 0-100 scale
    CurrentPrice     float64          // Deprecated (Lazy Tax)
    VaultDeposit     float64          // Monthly deposit (e.g., $20)
    EarnedRefund     float64          // Amount earned back
    TribeID          *uuid.UUID       // Optional tribe membership
    SignedContract   bool             // Contract acceptance
    CreatedAt        time.Time        // Account creation
}
```

**Subscription Tiers**:

- `free`: Basic features
- `premium`: Advanced analytics
- `elite`: Full features + personalized coaching

**Methods**:

- `IsPremium()`: Returns true if premium or elite
- `IsElite()`: Returns true if elite tier

### 2. FastingSession (`internal/core/domain/fasting.go`)

```go
type FastingSession struct {
    ID        uuid.UUID       // Session identifier
    UserID    uuid.UUID       // Owner
    StartTime time.Time       // When fast started
    EndTime   *time.Time      // When fast ended (nil if active)
    GoalHours int             // Target duration
    PlanType  FastingPlanType // Type of fast
    Status    FastingStatus   // active, completed, cancelled
}
```

**Fasting Plans**:

- `beginner`: Flexible beginner plan
- `16_8`: 16 hours fasting, 8 hours eating
- `18_6`: 18 hours fasting, 6 hours eating
- `omad`: One meal a day (~23 hours)
- `24h`: 24-hour fast
- `36h`: 36-hour fast
- `extended`: Multi-day fasts

**Status Types**:

- `active`: Currently fasting
- `completed`: Successfully finished
- `cancelled`: Stopped early

### 3. Meal (`internal/core/domain/meal.go`)

```go
type Meal struct {
    ID          uuid.UUID  // Meal identifier
    UserID      uuid.UUID  // Owner
    Image       string     // Base64 encoded image
    Description string     // User description
    LoggedAt    time.Time  // When logged
    Calories    int        // Estimated calories
    Analysis    string     // AI analysis text
    IsKeto      bool       // Keto-friendly flag
    IsAuthentic bool       // Real food vs screenshot
}
```

**AI Analysis**:

- Uses DeepSeek LLM to analyze meal description
- Estimates carb content
- Determines keto-friendliness (< 10g net carbs)
- Checks authenticity

### 4. Recipe (`internal/core/domain/recipe.go`)

```go
type Recipe struct {
    ID           string     // Recipe identifier
    Title        string     // Recipe name
    Description  string     // Brief description
    Ingredients  []string   // List of ingredients
    Instructions []string   // Step-by-step instructions
    Diet         DietType   // vegan, vegetarian, normal
    IsSimple     bool       // Simple recipe flag
    Calories     int        // Total calories
    Carbs        int        // Net carbs
    Image        string     // Image URL
}
```

**Diet Types**:

- `vegan`: Plant-based only
- `vegetarian`: Includes eggs/dairy
- `normal`: Includes meat

### 5. TelemetryData (`internal/core/domain/telemetry.go`)

```go
type TelemetryData struct {
    ID        uuid.UUID  // Data point identifier
    UserID    uuid.UUID  // Owner
    Type      string     // weight, steps, hydration
    Value     float64    // Measurement value
    Unit      string     // lbs, count, glasses
    Timestamp time.Time  // When recorded
    Source    string     // manual, device
}
```

**Telemetry Types**:

- `weight`: Body weight tracking
- `steps`: Daily step count
- `hydration`: Water intake

### 6. Tribe (`internal/core/domain/tribe.go`)

```go
type Tribe struct {
    ID               uuid.UUID   // Tribe identifier
    Name             string      // Tribe name
    Description      string      // Tribe description
    LeaderID         uuid.UUID   // Tribe leader
    Members          []uuid.UUID // Member IDs
    CollectiveDiscipline float64 // Average discipline
    CreatedAt        time.Time   // Creation date
}
```

**Tribe Features**:

- Group fasting challenges
- Shared accountability
- Collective discipline tracking

---

## Backend Services

### 1. AuthService (`internal/core/services/auth_service.go`)

**Purpose**: User authentication and authorization

**Methods**:

```go
Register(ctx, email, password) (*User, error)
Login(ctx, email, password) (token, refreshToken string, error)
ValidateToken(tokenString) (uuid.UUID, error)
```

**JWT Configuration**:

- Secret: `"your-secret-key"` (hardcoded - needs env var!)
- Expiration: 24 hours
- Algorithm: HS256

**Critical Issue**: Passwords stored in plain text! Needs bcrypt hashing.

### 2. FastingService (`internal/core/services/fasting_service.go`)

**Purpose**: Fasting session management

**Methods**:

```go
StartFast(ctx, userID, planType, goalHours, startTime) (*FastingSession, error)
StopFast(ctx, userID) (*FastingSession, error)
GetCurrentFast(ctx, userID) (*FastingSession, error)
```

**Business Logic**:

- Creates new fasting session
- Calculates duration on stop
- Updates discipline index:
  - Success: +5 discipline
  - Failure: -10 discipline
- Recalculates user price

### 3. PricingService (`internal/core/services/pricing_service.go`)

**Purpose**: Commitment vault and pricing calculations

**Constants**:

```go
MonthlyCharge = 30.0  // Base subscription
BaseFee       = 10.0  // Minimum fee
VaultDeposit  = 20.0  // Initial deposit
DailyMax      = 2.0   // Max earnings per day
```

**Methods**:

```go
CalculateVaultStatus(user) (deposit, earned, potentialRefund float64)
AddDailyEarnings(ctx, user, amount)
UpdateDisciplineIndex(ctx, user, completedFast, verifiedKetosis)
```

**Discipline Updates**:

- Completed fast: +1
- Verified ketosis: +2
- Clamped to 0-100 range

### 4. CortexService (`internal/core/services/cortex_service.go`)

**Purpose**: AI coaching and insights

**Methods**:

```go
Chat(ctx, userID, message) (string, error)
GenerateInsight(ctx, userID, fastingHours) (string, error)
AnalyzeMeal(ctx, imageBase64, description) (analysis string, isAuthentic, isKeto bool, error)
```

**Cortex Personality**:

- Ruthless but fair
- Discipline-based tone:
  - Low discipline: Tougher
  - High discipline: Encouraging but demanding
- Concise responses (<50 words)

**Meal Analysis Prompt**:

```
Analyze this meal based on: "{description}".
1. Is this real food or fake? (Authenticity)
2. Estimate carbs. Keto-friendly (under 10g net carbs)?

Output format:
Analysis: [description and carb estimate]
Authenticity: [Verified/Suspicious]
Keto-Friendly: [Yes/No]
```

**Response Parsing**:

- Checks for "keto-friendly: no" → `isKeto = false`
- Checks for "suspicious" or "fake" → `isAuthentic = false`

### 5. MealService (`internal/core/services/meal_service.go`)

**Purpose**: Meal logging and tracking

**Methods**:

```go
LogMeal(ctx, userID, image, description) (*Meal, error)
GetMeals(ctx, userID) ([]Meal, error)
```

**Workflow**:

1. Receives image (base64) and description
2. Calls `CortexService.AnalyzeMeal()`
3. Stores analysis results
4. Returns meal with AI insights

**Rewards**:

- $0.50 per meal logged
- Max 3 meals per day
- Capped at vault deposit

### 6. RecipeService (`internal/core/services/recipe_service.go`)

**Purpose**: Recipe management and filtering

**Methods**:

```go
GetRecipes(ctx, diet DietType) ([]Recipe, error)
```

**Filtering Logic**:

- `all` or `normal`: Returns all recipes
- `vegetarian`: Returns vegetarian + vegan
- `vegan`: Returns vegan only

**Seeded Recipes** (5 total):

1. Avocado & Egg Breakfast Bowl (vegetarian, simple)
2. Keto Chicken Salad (normal, simple)
3. Zucchini Noodles with Pesto (vegetarian)
4. Vegan Keto Tofu Stir-fry (vegan)
5. Simple Bulletproof Coffee (vegetarian, simple)

### 7. TelemetryService (`internal/core/services/telemetry_service.go`)

**Purpose**: Health metrics tracking

**Methods**:

```go
LogManualData(ctx, userID, metricType, value, unit) error
GetLatestMetric(ctx, userID, metricType) (*TelemetryData, error)
GetWeeklyStats(ctx, userID, metricType) ([]WeeklyStat, error)
```

**Supported Metrics**:

- Weight (lbs)
- Steps (count)
- Hydration (glasses)

### 8. TribeService (`internal/core/services/tribe_service.go`)

**Purpose**: Group management

**Methods**:

```go
CreateTribe(ctx, name, description, leaderID) (*Tribe, error)
JoinTribe(ctx, tribeID, userID) error
GetTribe(ctx, tribeID) (*Tribe, error)
UpdateCollectiveDiscipline(ctx, tribeID) error
```

**Collective Discipline**:

- Average of all member discipline indices
- Updated when members complete fasts

---

## API Endpoints

### Authentication

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
  "vault_deposit": 0,
  "earned_refund": 0
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

### Fasting

**POST /api/v1/fasting/start**

```json
Request:
{
  "plan_type": "16_8",
  "goal_hours": 16,
  "start_time": "2025-11-26T20:00:00Z" // optional
}

Response: 201 Created
{
  "id": "uuid",
  "user_id": "uuid",
  "plan_type": "16_8",
  "start_time": "2025-11-26T20:00:00Z",
  "goal_hours": 16,
  "status": "active"
}
```

**POST /api/v1/fasting/stop**

```json
Response: 200 OK
{
  "id": "uuid",
  "end_time": "2025-11-27T12:00:00Z",
  "status": "completed"
}
```

**GET /api/v1/fasting/current**

```json
Response: 200 OK
{
  "id": "uuid",
  "start_time": "2025-11-26T20:00:00Z",
  "goal_hours": 16,
  "status": "active"
}
```

### Meals

**POST /api/v1/meals/**

```json
Request:
{
  "image": "data:image/jpeg;base64,/9j/4AAQ...",
  "description": "Grilled chicken salad"
}

Response: 201 Created
{
  "id": "uuid",
  "logged_at": "2025-11-26T12:30:00Z",
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
    "description": "Grilled chicken salad",
    "logged_at": "2025-11-26T12:30:00Z",
    "is_keto": true
  }
]
```

### Recipes

**GET /api/v1/recipes/?diet=vegan**

```json
Response: 200 OK
[
  {
    "id": "4",
    "title": "Vegan Keto Tofu Stir-fry",
    "description": "Crispy tofu with low-carb veggies",
    "diet": "vegan",
    "is_simple": false,
    "calories": 320,
    "carbs": 9
  }
]
```

### Telemetry

**POST /api/v1/telemetry/manual**

```json
Request:
{
  "type": "weight",
  "value": 175.5,
  "unit": "lbs"
}

Response: 201 Created
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

### Cortex

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

**POST /api/v1/cortex/insight**

```json
Request:
{
  "fasting_hours": 16
}

Response: 200 OK
{
  "insight": "Your body is in ketosis. Autophagy is active..."
}
```

### Tribes

**POST /api/v1/tribes/**

```json
Request:
{
  "name": "Morning Warriors",
  "description": "Early morning fasters"
}

Response: 201 Created
{
  "id": "uuid",
  "name": "Morning Warriors",
  "leader_id": "uuid",
  "collective_discipline": 0
}
```

---

## Frontend Components

### Pages

**1. Dashboard (`frontend/src/pages/Dashboard.tsx`)**

- Fasting timer display
- Start/stop fast controls
- Current price display
- Discipline index gauge
- Quick action buttons

**2. Progress (`frontend/src/pages/Progress.tsx`)**

- Weight tracking chart
- Hydration tracker (8 glasses)
- Ketosis level gauge
- Meal logging with camera
- Recent meals list with AI analysis

**3. Community (`frontend/src/pages/Community.tsx`)**

- Social feed (4 tabs):
  - Feed: User activities
  - Leaderboard: Weekly rankings
  - Recipes: Keto recipe database
  - Knowledge: Educational content

**4. Login/Register (`frontend/src/pages/Login.tsx`, `Register.tsx`)**

- Email/password authentication
- JWT token storage
- Auto-redirect on success

**5. VaultIntro (`frontend/src/pages/VaultIntro.tsx`)**

- Commitment Vault explanation
- Deposit flow
- Contract signing

**6. Tribe (`frontend/src/pages/Tribe.tsx`)**

- Tribe overview
- Member list
- Collective discipline
- Join/create tribes

### Key Components

**FastingTimer (`frontend/src/components/FastingTimer.tsx`)**

- Real-time countdown
- Progress ring visualization
- Goal tracking
- Start/stop controls

**CortexWidget (`frontend/src/components/CortexWidget.tsx`)**

- Chat interface
- Message history
- Context-aware responses
- Floating widget

**VaultStatus (`frontend/src/components/VaultStatus.tsx`)**

- Deposit amount
- Earned refund
- Potential refund
- Progress bar

**MetricCard (`frontend/src/components/bio/MetricCard.tsx`)**

- Telemetry data display
- Trend visualization
- Manual logging

**KnowledgeHub (`frontend/src/components/community/KnowledgeHub.tsx`)**

- Educational articles
- Fasting science
- Keto guides

---

## Data Flow

### 1. User Registration Flow

```
User → Frontend (Register.tsx)
  ↓
  POST /api/v1/auth/register
  ↓
Backend (handlers.go → AuthService)
  ↓
  Create User in Repository
  ↓
  Return User object
  ↓
Frontend stores in state
```

### 2. Fasting Session Flow

```
User clicks "Start Fast"
  ↓
Frontend (Dashboard.tsx)
  ↓
  POST /api/v1/fasting/start
  ↓
Backend (FastingService)
  ↓
  Create FastingSession
  ↓
  Save to Repository
  ↓
  Return session
  ↓
Frontend updates timer
  ↓
User clicks "Stop Fast"
  ↓
  POST /api/v1/fasting/stop
  ↓
Backend calculates duration
  ↓
  Update discipline (+5 or -10)
  ↓
  Recalculate price
  ↓
  Return updated session
```

### 3. Meal Logging Flow

```
User uploads photo
  ↓
Frontend (Progress.tsx)
  ↓
  Convert to base64
  ↓
  POST /api/v1/meals/
  ↓
Backend (MealService)
  ↓
  Call CortexService.AnalyzeMeal()
  ↓
DeepSeek LLM analyzes description
  ↓
  Parse response
  ↓
  Save meal with analysis
  ↓
  Add $0.50 to earned refund
  ↓
  Return meal
  ↓
Frontend displays analysis
```

### 4. Recipe Discovery Flow

```
User navigates to Community → Recipes
  ↓
Frontend (Community.tsx)
  ↓
  GET /api/v1/recipes/?diet=vegan
  ↓
Backend (RecipeService)
  ↓
  Filter recipes by diet
  ↓
  Return filtered list
  ↓
Frontend displays recipe cards
```

---

## Configuration

### Environment Variables

**Backend**:

```bash
# Required
DEEPSEEK_API_KEY=sk-xxxxx
JWT_SECRET=your-secret-key

# Optional
DSN=user:pass@tcp(db:3306)/fastinghero
PORT=8080
ALLOWED_ORIGIN=https://yourdomain.com
```

**Frontend**:

```bash
VITE_API_URL=http://localhost:8080/api/v1
```

### Build Commands

**Backend**:

```bash
# Development
go run cmd/server/main.go

# Production
go build -o fastinghero cmd/server/main.go
./fastinghero
```

**Frontend**:

```bash
# Development
npm run dev

# Production
npm run build
npm run preview
```

---

## Deployment

### Docker Compose

```yaml
version: '3.8'
services:
  db:
    image: mariadb:10.11
    environment:
      MYSQL_DATABASE: fastinghero
      MYSQL_USER: fhuser
      MYSQL_PASSWORD: fhpass
    
  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      DSN: "fhuser:fhpass@tcp(db:3306)/fastinghero"
      JWT_SECRET: "${JWT_SECRET}"
      DEEPSEEK_API_KEY: "${DEEPSEEK_API_KEY}"
    
  frontend:
    build: ./frontend
    ports:
      - "80:80"
```

---

## Production Readiness

### Critical Issues

1. **Password Security** ⚠️
   - Currently stored in plain text
   - **Fix**: Implement bcrypt hashing

2. **JWT Secret** ⚠️
   - Hardcoded in source
   - **Fix**: Use environment variable

3. **Database** ⚠️
   - In-memory (data lost on restart)
   - **Fix**: Migrate to MariaDB/PostgreSQL

4. **Payment Integration** ⚠️
   - Commitment Vault is simulated
   - **Fix**: Integrate Stripe/PayPal

5. **Image Analysis** ⚠️
   - DeepSeek V2 is text-only
   - **Fix**: Upgrade to vision-capable model

### Recommended Timeline

- **Phase 1** (2-3 weeks): Security & Database
- **Phase 2** (2-3 weeks): Payment Integration
- **Phase 3** (1-2 weeks): Testing & Polish
- **Phase 4** (1 week): Soft Launch
- **Phase 5**: Public Launch

**Total**: 6-9 weeks to production

---

**End of Documentation**
