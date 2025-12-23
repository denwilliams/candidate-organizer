package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/candidate-organizer/backend/internal/api"
	"github.com/candidate-organizer/backend/internal/config"
	"github.com/candidate-organizer/backend/internal/database"
	"github.com/candidate-organizer/backend/internal/repository"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	dbWrapper, err := database.New(cfg.DatabaseURL, cfg.PostgresSchema)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbWrapper.Close()

	log.Println("Successfully connected to database")

	// Get the underlying sql.DB
	db := dbWrapper.DB

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(db)
	jobRepo := repository.NewPostgresJobRepository(db)
	candidateRepo := repository.NewPostgresCandidateRepository(db)
	commentRepo := repository.NewPostgresCommentRepository(db)
	attributeRepo := repository.NewPostgresAttributeRepository(db)

	// Initialize API server
	server := api.NewServer(cfg, userRepo, jobRepo, candidateRepo, commentRepo, attributeRepo)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, server.Router()); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
