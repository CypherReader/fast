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

// ============== WEEKLY REPORT TESTS ==============

func TestProgressAnalyzer_GenerateWeeklyReport_Success(t *testing.T) {
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)
	mockCortex := new(MockCortexServiceForMeal)

	analyzer := NewProgressAnalyzer(mockFastingRepo, mockUserRepo, mockCortex)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:              userID,
		Name:            "Test User",
		DisciplineIndex: 75,
	}

	// Create sessions from this week
	now := time.Now()
	sessions := []domain.FastingSession{
		{
			ID:                  uuid.New(),
			Status:              domain.StatusCompleted,
			StartTime:           now.Add(-48 * time.Hour),
			ActualDurationHours: 16,
		},
		{
			ID:                  uuid.New(),
			Status:              domain.StatusCompleted,
			StartTime:           now.Add(-24 * time.Hour),
			ActualDurationHours: 18,
		},
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockFastingRepo.On("FindByUserID", ctx, userID).Return(sessions, nil)
	mockCortex.On("Chat", ctx, userID, mock.AnythingOfType("string")).Return("Great week! Keep it up.", nil)

	report, err := analyzer.GenerateWeeklyReport(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 2, report.FastsCompleted)
	assert.Equal(t, 34.0, report.TotalFastingHours) // 16 + 18
	assert.Equal(t, 17.0, report.AverageDuration)   // 34 / 2
	assert.Equal(t, 18.0, report.LongestFast)
}

func TestProgressAnalyzer_GenerateWeeklyReport_NoSessions(t *testing.T) {
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)
	mockCortex := new(MockCortexServiceForMeal)

	analyzer := NewProgressAnalyzer(mockFastingRepo, mockUserRepo, mockCortex)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{ID: userID, Name: "Test User"}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockFastingRepo.On("FindByUserID", ctx, userID).Return([]domain.FastingSession{}, nil)
	mockCortex.On("Chat", ctx, userID, mock.AnythingOfType("string")).Return("No fasts this week", nil)

	report, err := analyzer.GenerateWeeklyReport(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 0, report.FastsCompleted)
	assert.Equal(t, 0.0, report.AverageDuration)
}

func TestProgressAnalyzer_GenerateWeeklyReport_UserNotFound(t *testing.T) {
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)
	mockCortex := new(MockCortexServiceForMeal)

	analyzer := NewProgressAnalyzer(mockFastingRepo, mockUserRepo, mockCortex)
	ctx := context.Background()
	userID := uuid.New()

	mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("not found"))

	report, err := analyzer.GenerateWeeklyReport(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, report)
}

func TestProgressAnalyzer_GenerateWeeklyReport_SessionsFetchError(t *testing.T) {
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)
	mockCortex := new(MockCortexServiceForMeal)

	analyzer := NewProgressAnalyzer(mockFastingRepo, mockUserRepo, mockCortex)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{ID: userID}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockFastingRepo.On("FindByUserID", ctx, userID).Return(nil, errors.New("db error"))

	report, err := analyzer.GenerateWeeklyReport(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, report)
}

// ============== GENERATE RECOMMENDATIONS TESTS ==============

func TestGenerateRecommendations_NoFasts(t *testing.T) {
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)
	mockCortex := new(MockCortexServiceForMeal)

	analyzer := NewProgressAnalyzer(mockFastingRepo, mockUserRepo, mockCortex)

	recs := analyzer.generateRecommendations(0, 0, map[string]int{})

	assert.NotEmpty(t, recs)
	assert.Contains(t, recs[0], "first fast")
}

func TestGenerateRecommendations_FewFasts(t *testing.T) {
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)
	mockCortex := new(MockCortexServiceForMeal)

	analyzer := NewProgressAnalyzer(mockFastingRepo, mockUserRepo, mockCortex)

	recs := analyzer.generateRecommendations(2, 16.0, map[string]int{"Monday": 1, "Tuesday": 1})

	assert.NotEmpty(t, recs)
	// Should recommend 3-4 fasts
	found := false
	for _, rec := range recs {
		if contains(rec, "3-4") {
			found = true
		}
	}
	assert.True(t, found)
}

func TestGenerateRecommendations_ShortFasts(t *testing.T) {
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)
	mockCortex := new(MockCortexServiceForMeal)

	analyzer := NewProgressAnalyzer(mockFastingRepo, mockUserRepo, mockCortex)

	recs := analyzer.generateRecommendations(5, 12.0, map[string]int{"Monday": 2, "Wednesday": 2, "Friday": 1})

	assert.NotEmpty(t, recs)
	// Should recommend extending fasts
}

// ============== PREDICT GOAL ACHIEVEMENT TESTS ==============

func TestPredictGoalAchievement_NoFasts(t *testing.T) {
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)
	mockCortex := new(MockCortexServiceForMeal)

	analyzer := NewProgressAnalyzer(mockFastingRepo, mockUserRepo, mockCortex)

	user := &domain.User{ID: uuid.New()}
	result := analyzer.predictGoalAchievement(user, 0, 0)

	assert.Empty(t, result) // No prediction for 0 fasts
}

func TestPredictGoalAchievement_WithFasts(t *testing.T) {
	mockFastingRepo := new(MockFastingRepository)
	mockUserRepo := new(MockUserRepository)
	mockCortex := new(MockCortexServiceForMeal)

	analyzer := NewProgressAnalyzer(mockFastingRepo, mockUserRepo, mockCortex)

	user := &domain.User{ID: uuid.New(), TargetWeightLbs: 170}
	result := analyzer.predictGoalAchievement(user, 16.0, 5)

	assert.NotEmpty(t, result) // Should have a prediction date
}

// Helper function for tests
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
