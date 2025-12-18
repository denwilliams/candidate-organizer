package repository

import (
	"context"
	"database/sql"

	"github.com/candidate-organizer/backend/internal/models"
)

// JobRepository defines the interface for job posting operations
type JobRepository interface {
	Create(ctx context.Context, job *models.JobPosting) error
	GetByID(ctx context.Context, id string) (*models.JobPosting, error)
	List(ctx context.Context, limit, offset int) ([]*models.JobPosting, error)
	Update(ctx context.Context, job *models.JobPosting) error
	Delete(ctx context.Context, id string) error
}

// PostgresJobRepository implements JobRepository for PostgreSQL
type PostgresJobRepository struct {
	db *sql.DB
}

// NewPostgresJobRepository creates a new PostgresJobRepository
func NewPostgresJobRepository(db *sql.DB) *PostgresJobRepository {
	return &PostgresJobRepository{db: db}
}

func (r *PostgresJobRepository) Create(ctx context.Context, job *models.JobPosting) error {
	query := `
		INSERT INTO job_postings (title, description, requirements, location, salary_range, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		job.Title, job.Description, job.Requirements, job.Location,
		job.SalaryRange, job.Status, job.CreatedBy,
	).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)
}

func (r *PostgresJobRepository) GetByID(ctx context.Context, id string) (*models.JobPosting, error) {
	query := `
		SELECT id, title, description, requirements, location, salary_range, status, created_at, updated_at, created_by
		FROM job_postings
		WHERE id = $1
	`
	job := &models.JobPosting{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&job.ID, &job.Title, &job.Description, &job.Requirements, &job.Location,
		&job.SalaryRange, &job.Status, &job.CreatedAt, &job.UpdatedAt, &job.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return job, err
}

func (r *PostgresJobRepository) List(ctx context.Context, limit, offset int) ([]*models.JobPosting, error) {
	query := `
		SELECT id, title, description, requirements, location, salary_range, status, created_at, updated_at, created_by
		FROM job_postings
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.JobPosting
	for rows.Next() {
		job := &models.JobPosting{}
		if err := rows.Scan(
			&job.ID, &job.Title, &job.Description, &job.Requirements, &job.Location,
			&job.SalaryRange, &job.Status, &job.CreatedAt, &job.UpdatedAt, &job.CreatedBy,
		); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, rows.Err()
}

func (r *PostgresJobRepository) Update(ctx context.Context, job *models.JobPosting) error {
	query := `
		UPDATE job_postings
		SET title = $1, description = $2, requirements = $3, location = $4, salary_range = $5, status = $6
		WHERE id = $7
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		job.Title, job.Description, job.Requirements, job.Location,
		job.SalaryRange, job.Status, job.ID,
	).Scan(&job.UpdatedAt)
}

func (r *PostgresJobRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM job_postings WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
