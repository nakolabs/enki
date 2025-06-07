package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type User struct {
	ID          uuid.UUID      `db:"id"`
	Email       string         `db:"email"`
	Name        string         `db:"name"`
	Password    string         `db:"password"`
	IsVerified  bool           `db:"is_verified"`
	Role        string         `db:"role"`
	Phone       string         `db:"phone"`
	DateOfBirth string         `db:"date_of_birth"`
	Gender      string         `db:"gender"`
	Address     string         `db:"address"`
	City        string         `db:"city"`
	Country     string         `db:"country"`
	Avatar      string         `db:"avatar"`
	Bio         string         `db:"bio"`
	ParentName  string         `db:"parent_name"`
	ParentPhone string         `db:"parent_phone"`
	ParentEmail string         `db:"parent_email"`
	IsDeleted   bool           `db:"is_deleted"`
	CreatedAt   int64          `db:"created_at"`
	UpdatedAt   int64          `db:"updated_at"`
	DeletedAt   int64          `db:"deleted_at"`
	DeletedBy   sql.NullString `db:"deleted_by"`
}

type UserVerifyEmailToken struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserForgotPasswordToken struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

const (
	VerifyEmailTokenKey    = "verify_email_token"
	ForgotPasswordTokenKey = "forgot_password"
)

func (r *repository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	u := &User{}
	err := r.db.GetContext(ctx, u, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		log.Err(err).Str("user", id.String()).Msg("failed to get user by id")
		return nil, err
	}

	return u, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	err := r.db.GetContext(ctx, u, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		log.Err(err).Str("user", email).Msg("failed to get user by email")
		return nil, err
	}
	return u, nil
}

func (r *repository) CreateUser(ctx context.Context, u *User) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO users (email, name, password, is_verified, role) VALUES (:email, :name, :password, :is_verified, :role)`, u)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) CreateVerifyEmailToken(ctx context.Context, u *UserVerifyEmailToken) error {
	key := VerifyEmailTokenKey + ":" + u.Email
	value, err := json.Marshal(u)
	if err != nil {
		log.Err(err).Str("user", u.Email).Msg("failed to marshal verify email token")
		return err
	}

	err = r.rdb.Set(ctx, key, value, time.Minute*30).Err()
	if err != nil {
		log.Err(err).Str("user", u.Email).Msg("failed to set verify email token")
		return err
	}

	return nil
}

func (r *repository) VerifyEmailToken(ctx context.Context, email string) (*UserVerifyEmailToken, error) {
	u := &UserVerifyEmailToken{}
	key := VerifyEmailTokenKey + ":" + email
	res, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		log.Err(err).Str("user", u.Email).Msg("failed get verify email token")
		return nil, err
	}

	err = json.Unmarshal([]byte(res), u)
	if err != nil {
		log.Err(err).Str("user", u.Email).Msg("failed to unmarshal verify email token")
		return nil, err
	}

	return u, nil
}

func (r *repository) CreateForgotPasswordToken(ctx context.Context, u *UserForgotPasswordToken) error {
	key := ForgotPasswordTokenKey + ":" + u.Email
	value, err := json.Marshal(u)
	if err != nil {
		log.Err(err).Str("user", u.Email).Msg("failed to marshal forgot password token")
		return err
	}

	err = r.rdb.Set(ctx, key, value, time.Minute*30).Err()
	if err != nil {
		log.Err(err).Str("user", u.Email).Msg("failed to set forgot password token")
		return err
	}

	return nil
}

func (r *repository) VerifyForgotPasswordToken(ctx context.Context, email string) (*UserForgotPasswordToken, error) {
	u := &UserForgotPasswordToken{}
	key := ForgotPasswordTokenKey + ":" + email
	res, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		log.Err(err).Str("user", u.Email).Msg("failed get verify forgot password token")
		return nil, err
	}

	err = json.Unmarshal([]byte(res), u)
	if err != nil {
		log.Err(err).Str("user", u.Email).Msg("failed to unmarshal verify forgot password token")
		return nil, err
	}

	return u, nil
}

func (r *repository) UpdatePassword(ctx context.Context, email, password string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET password = $1 WHERE email = $2", password, email)
	if err != nil {
		log.Err(err).Str("user", email).Msg("failed to update password")
		return err
	}
	return nil
}

func (r *repository) Redis() *redis.Client {
	return r.rdb
}
