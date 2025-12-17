package services

import (
	"context"
	"fastinghero/internal/core/ports"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CortexService struct {
	llm         ports.LLMProvider
	fastingRepo ports.FastingRepository
	userRepo    ports.UserRepository
}

func NewCortexService(llm ports.LLMProvider, fastingRepo ports.FastingRepository, userRepo ports.UserRepository) *CortexService {
	return &CortexService{
		llm:         llm,
		fastingRepo: fastingRepo,
		userRepo:    userRepo,
	}
}

func (s *CortexService) Chat(ctx context.Context, userID uuid.UUID, message string) (string, error) {
	// 1. Fetch User Context
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user: %w", err)
	}

	// 2. Fetch Fasting Context
	activeFast, err := s.fastingRepo.FindActiveByUserID(ctx, userID)
	isFasting := activeFast != nil
	fastingDuration := ""
	if isFasting {
		// Calculate duration (mocked for now or calculated if we had time lib)
		fastingDuration = "currently fasting"
	} else {
		fastingDuration = "not currently fasting"
	}

	// 3. Construct System Prompt
	systemPrompt := fmt.Sprintf(`You are Cortex, a ruthless but fair AI fasting coach. 
	The user has a Discipline Index of %.1f/100. 
	Current Status: %s.
	
	Your goal is to motivate them to stay on track with their fasting goals. 
	If their discipline is low, be tougher. If high, be encouraging but demanding.
	Keep responses concise (under 50 words) and impactful. Do not be polite. Be effective.`,
		user.DisciplineIndex, fastingDuration)

	// 4. Call LLM
	response, err := s.llm.GenerateResponse(ctx, message, systemPrompt)
	if err != nil {
		return "", fmt.Errorf("llm error: %w", err)
	}

	return response, nil
}

func (s *CortexService) GenerateInsight(ctx context.Context, userID uuid.UUID, fastingHours float64) (string, error) {
	// 1. Fetch User Context (optional, but good for personalization)
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user: %w", err)
	}

	// 2. Construct System Prompt
	systemPrompt := fmt.Sprintf(`You are a biological narrator for a fasting app. 
	The user has been fasting for %.1f hours. 
	User Discipline Index: %.1f/100.
	
	Your task is to describe the physiological processes happening right now (e.g., autophagy, ketosis, glycogen depletion).
	Be scientific but motivating. 
	Output format: A single concise paragraph. Max 50 words.`,
		fastingHours, user.DisciplineIndex)

	// 3. Call LLM
	// We use a generic prompt for the user message since the system prompt contains all context
	userMessage := "Describe my current biological status."
	response, err := s.llm.GenerateResponse(ctx, userMessage, systemPrompt)
	if err != nil {
		return "", fmt.Errorf("llm error: %w", err)
	}

	return response, nil
}

func (s *CortexService) AnalyzeMeal(ctx context.Context, imageBase64, description string) (string, bool, bool, error) {
	// 1. Construct Prompt
	// Since DeepSeek V2 is text-only, we rely heavily on the user's description for now.
	// In a real multimodal scenario, the image would be primary.
	prompt := fmt.Sprintf(`Analyze this meal based on the user's description: "%s".
	1. Is this a real photo of food taken by a camera, or does it look like a screen capture/fake? (Authenticity)
	2. Estimate the carb content based on the description. Is it Keto-friendly (under 10g net carbs)?
	
	Output format:
	Analysis: [Brief description of food and carb estimate]
	Authenticity: [Verified/Suspicious]
	Keto-Friendly: [Yes/No]
	`, description)

	// 2. Call LLM
	// We pass the image still, in case the adapter supports it or for future proofing
	response, err := s.llm.AnalyzeImage(ctx, imageBase64, prompt)
	if err != nil {
		return "", false, false, fmt.Errorf("llm error: %w", err)
	}

	// 3. Parse Response (Simple string parsing for MVP)
	isAuthentic := true
	isKeto := true

	// Basic parsing logic
	lowerResp := strings.ToLower(response)
	if strings.Contains(lowerResp, "keto-friendly: no") {
		isKeto = false
	}
	if strings.Contains(lowerResp, "authenticity: suspicious") || strings.Contains(lowerResp, "fake") {
		isAuthentic = false
	}

	return response, isAuthentic, isKeto, nil
}

