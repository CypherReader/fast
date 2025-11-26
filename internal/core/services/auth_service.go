package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo        ports.UserRepository
	referralService ports.ReferralService
	jwtSecret       []byte
}

func NewAuthService(userRepo ports.UserRepository, referralService ports.ReferralService, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		referralService: referralService,
		jwtSecret:       []byte(jwtSecret),
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string, referralCode string) (*domain.User, error) {
	// 1. Check if user exists
	if _, err := s.userRepo.FindByEmail(ctx, email); err == nil {
		return nil, errors.New("email already in use")
	}

	// 2. Hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Create User
	user := &domain.User{
		ID:               uuid.New(),
		Email:            email,
		PasswordHash:     string(hashedBytes),
		SubscriptionTier: domain.TierFree,
		DisciplineIndex:  0,
		CurrentPrice:     50.0,
		SignedContract:   false,
		CreatedAt:        time.Now(),
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	// 4. Track Referral
	if referralCode != "" && s.referralService != nil {
		if err := s.referralService.TrackReferral(ctx, referralCode, user.ID); err != nil {
			// Log error but don't fail registration
			// In a real app, use a logger
		}
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	// 1. Find User
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// 2. Check Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// 3. Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	return tokenString, "mock_refresh_token", nil
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user id in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	// Ideally, we might cache this or just trust the token if it contains enough info
	// For now, let's verify user still exists
	return s.userRepo.FindByID(ctx, userID)
}

func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}
