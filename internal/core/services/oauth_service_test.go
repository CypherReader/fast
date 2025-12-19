package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============== OAUTH SERVICE TESTS ==============

func TestOAuthService_GenerateStateToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOAuthService(mockUserRepo, "secret", "clientid", "clientsecret", "http://redirect")

	token, err := service.GenerateStateToken()

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Len(t, token, 44) // base64 encoded 32 bytes
}

func TestOAuthService_GenerateStateToken_Unique(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOAuthService(mockUserRepo, "secret", "clientid", "clientsecret", "http://redirect")

	token1, _ := service.GenerateStateToken()
	token2, _ := service.GenerateStateToken()

	assert.NotEqual(t, token1, token2) // Tokens should be unique
}

func TestOAuthService_GetAuthURL(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOAuthService(mockUserRepo, "secret", "test-client-id", "clientsecret", "http://example.com/callback")

	state := "test-state-123"
	url := service.GetAuthURL(state)

	assert.Contains(t, url, "accounts.google.com")
	assert.Contains(t, url, "test-client-id")
	assert.Contains(t, url, "test-state-123")
	assert.Contains(t, url, "http://example.com/callback")
}

func TestOAuthService_GetAuthURL_ContainsScopes(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOAuthService(mockUserRepo, "secret", "clientid", "clientsecret", "http://redirect")

	url := service.GetAuthURL("state123")

	assert.Contains(t, url, "scope=openid")
	assert.Contains(t, url, "email")
	assert.Contains(t, url, "profile")
}

func TestNewOAuthService(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOAuthService(mockUserRepo, "jwt-secret", "client-id", "client-secret", "http://redirect")

	assert.NotNil(t, service)
	assert.Equal(t, "client-id", service.clientID)
	assert.Equal(t, "client-secret", service.clientSecret)
	assert.Equal(t, "http://redirect", service.redirectURL)
}

// ============== findOrCreateUser TESTS ==============

func TestOAuthService_FindOrCreateUser_ExistingUser_NoGoogleID(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOAuthService(mockUserRepo, "secret", "clientid", "clientsecret", "http://redirect")
	ctx := context.Background()

	existingUser := &domain.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Name:     "Test User",
		GoogleID: "", // No Google ID yet
	}

	googleUser := &GoogleUserInfo{
		ID:            "google123",
		Email:         "test@example.com",
		VerifiedEmail: true,
		Name:          "Google Name",
		Picture:       "https://picture.url",
	}

	mockUserRepo.On("FindByEmail", ctx, "test@example.com").Return(existingUser, nil)
	mockUserRepo.On("Save", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := service.findOrCreateUser(ctx, googleUser)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "google123", user.GoogleID)
	assert.Equal(t, "https://picture.url", user.ProfilePictureURL)
	assert.Equal(t, "google", user.OAuthProvider)
	mockUserRepo.AssertExpectations(t)
}

func TestOAuthService_FindOrCreateUser_ExistingUser_HasGoogleID(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOAuthService(mockUserRepo, "secret", "clientid", "clientsecret", "http://redirect")
	ctx := context.Background()

	existingUser := &domain.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Name:     "Test User",
		GoogleID: "already-linked", // Already has Google ID
	}

	googleUser := &GoogleUserInfo{
		ID:            "google123",
		Email:         "test@example.com",
		VerifiedEmail: true,
		Name:          "Google Name",
	}

	mockUserRepo.On("FindByEmail", ctx, "test@example.com").Return(existingUser, nil)
	// Save should NOT be called since GoogleID is already set

	user, err := service.findOrCreateUser(ctx, googleUser)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "already-linked", user.GoogleID) // Should keep original
	mockUserRepo.AssertNotCalled(t, "Save")
}

func TestOAuthService_FindOrCreateUser_NewUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOAuthService(mockUserRepo, "secret", "clientid", "clientsecret", "http://redirect")
	ctx := context.Background()

	googleUser := &GoogleUserInfo{
		ID:            "google123",
		Email:         "newuser@example.com",
		VerifiedEmail: true,
		Name:          "New User",
		Picture:       "https://picture.url",
	}

	mockUserRepo.On("FindByEmail", ctx, "newuser@example.com").Return(nil, assert.AnError)
	mockUserRepo.On("Save", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := service.findOrCreateUser(ctx, googleUser)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "google123", user.GoogleID)
	assert.Equal(t, "newuser@example.com", user.Email)
	assert.Equal(t, "New User", user.Name)
	assert.Equal(t, "google", user.OAuthProvider)
	assert.Equal(t, domain.TierFree, user.SubscriptionTier)
	assert.Empty(t, user.PasswordHash) // OAuth users have no password
	mockUserRepo.AssertExpectations(t)
}

func TestOAuthService_FindOrCreateUser_SaveError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOAuthService(mockUserRepo, "secret", "clientid", "clientsecret", "http://redirect")
	ctx := context.Background()

	googleUser := &GoogleUserInfo{
		ID:            "google123",
		Email:         "newuser@example.com",
		VerifiedEmail: true,
		Name:          "New User",
	}

	mockUserRepo.On("FindByEmail", ctx, "newuser@example.com").Return(nil, assert.AnError)
	mockUserRepo.On("Save", ctx, mock.AnythingOfType("*domain.User")).Return(assert.AnError)

	user, err := service.findOrCreateUser(ctx, googleUser)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockUserRepo.AssertExpectations(t)
}
