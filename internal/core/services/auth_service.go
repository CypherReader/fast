package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"strings"
	"time"

	"fastinghero/pkg/logger"

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

func (s *AuthService) Register(ctx context.Context, email, password, name, referralCode string) (*domain.User, error) {
	email = strings.ToLower(email)

	// 1. Validate email format
	if !isValidEmail(email) {
		logger.Error().Str("email", email).Msg("Registration failed: invalid email format")
		return nil, errors.New("invalid email format")
	}

	// 2. Check if user exists
	if _, err := s.userRepo.FindByEmail(ctx, email); err == nil {
		logger.Error().Str("email", email).Msg("Registration failed: email already in use")
		return nil, errors.New("email already in use")
	}

	// 3. Enforce password policy
	if err := validatePasswordStrength(password); err != nil {
		logger.Error().Err(err).Msg("Registration failed: weak password")
		return nil, err
	}

	// 4. Hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(err).Msg("Registration failed: password hashing error")
		return nil, err
	}

	// 5. Create User
	user := &domain.User{
		ID:               uuid.New(),
		Email:            email,
		Name:             name,
		PasswordHash:     string(hashedBytes),
		SubscriptionTier: domain.TierFree,
		DisciplineIndex:  0,
		CurrentPrice:     50.0,
		SignedContract:   false,
		CreatedAt:        time.Now(),
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		logger.Error().Err(err).Msg("Registration failed: database save error")
		return nil, err
	}

	// 6. Track Referral
	if referralCode != "" && s.referralService != nil {
		if err := s.referralService.TrackReferral(ctx, referralCode, user.ID); err != nil {
			// Log error but don't fail registration
			logger.Error().Err(err).Msg("Failed to track referral")
		}
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, *domain.User, error) {
	email = strings.ToLower(email)

	// 1. Find User
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		logger.Error().Str("email", email).Err(err).Msg("Login failed: user not found")
		return "", "", nil, errors.New("invalid credentials")
	}

	// 2. Check Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		logger.Error().Str("email", email).Msg("Login failed: invalid password")
		return "", "", nil, errors.New("invalid credentials")
	}

	// 3. Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		logger.Error().Err(err).Msg("Login failed: token generation error")
		return "", "", nil, err
	}

	return tokenString, "mock_refresh_token", user, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Strictly enforce HS256 algorithm only
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Validate required claims exist
	userIDStr, ok := claims["user_id"].(string)
	if !ok || userIDStr == "" {
		return nil, errors.New("invalid user id in token")
	}

	// Validate expiration explicitly (defense in depth)
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, errors.New("token expired")
		}
	} else {
		return nil, errors.New("missing expiration claim")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	// Verify user still exists
	return s.userRepo.FindByID(ctx, userID)
}

func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}
