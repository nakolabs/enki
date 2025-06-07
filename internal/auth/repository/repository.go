package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	SuperAdmin  = "super_admin"
	Admin       = "admin"
	Teacher     = "teacher"
	Student     = "student"
	HeadTeacher = "head_teacher"
)

var Role = map[string]string{
	SuperAdmin:  SuperAdmin,
	Admin:       Admin,
	Teacher:     Teacher,
	Student:     Student,
	HeadTeacher: HeadTeacher,
}

type Repository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, u *User) error
	CreateVerifyEmailToken(ctx context.Context, u *UserVerifyEmailToken) error
	VerifyEmailToken(ctx context.Context, email string) (*UserVerifyEmailToken, error)
	Redis() *redis.Client
	GetFirstUserSchoolRoleByUserID(ctx context.Context, userID uuid.UUID) (*UserSchoolRole, error)
	CreateForgotPasswordToken(ctx context.Context, u *UserForgotPasswordToken) error
	VerifyForgotPasswordToken(ctx context.Context, email string) (*UserForgotPasswordToken, error)
	UpdatePassword(ctx context.Context, email, password string) error
	UpdateUser(ctx context.Context, profile *User) (*User, error)
}

type repository struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func New(db *sqlx.DB, rdb *redis.Client) Repository {
	return &repository{db: db, rdb: rdb}
}

func (r *repository) UpdateUser(ctx context.Context, user *User) (*User, error) {
	query := `UPDATE profiles SET 
			  name = :name,
			  phone = :phone, 
			  date_of_birth = :date_of_birth, 
			  gender = :gender, 
			  address = :address, 
			  city = :city, 
			  country = :country, 
			  avatar = :avatar, 
			  bio = :bio, 
			  parent_name = :parent_name, 
			  parent_phone = :parent_phone, 
			  parent_email = :parent_email, 
			  updated_at = :updated_at 
			  WHERE id = :id`

	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return nil, err
	}

	return r.GetUserByID(ctx, user.ID)
}
