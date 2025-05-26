package repository

import (
	"context"
	"enuma-elish/internal/class/service/data/request"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Class struct {
	ID        uuid.UUID `db:"id"`
	SchoolID  uuid.UUID `db:"school_id"`
	Name      string    `db:"name"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

type ClassStudent struct {
	ID        uuid.UUID `db:"id"`
	StudentID uuid.UUID `db:"student_id"`
	ClassID   uuid.UUID `db:"class_id"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

type ClassTeacher struct {
	ID        uuid.UUID `db:"id"`
	TeacherID uuid.UUID `db:"teacher_id"`
	ClassID   uuid.UUID `db:"class_id"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

type ClassSubject struct {
	ID        uuid.UUID `db:"id"`
	ClassID   uuid.UUID `db:"class_id"`
	SubjectID uuid.UUID `db:"subject_id"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

type Student struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Email      string    `db:"email"`
	IsVerified bool      `db:"is_verified"`
	CreatedAt  int64     `db:"created_at"`
	UpdatedAt  int64     `db:"updated_at"`
}

type Teacher struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Email      string    `db:"email"`
	IsVerified bool      `db:"is_verified"`
	CreatedAt  int64     `db:"created_at"`
	UpdatedAt  int64     `db:"updated_at"`
}

type Subject struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	SchoolID  uuid.UUID `db:"school_id"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

type Repository interface {
	CreateClass(ctx context.Context, class Class) error
	GetClassByID(ctx context.Context, classID uuid.UUID) (*Class, error)
	GetListClasses(ctx context.Context, httpQuery request.GetListClassQuery) ([]Class, int, error)
	UpdateClass(ctx context.Context, class Class) error
	DeleteClass(ctx context.Context, classID uuid.UUID) error
	AddStudentsToClass(ctx context.Context, classID uuid.UUID, studentIDs []uuid.UUID) error
	AddTeachersToClass(ctx context.Context, classID uuid.UUID, teacherIDs []uuid.UUID) error
	AddSubjectsToClass(ctx context.Context, classID uuid.UUID, subjectIDs []uuid.UUID) error
	GetStudentsByClassID(ctx context.Context, classID uuid.UUID, query request.GetStudentsByClassQuery) ([]Student, int, error)
	GetTeachersByClassID(ctx context.Context, classID uuid.UUID, query request.GetTeachersByClassQuery) ([]Teacher, int, error)
	GetSubjectsByClassID(ctx context.Context, classID uuid.UUID, query request.GetSubjectsByClassQuery) ([]Subject, int, error)
	RemoveTeachersFromClass(ctx context.Context, classID uuid.UUID, teacherIDs []uuid.UUID) error
	RemoveStudentsFromClass(ctx context.Context, classID uuid.UUID, studentIDs []uuid.UUID) error
	RemoveSubjectsFromClass(ctx context.Context, classID uuid.UUID, subjectIDs []uuid.UUID) error
}

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateClass(ctx context.Context, class Class) error {
	query := `INSERT INTO class (id, school_id, name, created_at, updated_at) 
			  VALUES (:id, :school_id, :name, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, class)
	return err
}

func (r *repository) GetClassByID(ctx context.Context, classID uuid.UUID) (*Class, error) {
	class := &Class{}
	err := r.db.GetContext(ctx, class, "SELECT * FROM class WHERE id = $1", classID)
	if err != nil {
		return nil, err
	}
	return class, nil
}

