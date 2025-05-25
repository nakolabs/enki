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
	CreateVerifyEmailToken(ctx context.Context, u *UserVerifyEmailToken) error
	VerifyEmailToken(ctx context.Context, email string) (*UserVerifyEmailToken, error)
	Redis() *redis.Client
	GetFirstUserSchoolRolByUserID(ctx context.Context, userID uuid.UUID) (*UserSchoolRole, error)
	CreateForgotPasswordToken(ctx context.Context, u *UserForgotPasswordToken) error
	VerifyForgotPasswordToken(ctx context.Context, email string) (*UserForgotPasswordToken, error)
	UpdatePassword(ctx context.Context, email, password string) error
}
type repository struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func New(db *sqlx.DB, rdb *redis.Client) Repository {
	return &repository{db: db, rdb: rdb}
}
