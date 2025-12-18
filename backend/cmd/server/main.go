package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/candidate-organizer/backend/internal/api"
	"github.com/candidate-organizer/backend/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize API server
	server := api.NewServer(cfg)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, server.Router()); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