// GetFastingMilestoneInsight returns structured insight based on fasting duration
func (s *CortexService) GetFastingMilestoneInsight(ctx context.Context, userID uuid.UUID, hours float64) (map[string]interface{}, error) {
	// Determine milestone
	milestone := getMilestone(hours)

	// Construct prompt for structured response
	systemPrompt := `You are a fasting science expert. Provide insights about fasting in JSON format.
	Be scientific but motivating. Keep each field concise.`

	userMessage := fmt.Sprintf(`The user has been fasting for %.1f hours (milestone: %s).
	Provide a JSON response with:
	{
		"insight": "2-3 sentence description of what's happening in the body now",
		"benefits": ["benefit1", "benefit2", "benefit3"],
		"motivation": "One powerful motivational quote (under 15 words)"
	}
	
	Focus on the science at this milestone. Be specific about biological processes.`, hours, milestone)

	// Call LLM
	response, err := s.llm.GenerateResponse(ctx, userMessage, systemPrompt)
	if err != nil {
		return nil, fmt.Errorf("llm error: %w", err)
	}

	// For MVP, return a structured response manually parsed
	// In production, you'd parse the JSON response from DeepSeek
	result := map[string]interface{}{
		"hours":      hours,
		"milestone":  milestone,
		"insight":    response,
		"benefits":   extractBenefits(response, milestone),
		"motivation": extractMotivation(response),
	}

	return result, nil
}

// getMilestone identifies the fasting milestone
func getMilestone(hours float64) string {
	switch {
	case hours < 4:
		return "Early Stage"
	case hours >= 4 && hours < 8:
		return "4h - Glycogen Depletion"
	case hours >= 8 && hours < 12:
		return "8h - Fat Adaptation"
	case hours >= 12 && hours < 16:
		return "12h - Ketosis Begins"
	case hours >= 16 && hours < 20:
		return "16h - Peak Ketosis"
	case hours >= 20 && hours < 24:
		return "20h - Deep Autophagy"
	case hours >= 24:
		return "24h+ - Extended Fasting"
	default:
		return "Active Fasting"
	}
}

// extractBenefits extracts or generates benefits based on milestone
func extractBenefits(response string, milestone string) []string {
	// For MVP, return milestone-specific benefits
	// In production, parse from AI response
	benefitsMap := map[string][]string{
		"4h - Glycogen Depletion": {"Insulin levels dropping", "Starting fat burning", "Digestive rest"},
		"8h - Fat Adaptation":     {"Increased fat oxidation", "Stable energy", "Mental clarity"},
		"12h - Ketosis Begins":    {"Ketone production starts", "Enhanced focus", "Autophagy initiating"},
		"16h - Peak Ketosis":      {"Maximum fat burning", "Deep autophagy", "HGH boost"},
		"20h - Deep Autophagy":    {"Cellular renewal peak", "Anti-aging benefits", "Immune system reset"},
		"24h+ - Extended Fasting": {"Maximum autophagy", "Stem cell activation", "Deep healing"},
	}

	if benefits, ok := benefitsMap[milestone]; ok {
		return benefits
	}
	return []string{"Fat burning", "Mental clarity", "Cellular repair"}
}

// extractMotivation generates motivational message
func extractMotivation(response string) string {
	// For MVP, return generic but try to extract from response
	// In production, parse from AI response
	if strings.Contains(strings.ToLower(response), "keep") || strings.Contains(strings.ToLower(response), "you") {
		// Try to find a motivational sentence in the response
		sentences := strings.Split(response, ".")
		for _, sentence := range sentences {
			if len(strings.TrimSpace(sentence)) > 10 && len(strings.TrimSpace(sentence)) < 100 {
				return strings.TrimSpace(sentence)
			}
		}
	}
	return "Every hour fasting is a victory for your health!"
}

