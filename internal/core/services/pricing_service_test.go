package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"testing"
)

func TestCalculatePrice(t *testing.T) {
	service := NewPricingService()
	ctx := context.Background()

	tests := []struct {
		name            string
		disciplineIndex float64
		expectedPrice   float64
	}{
		{"Zero Discipline", 0.0, 50.0},
		{"Full Discipline", 100.0, 1.0},
		{"Half Discipline", 50.0, 25.5}, // 50 - (50 * 0.49) = 50 - 24.5 = 25.5
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &domain.User{DisciplineIndex: tt.disciplineIndex}
			price := service.CalculatePrice(ctx, user)
			if price != tt.expectedPrice {
				t.Errorf("expected price %f, got %f", tt.expectedPrice, price)
			}
		})
	}
}

func TestUpdateDisciplineIndex(t *testing.T) {
	service := NewPricingService()
	ctx := context.Background()

	user := &domain.User{
		DisciplineIndex: 0,
		CurrentPrice:    50.0,
	}

	// Initial check
	if service.CalculatePrice(ctx, user) != 50.0 {
		t.Errorf("initial price should be 50.0")
	}

	// Simulate a completed fast (+1)
	service.UpdateDisciplineIndex(ctx, user, true, false)
	if user.DisciplineIndex != 1.0 {
		t.Errorf("expected index 1.0, got %f", user.DisciplineIndex)
	}

	// Simulate verified ketosis (+2)
	service.UpdateDisciplineIndex(ctx, user, false, true)
	if user.DisciplineIndex != 3.0 {
		t.Errorf("expected index 3.0, got %f", user.DisciplineIndex)
	}
}
