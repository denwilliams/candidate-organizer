package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/candidate-organizer/backend/internal/auth"
	"github.com/candidate-organizer/backend/internal/errors"
	"github.com/candidate-organizer/backend/internal/models"
	"github.com/candidate-organizer/backend/internal/repository"
)

// AuthMiddleware is middleware for authenticating requests
type AuthMiddleware struct {
	jwtManager *auth.JWTManager
	userRepo   repository.UserRepository
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtManager *auth.JWTManager, userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		userRepo:   userRepo,
	}
}

// Authenticate validates the JWT token and adds user to context
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get token from cookie first
		token := ""
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			token = cookie.Value
		}

		// If no cookie, try Authorization header
		if token == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				// Bearer token format
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
				}
			}
		}

		// If still no token, return unauthorized
		if token == "" {
			errors.WriteError(w, errors.NewUnauthorizedError("No authentication token provided"))
			return
		}

		// Validate token
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			errors.WriteError(w, errors.NewUnauthorizedError("Invalid or expired token"))
			return
		}

		// Get user from database
		user, err := m.userRepo.GetByID(r.Context(), claims.UserID)
		if err != nil {
			errors.WriteError(w, errors.NewUnauthorizedError("User not found"))
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin ensures the authenticated user is an admin
func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(*models.User)
		if !ok {
			errors.WriteError(w, errors.NewUnauthorizedError("User not found in context"))
			return
		}

		if user.Role != "admin" {
			errors.WriteError(w, errors.NewForbiddenError("Admin access required"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
