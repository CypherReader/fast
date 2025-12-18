package http

import (
	"context"
	"fastinghero/internal/core/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTribeService is a mock implementation of ports.TribeService
type MockTribeService struct {
	mock.Mock
}

func (m *MockTribeService) CreateTribe(ctx context.Context, userID string, req domain.CreateTribeRequest) (*domain.Tribe, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tribe), args.Error(1)
}

func (m *MockTribeService) GetTribe(ctx context.Context, tribeID string, currentUserID *string) (*domain.Tribe, error) {
	args := m.Called(ctx, tribeID, currentUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tribe), args.Error(1)
}

func (m *MockTribeService) ListTribes(ctx context.Context, query domain.ListTribesQuery, currentUserID *string) ([]domain.Tribe, int, error) {
	args := m.Called(ctx, query, currentUserID)
	return args.Get(0).([]domain.Tribe), args.Int(1), args.Error(2)
}

func (m *MockTribeService) JoinTribe(ctx context.Context, tribeID, userID string) error {
	args := m.Called(ctx, tribeID, userID)
	return args.Error(0)
}

func (m *MockTribeService) LeaveTribe(ctx context.Context, tribeID, userID string) error {
	args := m.Called(ctx, tribeID, userID)
	return args.Error(0)
}

func (m *MockTribeService) GetTribeMembers(ctx context.Context, tribeID string, limit, offset int) ([]domain.TribeMember, error) {
	args := m.Called(ctx, tribeID, limit, offset)
	return args.Get(0).([]domain.TribeMember), args.Error(1)
}

func (m *MockTribeService) GetMyTribes(ctx context.Context, userID string) ([]domain.Tribe, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Tribe), args.Error(1)
}

func (m *MockTribeService) GetTribeStats(ctx context.Context, tribeID string) (*domain.TribeStats, error) {
	args := m.Called(ctx, tribeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TribeStats), args.Error(1)
}

func TestRegisterTribesRoutes_JoinTribe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockTribeService)
	handler := NewTribeHandler(mockService)

	router := gin.New()
	api := router.Group("/api/v1")

	// Mock Auth Middleware that sets user_id
	authMiddleware := func(c *gin.Context) {
		c.Set("user_id", uuid.New())
		c.Next()
	}

	RegisterTribesRoutes(api, handler, authMiddleware)

	// Test Case: Join Tribe
	t.Run("Join Tribe Endpoint Exists", func(t *testing.T) {
		tribeID := "123"
		mockService.On("JoinTribe", mock.Anything, tribeID, mock.AnythingOfType("string")).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/tribes/123/join", nil)
		router.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusNotFound, w.Code, "Endpoint returns 404 Not Found")
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
