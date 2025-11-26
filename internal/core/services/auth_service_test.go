package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"testing"

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
	password := "password123"

	// Expect FindByEmail to return error (user not found)
	mockRepo.On("FindByEmail", ctx, email).Return(nil, assert.AnError)

	// Expect Save to be called with a user object
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.User")).Run(func(args mock.Arguments) {
		user := args.Get(1).(*domain.User)
		assert.NotEqual(t, password, user.PasswordHash, "Password should be hashed")
		assert.NotEmpty(t, user.PasswordHash, "Password hash should not be empty")
	}).Return(nil)

	user, err := authService.Register(ctx, email, password, "")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	email := "test@example.com"
	password := "password123"

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
	}

	mockRepo.On("FindByEmail", ctx, email).Return(existingUser, nil)

	token, _, err := authService.Login(ctx, email, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, nil, "test-secret")
	ctx := context.Background()

	email := "test@example.com"
	password := "password123"
	wrongPassword := "wrongpassword"

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
	}

	mockRepo.On("FindByEmail", ctx, email).Return(existingUser, nil)

	token, _, err := authService.Login(ctx, email, wrongPassword)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}
