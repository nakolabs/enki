package repository

import (
	"context"
	"database/sql"
	commonHttp "enuma-elish/pkg/http"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateTeachers(ctx context.Context, teachers []User, schoolID uuid.UUID) error
	UpdateTeacher(ctx context.Context, teacher User) error
	GetListTeachers(ctx context.Context, schoolID uuid.UUID, httpQuery commonHttp.Query) ([]User, int, error)
	CreateTeacherVerifyToken(ctx context.Context, email string) (string, error)
	VerifyEmailToken(ctx context.Context, email string) (string, error)
	GetTeacherByEmail(ctx context.Context, email string) (*User, error)
	Redis() *redis.Client
	Tx(ctx context.Context, options *sql.TxOptions) (*sqlx.Tx, error)
}

type repository struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func New(db *sqlx.DB, rdb *redis.Client) Repository {
	return &repository{
		db:  db,
		rdb: rdb,
	}
}
