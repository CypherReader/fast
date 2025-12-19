package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of ports.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByReferralCode(ctx context.Context, code string) (*domain.User, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestAuthService_Register_HashesPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	email := "test@example.com"
	password := "SecurePass123!" // Meets 12+ char requirement

	// Expect FindByEmail to return error (user not found)
	mockRepo.On("FindByEmail", ctx, email).Return(nil, assert.AnError)

	// Expect Save to be called with a user object
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.User")).Run(func(args mock.Arguments) {
		user := args.Get(1).(*domain.User)
		assert.NotEqual(t, password, user.PasswordHash, "Password should be hashed")
		assert.NotEmpty(t, user.PasswordHash, "Password hash should not be empty")
	}).Return(nil)

	user, err := authService.Register(ctx, email, password, "Test User", "")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	email := "test@example.com"
	password := "SecurePass123!"

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
	}

	mockRepo.On("FindByEmail", ctx, email).Return(existingUser, nil)

	token, _, _, err := authService.Login(ctx, email, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	email := "test@example.com"
	password := "SecurePass123!"
	wrongPassword := "WrongPass123!"

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
	}

	mockRepo.On("FindByEmail", ctx, email).Return(existingUser, nil)

	token, _, _, err := authService.Login(ctx, email, wrongPassword)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	email := "notfound@example.com"
	password := "SecurePass123!"

	mockRepo.On("FindByEmail", ctx, email).Return(nil, assert.AnError)

	token, _, _, err := authService.Login(ctx, email, password)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_InvalidEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	invalidEmails := []string{
		"notanemail",
		"missing@domain",
		"@nodomain.com",
		"spaces in@email.com",
		"",
	}

	for _, email := range invalidEmails {
		user, err := authService.Register(ctx, email, "SecurePass123!", "Test User", "")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid email format")
	}
}

func TestAuthService_Register_WeakPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	email := "valid@example.com"

	// Expect FindByEmail to return user not found
	mockRepo.On("FindByEmail", ctx, email).Return(nil, assert.AnError)

	weakPasswords := []string{
		"short",         // Too short
		"nouppercase1!", // No uppercase
		"NOLOWERCASE1!", // No lowercase
	}

	for _, password := range weakPasswords {
		user, err := authService.Register(ctx, email, password, "Test User", "")

		assert.Error(t, err)
		assert.Nil(t, user)
	}
}

func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	email := "existing@example.com"
	existingUser := &domain.User{
		ID:    uuid.New(),
		Email: email,
	}

	mockRepo.On("FindByEmail", ctx, email).Return(existingUser, nil)

	user, err := authService.Register(ctx, email, "SecurePass123!", "Test User", "")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email already in use", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	userID := uuid.New()
	user := &domain.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Create a valid token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("test-secret"))

	mockRepo.On("FindByID", ctx, userID).Return(user, nil)

	validatedUser, err := authService.ValidateToken(ctx, tokenString)

	assert.NoError(t, err)
	assert.NotNil(t, validatedUser)
	assert.Equal(t, userID, validatedUser.ID)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken_ExpiredToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	userID := uuid.New()

	// Create an expired token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
	})
	tokenString, _ := token.SignedString([]byte("test-secret"))

	validatedUser, err := authService.ValidateToken(ctx, tokenString)

	assert.Error(t, err)
	assert.Nil(t, validatedUser)
}

func TestAuthService_ValidateToken_InvalidSignature(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	userID := uuid.New()

	// Create a token with different secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("wrong-secret"))

	validatedUser, err := authService.ValidateToken(ctx, tokenString)

	assert.Error(t, err)
	assert.Nil(t, validatedUser)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestAuthService_ValidateToken_MalformedToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	validatedUser, err := authService.ValidateToken(ctx, "not.a.valid.token")

	assert.Error(t, err)
	assert.Nil(t, validatedUser)
}

func TestAuthService_ValidateToken_MissingUserID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	// Create a token without user_id
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("test-secret"))

	validatedUser, err := authService.ValidateToken(ctx, tokenString)

	assert.Error(t, err)
	assert.Nil(t, validatedUser)
	assert.Contains(t, err.Error(), "invalid user id")
}

func TestAuthService_GetUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	userID := uuid.New()
	user := &domain.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockRepo.On("FindByID", ctx, userID).Return(user, nil)

	result, err := authService.GetUserByID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.ID)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	userID := uuid.New()

	mockRepo.On("FindByID", ctx, userID).Return(nil, assert.AnError)

	result, err := authService.GetUserByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
