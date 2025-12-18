package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/candidate-organizer/backend/internal/api/handlers"
	appmiddleware "github.com/candidate-organizer/backend/internal/api/middleware"
	"github.com/candidate-organizer/backend/internal/auth"
	"github.com/candidate-organizer/backend/internal/config"
	"github.com/candidate-organizer/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server represents the API server
type Server struct {
	config         *config.Config
	userRepo       repository.UserRepository
	jobRepo        repository.JobRepository
	candidateRepo  repository.CandidateRepository
	commentRepo    repository.CommentRepository
	attributeRepo  repository.AttributeRepository
	authHandler    *handlers.AuthHandler
	authMiddleware *appmiddleware.AuthMiddleware
}

// NewServer creates a new API server
func NewServer(
	cfg *config.Config,
	userRepo repository.UserRepository,
	jobRepo repository.JobRepository,
	candidateRepo repository.CandidateRepository,
	commentRepo repository.CommentRepository,
	attributeRepo repository.AttributeRepository,
) *Server {
	// Create auth handler
	authHandler := handlers.NewAuthHandler(userRepo, cfg)

	// Create JWT manager and auth middleware
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, 24*time.Hour)
	authMiddleware := appmiddleware.NewAuthMiddleware(jwtManager, userRepo)

	return &Server{
		config:         cfg,
		userRepo:       userRepo,
		jobRepo:        jobRepo,
		candidateRepo:  candidateRepo,
		commentRepo:    commentRepo,
		attributeRepo:  attributeRepo,
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
	}
}

// Router sets up and returns the HTTP router
func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{s.config.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", s.handleHealth)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes (public and protected combined)
		r.Route("/auth", func(r chi.Router) {
			// Public auth routes (no middleware)
			r.Get("/google", s.authHandler.GoogleLogin)
			r.Get("/callback", s.authHandler.GoogleCallback)

			// Protected auth routes (with middleware)
			r.Group(func(r chi.Router) {
				r.Use(s.authMiddleware.Authenticate)
				r.Post("/refresh", s.authHandler.RefreshToken)
				r.Post("/logout", s.authHandler.Logout)
				r.Get("/me", s.authHandler.GetProfile)
			})
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			// Apply auth middleware to all routes in this group
			r.Use(s.authMiddleware.Authenticate)

			// User routes
			r.Route("/users", func(r chi.Router) {
				r.Get("/", s.handleListUsers)
				r.Post("/{id}/promote", s.handlePromoteUser)
			})

			// Job posting routes
			r.Route("/jobs", func(r chi.Router) {
				r.Get("/", s.handleListJobs)
				r.Post("/", s.handleCreateJob)
				r.Get("/{id}", s.handleGetJob)
				r.Put("/{id}", s.handleUpdateJob)
				r.Delete("/{id}", s.handleDeleteJob)
			})

			// Candidate routes
			r.Route("/candidates", func(r chi.Router) {
				r.Get("/", s.handleListCandidates)
				r.Post("/", s.handleCreateCandidate)
				r.Post("/upload", s.handleUploadResume)
				r.Get("/{id}", s.handleGetCandidate)
				r.Put("/{id}", s.handleUpdateCandidate)
				r.Delete("/{id}", s.handleDeleteCandidate)
				r.Put("/{id}/status", s.handleUpdateCandidateStatus)

				// Candidate attributes
				r.Post("/{id}/attributes", s.handleAddAttribute)
				r.Put("/{id}/attributes/{attrId}", s.handleUpdateAttribute)
				r.Delete("/{id}/attributes/{attrId}", s.handleDeleteAttribute)

				// Comments
				r.Get("/{id}/comments", s.handleListComments)
				r.Post("/{id}/comments", s.handleAddComment)
				r.Put("/{id}/comments/{commentId}", s.handleUpdateComment)
				r.Delete("/{id}/comments/{commentId}", s.handleDeleteComment)

				// AI features
				r.Post("/{id}/summary", s.handleGenerateSummary)
			})

			// AI chat
			r.Post("/chat", s.handleChat)
		})
	})

	// Serve static files from frontend build
	// This should come after all API routes
	staticDir := getStaticDir()
	if staticDir != "" {
		fileServer := http.FileServer(http.Dir(staticDir))
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			// Serve static files, but fallback to index.html for client-side routing
			path := filepath.Join(staticDir, r.URL.Path)

			// Check if the file exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				// If it's not a file, check if it's a directory with index.html
				indexPath := filepath.Join(path, "index.html")
				if _, err := os.Stat(indexPath); err == nil {
					http.ServeFile(w, r, indexPath)
					return
				}

				// If no file or directory exists, serve the root index.html for client-side routing
				rootIndex := filepath.Join(staticDir, "index.html")
				if _, err := os.Stat(rootIndex); err == nil {
					http.ServeFile(w, r, rootIndex)
					return
				}
			}

			// Remove leading slash for file server
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/")
			fileServer.ServeHTTP(w, r)
		})
	}

	return r
}