// CravingResponse contains structured help for hunger cravings
type CravingResponse struct {
	ImmediateAction   string   `json:"immediate_action"`
	DistractionIdea   string   `json:"distraction_idea"`
	BiologicalFact    string   `json:"biological_fact"`
	Motivation        string   `json:"motivation"`
	TimeRemaining     string   `json:"time_remaining,omitempty"`
	SupportStrategies []string `json:"support_strategies"`
}

// GetCravingHelp provides personalized support for hunger cravings
func (s *CortexService) GetCravingHelp(ctx context.Context, userID uuid.UUID, cravingDescription string) (*CravingResponse, error) {
	// 1. Fetch user context
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// 2. Fetch active fast context
	activeFast, _ := s.fastingRepo.FindActiveByUserID(ctx, userID)
	if activeFast == nil {
		// User not currently fasting
		return &CravingResponse{
			ImmediateAction: "Start a fast first! You can't beat cravings you don't face.",
			Motivation:      "Set your goal and commit.",
		}, nil
	}

	// 3. Calculate fast duration
	fastDuration := time.Since(activeFast.StartTime).Hours()
	hoursRemaining := float64(activeFast.GoalHours) - fastDuration

	// 4. Construct AI prompt
	systemPrompt := fmt.Sprintf(`You are an emergency fasting coach. The user is %.1f hours into a %d-hour fast and experiencing cravings.
Discipline Score: %.1f/100
Craving: %s

Respond in this EXACT format:
IMMEDIATE: [One 20-second action they can do RIGHT NOW]
DISTRACTION: [One 5-minute activity to redirect their focus]
SCIENCE: [One biological fact about what's happening in their body at this stage]
MOTIVATION: [Powerful one-liner under 15 words]

Be firm, direct, and supportive. No fluff. Total response under 100 words.`,
		fastDuration, activeFast.GoalHours, user.DisciplineIndex, cravingDescription)

	userMessage := "Help me fight this craving."

	// 5. Call LLM
	response, err := s.llm.GenerateResponse(ctx, userMessage, systemPrompt)
	if err != nil {
		// Fallback response if AI fails
		return &CravingResponse{
			ImmediateAction:   "Drink 16oz of water RIGHT NOW. Set a 5-minute timer.",
			DistractionIdea:   "Take a brisk 5-minute walk or do 20 pushups.",
			BiologicalFact:    fmt.Sprintf("At %.0f hours, your body is actively burning fat and producing ketones for energy.", fastDuration),
			Motivation:        "You're stronger than this craving.",
			TimeRemaining:     fmt.Sprintf("%.1f hours until your goal", hoursRemaining),
			SupportStrategies: []string{"Drink water", "Move your body", "Call a tribe member"},
		}, nil
	}

	// 6. Parse AI response
	cravingResp := parseCravingResponse(response, fastDuration, hoursRemaining)

	return cravingResp, nil
}

// parseCravingResponse extracts structured data from AI response
func parseCravingResponse(aiResponse string, fastDuration, hoursRemaining float64) *CravingResponse {
	lines := strings.Split(aiResponse, "\n")
	resp := &CravingResponse{
		TimeRemaining:     fmt.Sprintf("%.1f hours until your goal", hoursRemaining),
		SupportStrategies: []string{"Drink water", "5-minute walk", "Deep breathing"},
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "IMMEDIATE:") {
			resp.ImmediateAction = strings.TrimSpace(strings.TrimPrefix(line, "IMMEDIATE:"))
		} else if strings.HasPrefix(line, "DISTRACTION:") {
			resp.DistractionIdea = strings.TrimSpace(strings.TrimPrefix(line, "DISTRACTION:"))
		} else if strings.HasPrefix(line, "SCIENCE:") {
			resp.BiologicalFact = strings.TrimSpace(strings.TrimPrefix(line, "SCIENCE:"))
		} else if strings.HasPrefix(line, "MOTIVATION:") {
			resp.Motivation = strings.TrimSpace(strings.TrimPrefix(line, "MOTIVATION:"))
		}
	}

	// Fallbacks if parsing fails
	if resp.ImmediateAction == "" {
		resp.ImmediateAction = "Drink 16oz water immediately and wait 5 minutes."
	}
	if resp.Motivation == "" {
		resp.Motivation = "You're stronger than this."
	}

	return resp
}
