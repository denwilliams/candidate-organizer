package api

import (
	"encoding/json"
	"net/http"

	"github.com/candidate-organizer/backend/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server represents the API server
type Server struct {
	config *config.Config
}

// NewServer creates a new API server
func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
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
		// Auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Get("/google", s.handleGoogleLogin)
			r.Get("/callback", s.handleGoogleCallback)
			r.Post("/refresh", s.handleRefreshToken)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			// r.Use(s.authMiddleware) // TODO: Implement auth middleware

			// User routes
			r.Route("/users", func(r chi.Router) {
				r.Get("/", s.handleListUsers)
				r.Get("/me", s.handleGetCurrentUser)
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

	return r
}

// handleHealth is a health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// Placeholder handlers - to be implemented
func (s *Server) handleGoogleLogin(w http.ResponseWriter, r *http.Request)        { notImplemented(w) }
func (s *Server) handleGoogleCallback(w http.ResponseWriter, r *http.Request)     { notImplemented(w) }
func (s *Server) handleRefreshToken(w http.ResponseWriter, r *http.Request)       { notImplemented(w) }
func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request)          { notImplemented(w) }
func (s *Server) handleGetCurrentUser(w http.ResponseWriter, r *http.Request)     { notImplemented(w) }
func (s *Server) handlePromoteUser(w http.ResponseWriter, r *http.Request)        { notImplemented(w) }
func (s *Server) handleListJobs(w http.ResponseWriter, r *http.Request)           { notImplemented(w) }
func (s *Server) handleCreateJob(w http.ResponseWriter, r *http.Request)          { notImplemented(w) }
func (s *Server) handleGetJob(w http.ResponseWriter, r *http.Request)             { notImplemented(w) }
func (s *Server) handleUpdateJob(w http.ResponseWriter, r *http.Request)          { notImplemented(w) }
func (s *Server) handleDeleteJob(w http.ResponseWriter, r *http.Request)          { notImplemented(w) }
func (s *Server) handleListCandidates(w http.ResponseWriter, r *http.Request)     { notImplemented(w) }
func (s *Server) handleCreateCandidate(w http.ResponseWriter, r *http.Request)    { notImplemented(w) }
func (s *Server) handleUploadResume(w http.ResponseWriter, r *http.Request)       { notImplemented(w) }
func (s *Server) handleGetCandidate(w http.ResponseWriter, r *http.Request)       { notImplemented(w) }
func (s *Server) handleUpdateCandidate(w http.ResponseWriter, r *http.Request)    { notImplemented(w) }
func (s *Server) handleDeleteCandidate(w http.ResponseWriter, r *http.Request)    { notImplemented(w) }
func (s *Server) handleUpdateCandidateStatus(w http.ResponseWriter, r *http.Request) { notImplemented(w) }
func (s *Server) handleAddAttribute(w http.ResponseWriter, r *http.Request)       { notImplemented(w) }
func (s *Server) handleUpdateAttribute(w http.ResponseWriter, r *http.Request)    { notImplemented(w) }
func (s *Server) handleDeleteAttribute(w http.ResponseWriter, r *http.Request)    { notImplemented(w) }
func (s *Server) handleListComments(w http.ResponseWriter, r *http.Request)       { notImplemented(w) }
func (s *Server) handleAddComment(w http.ResponseWriter, r *http.Request)         { notImplemented(w) }
func (s *Server) handleUpdateComment(w http.ResponseWriter, r *http.Request)      { notImplemented(w) }
func (s *Server) handleDeleteComment(w http.ResponseWriter, r *http.Request)      { notImplemented(w) }
func (s *Server) handleGenerateSummary(w http.ResponseWriter, r *http.Request)    { notImplemented(w) }
func (s *Server) handleChat(w http.ResponseWriter, r *http.Request)               { notImplemented(w) }

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func notImplemented(w http.ResponseWriter) {
	respondJSON(w, http.StatusNotImplemented, map[string]string{"error": "Not implemented yet"})
}