func (r *repository) GetListClasses(ctx context.Context, httpQuery request.GetListClassQuery) ([]Class, int, error) {
	var classes []Class
	selectQuery := "SELECT * FROM class WHERE school_id = $1"
	countQuery := "SELECT COUNT(*) FROM class WHERE school_id = $1"

	filterParams := []interface{}{httpQuery.SchoolID}
	filterQuery := ""

	if httpQuery.Search != "" {
		filterQuery += " AND name ILIKE $2"
		filterParams = append(filterParams, "%"+httpQuery.Search+"%")
	}

	if httpQuery.StartDate > 0 && httpQuery.EndDate > 0 {
		if httpQuery.Search != "" {
			filterQuery += " AND created_at BETWEEN $3 AND $4"
			filterParams = append(filterParams, httpQuery.StartDate, httpQuery.EndDate)
		} else {
			filterQuery += " AND created_at BETWEEN $2 AND $3"
			filterParams = append(filterParams, httpQuery.StartDate, httpQuery.EndDate)
		}
	}

	orderQuery := fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d",
		httpQuery.OrderBy, httpQuery.Order,
		len(filterParams)+1, len(filterParams)+2)

	limitParams := []interface{}{httpQuery.PageSize, httpQuery.GetOffset()}
	selectParams := append(filterParams, limitParams...)

	err := r.db.SelectContext(ctx, &classes, selectQuery+filterQuery+orderQuery, selectParams...)
	if err != nil {
		return nil, 0, err
	}

	var total int
	err = r.db.GetContext(ctx, &total, countQuery+filterQuery, filterParams...)
	if err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}

