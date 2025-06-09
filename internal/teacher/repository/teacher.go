package repository

import (
	"context"
	"database/sql"
	"enuma-elish/internal/teacher/service/data/request"
	"enuma-elish/pkg/jwt"
	"fmt"
	"time"

	commonError "enuma-elish/pkg/error"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

const TeacherVerifyEmailTokenKey = "teacher:verify:email"

type User struct {
	ID         uuid.UUID      `db:"id"`
	Name       string         `db:"name"`
	Email      string         `db:"email"`
	Password   string         `db:"password"`
	IsVerified bool           `db:"is_verified"`
	CreatedAt  int64          `db:"created_at"`
	CreatedBy  uuid.UUID      `db:"created_by"`
	UpdatedAt  int64          `db:"updated_at"`
	UpdatedBy  sql.NullString `db:"updated_by"`
	DeletedAt  int64          `db:"deleted_at"`
	DeletedBy  sql.NullString `db:"deleted_by"`
}

type UserSchoolRole struct {
	ID        uuid.UUID      `db:"id"`
	UserID    uuid.UUID      `db:"user_id"`
	SchoolID  uuid.UUID      `db:"school_id"`
	RoleID    string         `db:"role_id"`
	IsDeleted bool           `db:"is_deleted"`
	CreatedAt int64          `db:"created_at"`
	CreatedBy uuid.UUID      `db:"created_by"`
	UpdatedAt int64          `db:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by"`
	DeletedAt int64          `db:"deleted_at"`
	DeletedBy sql.NullString `db:"deleted_by"`
}

type Subject struct {
	ID        uuid.UUID      `db:"id"`
	Name      string         `db:"name"`
	SchoolID  uuid.UUID      `db:"school_id"`
	CreatedAt int64          `db:"created_at"`
	CreatedBy uuid.UUID      `db:"created_by"`
	UpdatedAt int64          `db:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by"`
	DeletedAt int64          `db:"deleted_at"`
	DeletedBy sql.NullString `db:"deleted_by"`
}

type Class struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	SchoolID  uuid.UUID `db:"school_id"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

type TeacherSubject struct {
	ID        uuid.UUID `db:"id"`
	TeacherID uuid.UUID `db:"teacher_id"`
	SubjectID uuid.UUID `db:"subject_id"`
	CreatedAt int64     `db:"created_at"`
	CreatedBy uuid.UUID `db:"created_by"`
	UpdatedAt int64     `db:"updated_at"`
	IsDeleted bool      `db:"is_deleted"`
}

type TeacherClass struct {
	ID        uuid.UUID `db:"id"`
	TeacherID uuid.UUID `db:"teacher_id"`
	ClassID   uuid.UUID `db:"class_id"`
	CreatedAt int64     `db:"created_at"`
	CreatedBy uuid.UUID `db:"created_by"`
	UpdatedAt int64     `db:"updated_at"`
	IsDeleted bool      `db:"is_deleted"`
}

func (r *repository) CreateTeachers(ctx context.Context, teachers []User, schoolID uuid.UUID) error {

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	insertTeacher := `INSERT INTO users (id, email, name, password, is_verified, created_at, created_by, updated_at)
	 VALUES (:id, :email, :name, :password, :is_verified, :created_at, :created_by, :updated_at) ON CONFLICT (email, is_deleted) DO NOTHING`

	_, err = tx.NamedExecContext(ctx, insertTeacher, teachers)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert teachers")
		return err
	}

	var users []User
	var emails []string
	for _, teacher := range teachers {
		emails = append(emails, teacher.Email)
	}

	q, args, err := sqlx.In("SELECT id, name, email, password, created_at, created_by, updated_at, deleted_at, deleted_by FROM users WHERE email = (?)", emails)
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
			RoleID:    "teacher",
			SchoolID:  schoolID,
			CreatedAt: now,
			CreatedBy: v.CreatedBy,
			UpdatedAt: 0,
		})
	}

	insertUserSchoolRole := `INSERT INTO user_school_role (id, user_id, school_id, role_id, created_at, created_by, updated_at)
	 VALUES (:id, :user_id, :school_id, :role_id, :created_at, :created_by, :updated_at) ON CONFLICT (user_id, school_id, role_id, is_deleted) DO NOTHING`
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

func (r *repository) CreateTeachersWithAssignments(ctx context.Context, teachers []User, schoolID uuid.UUID, teacherSubjects []TeacherSubject, teacherClasses []TeacherClass) error {
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

	// Insert teachers
	insertTeacher := `INSERT INTO users (id, email, name, password, is_verified, created_at, created_by, updated_at)
	 VALUES (:id, :email, :name, :password, :is_verified, :created_at, :created_by, :updated_at)`

	_, err = tx.NamedExecContext(ctx, insertTeacher, teachers)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert teachers")
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return commonError.ErrrTeacherAlreadyExists
		}
		return err
	}

	// Get created teachers
	var users []User
	var emails []string
	for _, teacher := range teachers {
		emails = append(emails, teacher.Email)
	}

	q, args, err := sqlx.In("SELECT id, name, email, password, created_at, created_by, updated_at, deleted_at, deleted_by FROM users WHERE email IN (?)", emails)
	if err != nil {
		return err
	}

	err = tx.SelectContext(ctx, &users, tx.Rebind(q), args...)
	if err != nil {
		return err
	}

	// Insert user school roles
	var userSchoolRole []UserSchoolRole
	now := time.Now().UnixMilli()
	for _, v := range users {
		userSchoolRole = append(userSchoolRole, UserSchoolRole{
			ID:        uuid.New(),
			UserID:    v.ID,
			RoleID:    "teacher",
			SchoolID:  schoolID,
			CreatedAt: now,
			CreatedBy: v.CreatedBy,
			UpdatedAt: 0,
		})
	}

	insertUserSchoolRole := `INSERT INTO user_school_role (id, user_id, school_id, role_id, created_at, created_by, updated_at)
	 VALUES (:id, :user_id, :school_id, :role_id, :created_at, :created_by, :updated_at) ON CONFLICT (user_id, school_id, role_id, is_deleted) DO NOTHING`
	_, err = tx.NamedExecContext(ctx, insertUserSchoolRole, userSchoolRole)
	if err != nil {
		return err
	}

	// Insert subject assignments
	if len(teacherSubjects) > 0 {
		insertSubjectAssignment := `INSERT INTO teacher_subject (id, teacher_id, subject_id, created_at, created_by, updated_at, is_deleted) 
						VALUES (:id, :teacher_id, :subject_id, :created_at, :created_by, :updated_at, :is_deleted) 
						ON CONFLICT (teacher_id, subject_id, is_deleted) DO NOTHING`

		_, err = tx.NamedExecContext(ctx, insertSubjectAssignment, teacherSubjects)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert teacher subject assignments")
			return err
		}
	}

	// Insert class assignments
	if len(teacherClasses) > 0 {
		insertClassAssignment := `INSERT INTO class_teacher (id, teacher_id, class_id, created_at, created_by, updated_at, is_deleted) 
						VALUES (:id, :teacher_id, :class_id, :created_at, :created_by, :updated_at, :is_deleted) 
						ON CONFLICT (teacher_id, class_id, is_deleted) DO NOTHING`

		_, err = tx.NamedExecContext(ctx, insertClassAssignment, teacherClasses)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert teacher class assignments")
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

	filterParams := make([]interface{}, 0)
	filterQuery := ""
	filterParams = append(filterParams, httpQuery.SchoolID, "teacher")
	filterQuery += " AND user_school_role.school_id = ? AND user_school_role.role_id = ? "

	if httpQuery.Search != "" {
		first := false
		for _, v := range httpQuery.SearchBy {
			if !first {
				filterQuery += " AND ("
				first = true
			} else {
				filterQuery += " OR "
			}
			filterQuery += fmt.Sprintf(" %s LIKE ? ", v)
			filterParams = append(filterParams, "%"+httpQuery.Search+"%")
		}
		filterQuery += " ) "
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

	if httpQuery.ClassID != "" {
		filterQuery += " AND EXISTS (SELECT 1 FROM class_teacher ct WHERE ct.teacher_id = users.id AND ct.class_id = ? AND ct.is_deleted = false) "
		filterParams = append(filterParams, httpQuery.ClassID)
	}

	if httpQuery.SubjectID != "" {
		filterQuery += " AND EXISTS (SELECT 1 FROM teacher_subject ts WHERE ts.teacher_id = users.id AND ts.subject_id = ? AND ts.is_deleted = false) "
		filterParams = append(filterParams, httpQuery.SubjectID)
	}

	if httpQuery.IsVerified != "" {
		isVerified := httpQuery.IsVerified == "true"
		filterQuery += " AND users.is_verified = ? "
		filterParams = append(filterParams, isVerified)
	}

	limitOrderQuery := fmt.Sprintf(" ORDER BY %s %s LIMIT ? OFFSET ? ", httpQuery.OrderBy, httpQuery.Order)
	limitOrderParams := []interface{}{httpQuery.PageSize, httpQuery.GetOffset()}

	selectParams := append(filterParams, limitOrderParams...)
	err := r.db.SelectContext(ctx, &teachers, r.db.Rebind(selectTeacher+filterQuery+limitOrderQuery), selectParams...)
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
	err := r.db.GetContext(ctx, teacher, "SELECT id, name, email, password, created_at, created_by, updated_at, deleted_at, deleted_by FROM users WHERE email = $1", email)
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

	_, err = tx.ExecContext(ctx, "DELETE FROM user_school_role WHERE user_id = $1 AND school_id = $2", teacherID, schoolID)
	if err != nil {
		return err
	}

	var count int
	err = tx.GetContext(ctx, &count, "SELECT COUNT(*) FROM user_school_role WHERE user_id = $1", teacherID)
	if err != nil {
		return err
	}

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
	err := r.db.GetContext(ctx, teacher, "SELECT  id, name, email, password, created_at, created_by, updated_at, deleted_at, deleted_by FROM users WHERE id = $1", teacherID)
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

func (r *repository) GetTeacherSubjects(ctx context.Context, teacherID uuid.UUID) ([]Subject, error) {

	payload, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to extract JWT payload from context")
		return nil, err
	}
	query := `
		SELECT s.id, s.name, s.school_id, s.created_at, s.updated_at
		FROM subject s
		INNER JOIN teacher_subject ts ON s.id = ts.subject_id
		WHERE ts.teacher_id = $1 AND s.school_id = $2
	`
	var subjects []Subject
	err = r.db.SelectContext(ctx, &subjects, query, teacherID, payload.User.SchoolID)
	if err != nil {
		return nil, err
	}
	return subjects, nil
}

func (r *repository) GetTeacherClasses(ctx context.Context, teacherID uuid.UUID) ([]Class, error) {
	payload, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to extract JWT payload from context")
		return nil, err
	}
	query := `
		SELECT c.id, c.name, c.school_id, c.created_at, c.updated_at
		FROM class c
		INNER JOIN class_teacher ct ON c.id = ct.class_id
		WHERE ct.teacher_id = $1 AND c.school_id = $2 AND ct.is_deleted = false
	`
	var classes []Class
	err = r.db.SelectContext(ctx, &classes, query, teacherID, payload.User.SchoolID)
	if err != nil {
		return nil, err
	}
	return classes, nil
}

func (r *repository) GetTeacherAssignments(ctx context.Context, teacherIDs []uuid.UUID, schoolID uuid.UUID) (map[uuid.UUID][]Subject, map[uuid.UUID][]Class, error) {
	// Get teacher subjects
	subjectQuery := `
		SELECT ts.teacher_id, s.id, s.name, s.school_id, s.created_at, s.updated_at
		FROM teacher_subject ts
		INNER JOIN subject s ON ts.subject_id = s.id
		WHERE ts.teacher_id = ANY($1) AND s.school_id = $2 AND ts.is_deleted = false
	`

	var teacherSubjects []struct {
		TeacherID uuid.UUID `db:"teacher_id"`
		Subject
	}

	err := r.db.SelectContext(ctx, &teacherSubjects, subjectQuery, pq.Array(teacherIDs), schoolID)
	if err != nil {
		return nil, nil, err
	}

	// Get teacher classes
	classQuery := `
		SELECT ct.teacher_id, c.id, c.name, c.school_id, c.created_at, c.updated_at
		FROM class_teacher ct
		INNER JOIN class c ON ct.class_id = c.id
		WHERE ct.teacher_id = ANY($1) AND c.school_id = $2 AND ct.is_deleted = false
	`

	var teacherClasses []struct {
		TeacherID uuid.UUID `db:"teacher_id"`
		Class
	}

	err = r.db.SelectContext(ctx, &teacherClasses, classQuery, pq.Array(teacherIDs), schoolID)
	if err != nil {
		return nil, nil, err
	}

	// Organize data by teacher ID
	subjectMap := make(map[uuid.UUID][]Subject)
	classMap := make(map[uuid.UUID][]Class)

	for _, ts := range teacherSubjects {
		subjectMap[ts.TeacherID] = append(subjectMap[ts.TeacherID], ts.Subject)
	}

	for _, tc := range teacherClasses {
		classMap[tc.TeacherID] = append(classMap[tc.TeacherID], tc.Class)
	}

	return subjectMap, classMap, nil
}

func (r *repository) GetTeacherStatistics(ctx context.Context, schoolID string) (int, int, int, int, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN users.is_verified = true THEN 1 END) as verified,
			COUNT(CASE WHEN users.is_verified = false THEN 1 END) as pending,
			COUNT(CASE WHEN users.is_verified = true AND EXISTS(
				SELECT 1 FROM teacher_subject ts WHERE ts.teacher_id = users.id AND ts.is_deleted = false
				UNION
				SELECT 1 FROM class_teacher ct WHERE ct.teacher_id = users.id AND ct.is_deleted = false
			) THEN 1 END) as active
		FROM users 
		JOIN user_school_role ON users.id = user_school_role.user_id 
		WHERE user_school_role.school_id = $1 AND user_school_role.role_id = 'teacher' AND user_school_role.is_deleted = false
	`

	var stats struct {
		Total    int `db:"total"`
		Verified int `db:"verified"`
		Pending  int `db:"pending"`
		Active   int `db:"active"`
	}

	err := r.db.GetContext(ctx, &stats, query, schoolID)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return stats.Total, stats.Verified, stats.Pending, stats.Active, nil
}

func (r *repository) AssignSubjectsToTeachers(ctx context.Context, teacherSubjects []TeacherSubject) error {
	if len(teacherSubjects) == 0 {
		return nil
	}

	insertQuery := `INSERT INTO teacher_subject (id, teacher_id, subject_id, created_at, created_by, updated_at, is_deleted) 
					VALUES (:id, :teacher_id, :subject_id, :created_at, :created_by, :updated_at, :is_deleted) 
					ON CONFLICT (teacher_id, subject_id, is_deleted) DO NOTHING`

	_, err := r.db.NamedExecContext(ctx, insertQuery, teacherSubjects)
	return err
}

func (r *repository) AssignClassesToTeachers(ctx context.Context, teacherClasses []TeacherClass) error {
	if len(teacherClasses) == 0 {
		return nil
	}

	insertQuery := `INSERT INTO class_teacher (id, teacher_id, class_id, created_at, created_by, updated_at, is_deleted) 
					VALUES (:id, :teacher_id, :class_id, :created_at, :created_by, :updated_at, :is_deleted) 
					ON CONFLICT (teacher_id, class_id, is_deleted) DO NOTHING`

	_, err := r.db.NamedExecContext(ctx, insertQuery, teacherClasses)
	return err
}
