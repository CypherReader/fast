# FastingHero

## Configuration

Copy `.env.example` to `.env` and update the values:

```bash
cp .env.example .env
```

**Important**: You must generate a secure `JWT_SECRET` for production:

```bash
openssl rand -base64 32
```

 (Neuro-Fast)

> "Put Your Money Where Your Goals Are."

FastingHero is a high-stakes accountability application designed to gamify discipline. It uses a unique "Commitment Vault" model where you deposit funds upfront and earn them back through disciplined actions like fasting, logging meals, and staying active.

## 🧠 Core Features

### 1. The Commitment Vault (Financial Accountability)

- **The Deal**: $30/month subscription ($10 Base Fee + $20 Vault Deposit).
- **The Mechanism**: You deposit $20 into your Vault at the start of the month.
- **Earn It Back**: Every disciplined action unlocks a refund from your Vault.
  - **Log a Meal**: +$0.50
  - **Hit Step Goal**: +$0.50
  - **Perfect Week Streak**: +$5.00
- **Philosophy**: Loss aversion is a powerful motivator. Don't lose your own money.

### 2. Bio-Narrative Engine (Cortex 2.0)

- **Biological Timeline**: Visualizes your journey through Digestion, Glycogen Depletion, Ketosis, and Autophagy.
- **Real-Time Insights**: Generates scientifically accurate physiological facts based on your exact fasting duration.
- **Fasting Clock**: A beautiful, animated interface to track your fasts.

### 3. The Tribe

- **Social Accountability**: See where you rank against friends.
- **Nudge System**: Pay "Discipline Points" to nudge friends who are slacking off.

### 4. Cortex AI Coach

- **Always On**: A chat interface powered by DeepSeek to answer questions and provide tough love.

## 🛠️ Tech Stack

- **Backend**: Go (Golang) with Gin Framework.
- **Architecture**: Hexagonal Architecture (Ports & Adapters).
- **Frontend**: React, Vite, TypeScript, Tailwind CSS, shadcn/ui.
- **Database**: MariaDB (with In-Memory fallback for dev).
- **AI**: DeepSeek API.

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- DeepSeek API Key

### Backend Setup

```bash
# Set API Key
$env:DEEPSEEK_API_KEY="your-key-here"

# Run Server
go run cmd/server/main.go
```

Server runs on `http://localhost:8080`.

### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

App runs on `http://localhost:5173`.

## 🧪 Verification

- **Unit Tests**: `go test ./...`
- **Manual Verification**: Check `walkthrough.md` for detailed verification logs and screenshots.
