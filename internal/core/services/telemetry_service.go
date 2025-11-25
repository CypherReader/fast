package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

type TelemetryService struct {
	repo ports.TelemetryRepository
}

func NewTelemetryService(repo ports.TelemetryRepository) *TelemetryService {
	return &TelemetryService{
		repo: repo,
	}
}

func (s *TelemetryService) ConnectDevice(ctx context.Context, userID uuid.UUID, source domain.TelemetrySource) (*domain.DeviceConnection, error) {
	// In a real app, this would handle OAuth flow
	conn := &domain.DeviceConnection{
		ID:          uuid.New(),
		UserID:      userID,
		Source:      source,
		ConnectedAt: time.Now(),
		Status:      "connected",
	}
	if err := s.repo.SaveConnection(ctx, conn); err != nil {
		return nil, err
	}
	return conn, nil
}

func (s *TelemetryService) DisconnectDevice(ctx context.Context, userID uuid.UUID, source domain.TelemetrySource) error {
	conn, err := s.repo.GetConnection(ctx, userID, source)
	if err != nil {
		return err
	}
	conn.Status = "disconnected"
	return s.repo.SaveConnection(ctx, conn)
}

func (s *TelemetryService) SyncData(ctx context.Context, userID uuid.UUID, source domain.TelemetrySource) error {
	// Mock sync logic
	conn, err := s.repo.GetConnection(ctx, userID, source)
	if err != nil {
		return errors.New("device not connected")
	}

	now := time.Now()
	conn.LastSyncAt = &now
	if err := s.repo.SaveConnection(ctx, conn); err != nil {
		return err
	}

	// Mock data ingestion
	// In reality, we'd fetch from Garmin/Apple API here
	data := &domain.TelemetryData{
		ID:         uuid.New(),
		UserID:     userID,
		Source:     source,
		Type:       domain.MetricSteps,
		Value:      8500, // Mock value
		Unit:       "steps",
		Timestamp:  now,
		IsManual:   false,
		TrustScore: 1.0,
	}
	return s.repo.SaveData(ctx, data)
}

func (s *TelemetryService) GetDeviceStatus(ctx context.Context, userID uuid.UUID) ([]domain.DeviceConnection, error) {
	return s.repo.ListConnections(ctx, userID)
}

func (s *TelemetryService) LogManualData(ctx context.Context, userID uuid.UUID, metricType domain.MetricType, value float64, unit string) (*domain.TelemetryData, error) {
	data := &domain.TelemetryData{
		ID:         uuid.New(),
		UserID:     userID,
		Source:     domain.TelemetrySourceManual,
		Type:       metricType,
		Value:      value,
		Unit:       unit,
		Timestamp:  time.Now(),
		IsManual:   true,
		TrustScore: 0.5, // Lower trust for manual entry
	}
	if err := s.repo.SaveData(ctx, data); err != nil {
		return nil, err
	}
	return data, nil
}
