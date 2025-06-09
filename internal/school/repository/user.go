package repository

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `db:"id"`
	Name  string    `db:"name"`
	Email string    `db:"email"`
}

func (r *repository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := "SELECT id, name, email FROM users WHERE id = $1"

	var user User
	err := r.db.GetContext(ctx, &user, query, id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
