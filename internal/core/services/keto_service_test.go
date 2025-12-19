package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockKetoRepository is a mock of ports.KetoRepository
type MockKetoRepository struct {
	mock.Mock
}

func (m *MockKetoRepository) Save(ctx context.Context, entry *domain.KetoEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

func (m *MockKetoRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.KetoEntry, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.KetoEntry), args.Error(1)
}

// ============== LOG ENTRY TESTS ==============

func TestKetoService_LogEntry_Success_SoftData(t *testing.T) {
	mockKetoRepo := new(MockKetoRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewKetoService(mockKetoRepo, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	// Soft data - no ketone/acetone levels (no premium required)
	entry := domain.KetoEntry{
		Source: domain.SourceManual,
	}

	mockKetoRepo.On("Save", ctx, mock.AnythingOfType("*domain.KetoEntry")).Return(nil)

	err := service.LogEntry(ctx, userID, entry)

	assert.NoError(t, err)
	mockKetoRepo.AssertExpectations(t)
}

func TestKetoService_LogEntry_HardData_PremiumUser(t *testing.T) {
	mockKetoRepo := new(MockKetoRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewKetoService(mockKetoRepo, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	ketoneLevel := 1.5
	entry := domain.KetoEntry{
		KetoneLevel: &ketoneLevel,
	}

	user := &domain.User{
		ID:                 userID,
		SubscriptionTier:   domain.TierVault,
		SubscriptionStatus: domain.SubStatusActive,
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockKetoRepo.On("Save", ctx, mock.AnythingOfType("*domain.KetoEntry")).Return(nil)

	err := service.LogEntry(ctx, userID, entry)

	assert.NoError(t, err)
}

func TestKetoService_LogEntry_HardData_NonPremiumUser(t *testing.T) {
	mockKetoRepo := new(MockKetoRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewKetoService(mockKetoRepo, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	ketoneLevel := 1.5
	entry := domain.KetoEntry{
		KetoneLevel: &ketoneLevel,
	}

	user := &domain.User{
		ID:                 userID,
		SubscriptionTier:   domain.TierFree,
		SubscriptionStatus: domain.SubStatusNone,
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)

	err := service.LogEntry(ctx, userID, entry)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "premium subscription required")
}

func TestKetoService_LogEntry_HardData_UserNotFound(t *testing.T) {
	mockKetoRepo := new(MockKetoRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewKetoService(mockKetoRepo, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	acetoneLevel := 0.5
	entry := domain.KetoEntry{
		AcetoneLevel: &acetoneLevel,
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("not found"))

	err := service.LogEntry(ctx, userID, entry)

	assert.Error(t, err)
}

func TestKetoService_LogEntry_SaveError(t *testing.T) {
	mockKetoRepo := new(MockKetoRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewKetoService(mockKetoRepo, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	entry := domain.KetoEntry{
		Source: domain.SourceManual,
	}

	mockKetoRepo.On("Save", ctx, mock.AnythingOfType("*domain.KetoEntry")).Return(errors.New("db error"))

	err := service.LogEntry(ctx, userID, entry)

	assert.Error(t, err)
}
