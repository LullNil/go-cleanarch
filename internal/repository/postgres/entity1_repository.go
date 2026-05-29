package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/LullNil/go-cleanarch/domain"
	"github.com/LullNil/go-cleanarch/domain/entity1"

	"github.com/lib/pq"
)

type entity1Repo struct {
	db *sql.DB
}

// NewEntity1Repository creates a new Entity1 PostgreSQL repository.
func NewEntity1Repository(db *sql.DB) *entity1Repo {
	return &entity1Repo{db: db}
}

// Ensure interface compliance
var _ entity1.Repository = (*entity1Repo)(nil)

// Save inserts a new Entity1 into the database.
func (r *entity1Repo) Save(ctx context.Context, e *entity1.Entity1) (int64, error) {
	const op = "repository.postgres.entity1.Save"

	query := `
		INSERT INTO entity1 (field1, field2, field3)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(ctx, query,
		e.Field1,
		e.Field2,
		e.Field3,
	).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // unique_violation
				return 0, domain.ErrAlreadyExists
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// GetByID retrieves an Entity1 by its ID.
func (r *entity1Repo) GetByID(ctx context.Context, id int64) (*entity1.Entity1, error) {
	const op = "repository.postgres.entity1.GetByID"

	query := `
		SELECT id, field1, field2, field3
		FROM entity1
		WHERE id = $1
	`

	var e entity1.Entity1
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID,
		&e.Field1,
		&e.Field2,
		&e.Field3,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &e, nil
}

// Update updates existing Entity1 in the database.
func (r *entity1Repo) Update(ctx context.Context, e *entity1.Entity1) error {
	const op = "repository.postgres.entity1.Update"

	query := `
		UPDATE entity1
		SET field1 = $1, field2 = $2, field3 = $3
		WHERE id = $4
	`

	res, err := r.db.ExecContext(ctx, query,
		e.Field1,
		e.Field2,
		e.Field3,
		e.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to check affected rows: %w", op, err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Delete deletes an Entity1 by its ID.
func (r *entity1Repo) Delete(ctx context.Context, id int64) error {
	const op = "repository.postgres.entity1.Delete"

	query := `DELETE FROM entity1 WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to check affected rows: %w", op, err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
