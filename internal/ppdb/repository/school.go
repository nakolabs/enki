package repository

import (
	"context"
	"github.com/google/uuid"
)

type School struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Level     string    `db:"level"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

func (r *repository) CreateSchool(ctx context.Context, school School) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO school (name, level) VALUES (:name, :level)`, school)
	if err != nil {
		return err
	}

	return nil
}
