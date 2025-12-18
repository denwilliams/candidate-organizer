package repository

import (
	"context"
	"database/sql"

	"github.com/candidate-organizer/backend/internal/models"
)

// CommentRepository defines the interface for comment operations
type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	GetByID(ctx context.Context, id string) (*models.Comment, error)
	ListByCandidate(ctx context.Context, candidateID string) ([]*models.Comment, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, id string) error
}

// PostgresCommentRepository implements CommentRepository for PostgreSQL
type PostgresCommentRepository struct {
	db *sql.DB
}

// NewPostgresCommentRepository creates a new PostgresCommentRepository
func NewPostgresCommentRepository(db *sql.DB) *PostgresCommentRepository {
	return &PostgresCommentRepository{db: db}
}

func (r *PostgresCommentRepository) Create(ctx context.Context, comment *models.Comment) error {
	query := `
		INSERT INTO comments (candidate_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		comment.CandidateID, comment.UserID, comment.Content,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
}

func (r *PostgresCommentRepository) GetByID(ctx context.Context, id string) (*models.Comment, error) {
	query := `
		SELECT c.id, c.candidate_id, c.user_id, u.name as user_name, c.content, c.created_at, c.updated_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = $1
	`
	comment := &models.Comment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&comment.ID, &comment.CandidateID, &comment.UserID,
		&comment.UserName, &comment.Content,
		&comment.CreatedAt, &comment.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return comment, err
}

func (r *PostgresCommentRepository) ListByCandidate(ctx context.Context, candidateID string) ([]*models.Comment, error) {
	query := `
		SELECT c.id, c.candidate_id, c.user_id, u.name as user_name, c.content, c.created_at, c.updated_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.candidate_id = $1
		ORDER BY c.created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, candidateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		if err := rows.Scan(
			&comment.ID, &comment.CandidateID, &comment.UserID,
			&comment.UserName, &comment.Content,
			&comment.CreatedAt, &comment.UpdatedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}

func (r *PostgresCommentRepository) Update(ctx context.Context, comment *models.Comment) error {
	query := `
		UPDATE comments
		SET content = $1
		WHERE id = $2
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query, comment.Content, comment.ID).
		Scan(&comment.UpdatedAt)
}

func (r *PostgresCommentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
