package services

import (
	"context"
	"fastinghero/internal/core/ports"
	"fmt"

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
