package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateSchool(ctx context.Context, school School) error
}

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}
