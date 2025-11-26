package services

import (
	"context"
	"fastinghero/internal/core/ports"
	"fmt"
	"strings"

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