func (r *repository) UpdateClass(ctx context.Context, class Class) error {
	query := `UPDATE class SET name = :name, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, class)
	return err
}

func (r *repository) DeleteClass(ctx context.Context, classID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM class WHERE id = $1", classID)
	return err
}

func (r *repository) AddStudentsToClass(ctx context.Context, classID uuid.UUID, studentIDs []uuid.UUID) error {
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

	now := time.Now().UnixMilli()
	var classStudents []ClassStudent

	for _, studentID := range studentIDs {
		classStudents = append(classStudents, ClassStudent{
			ID:        uuid.New(),
			StudentID: studentID,
			ClassID:   classID,
			CreatedAt: now,
			UpdatedAt: 0,
		})
	}

	insertQuery := `INSERT INTO class_student (id, student_id, class_id, created_at, updated_at) 
					VALUES (:id, :student_id, :class_id, :created_at, :updated_at) 
					ON CONFLICT (student_id, class_id) DO NOTHING`

	_, err = tx.NamedExecContext(ctx, insertQuery, classStudents)
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

func (r *repository) AddTeachersToClass(ctx context.Context, classID uuid.UUID, teacherIDs []uuid.UUID) error {
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

	now := time.Now().UnixMilli()
	var classTeachers []ClassTeacher

	for _, teacherID := range teacherIDs {
		classTeachers = append(classTeachers, ClassTeacher{
			ID:        uuid.New(),
			TeacherID: teacherID,
			ClassID:   classID,
			CreatedAt: now,
			UpdatedAt: 0,
		})
	}

	insertQuery := `INSERT INTO class_teacher (id, teacher_id, class_id, created_at, updated_at) 
					VALUES (:id, :teacher_id, :class_id, :created_at, :updated_at) 
					ON CONFLICT (teacher_id, class_id) DO NOTHING`

	_, err = tx.NamedExecContext(ctx, insertQuery, classTeachers)
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

func (r *repository) AddSubjectsToClass(ctx context.Context, classID uuid.UUID, subjectIDs []uuid.UUID) error {
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

	now := time.Now().UnixMilli()
	var classSubjects []ClassSubject

	for _, subjectID := range subjectIDs {
		classSubjects = append(classSubjects, ClassSubject{
			ID:        uuid.New(),
			ClassID:   classID,
			SubjectID: subjectID,
			CreatedAt: now,
			UpdatedAt: 0,
		})
	}

	insertQuery := `INSERT INTO class_subject (id, class_id, subject_id, created_at, updated_at) 
					VALUES (:id, :class_id, :subject_id, :created_at, :updated_at) 
					ON CONFLICT (class_id, subject_id) DO NOTHING`

	_, err = tx.NamedExecContext(ctx, insertQuery, classSubjects)
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

func (r *repository) GetStudentsByClassID(ctx context.Context, classID uuid.UUID, query request.GetStudentsByClassQuery) ([]Student, int, error) {
	baseQuery := `
		SELECT u.id, u.name, u.email, u.is_verified, u.created_at, u.updated_at
		FROM users u
		INNER JOIN class_student cs ON u.id = cs.student_id
		WHERE cs.class_id = $1
	`

	countQuery := `
		SELECT COUNT(*)
		FROM users u
		INNER JOIN class_student cs ON u.id = cs.student_id
		WHERE cs.class_id = $1
	`

	var students []Student
	limitOrderQuery := fmt.Sprintf(" ORDER BY %s %s LIMIT $2 OFFSET $3", query.OrderBy, query.Order)

	err := r.db.SelectContext(ctx, &students, baseQuery+limitOrderQuery, classID, query.PageSize, query.GetOffset())
	if err != nil {
		return nil, 0, err
	}

	var total int
	err = r.db.GetContext(ctx, &total, countQuery, classID)
	if err != nil {
		return nil, 0, err
	}

	return students, total, nil
}

func (r *repository) GetTeachersByClassID(ctx context.Context, classID uuid.UUID, query request.GetTeachersByClassQuery) ([]Teacher, int, error) {
	baseQuery := `
		SELECT u.id, u.name, u.email, u.is_verified, u.created_at, u.updated_at
		FROM users u
		INNER JOIN class_teacher ct ON u.id = ct.teacher_id
		WHERE ct.class_id = $1
	`

	countQuery := `
		SELECT COUNT(*)
		FROM users u
		INNER JOIN class_teacher ct ON u.id = ct.teacher_id
		WHERE ct.class_id = $1
	`

	var teachers []Teacher
	limitOrderQuery := fmt.Sprintf(" ORDER BY %s %s LIMIT $2 OFFSET $3", query.OrderBy, query.Order)

	err := r.db.SelectContext(ctx, &teachers, baseQuery+limitOrderQuery, classID, query.PageSize, query.GetOffset())
	if err != nil {
		return nil, 0, err
	}

	var total int
	err = r.db.GetContext(ctx, &total, countQuery, classID)
	if err != nil {
		return nil, 0, err
	}

	return teachers, total, nil
}

func (r *repository) GetSubjectsByClassID(ctx context.Context, classID uuid.UUID, query request.GetSubjectsByClassQuery) ([]Subject, int, error) {
	baseQuery := `
		SELECT s.id, s.name, s.school_id, s.created_at, s.updated_at
		FROM subject s
		INNER JOIN class_subject cs ON s.id = cs.subject_id
		WHERE cs.class_id = $1
	`

	countQuery := `
		SELECT COUNT(*)
		FROM subject s
		INNER JOIN class_subject cs ON s.id = cs.subject_id
		WHERE cs.class_id = $1
	`

	var subjects []Subject
	limitOrderQuery := fmt.Sprintf(" ORDER BY %s %s LIMIT $2 OFFSET $3", query.OrderBy, query.Order)

	err := r.db.SelectContext(ctx, &subjects, baseQuery+limitOrderQuery, classID, query.PageSize, query.GetOffset())
	if err != nil {
		return nil, 0, err
	}

	var total int
	err = r.db.GetContext(ctx, &total, countQuery, classID)
	if err != nil {
		return nil, 0, err
	}

	return subjects, total, nil
}

func (r *repository) RemoveTeachersFromClass(ctx context.Context, classID uuid.UUID, teacherIDs []uuid.UUID) error {
	query := `DELETE FROM class_teacher WHERE class_id = $1 AND teacher_id = ANY($2)`
	_, err := r.db.ExecContext(ctx, query, classID, pq.Array(teacherIDs))
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) RemoveStudentsFromClass(ctx context.Context, classID uuid.UUID, studentIDs []uuid.UUID) error {
	query := `DELETE FROM class_student WHERE class_id = $1 AND student_id = ANY($2)`
	_, err := r.db.ExecContext(ctx, query, classID, pq.Array(studentIDs))
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) RemoveSubjectsFromClass(ctx context.Context, classID uuid.UUID, subjectIDs []uuid.UUID) error {
	query := `DELETE FROM class_subject WHERE class_id = $1 AND subject_id = ANY($2)`
	_, err := r.db.ExecContext(ctx, query, classID, pq.Array(subjectIDs))
	if err != nil {
		return err
	}
	return nil
}
