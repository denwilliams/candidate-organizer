package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/candidate-organizer/backend/internal/models"
)

// CandidateRepository defines the interface for candidate operations
type CandidateRepository interface {
	Create(ctx context.Context, candidate *models.Candidate) error
	GetByID(ctx context.Context, id string) (*models.Candidate, error)
	List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*models.Candidate, error)
	Update(ctx context.Context, candidate *models.Candidate) error
	UpdateStatus(ctx context.Context, id, status string) error
	Delete(ctx context.Context, id string) error
}

// PostgresCandidateRepository implements CandidateRepository for PostgreSQL
type PostgresCandidateRepository struct {
	db *sql.DB
}

// NewPostgresCandidateRepository creates a new PostgresCandidateRepository
func NewPostgresCandidateRepository(db *sql.DB) *PostgresCandidateRepository {
	return &PostgresCandidateRepository{db: db}
}

func (r *PostgresCandidateRepository) Create(ctx context.Context, candidate *models.Candidate) error {
	parsedDataJSON, err := json.Marshal(candidate.ParsedData)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO candidates (name, email, phone, resume_url, parsed_data, status, salary_expectation, job_posting_id, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		candidate.Name, candidate.Email, candidate.Phone, candidate.ResumeURL,
		parsedDataJSON, candidate.Status, candidate.SalaryExpectation,
		nullStringOrNil(candidate.JobPostingID), candidate.CreatedBy,
	).Scan(&candidate.ID, &candidate.CreatedAt, &candidate.UpdatedAt)
}

func (r *PostgresCandidateRepository) GetByID(ctx context.Context, id string) (*models.Candidate, error) {
	query := `
		SELECT id, name, email, phone, resume_url, parsed_data, status, salary_expectation, job_posting_id, created_at, updated_at, created_by
		FROM candidates
		WHERE id = $1
	`
	candidate := &models.Candidate{}
	var parsedDataJSON []byte
	var jobPostingID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&candidate.ID, &candidate.Name, &candidate.Email, &candidate.Phone,
		&candidate.ResumeURL, &parsedDataJSON, &candidate.Status,
		&candidate.SalaryExpectation, &jobPostingID,
		&candidate.CreatedAt, &candidate.UpdatedAt, &candidate.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if jobPostingID.Valid {
		candidate.JobPostingID = jobPostingID.String
	}

	if len(parsedDataJSON) > 0 {
		if err := json.Unmarshal(parsedDataJSON, &candidate.ParsedData); err != nil {
			return nil, err
		}
	}

	return candidate, nil
}

func (r *PostgresCandidateRepository) List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*models.Candidate, error) {
	query := `
		SELECT id, name, email, phone, resume_url, parsed_data, status, salary_expectation, job_posting_id, created_at, updated_at, created_by
		FROM candidates
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidates []*models.Candidate
	for rows.Next() {
		candidate := &models.Candidate{}
		var parsedDataJSON []byte
		var jobPostingID sql.NullString

		if err := rows.Scan(
			&candidate.ID, &candidate.Name, &candidate.Email, &candidate.Phone,
			&candidate.ResumeURL, &parsedDataJSON, &candidate.Status,
			&candidate.SalaryExpectation, &jobPostingID,
			&candidate.CreatedAt, &candidate.UpdatedAt, &candidate.CreatedBy,
		); err != nil {
			return nil, err
		}

		if jobPostingID.Valid {
			candidate.JobPostingID = jobPostingID.String
		}

		if len(parsedDataJSON) > 0 {
			if err := json.Unmarshal(parsedDataJSON, &candidate.ParsedData); err != nil {
				return nil, err
			}
		}

		candidates = append(candidates, candidate)
	}
	return candidates, rows.Err()
}

func (r *PostgresCandidateRepository) Update(ctx context.Context, candidate *models.Candidate) error {
	parsedDataJSON, err := json.Marshal(candidate.ParsedData)
	if err != nil {
		return err
	}

	query := `
		UPDATE candidates
		SET name = $1, email = $2, phone = $3, resume_url = $4, parsed_data = $5, status = $6, salary_expectation = $7, job_posting_id = $8
		WHERE id = $9
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		candidate.Name, candidate.Email, candidate.Phone, candidate.ResumeURL,
		parsedDataJSON, candidate.Status, candidate.SalaryExpectation,
		nullStringOrNil(candidate.JobPostingID), candidate.ID,
	).Scan(&candidate.UpdatedAt)
}

func (r *PostgresCandidateRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `UPDATE candidates SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *PostgresCandidateRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM candidates WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Helper function to handle nullable strings
func nullStringOrNil(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
