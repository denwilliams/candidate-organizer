package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/candidate-organizer/backend/internal/auth"
	"github.com/candidate-organizer/backend/internal/config"
	"github.com/candidate-organizer/backend/internal/errors"
	"github.com/candidate-organizer/backend/internal/models"
	"github.com/candidate-organizer/backend/internal/repository"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	userRepo    repository.UserRepository
	oauthConfig *auth.OAuthConfig
	jwtManager  *auth.JWTManager
	frontendURL string
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	userRepo repository.UserRepository,
	cfg *config.Config,
) *AuthHandler {
	oauthConfig := auth.NewGoogleOAuthConfig(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
		cfg.WorkspaceDomain,
	)

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, 24*time.Hour)

	return &AuthHandler{
		userRepo:    userRepo,
		oauthConfig: oauthConfig,
		jwtManager:  jwtManager,
		frontendURL: cfg.FrontendURL,
	}
}

// GoogleLogin initiates the Google OAuth flow
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Generate state token for CSRF protection
	state, err := generateStateToken()
	if err != nil {
		errors.WriteError(w, errors.NewInternalServerError("Failed to generate state token", err))
		return
	}

	// Store state in cookie for validation in callback
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   600, // 10 minutes
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to Google OAuth
	url := h.oauthConfig.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the OAuth callback from Google
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state token
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		h.redirectToFrontendWithError(w, r, "Missing state cookie")
		return
	}

	state := r.URL.Query().Get("state")
	if state != stateCookie.Value {
		h.redirectToFrontendWithError(w, r, "Invalid state token")
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Get authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		h.redirectToFrontendWithError(w, r, "Missing authorization code")
		return
	}

	// Exchange code for token
	ctx := r.Context()
	token, err := h.oauthConfig.ExchangeCode(ctx, code)
	if err != nil {
		h.redirectToFrontendWithError(w, r, "Failed to exchange code for token")
		return
	}

	// Get user info from Google
	userInfo, err := h.oauthConfig.GetUserInfo(ctx, token)
	if err != nil {
		h.redirectToFrontendWithError(w, r, "Failed to get user info")
		return
	}

	// Validate workspace domain
	if err := h.oauthConfig.ValidateWorkspaceDomain(userInfo); err != nil {
		h.redirectToFrontendWithError(w, r, "Unauthorized: "+err.Error())
		return
	}

	// Get or create user
	user, err := h.getOrCreateUser(ctx, userInfo)
	if err != nil {
		h.redirectToFrontendWithError(w, r, "Failed to create user")
		return
	}

	// Generate JWT token
	jwtToken, err := h.jwtManager.GenerateToken(user)
	if err != nil {
		h.redirectToFrontendWithError(w, r, "Failed to generate token")
		return
	}

	// Set JWT token in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		Path:     "/",
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	})

	// Redirect to frontend with success
	redirectURL := fmt.Sprintf("%s/auth/callback?success=true", h.frontendURL)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		errors.WriteError(w, errors.NewUnauthorizedError("User not found in context"))
		return
	}

	errors.WriteJSON(w, http.StatusOK, user)
}

// Logout logs out the user by clearing the auth token cookie
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear auth token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	errors.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

// getOrCreateUser gets an existing user or creates a new one
func (h *AuthHandler) getOrCreateUser(ctx context.Context, userInfo *auth.GoogleUserInfo) (*models.User, error) {
	// Try to get existing user
	user, err := h.userRepo.GetByEmail(ctx, userInfo.Email)
	if err != nil {
		return nil, err
	}

	// If user exists, return it
	if user != nil {
		return user, nil
	}

	// Check if this is the first user
	isFirst, err := h.userRepo.IsFirstUser(ctx)
	if err != nil {
		return nil, err
	}

	// Create new user
	role := "user"
	if isFirst {
		role = "admin" // First user becomes admin
	}

	// Extract workspace domain from email or use HD field
	workspaceDomain := userInfo.HD
	if workspaceDomain == "" {
		// Fallback to email domain
		emailParts := []rune(userInfo.Email)
		for i, c := range emailParts {
			if c == '@' && i < len(emailParts)-1 {
				workspaceDomain = string(emailParts[i+1:])
				break
			}
		}
	}

	newUser := &models.User{
		Email:           userInfo.Email,
		Name:            userInfo.Name,
		Role:            role,
		WorkspaceDomain: workspaceDomain,
	}

	if err := h.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// redirectToFrontendWithError redirects to the frontend with an error message
func (h *AuthHandler) redirectToFrontendWithError(w http.ResponseWriter, r *http.Request, errMsg string) {
	redirectURL := fmt.Sprintf("%s/auth/callback?error=%s", h.frontendURL, errMsg)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// generateStateToken generates a random state token for CSRF protection
func generateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RefreshToken refreshes the JWT token
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		errors.WriteError(w, errors.NewUnauthorizedError("User not found in context"))
		return
	}

	// Generate new JWT token
	jwtToken, err := h.jwtManager.GenerateToken(user)
	if err != nil {
		errors.WriteError(w, errors.NewInternalServerError("Failed to generate token", err))
		return
	}

	// Return token in response
	errors.WriteJSON(w, http.StatusOK, map[string]string{
		"token": jwtToken,
	})
}

// TokenResponse represents the response for token-related endpoints
type TokenResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// GetToken returns the current user's token and profile
func (h *AuthHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		errors.WriteError(w, errors.NewUnauthorizedError("User not found in context"))
		return
	}

	// Get token from cookie
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		errors.WriteError(w, errors.NewUnauthorizedError("No auth token found"))
		return
	}

	response := TokenResponse{
		Token: cookie.Value,
		User:  user,
	}

	json.NewEncoder(w).Encode(response)
}