// getStaticDir returns the directory containing static frontend files
func getStaticDir() string {
	// Check for STATIC_DIR environment variable first
	if dir := os.Getenv("STATIC_DIR"); dir != "" {
		return dir
	}

	// Default locations to check
	dirs := []string{
		"./static",
		"../static",
		"./frontend/out",
		"../frontend/out",
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
	}

	return ""
}

// handleHealth is a health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// Placeholder handlers - to be implemented
func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request)               { notImplemented(w) }
func (s *Server) handlePromoteUser(w http.ResponseWriter, r *http.Request)             { notImplemented(w) }
func (s *Server) handleListJobs(w http.ResponseWriter, r *http.Request)                { notImplemented(w) }
func (s *Server) handleCreateJob(w http.ResponseWriter, r *http.Request)               { notImplemented(w) }
func (s *Server) handleGetJob(w http.ResponseWriter, r *http.Request)                  { notImplemented(w) }
func (s *Server) handleUpdateJob(w http.ResponseWriter, r *http.Request)               { notImplemented(w) }
func (s *Server) handleDeleteJob(w http.ResponseWriter, r *http.Request)               { notImplemented(w) }
func (s *Server) handleListCandidates(w http.ResponseWriter, r *http.Request)          { notImplemented(w) }
func (s *Server) handleCreateCandidate(w http.ResponseWriter, r *http.Request)         { notImplemented(w) }
func (s *Server) handleUploadResume(w http.ResponseWriter, r *http.Request)            { notImplemented(w) }
func (s *Server) handleGetCandidate(w http.ResponseWriter, r *http.Request)            { notImplemented(w) }
func (s *Server) handleUpdateCandidate(w http.ResponseWriter, r *http.Request)         { notImplemented(w) }
func (s *Server) handleDeleteCandidate(w http.ResponseWriter, r *http.Request)         { notImplemented(w) }
func (s *Server) handleUpdateCandidateStatus(w http.ResponseWriter, r *http.Request)   { notImplemented(w) }
func (s *Server) handleAddAttribute(w http.ResponseWriter, r *http.Request)            { notImplemented(w) }
func (s *Server) handleUpdateAttribute(w http.ResponseWriter, r *http.Request)         { notImplemented(w) }
func (s *Server) handleDeleteAttribute(w http.ResponseWriter, r *http.Request)         { notImplemented(w) }
func (s *Server) handleListComments(w http.ResponseWriter, r *http.Request)            { notImplemented(w) }
func (s *Server) handleAddComment(w http.ResponseWriter, r *http.Request)              { notImplemented(w) }
func (s *Server) handleUpdateComment(w http.ResponseWriter, r *http.Request)           { notImplemented(w) }
func (s *Server) handleDeleteComment(w http.ResponseWriter, r *http.Request)           { notImplemented(w) }
func (s *Server) handleGenerateSummary(w http.ResponseWriter, r *http.Request)         { notImplemented(w) }
func (s *Server) handleChat(w http.ResponseWriter, r *http.Request)                    { notImplemented(w) }

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func notImplemented(w http.ResponseWriter) {
	respondJSON(w, http.StatusNotImplemented, map[string]string{"error": "Not implemented yet"})
}
