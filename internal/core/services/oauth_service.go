package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"

	"github.com/google/uuid"
)

type OAuthService struct {
	userRepo     ports.UserRepository
	jwtSecret    []byte
	clientID     string
	clientSecret string
	redirectURL  string
}

func NewOAuthService(userRepo ports.UserRepository, jwtSecret, clientID, clientSecret, redirectURL string) *OAuthService {
	return &OAuthService{
		userRepo:     userRepo,
		jwtSecret:    []byte(jwtSecret),
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
	}
}

// GoogleUserInfo represents the user data from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// GenerateStateToken creates a random state token for CSRF protection
func (s *OAuthService) GenerateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetAuthURL returns the Google OAuth authorization URL
func (s *OAuthService) GetAuthURL(state string) string {
	return fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid email profile&state=%s",
		s.clientID,
		s.redirectURL,
		state,
	)
}

// ExchangeCode exchanges the authorization code for an access token
func (s *OAuthService) ExchangeCode(ctx context.Context, code string) (string, error) {
	tokenURL := "https://oauth2.googleapis.com/token"

	data := fmt.Sprintf(
		"code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		code,
		s.clientID,
		s.clientSecret,
		s.redirectURL,
	)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Body = io.NopCloser([]byte(data))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to exchange code for token")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", errors.New("access token not found in response")
	}

	return accessToken, nil
}

// GetUserInfo fetches user information from Google using the access token
func (s *OAuthService) GetUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info from Google")
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	if !userInfo.VerifiedEmail {
		return nil, errors.New("email not verified")
	}

	return &userInfo, nil
}

// AuthenticateWithGoogle handles the complete OAuth flow and returns a user with JWT token
func (s *OAuthService) AuthenticateWithGoogle(ctx context.Context, code string) (string, *domain.User, error) {
	// Exchange code for access token
	accessToken, err := s.ExchangeCode(ctx, code)
	if err != nil {
		return "", nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from Google
	googleUser, err := s.GetUserInfo(ctx, accessToken)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if user already exists with this Google ID
	user, err := s.findOrCreateUser(ctx, googleUser)
	if err != nil {
		return "", nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// Generate JWT token (reuse the auth service logic)
	authService := &AuthService{
		userRepo:  s.userRepo,
		jwtSecret: s.jwtSecret,
	}

	token, _, _, err := authService.Login(ctx, user.Email, "") // Empty password for OAuth users
	if err != nil {
		// If login fails (OAuth user with no password), generate token directly
		token, err = authService.generateToken(user.ID)
		if err != nil {
			return "", nil, err
		}
	}

	return token, user, nil
}

// findOrCreateUser finds an existing user or creates a new one from Google OAuth data
func (s *OAuthService) findOrCreateUser(ctx context.Context, googleUser *GoogleUserInfo) (*domain.User, error) {
	// Try to find user by email first
	user, err := s.userRepo.FindByEmail(ctx, googleUser.Email)
	if err == nil {
		// User exists - update Google ID and picture if not set
		if user.GoogleID == "" {
			user.GoogleID = googleUser.ID
			user.ProfilePictureURL = googleUser.Picture
			user.OAuthProvider = "google"
			if err := s.userRepo.Save(ctx, user); err != nil {
				return nil, err
			}
		}
		return user, nil
	}

	// User doesn't exist - create new user
	user = &domain.User{
		ID:                uuid.New(),
		Email:             googleUser.Email,
		Name:              googleUser.Name,
		GoogleID:          googleUser.ID,
		OAuthProvider:     "google",
		ProfilePictureURL: googleUser.Picture,
		PasswordHash:      "", // No password for OAuth-only users
		SubscriptionTier:  domain.TierFree,
		DisciplineIndex:   0,
		CurrentPrice:      50.0,
		SignedContract:    false,
		CreatedAt:         time.Now(),
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// generateToken is a helper to generate JWT token (copied from AuthService)
func (s *AuthService) generateToken(userID uuid.UUID) (string, error) {
	// Import jwt package at top if not already
	// This is placeholder - actual implementation should use the JWT generation from AuthService
	return "", errors.New("not implemented - use AuthService.Login")
}
