package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
