package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"time"
)

type Repository interface {
	CreateStudent(ctx context.Context, schoolID uuid.UUID, u []User) error
	CreateStudentVerifyEmailToken(ctx context.Context, email string) (string, error)
	VerifyEmailToken(ctx context.Context, email string) (string, error)
	GetStudentByEmail(ctx context.Context, email string) (*User, error)
	Redis() *redis.Client
	UpdateTeacher(ctx context.Context, teacher User) error
}

type repository struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func New(db *sqlx.DB, rdb *redis.Client) Repository {
	return &repository{db: db, rdb: rdb}
}

type User struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Password   string    `db:"password"`
	Email      string    `db:"email"`
	IsVerified bool      `db:"is_verified"`
	CreatedAt  int64     `db:"created_at"`
	UpdatedAt  int64     `db:"updated_at"`
}

type UserSchoolRole struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	RoleID    uuid.UUID `db:"role_id"`
	SchoolID  uuid.UUID `db:"school_id"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

const StudentVerifyEmailTokenKey = "student:verify:email"

func (r *repository) CreateStudent(ctx context.Context, schoolID uuid.UUID, u []User) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func(tx *sqlx.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Error().Err(err).Msg("error rolling back transaction")
		}
	}(tx)

	insertUserQuery := "INSERT INTO users (id, name, password, email, is_verified, created_at, updated_at) VALUES (:id, :name, :password, :email, :is_verified, :created_at, :updated_at) ON CONFLICT (email) DO NOTHING "
	insertUserSchoolRoleQuery := "INSERT INTO user_school_role (id, user_id, role_id, school_id, created_at, updated_at) VALUES (:id, :user_id, :role_id, :school_id, :created_at, :updated_at) ON CONFLICT (user_id, school_id, role_id) DO NOTHING"

	_, err = tx.NamedExecContext(ctx, insertUserQuery, u)
	if err != nil {
		return err
	}

	emails := make([]interface{}, len(u))
	for i, v := range u {
		emails[i] = v.Email
	}

	selectUserQuery, args, err := sqlx.In("SELECT * FROM users WHERE email IN (?)", emails)
	if err != nil {
		return err
	}

	var users []User
	err = tx.SelectContext(ctx, &users, tx.Rebind(selectUserQuery), args...)
	if err != nil {
		return err
	}

	var studentRoleID uuid.UUID
	err = tx.GetContext(ctx, &studentRoleID, "SELECT id FROM role WHERE name = 'student'")

	var userSchoolRole []UserSchoolRole
	now := time.Now().UnixMilli()
	for _, user := range users {
		userSchoolRole = append(userSchoolRole, UserSchoolRole{
			ID:        user.ID,
			UserID:    user.ID,
			SchoolID:  schoolID,
			RoleID:    studentRoleID,
			CreatedAt: now,
			UpdatedAt: 0,
		})
	}
	_, err = tx.NamedExecContext(ctx, insertUserSchoolRoleQuery, userSchoolRole)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CreateStudentVerifyEmailToken(ctx context.Context, email string) (string, error) {
	token := uuid.New().String()
	err := r.rdb.Set(ctx, StudentVerifyEmailTokenKey+":"+email, token, time.Hour*24).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *repository) VerifyEmailToken(ctx context.Context, email string) (string, error) {
	key := StudentVerifyEmailTokenKey + ":" + email
	token := ""
	err := r.rdb.Get(ctx, key).Scan(&token)
	if err != nil {
		log.Err(err).Str("user", email).Msg("failed get verify email token")
		return token, err
	}

	return token, nil
}

func (r *repository) GetStudentByEmail(ctx context.Context, email string) (*User, error) {
	teacher := &User{}
	err := r.db.GetContext(ctx, teacher, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	return teacher, nil
}

func (r *repository) Redis() *redis.Client {
	return r.rdb
}

func (r *repository) UpdateTeacher(ctx context.Context, teacher User) error {
	updateTeacher := "UPDATE users SET name = :name, email = :email, password = :password, is_verified = :is_verified, updated_at = :updated_at WHERE id = :id"
	_, err := r.db.NamedExecContext(ctx, updateTeacher, teacher)
	if err != nil {
		return err
	}
	return nil
}
