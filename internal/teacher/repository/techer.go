package repository

import (
	"context"
	"database/sql"
	"enuma-elish/internal/teacher/service/data/request"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const TeacherVerifyEmailTokenKey = "teacher:verify:email"

type User struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	IsVerified bool      `db:"is_verified"`
	CreatedAt  int64     `db:"created_at"`
	UpdatedAt  int64     `db:"updated_at"`
}

type Role struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

type UserSchoolRole struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	SchoolID  uuid.UUID `db:"school_id"`
	RoleID    uuid.UUID `db:"role_id"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

func (r *repository) CreateTeachers(ctx context.Context, teachers []User, schoolID uuid.UUID) error {

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	insertTeacher := "INSERT INTO users (id, email, name, password, is_verified, created_at, updated_at) VALUES (:id, :email, :name, :password, :is_verified, :created_at, :updated_at) ON CONFLICT (email) DO NOTHING"
	_, err = tx.NamedExecContext(ctx, insertTeacher, teachers)
	if err != nil {
		return err
	}

	getRole := "SELECT id, name, created_at, updated_at FROM role WHERE name = 'teacher'"
	role := &Role{}
	err = tx.GetContext(ctx, role, getRole)
	if err != nil {
		return err
	}

	var users []User
	var emails []string
	for _, teacher := range teachers {
		emails = append(emails, teacher.Email)
	}

	q, args, err := sqlx.In("SELECT * FROM users WHERE email = (?)", emails)
	if err != nil {
		return err
	}

	err = tx.SelectContext(ctx, &users, r.db.Rebind(q), args...)
	if err != nil {
		return err
	}

	var userSchoolRole []UserSchoolRole
	now := time.Now().UnixMilli()
	for _, v := range users {
		userSchoolRole = append(userSchoolRole, UserSchoolRole{
			ID:        uuid.New(),
			UserID:    v.ID,
			RoleID:    role.ID,
			SchoolID:  schoolID,
			CreatedAt: now,
			UpdatedAt: 0,
		})
	}

	insertUserSchoolRole := "INSERT INTO user_school_role (id, user_id, school_id, role_id, created_at, updated_at) VALUES (:id, :user_id, :school_id, :role_id, :created_at, :updated_at) ON CONFLICT (user_id, school_id, role_id) DO NOTHING"
	_, err = tx.NamedExecContext(ctx, insertUserSchoolRole, userSchoolRole)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateTeacher(ctx context.Context, teacher User) error {
	updateTeacher := "UPDATE users SET name = :name, email = :email, password = :password, is_verified = :is_verified, updated_at = :updated_at WHERE id = :id"
	_, err := r.db.NamedExecContext(ctx, updateTeacher, teacher)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetListTeachers(ctx context.Context, httpQuery request.GetListTeacherQuery) ([]User, int, error) {
	var teachers []User
	selectTeacher := "SELECT users.id, name, email, is_verified, user_school_role.created_at, user_school_role.updated_at FROM users JOIN user_school_role on users.id = user_school_role.user_id WHERE true"
	countQuery := "SELECT COUNT(*) FROM users JOIN user_school_role on users.id = user_school_role.user_id WHERE true "

	var teacherRoleID uuid.UUID
	err := r.db.GetContext(ctx, &teacherRoleID, "SELECT id FROM role WHERE name = 'teacher'")
	if err != nil {
		return nil, 0, err
	}

	filterParams := make([]interface{}, 0)
	filterQuery := ""
	filterParams = append(filterParams, httpQuery.SchoolID, teacherRoleID)
	filterQuery += " AND user_school_role.school_id = ? AND user_school_role.role_id = ? "

	if httpQuery.Search != "" {
		for _, v := range httpQuery.SearchBy {
			filterQuery += fmt.Sprintf(" OR %s LIKE ? ", v)
			filterParams = append(filterParams, httpQuery.Search)
		}
	}

	if httpQuery.StartDate > 0 && httpQuery.EndDate > 0 {
		filterParams = append(filterParams, httpQuery.DateRange.StartDate, httpQuery.DateRange.EndDate)
		filterQuery += " AND created_at BETWEEN ? AND ? "
	}

	if httpQuery.StartDate > 0 && httpQuery.EndDate <= 0 {
		filterParams = append(filterParams, httpQuery.DateRange.StartDate)
		filterQuery += " AND created_at >= ? "
	}

	if httpQuery.StartDate <= 0 && httpQuery.EndDate > 0 {
		filterParams = append(filterParams, httpQuery.DateRange.EndDate)
		filterQuery += " AND created_at <= ? "
	}

	limitOrderQuery := fmt.Sprintf(" ORDER BY %s %s LIMIT ? OFFSET ? ", httpQuery.OrderBy, httpQuery.Order)
	limitOrderParams := []interface{}{httpQuery.PageSize, httpQuery.GetOffset()}

	selectParams := append(filterParams, limitOrderParams...)
	err = r.db.SelectContext(ctx, &teachers, r.db.Rebind(selectTeacher+filterQuery+limitOrderQuery), selectParams...)
	if err != nil {
		return nil, 9, err
	}

	total := 0
	err = r.db.GetContext(ctx, &total, r.db.Rebind(countQuery+filterQuery), filterParams...)
	if err != nil {
		return nil, 0, err
	}

	return teachers, total, nil
}

func (r *repository) CreateTeacherVerifyToken(ctx context.Context, email string) (string, error) {
	token := uuid.New().String()
	err := r.rdb.Set(ctx, TeacherVerifyEmailTokenKey+":"+email, token, time.Hour*24).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *repository) VerifyEmailToken(ctx context.Context, email string) (string, error) {
	key := TeacherVerifyEmailTokenKey + ":" + email
	token := ""
	err := r.rdb.Get(ctx, key).Scan(&token)
	if err != nil {
		log.Err(err).Str("user", email).Msg("failed get verify email token")
		return token, err
	}

	return token, nil
}

func (r *repository) GetTeacherByEmail(ctx context.Context, email string) (*User, error) {
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

func (r *repository) Tx(ctx context.Context, options *sql.TxOptions) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, options)
}

func (r *repository) DeleteTeacher(ctx context.Context, teacherID uuid.UUID, schoolID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func(tx *sqlx.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Err(err).Msg("error rolling back transaction")
		}
	}(tx)

	// Delete from user_school_role first (foreign key constraint)
	_, err = tx.ExecContext(ctx, "DELETE FROM user_school_role WHERE user_id = $1 AND school_id = $2", teacherID, schoolID)
	if err != nil {
		return err
	}

	// Check if user has other school associations
	var count int
	err = tx.GetContext(ctx, &count, "SELECT COUNT(*) FROM user_school_role WHERE user_id = $1", teacherID)
	if err != nil {
		return err
	}

	// If no other associations, delete the user
	if count == 0 {
		_, err = tx.ExecContext(ctx, "DELETE FROM users WHERE id = $1", teacherID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetTeacherByID(ctx context.Context, teacherID uuid.UUID) (*User, error) {
	teacher := &User{}
	err := r.db.GetContext(ctx, teacher, "SELECT * FROM users WHERE id = $1", teacherID)
	if err != nil {
		return nil, err
	}
	return teacher, nil
}

func (r *repository) UpdateTeacherClass(ctx context.Context, teacherID, oldClassID, newClassID uuid.UUID) error {
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

	// Update the class_teacher record
	updateQuery := `UPDATE class_teacher 
					SET class_id = $1, updated_at = $2 
					WHERE teacher_id = $3 AND class_id = $4`

	now := time.Now().UnixMilli()
	result, err := tx.ExecContext(ctx, updateQuery, newClassID, now, teacherID, oldClassID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no class assignment found for teacher %s in class %s", teacherID, oldClassID)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	committed = true

	return nil
}
