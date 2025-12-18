package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/candidate-organizer/backend/internal/api/handlers"
	appmiddleware "github.com/candidate-organizer/backend/internal/api/middleware"
	"github.com/candidate-organizer/backend/internal/auth"
	"github.com/candidate-organizer/backend/internal/config"
	"github.com/candidate-organizer/backend/internal/models"
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

			// User routes (admin only)
			r.Route("/users", func(r chi.Router) {
				r.Use(s.authMiddleware.RequireAdmin)
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

// User management handlers

// handleListUsers returns a list of all users (admin only)
func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.userRepo.List(r.Context())
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch users",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"users": users,
	})
}

// handlePromoteUser promotes a user to admin role (admin only)
func (s *Server) handlePromoteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "User ID is required",
		})
		return
	}

	// Promote the user to admin
	if err := s.userRepo.PromoteToAdmin(r.Context(), userID); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to promote user",
		})
		return
	}

	// Fetch the updated user to return
	user, err := s.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "User promoted but failed to fetch updated user",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User promoted to admin successfully",
		"user":    user,
	})
}

// Job posting handlers

// handleListJobs returns a list of job postings with pagination
func (s *Server) handleListJobs(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit := 20 // default
	offset := 0 // default

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := parseInt(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := parseInt(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	jobs, err := s.jobRepo.List(r.Context(), limit, offset)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch job postings",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"jobs":   jobs,
		"limit":  limit,
		"offset": offset,
	})
}

// handleCreateJob creates a new job posting
func (s *Server) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		Requirements string `json:"requirements"`
		Location     string `json:"location"`
		SalaryRange  string `json:"salary_range"`
		Status       string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// Validate required fields
	if req.Title == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Title is required",
		})
		return
	}

	// Set default status if not provided
	if req.Status == "" {
		req.Status = "draft"
	}

	// Validate status
	if req.Status != "draft" && req.Status != "open" && req.Status != "closed" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Status must be 'draft', 'open', or 'closed'",
		})
		return
	}

	// Get user from context
	user := r.Context().Value("user").(*models.User)

	job := &models.JobPosting{
		Title:        req.Title,
		Description:  req.Description,
		Requirements: req.Requirements,
		Location:     req.Location,
		SalaryRange:  req.SalaryRange,
		Status:       req.Status,
		CreatedBy:    user.ID,
	}

	if err := s.jobRepo.Create(r.Context(), job); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create job posting",
		})
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Job posting created successfully",
		"job":     job,
	})
}

// handleGetJob returns a single job posting by ID
func (s *Server) handleGetJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")
	if jobID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Job ID is required",
		})
		return
	}

	job, err := s.jobRepo.GetByID(r.Context(), jobID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch job posting",
		})
		return
	}

	if job == nil {
		respondJSON(w, http.StatusNotFound, map[string]string{
			"error": "Job posting not found",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"job": job,
	})
}

// handleUpdateJob updates an existing job posting
func (s *Server) handleUpdateJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")
	if jobID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Job ID is required",
		})
		return
	}

	// Check if job exists
	existingJob, err := s.jobRepo.GetByID(r.Context(), jobID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch job posting",
		})
		return
	}

	if existingJob == nil {
		respondJSON(w, http.StatusNotFound, map[string]string{
			"error": "Job posting not found",
		})
		return
	}

	var req struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		Requirements string `json:"requirements"`
		Location     string `json:"location"`
		SalaryRange  string `json:"salary_range"`
		Status       string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// Validate required fields
	if req.Title == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Title is required",
		})
		return
	}

	// Validate status
	if req.Status != "" && req.Status != "draft" && req.Status != "open" && req.Status != "closed" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Status must be 'draft', 'open', or 'closed'",
		})
		return
	}

	// Update job with new data
	existingJob.Title = req.Title
	existingJob.Description = req.Description
	existingJob.Requirements = req.Requirements
	existingJob.Location = req.Location
	existingJob.SalaryRange = req.SalaryRange
	if req.Status != "" {
		existingJob.Status = req.Status
	}

	if err := s.jobRepo.Update(r.Context(), existingJob); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to update job posting",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Job posting updated successfully",
		"job":     existingJob,
	})
}

// handleDeleteJob deletes a job posting
func (s *Server) handleDeleteJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")
	if jobID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Job ID is required",
		})
		return
	}

	// Check if job exists
	job, err := s.jobRepo.GetByID(r.Context(), jobID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch job posting",
		})
		return
	}

	if job == nil {
		respondJSON(w, http.StatusNotFound, map[string]string{
			"error": "Job posting not found",
		})
		return
	}

	if err := s.jobRepo.Delete(r.Context(), jobID); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete job posting",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Job posting deleted successfully",
	})
}

// Placeholder handlers - to be implemented
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

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
