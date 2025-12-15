package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/suipic/backend/models"
)

type PostgresCommentRepository struct {
	db *sql.DB
}

func NewPostgresCommentRepository(db *sql.DB) *PostgresCommentRepository {
	return &PostgresCommentRepository{db: db}
}

func (r *PostgresCommentRepository) Create(ctx context.Context, comment *models.Comment) error {
	query := `
		INSERT INTO comments (photo_id, user_id, parent_comment_id, text, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		comment.PhotoID,
		comment.UserID,
		comment.ParentCommentID,
		comment.Text,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

func (r *PostgresCommentRepository) GetByID(ctx context.Context, id int) (*models.Comment, error) {
	query := `
		SELECT id, photo_id, user_id, parent_comment_id, text, created_at, updated_at
		FROM comments
		WHERE id = $1
	`
	comment := &models.Comment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.PhotoID,
		&comment.UserID,
		&comment.ParentCommentID,
		&comment.Text,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get comment by id: %w", err)
	}

	return comment, nil
}

func (r *PostgresCommentRepository) Update(ctx context.Context, comment *models.Comment) error {
	query := `
		UPDATE comments
		SET photo_id = $1, user_id = $2, parent_comment_id = $3, text = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		comment.PhotoID,
		comment.UserID,
		comment.ParentCommentID,
		comment.Text,
		comment.ID,
	).Scan(&comment.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("comment not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

func (r *PostgresCommentRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM comments WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

func (r *PostgresCommentRepository) List(ctx context.Context, limit, offset int) ([]*models.Comment, error) {
	query := `
		SELECT id, photo_id, user_id, parent_comment_id, text, created_at, updated_at
		FROM comments
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.PhotoID,
			&comment.UserID,
			&comment.ParentCommentID,
			&comment.Text,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %w", err)
	}

	return comments, nil
}

func (r *PostgresCommentRepository) GetByPhoto(ctx context.Context, photoID int) ([]*models.Comment, error) {
	query := `
		SELECT id, photo_id, user_id, parent_comment_id, text, created_at, updated_at
		FROM comments
		WHERE photo_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, photoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by photo: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.PhotoID,
			&comment.UserID,
			&comment.ParentCommentID,
			&comment.Text,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %w", err)
	}

	return comments, nil
}

func (r *PostgresCommentRepository) GetThreads(ctx context.Context, photoID int) ([]*models.Comment, error) {
	query := `
		SELECT id, photo_id, user_id, parent_comment_id, text, created_at, updated_at
		FROM comments
		WHERE photo_id = $1 AND parent_comment_id IS NULL
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, photoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment threads: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.PhotoID,
			&comment.UserID,
			&comment.ParentCommentID,
			&comment.Text,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %w", err)
	}

	return comments, nil
}

func (r *PostgresCommentRepository) GetReplies(ctx context.Context, parentCommentID int) ([]*models.Comment, error) {
	query := `
		SELECT id, photo_id, user_id, parent_comment_id, text, created_at, updated_at
		FROM comments
		WHERE parent_comment_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, parentCommentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment replies: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.PhotoID,
			&comment.UserID,
			&comment.ParentCommentID,
			&comment.Text,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %w", err)
	}

	return comments, nil
}
