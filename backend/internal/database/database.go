package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"

	_ "github.com/lib/pq"
)

// DB wraps the sql.DB connection
type DB struct {
	*sql.DB
}

// New creates a new database connection
func New(databaseURL, schema string) (*DB, error) {
	// Set default schema if not provided
	if schema == "" {
		schema = "public"
	}

	// Append search_path to the connection string
	connURL, err := appendSearchPath(databaseURL, schema)
	if err != nil {
		return nil, fmt.Errorf("error configuring connection URL: %w", err)
	}

	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Printf("Database connection established (schema: %s)", schema)

	return &DB{db}, nil
}

// appendSearchPath adds the search_path parameter to the connection string
func appendSearchPath(databaseURL, schema string) (string, error) {
	// Parse the URL
	u, err := url.Parse(databaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid database URL: %w", err)
	}

	// Get existing query parameters
	q := u.Query()

	// Add or update the search_path parameter
	// Use options parameter which is the standard way for PostgreSQL
	options := q.Get("options")
	searchPathOption := fmt.Sprintf("-c search_path=%s", schema)

	if options != "" {
		// Append to existing options
		options = options + " " + searchPathOption
	} else {
		options = searchPathOption
	}

	q.Set("options", options)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
