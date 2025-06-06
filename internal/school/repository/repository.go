package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateSchool(ctx context.Context, userID uuid.UUID, school School) error
	GetSchoolByID(ctx context.Context, id uuid.UUID) (*School, error)
	GetListSchool(ctx context.Context, userID uuid.UUID) ([]School, error)
	GetSchoolRoleByUserIDAndSchoolID(ctx context.Context, userID uuid.UUID, schoolID uuid.UUID) (*UserSchoolRole, error)
	DeleteSchool(ctx context.Context, schoolID uuid.UUID) error
	UpdateSchoolProfile(ctx context.Context, schoolID uuid.UUID, school School) (*School, error)
}

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}
