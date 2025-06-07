package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UserSchoolRole struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	RoleID    string    `db:"role_id"`
	SchoolID  uuid.UUID `db:"school_id"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

func (r *repository) GetFirstUserSchoolRoleByUserID(ctx context.Context, userID uuid.UUID) (*UserSchoolRole, error) {
	userSchoolRole := &UserSchoolRole{}
	err := r.db.GetContext(ctx, userSchoolRole, "SELECT id, user_id, school_id, role_id, created_at, updated_at FROM user_school WHERE user_id = $1 ORDER BY created_at ASC", userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to get user school role")
		return nil, err
	}
	return userSchoolRole, nil
}
