package ports

import (
	"context"
	"fastinghero/internal/core/domain"

	"github.com/google/uuid"
)

type TelemetryRepository interface {
	SaveData(ctx context.Context, data *domain.TelemetryData) error
	GetLatestMetric(ctx context.Context, userID uuid.UUID, metricType domain.MetricType) (*domain.TelemetryData, error)
	SaveConnection(ctx context.Context, conn *domain.DeviceConnection) error
	GetConnection(ctx context.Context, userID uuid.UUID, source domain.TelemetrySource) (*domain.DeviceConnection, error)
	ListConnections(ctx context.Context, userID uuid.UUID) ([]domain.DeviceConnection, error)
}

type TelemetryService interface {
	ConnectDevice(ctx context.Context, userID uuid.UUID, source domain.TelemetrySource) (*domain.DeviceConnection, error)
	DisconnectDevice(ctx context.Context, userID uuid.UUID, source domain.TelemetrySource) error
	SyncData(ctx context.Context, userID uuid.UUID, source domain.TelemetrySource) error
	GetDeviceStatus(ctx context.Context, userID uuid.UUID) ([]domain.DeviceConnection, error)
	LogManualData(ctx context.Context, userID uuid.UUID, metricType domain.MetricType, value float64, unit string) (*domain.TelemetryData, error)
}
