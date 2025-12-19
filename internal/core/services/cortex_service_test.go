package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLLMProvider is a mock of ports.LLMProvider
type MockLLMProvider struct {
	mock.Mock
}

func (m *MockLLMProvider) GenerateResponse(ctx context.Context, prompt, systemPrompt string) (string, error) {
	args := m.Called(ctx, prompt, systemPrompt)
	return args.String(0), args.Error(1)
}

func (m *MockLLMProvider) AnalyzeImage(ctx context.Context, imageBase64, prompt string) (string, error) {
	args := m.Called(ctx, imageBase64, prompt)
	return args.String(0), args.Error(1)
}

// ============== CHAT TESTS ==============

func TestCortexService_Chat_Success(t *testing.T) {
	mockLLM := new(MockLLMProvider)
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewCortexService(mockLLM, mockFastingRepo, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:              userID,
		Name:            "Test User",
		DisciplineIndex: 75,
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockFastingRepo.On("FindActiveByUserID", ctx, userID).Return(nil, nil)
	mockLLM.On("GenerateResponse", ctx, mock.Anything, mock.Anything).Return("This is a helpful response", nil)

	response, err := service.Chat(ctx, userID, "How can I stay motivated?")

	assert.NoError(t, err)
	assert.Contains(t, response, "helpful response")
}

func TestCortexService_Chat_LLMError(t *testing.T) {
	mockLLM := new(MockLLMProvider)
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewCortexService(mockLLM, mockFastingRepo, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{ID: userID, Name: "Test User"}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockFastingRepo.On("FindActiveByUserID", ctx, userID).Return(nil, nil)
	mockLLM.On("GenerateResponse", ctx, mock.Anything, mock.Anything).Return("", errors.New("LLM error"))

	response, err := service.Chat(ctx, userID, "Hello")

	assert.Error(t, err)
	assert.Empty(t, response)
}

// ============== GENERATE INSIGHT TESTS ==============

func TestCortexService_GenerateInsight_Success(t *testing.T) {
	mockLLM := new(MockLLMProvider)
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewCortexService(mockLLM, mockFastingRepo, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:              userID,
		Name:            "Test User",
		DisciplineIndex: 50,
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockLLM.On("GenerateResponse", ctx, mock.Anything, mock.Anything).Return("Great progress! You're entering ketosis.", nil)

	insight, err := service.GenerateInsight(ctx, userID, 14.5)

	assert.NoError(t, err)
	assert.Contains(t, insight, "ketosis")
}

// ============== ANALYZE MEAL TESTS ==============

func TestCortexService_AnalyzeMeal_Success(t *testing.T) {
	mockLLM := new(MockLLMProvider)
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewCortexService(mockLLM, mockFastingRepo, mockUserRepo)
	ctx := context.Background()

	// Mock returns: analysis, isAuthentic, isKetoFriendly
	mockLLM.On("AnalyzeImage", ctx, mock.Anything, mock.Anything).Return("This is a healthy keto meal with eggs and avocado. AUTHENTIC: true KETO: true", nil)

	analysis, isAuthentic, isKeto, err := service.AnalyzeMeal(ctx, "base64imagedata", "Eggs and avocado")

	assert.NoError(t, err)
	assert.NotEmpty(t, analysis)
	// Note: actual values depend on parsing logic
	_ = isAuthentic
	_ = isKeto
}

func TestCortexService_AnalyzeMeal_NoImage(t *testing.T) {
	mockLLM := new(MockLLMProvider)
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewCortexService(mockLLM, mockFastingRepo, mockUserRepo)
	ctx := context.Background()

	mockLLM.On("AnalyzeImage", ctx, "", mock.Anything).Return("Description-only analysis", nil)

	analysis, isAuthentic, isKeto, err := service.AnalyzeMeal(ctx, "", "Just a salad")

	assert.NoError(t, err)
	assert.NotEmpty(t, analysis)
	_ = isAuthentic
	_ = isKeto
}

// ============== GET FASTING MILESTONE INSIGHT TESTS ==============

func TestCortexService_GetFastingMilestoneInsight_Success(t *testing.T) {
	mockLLM := new(MockLLMProvider)
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewCortexService(mockLLM, mockFastingRepo, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{ID: userID, Name: "Test User"}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockLLM.On("GenerateResponse", ctx, mock.Anything, mock.Anything).Return("At 16 hours, you're entering ketosis. Your body is burning fat!", nil)

	insight, err := service.GetFastingMilestoneInsight(ctx, userID, 16.0)

	assert.NoError(t, err)
	assert.NotNil(t, insight)
	assert.NotEmpty(t, insight["milestone"])
}

// ============== GET CRAVING HELP TESTS ==============

func TestCortexService_GetCravingHelp_Success(t *testing.T) {
	mockLLM := new(MockLLMProvider)
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewCortexService(mockLLM, mockFastingRepo, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{ID: userID, Name: "Test User"}
	startTime := time.Now().Add(-10 * time.Hour)
	activeFast := &domain.FastingSession{
		ID:        uuid.New(),
		UserID:    userID,
		Status:    domain.StatusActive,
		StartTime: startTime,
		GoalHours: 16,
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockFastingRepo.On("FindActiveByUserID", ctx, userID).Return(activeFast, nil)
	mockLLM.On("GenerateResponse", ctx, mock.Anything, mock.Anything).Return("Drink a glass of water. This craving will pass in 15 minutes.", nil)

	help, err := service.GetCravingHelp(ctx, userID, "I really want pizza")

	assert.NoError(t, err)
	assert.NotNil(t, help)
}

func TestCortexService_GetCravingHelp_NoActiveFast(t *testing.T) {
	mockLLM := new(MockLLMProvider)
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewCortexService(mockLLM, mockFastingRepo, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{ID: userID, Name: "Test User"}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockFastingRepo.On("FindActiveByUserID", ctx, userID).Return(nil, nil)
	mockLLM.On("GenerateResponse", ctx, mock.Anything, mock.Anything).Return("General craving help...", nil)

	help, err := service.GetCravingHelp(ctx, userID, "I want ice cream")

	assert.NoError(t, err)
	assert.NotNil(t, help)
}

// ============== GET BREAKFAST RECOMMENDATIONS TESTS ==============

func TestCortexService_GetBreakFastRecommendations_ShortFast(t *testing.T) {
	mockLLM := new(MockLLMProvider)
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewCortexService(mockLLM, mockFastingRepo, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{ID: userID, Name: "Test User"}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockLLM.On("GenerateResponse", ctx, mock.Anything, mock.Anything).Return("For a short fast, you can break with a light meal.", nil)

	guide, err := service.GetBreakFastRecommendations(ctx, userID, 16.0)

	assert.NoError(t, err)
	assert.NotNil(t, guide)
	assert.Equal(t, 16.0, guide.FastDuration)
}

func TestCortexService_GetBreakFastRecommendations_LongFast(t *testing.T) {
	mockLLM := new(MockLLMProvider)
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewCortexService(mockLLM, mockFastingRepo, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{ID: userID, Name: "Test User"}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockLLM.On("GenerateResponse", ctx, mock.Anything, mock.Anything).Return("For an extended fast, reintroduce food slowly.", nil)

	guide, err := service.GetBreakFastRecommendations(ctx, userID, 48.0)

	assert.NoError(t, err)
	assert.NotNil(t, guide)
	assert.Equal(t, 48.0, guide.FastDuration)
}

// ============== GENERATE DAILY QUOTE TESTS ==============

func TestCortexService_GenerateDailyQuote_Success(t *testing.T) {
	mockLLM := new(MockLLMProvider)
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewCortexService(mockLLM, mockFastingRepo, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{ID: userID, Name: "Test User", DisciplineIndex: 80}
	sessions := []domain.FastingSession{
		{ID: uuid.New(), Status: domain.StatusCompleted},
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockFastingRepo.On("FindByUserID", ctx, userID).Return(sessions, nil)
	mockFastingRepo.On("FindActiveByUserID", ctx, userID).Return(nil, nil)
	mockLLM.On("GenerateResponse", ctx, mock.Anything, mock.Anything).Return("Every fast you complete is a victory. Keep going!", nil)

	quote, err := service.GenerateDailyQuote(ctx, userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, quote)
}

// ============== HELPER FUNCTION TESTS ==============

func TestGetMilestone(t *testing.T) {
	testCases := []struct {
		hours    float64
		expected string
	}{
		{2, "Early Stage"},
		{6, "4h - Glycogen Depletion"},
		{10, "8h - Fat Adaptation"},
		{14, "12h - Ketosis Begins"},
		{18, "16h - Peak Ketosis"},
		{22, "20h - Deep Autophagy"},
		{48, "24h+ - Extended Fasting"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := getMilestone(tc.hours)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestExtractMotivation(t *testing.T) {
	testCases := []struct {
		response string
		hasMotiv bool
	}{
		{"You're doing great! Keep it up!", true},
		{"", false},
	}

	for _, tc := range testCases {
		result := extractMotivation(tc.response)
		if tc.hasMotiv {
			assert.NotEmpty(t, result)
		}
	}
}

func TestExtractBenefits(t *testing.T) {
	response := "Fat burning is accelerating. Autophagy is starting."
	milestone := "Ketosis begins"

	benefits := extractBenefits(response, milestone)

	assert.NotEmpty(t, benefits)
}
