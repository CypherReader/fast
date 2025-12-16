package http

import (
	"fastinghero/internal/core/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	oauthService *services.OAuthService
}

func NewOAuthHandler(oauthService *services.OAuthService) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
	}
}

// HandleGoogleLogin initiates the Google OAuth flow
func (h *OAuthHandler) HandleGoogleLogin(c *gin.Context) {
	// Generate state token for CSRF protection
	state, err := h.oauthService.GenerateStateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state token"})
		return
	}

	// Store state in cookie for verification
	c.SetCookie(
		"oauth_state",
		state,
		600, // 10 minutes
		"/",
		"",
		false, // TODO: Set to true in production with HTTPS
		true,  // httpOnly
	)

	// Get OAuth URL and redirect
	authURL := h.oauthService.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// HandleGoogleCallback handles the OAuth callback from Google
func (h *OAuthHandler) HandleGoogleCallback(c *gin.Context) {
	// Verify state token (CSRF protection)
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state token"})
		return
	}

	stateParam := c.Query("state")
	if stateParam == "" || stateParam != stateCookie {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state token mismatch"})
		return
	}

	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// Get authorization code
	code := c.Query("code")
	if code == "" {
		// Check if there was an error from Google
		if errorParam := c.Query("error"); errorParam != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "OAuth error: " + errorParam})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not found"})
		return
	}

	// Authenticate with Google
	token, user, err := h.oauthService.AuthenticateWithGoogle(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
		return
	}

	// Set cookie for the JWT token (optional - or return in JSON)
	c.SetCookie(
		"access_token",
		token,
		24*3600, // 24 hours
		"/",
		"",
		false, // TODO: Set to true in production with HTTPS
		true,  // httpOnly
	)

	// Redirect to frontend with token and user info
	// Frontend will store the token and redirect to dashboard or onboarding
	frontendURL := "http://localhost:5173"

	// Check for production frontend URL
	if prodURL := c.Request.Header.Get("Origin"); prodURL != "" && prodURL != "http://localhost:5173" {
		frontendURL = prodURL
	}

	// Determine redirect path based on onboarding status
	redirectPath := "/dashboard"
	if !user.OnboardingCompleted {
		redirectPath = "/onboarding"
	}

	// Redirect to frontend with token as URL parameter
	// Frontend will extract token from URL and store it
	c.Redirect(http.StatusTemporaryRedirect, frontendURL+redirectPath+"?token="+token)
}
