package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"math"
)

const (
	BasePrice   = 50.0
	PriceFloor  = 1.0
	IndexFactor = 0.49
)

type PricingService struct{}

func NewPricingService() *PricingService {
	return &PricingService{}
}

func (s *PricingService) CalculatePrice(ctx context.Context, user *domain.User) float64 {
	// Formula: Price = BasePrice ($50) - (DisciplineIndex * 0.49)
	// Example: Index 100 -> 50 - 49 = $1
	// Example: Index 0 -> 50 - 0 = $50

	discount := user.DisciplineIndex * IndexFactor
	price := BasePrice - discount

	if price < PriceFloor {
		return PriceFloor
	}

	// Round to 2 decimal places
	return math.Round(price*100) / 100
}

func (s *PricingService) UpdateDisciplineIndex(ctx context.Context, user *domain.User, completedFast bool, verifiedKetosis bool) {
	// Simple algorithm for MVP
	// Completed fast: +1
	// Verified ketosis: +2
	// Missed fast (not implemented yet): -2

	if completedFast {
		user.DisciplineIndex += 1
	}
	if verifiedKetosis {
		user.DisciplineIndex += 2
	}

	if user.DisciplineIndex > 100 {
		user.DisciplineIndex = 100
	}
	if user.DisciplineIndex < 0 {
		user.DisciplineIndex = 0
	}

	user.CurrentPrice = s.CalculatePrice(ctx, user)
}
