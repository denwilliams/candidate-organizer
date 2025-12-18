package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleUserInfo represents the user information from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	HD            string `json:"hd"` // Hosted domain for Google Workspace
}

// OAuthConfig holds OAuth configuration
type OAuthConfig struct {
	Config          *oauth2.Config
	WorkspaceDomain string
}

// NewGoogleOAuthConfig creates a new Google OAuth configuration
func NewGoogleOAuthConfig(clientID, clientSecret, redirectURL, workspaceDomain string) *OAuthConfig {
	return &OAuthConfig{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		WorkspaceDomain: workspaceDomain,
	}
}

// GetAuthURL generates the OAuth authorization URL
func (c *OAuthConfig) GetAuthURL(state string) string {
	// Add hd parameter if workspace domain is configured
	opts := []oauth2.AuthCodeOption{}
	if c.WorkspaceDomain != "" {
		opts = append(opts, oauth2.SetAuthURLParam("hd", c.WorkspaceDomain))
	}
	return c.Config.AuthCodeURL(state, opts...)
}

// ExchangeCode exchanges the authorization code for a token
func (c *OAuthConfig) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.Config.Exchange(ctx, code)
}

// GetUserInfo retrieves user information from Google
func (c *OAuthConfig) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := c.Config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: status %d, body: %s", resp.StatusCode, string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// ValidateWorkspaceDomain checks if the user's email belongs to the configured workspace domain
func (c *OAuthConfig) ValidateWorkspaceDomain(userInfo *GoogleUserInfo) error {
	if c.WorkspaceDomain == "" {
		// No domain restriction
		return nil
	}

	// Check hd field first (Hosted Domain)
	if userInfo.HD != "" && userInfo.HD == c.WorkspaceDomain {
		return nil
	}

	// Fallback to email domain check
	parts := strings.Split(userInfo.Email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email format")
	}

	emailDomain := parts[1]
	if emailDomain != c.WorkspaceDomain {
		return fmt.Errorf("email domain %s does not match required workspace domain %s", emailDomain, c.WorkspaceDomain)
	}

	return nil
}
