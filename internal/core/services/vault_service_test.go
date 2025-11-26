package services

import (
	"fastinghero/internal/core/domain"
	"testing"
)

func TestVaultService_CalculateVaultStatus(t *testing.T) {
	service := NewVaultService(nil, nil)

	tests := []struct {
		name            string
		user            *domain.User
		expectedDeposit float64
		expectedEarned  float64
		expectedRefund  float64
	}{
		{
			name: "Full earnings",
			user: &domain.User{
				VaultDeposit: 20.0,
				EarnedRefund: 20.0,
			},
			expectedDeposit: 20.0,
			expectedEarned:  20.0,
			expectedRefund:  20.0,
		},
		{
			name: "Partial earnings",
			user: &domain.User{
				VaultDeposit: 20.0,
				EarnedRefund: 10.0,
			},
			expectedDeposit: 20.0,
			expectedEarned:  10.0,
			expectedRefund:  10.0,
		},
		{
			name: "Over earnings",
			user: &domain.User{
				VaultDeposit: 20.0,
				EarnedRefund: 25.0,
			},
			expectedDeposit: 20.0,
			expectedEarned:  25.0,
			expectedRefund:  20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deposit, earned, refund := service.CalculateVaultStatus(tt.user)
			if deposit != tt.expectedDeposit {
				t.Errorf("expected deposit %f, got %f", tt.expectedDeposit, deposit)
			}
			if earned != tt.expectedEarned {
				t.Errorf("expected earned %f, got %f", tt.expectedEarned, earned)
			}
			if refund != tt.expectedRefund {
				t.Errorf("expected refund %f, got %f", tt.expectedRefund, refund)
			}
		})
	}
}

func TestVaultService_CalculateDailyEarning(t *testing.T) {
	service := NewVaultService(nil, nil)

	tests := []struct {
		name            string
		disciplineIndex int
		expectedEarning float64
	}{
		{"Zero Discipline", 0, 0.0},
		{"Full Discipline", 100, 2.0},
		{"Half Discipline", 50, 1.0},
		{"High Discipline", 80, 1.6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			earning := service.CalculateDailyEarning(tt.disciplineIndex)
			if earning != tt.expectedEarning {
				t.Errorf("expected earning %f, got %f", tt.expectedEarning, earning)
			}
		})
	}
}
