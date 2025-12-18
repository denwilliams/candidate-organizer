package repository

import (
	"context"
	"database/sql"

	"github.com/candidate-organizer/backend/internal/models"
)

// AttributeRepository defines the interface for candidate attribute operations
type AttributeRepository interface {
	Create(ctx context.Context, attribute *models.CandidateAttribute) error
	GetByID(ctx context.Context, id string) (*models.CandidateAttribute, error)
	ListByCandidate(ctx context.Context, candidateID string) ([]*models.CandidateAttribute, error)
	Update(ctx context.Context, attribute *models.CandidateAttribute) error
	Delete(ctx context.Context, id string) error
	DeleteByKey(ctx context.Context, candidateID, attributeKey string) error
}

// PostgresAttributeRepository implements AttributeRepository for PostgreSQL
type PostgresAttributeRepository struct {
	db *sql.DB
}

// NewPostgresAttributeRepository creates a new PostgresAttributeRepository
func NewPostgresAttributeRepository(db *sql.DB) *PostgresAttributeRepository {
	return &PostgresAttributeRepository{db: db}
}

func (r *PostgresAttributeRepository) Create(ctx context.Context, attribute *models.CandidateAttribute) error {
	query := `
		INSERT INTO candidate_attributes (candidate_id, attribute_key, attribute_value)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		attribute.CandidateID, attribute.AttributeKey, attribute.AttributeValue,
	).Scan(&attribute.ID, &attribute.CreatedAt, &attribute.UpdatedAt)
}

func (r *PostgresAttributeRepository) GetByID(ctx context.Context, id string) (*models.CandidateAttribute, error) {
	query := `
		SELECT id, candidate_id, attribute_key, attribute_value, created_at, updated_at
		FROM candidate_attributes
		WHERE id = $1
	`
	attribute := &models.CandidateAttribute{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&attribute.ID, &attribute.CandidateID, &attribute.AttributeKey,
		&attribute.AttributeValue, &attribute.CreatedAt, &attribute.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return attribute, err
}

func (r *PostgresAttributeRepository) ListByCandidate(ctx context.Context, candidateID string) ([]*models.CandidateAttribute, error) {
	query := `
		SELECT id, candidate_id, attribute_key, attribute_value, created_at, updated_at
		FROM candidate_attributes
		WHERE candidate_id = $1
		ORDER BY attribute_key ASC
	`
	rows, err := r.db.QueryContext(ctx, query, candidateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attributes []*models.CandidateAttribute
	for rows.Next() {
		attribute := &models.CandidateAttribute{}
		if err := rows.Scan(
			&attribute.ID, &attribute.CandidateID, &attribute.AttributeKey,
			&attribute.AttributeValue, &attribute.CreatedAt, &attribute.UpdatedAt,
		); err != nil {
			return nil, err
		}
		attributes = append(attributes, attribute)
	}
	return attributes, rows.Err()
}

func (r *PostgresAttributeRepository) Update(ctx context.Context, attribute *models.CandidateAttribute) error {
	query := `
		UPDATE candidate_attributes
		SET attribute_key = $1, attribute_value = $2
		WHERE id = $3
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		attribute.AttributeKey, attribute.AttributeValue, attribute.ID,
	).Scan(&attribute.UpdatedAt)
}

func (r *PostgresAttributeRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM candidate_attributes WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresAttributeRepository) DeleteByKey(ctx context.Context, candidateID, attributeKey string) error {
	query := `DELETE FROM candidate_attributes WHERE candidate_id = $1 AND attribute_key = $2`
	_, err := r.db.ExecContext(ctx, query, candidateID, attributeKey)
	return err
}
