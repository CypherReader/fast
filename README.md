# FastingHero (Neuro-Fast)

**"We don't charge for features. We charge for lack of discipline."**

FastingHero is a high-stakes fasting application designed to gamify discipline. It uses a unique "Lazy Tax" model where your subscription price is inversely correlated with your discipline score.

## 🧠 Core Features

### 1. The Lazy Tax (Dynamic Pricing)

- **Base Price**: $50.00 / month.
- **Your Goal**: $0.00 / month.
- **Mechanism**: Every completed fast lowers your price. Quitting early raises it.
- **Philosophy**: Financial pain is a powerful motivator.

### 2. Bio-Narrative Engine (Cortex 2.0)

- **Biological Timeline**: Visualizes your journey through Digestion, Glycogen Depletion, Ketosis, and Autophagy.
- **DeepSeek AI Insights**: Generates real-time, scientifically accurate physiological facts based on your exact fasting duration.
- **"Focus" Mode**: A distraction-free timer interface.

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
