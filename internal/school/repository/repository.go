package repository

import (
	"context"
	"enuma-elish/internal/school/service/data/request"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateSchool(ctx context.Context, userID uuid.UUID, school School) error
	GetSchoolByID(ctx context.Context, id uuid.UUID) (*School, error)
	GetListSchool(ctx context.Context, userID uuid.UUID, httpQuery request.GetListSchoolQuery) ([]School, int, error)
	GetSchoolRoleByUserIDAndSchoolID(ctx context.Context, userID uuid.UUID, schoolID uuid.UUID) (*UserSchoolRole, error)
	DeleteSchool(ctx context.Context, schoolID uuid.UUID) error
	UpdateSchoolProfile(ctx context.Context, schoolID uuid.UUID, school School) (*School, error)
	GetSchoolStatistics(ctx context.Context, schoolID uuid.UUID) (*SchoolStatistics, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetSchoolCounts(ctx context.Context, schoolIDs []uuid.UUID) (map[uuid.UUID]SchoolCounts, error)
	GetListSchoolStatistics(ctx context.Context, userID uuid.UUID) (*ListSchoolStatistics, error)
}

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}
