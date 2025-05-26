package repository

import (
	"context"
	"enuma-elish/internal/student/service/data/request"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	CreateStudent(ctx context.Context, schoolID uuid.UUID, u []User) error
	CreateStudentVerifyEmailToken(ctx context.Context, email string) (string, error)
	VerifyEmailToken(ctx context.Context, email string) (string, error)
	GetStudentByEmail(ctx context.Context, email string) (*User, error)
	GetStudentByID(ctx context.Context, studentID uuid.UUID) (*User, error)
	Redis() *redis.Client
	UpdateStudent(ctx context.Context, student User) error
	DeleteStudent(ctx context.Context, studentID uuid.UUID, schoolID uuid.UUID) error
	GetListStudent(ctx context.Context, httpQuery request.GetListStudentQuery) ([]User, int, error)
	UpdateStudentClass(ctx context.Context, studentID, oldClassID, newClassID uuid.UUID) error
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

	// Rollback hanya jika terjadi error sebelum commit
	committed := false
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				log.Error().Err(err).Msg("error rolling back transaction")
			}
		}
	}()

	insertUserQuery := "INSERT INTO users (id, name, password, email, is_verified, created_at, updated_at) VALUES (:id, :name, :password, :email, :is_verified, :created_at, :updated_at) ON CONFLICT (email) DO NOTHING"
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
	if err != nil {
		return err
	}

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
	committed = true

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
	student := &User{}
	err := r.db.GetContext(ctx, student, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	return student, nil
}

func (r *repository) GetStudentByID(ctx context.Context, studentID uuid.UUID) (*User, error) {
	student := &User{}
	err := r.db.GetContext(ctx, student, "SELECT * FROM users WHERE id = $1", studentID)
	if err != nil {
		return nil, err
	}
	return student, nil
}

func (r *repository) Redis() *redis.Client {
	return r.rdb
}

func (r *repository) UpdateStudent(ctx context.Context, student User) error {
	updateStudent := "UPDATE users SET name = :name, email = :email, password = :password, is_verified = :is_verified, updated_at = :updated_at WHERE id = :id"
	_, err := r.db.NamedExecContext(ctx, updateStudent, student)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteStudent(ctx context.Context, studentID uuid.UUID, schoolID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	committed := false
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				log.Error().Err(err).Msg("error rolling back transaction")
			}
		}
	}()

	// Delete from user_school_role first (foreign key constraint)
	_, err = tx.ExecContext(ctx, "DELETE FROM user_school_role WHERE user_id = $1 AND school_id = $2", studentID, schoolID)
	if err != nil {
		return err
	}

	// Check if user has other school associations
	var count int
	err = tx.GetContext(ctx, &count, "SELECT COUNT(*) FROM user_school_role WHERE user_id = $1", studentID)
	if err != nil {
		return err
	}

	// If no other associations, delete the user
	if count == 0 {
		_, err = tx.ExecContext(ctx, "DELETE FROM users WHERE id = $1", studentID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	committed = true

	return nil
}

func (r *repository) GetListStudent(ctx context.Context, httpQuery request.GetListStudentQuery) ([]User, int, error) {
	var students []User
	selectStudent := "SELECT users.id, name, email, is_verified, user_school_role.created_at, user_school_role.updated_at FROM users JOIN user_school_role on users.id = user_school_role.user_id WHERE true"
	countQuery := "SELECT COUNT(*) FROM users JOIN user_school_role on users.id = user_school_role.user_id WHERE true"

	var studentRoleID uuid.UUID
	err := r.db.GetContext(ctx, &studentRoleID, "SELECT id FROM role WHERE name = 'student'")
	if err != nil {
		return nil, 0, err
	}

	filterParams := []interface{}{httpQuery.SchoolID, studentRoleID}
	filterQuery := " AND user_school_role.school_id = ? AND user_school_role.role_id = ?"

	if httpQuery.Search != "" && len(httpQuery.SearchBy) > 0 {
		filterQuery += " AND ("
		for i, v := range httpQuery.SearchBy {
			if i > 0 {
				filterQuery += " OR "
			}
			filterQuery += fmt.Sprintf("%s LIKE ?", v)
			filterParams = append(filterParams, "%"+httpQuery.Search+"%")
		}
		filterQuery += ")"
	}

	if httpQuery.StartDate > 0 && httpQuery.EndDate > 0 {
		filterParams = append(filterParams, httpQuery.DateRange.StartDate, httpQuery.DateRange.EndDate)
		filterQuery += " AND created_at BETWEEN ? AND ?"
	} else if httpQuery.StartDate > 0 {
		filterParams = append(filterParams, httpQuery.DateRange.StartDate)
		filterQuery += " AND created_at >= ?"
	} else if httpQuery.EndDate > 0 {
		filterParams = append(filterParams, httpQuery.DateRange.EndDate)
		filterQuery += " AND created_at <= ?"
	}

	limitOrderQuery := fmt.Sprintf(" ORDER BY %s %s LIMIT ? OFFSET ? ", httpQuery.OrderBy, httpQuery.Order)
	limitOrderParams := []interface{}{httpQuery.PageSize, httpQuery.GetOffset()}

	selectParams := append(filterParams, limitOrderParams...)
	err = r.db.SelectContext(ctx, &students, r.db.Rebind(selectStudent+filterQuery+limitOrderQuery), selectParams...)
	if err != nil {
		return nil, 0, err
	}

	total := 0
	err = r.db.GetContext(ctx, &total, r.db.Rebind(countQuery+filterQuery), filterParams...)
	if err != nil {
		return nil, 0, err
	}

	return students, total, nil
}
