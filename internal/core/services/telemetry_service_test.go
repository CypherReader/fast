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

// MockTelemetryRepository is a mock of ports.TelemetryRepository
type MockTelemetryRepository struct {
	mock.Mock
}

func (m *MockTelemetryRepository) SaveConnection(ctx context.Context, conn *domain.DeviceConnection) error {
	args := m.Called(ctx, conn)
	return args.Error(0)
}

func (m *MockTelemetryRepository) GetConnection(ctx context.Context, userID uuid.UUID, source domain.TelemetrySource) (*domain.DeviceConnection, error) {
	args := m.Called(ctx, userID, source)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DeviceConnection), args.Error(1)
}

func (m *MockTelemetryRepository) ListConnections(ctx context.Context, userID uuid.UUID) ([]domain.DeviceConnection, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.DeviceConnection), args.Error(1)
}

func (m *MockTelemetryRepository) SaveData(ctx context.Context, data *domain.TelemetryData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockTelemetryRepository) GetLatestMetric(ctx context.Context, userID uuid.UUID, metricType domain.MetricType) (*domain.TelemetryData, error) {
	args := m.Called(ctx, userID, metricType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TelemetryData), args.Error(1)
}

func (m *MockTelemetryRepository) GetWeeklyStats(ctx context.Context, userID uuid.UUID, metricType domain.MetricType) ([]domain.DailyStat, error) {
	args := m.Called(ctx, userID, metricType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.DailyStat), args.Error(1)
}

// ============== CONNECT DEVICE TESTS ==============

func TestTelemetryService_ConnectDevice_Success(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	service := NewTelemetryService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("SaveConnection", ctx, mock.AnythingOfType("*domain.DeviceConnection")).Return(nil)

	conn, err := service.ConnectDevice(ctx, userID, domain.SourceAppleHealth)

	assert.NoError(t, err)
	assert.NotNil(t, conn)
	assert.Equal(t, domain.SourceAppleHealth, conn.Source)
	assert.Equal(t, "connected", conn.Status)
}

func TestTelemetryService_ConnectDevice_SaveError(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	service := NewTelemetryService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("SaveConnection", ctx, mock.AnythingOfType("*domain.DeviceConnection")).Return(errors.New("db error"))

	conn, err := service.ConnectDevice(ctx, userID, domain.SourceGarmin)

	assert.Error(t, err)
	assert.Nil(t, conn)
}

// ============== DISCONNECT DEVICE TESTS ==============

func TestTelemetryService_DisconnectDevice_Success(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	service := NewTelemetryService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	existingConn := &domain.DeviceConnection{
		ID:     uuid.New(),
		UserID: userID,
		Source: domain.SourceAppleHealth,
		Status: "connected",
	}

	mockRepo.On("GetConnection", ctx, userID, domain.SourceAppleHealth).Return(existingConn, nil)
	mockRepo.On("SaveConnection", ctx, existingConn).Return(nil)

	err := service.DisconnectDevice(ctx, userID, domain.SourceAppleHealth)

	assert.NoError(t, err)
	assert.Equal(t, "disconnected", existingConn.Status)
}

func TestTelemetryService_DisconnectDevice_NotFound(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	service := NewTelemetryService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("GetConnection", ctx, userID, domain.SourceAppleHealth).Return(nil, errors.New("not found"))

	err := service.DisconnectDevice(ctx, userID, domain.SourceAppleHealth)

	assert.Error(t, err)
}

// ============== SYNC DATA TESTS ==============

func TestTelemetryService_SyncData_Success(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	service := NewTelemetryService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	existingConn := &domain.DeviceConnection{
		ID:     uuid.New(),
		UserID: userID,
		Source: domain.SourceGarmin,
		Status: "connected",
	}

	mockRepo.On("GetConnection", ctx, userID, domain.SourceGarmin).Return(existingConn, nil)
	mockRepo.On("SaveConnection", ctx, existingConn).Return(nil)
	mockRepo.On("SaveData", ctx, mock.AnythingOfType("*domain.TelemetryData")).Return(nil)

	err := service.SyncData(ctx, userID, domain.SourceGarmin)

	assert.NoError(t, err)
	assert.NotNil(t, existingConn.LastSyncAt)
}

func TestTelemetryService_SyncData_DeviceNotConnected(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	service := NewTelemetryService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("GetConnection", ctx, userID, domain.SourceAppleHealth).Return(nil, errors.New("not found"))

	err := service.SyncData(ctx, userID, domain.SourceAppleHealth)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "device not connected")
}

// ============== GET DEVICE STATUS TESTS ==============

func TestTelemetryService_GetDeviceStatus_Success(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	service := NewTelemetryService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	connections := []domain.DeviceConnection{
		{ID: uuid.New(), Source: domain.SourceAppleHealth, Status: "connected"},
		{ID: uuid.New(), Source: domain.SourceGarmin, Status: "disconnected"},
	}

	mockRepo.On("ListConnections", ctx, userID).Return(connections, nil)

	result, err := service.GetDeviceStatus(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

// ============== LOG MANUAL DATA TESTS ==============

func TestTelemetryService_LogManualData_Success(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	service := NewTelemetryService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("SaveData", ctx, mock.AnythingOfType("*domain.TelemetryData")).Return(nil)

	data, err := service.LogManualData(ctx, userID, domain.MetricSteps, 10000, "steps")

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, float64(10000), data.Value)
	assert.True(t, data.IsManual)
	assert.Equal(t, 0.5, data.TrustScore) // Manual has lower trust
}

func TestTelemetryService_LogManualData_Error(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	service := NewTelemetryService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("SaveData", ctx, mock.AnythingOfType("*domain.TelemetryData")).Return(errors.New("db error"))

	data, err := service.LogManualData(ctx, userID, domain.MetricSteps, 5000, "steps")

	assert.Error(t, err)
	assert.Nil(t, data)
}

// ============== GET LATEST METRIC TESTS ==============

func TestTelemetryService_GetLatestMetric_Success(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	service := NewTelemetryService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	latestData := &domain.TelemetryData{
		ID:        uuid.New(),
		Type:      domain.MetricSteps,
		Value:     8500,
		Timestamp: time.Now(),
	}

	mockRepo.On("GetLatestMetric", ctx, userID, domain.MetricSteps).Return(latestData, nil)

	result, err := service.GetLatestMetric(ctx, userID, domain.MetricSteps)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, float64(8500), result.Value)
}

// ============== GET WEEKLY STATS TESTS ==============

func TestTelemetryService_GetWeeklyStats_Success(t *testing.T) {
	mockRepo := new(MockTelemetryRepository)
	service := NewTelemetryService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	stats := []domain.DailyStat{
		{Date: "2024-01-01", Value: 8000},
		{Date: "2024-01-02", Value: 10000},
		{Date: "2024-01-03", Value: 7500},
	}

	mockRepo.On("GetWeeklyStats", ctx, userID, domain.MetricSteps).Return(stats, nil)

	result, err := service.GetWeeklyStats(ctx, userID, domain.MetricSteps)

	assert.NoError(t, err)
	assert.Len(t, result, 3)
}
