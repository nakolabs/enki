package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, u *User) error
	CreateVerifyEmailToken(ctx context.Context, u *UserVerifyEmailToken) (string, error)
	VerifyEmailToken(ctx context.Context, email string) (*UserVerifyEmailToken, error)
	Redis() *redis.Client
}
type repository struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func New(db *sqlx.DB, rdb *redis.Client) Repository {
	return &repository{db: db, rdb: rdb}
}
