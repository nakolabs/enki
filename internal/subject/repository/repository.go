package repository

import (
	"context"
	"database/sql"
	"enuma-elish/internal/subject/service/data/request"
	"enuma-elish/pkg/jwt"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Subject struct {
	ID        uuid.UUID      `db:"id"`
	SchoolID  uuid.UUID      `db:"school_id"`
	Name      string         `db:"name"`
	CreatedAt int64          `db:"created_at"`
	CreatedBy uuid.UUID      `db:"created_by"`
	UpdatedAt int64          `db:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by"`
	DeletedAt int64          `db:"deleted_at"`
	DeletedBy sql.NullString `db:"deleted_by"`
}

type TeacherSubject struct {
	ID        uuid.UUID      `db:"id"`
	TeacherID uuid.UUID      `db:"teacher_id"`
	SubjectID uuid.UUID      `db:"subject_id"`
	CreatedAt int64          `db:"created_at"`
	CreatedBy uuid.UUID      `db:"created_by"`
	UpdatedAt int64          `db:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by"`
	DeletedAt int64          `db:"deleted_at"`
	DeletedBy sql.NullString `db:"deleted_by"`
}

type Teacher struct {
	ID         uuid.UUID      `db:"id"`
	Name       string         `db:"name"`
	Email      string         `db:"email"`
	IsVerified bool           `db:"is_verified"`
	CreatedAt  int64          `db:"created_at"`
	CreatedBy  uuid.UUID      `db:"created_by"`
	UpdatedAt  int64          `db:"updated_at"`
	UpdatedBy  sql.NullString `db:"updated_by"`
	DeletedAt  int64          `db:"deleted_at"`
	DeletedBy  sql.NullString `db:"deleted_by"`
}

type Repository interface {
	CreateSubject(ctx context.Context, subject Subject) error
	GetSubjectByID(ctx context.Context, subjectID uuid.UUID) (*Subject, error)
	GetListSubjects(ctx context.Context, httpQuery request.GetListSubjectQuery) ([]Subject, int, error)
	UpdateSubject(ctx context.Context, subject Subject) error
	DeleteSubject(ctx context.Context, subjectID uuid.UUID) error
	AssignTeachersToSubject(ctx context.Context, subjectID uuid.UUID, teacherIDs []uuid.UUID) error
	GetTeachersBySubjectID(ctx context.Context, subjectID uuid.UUID, query request.GetTeachersBySubjectQuery) ([]Teacher, int, error)
	UpdateSubjectClass(ctx context.Context, subjectID, oldClassID, newClassID uuid.UUID) error
}

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateSubject(ctx context.Context, subject Subject) error {
	query := `INSERT INTO subject (id, school_id, name, created_at, created_by, updated_at) 
			  VALUES (:id, :school_id, :name, :created_at, :created_by, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, subject)
	return err
}

func (r *repository) GetSubjectByID(ctx context.Context, subjectID uuid.UUID) (*Subject, error) {
	subject := &Subject{}
	err := r.db.GetContext(ctx, subject, "SELECT * FROM subject WHERE id = $1", subjectID)
	if err != nil {
		return nil, err
	}
	return subject, nil
}

func (r *repository) GetListSubjects(ctx context.Context, httpQuery request.GetListSubjectQuery) ([]Subject, int, error) {
	var subjects []Subject
	selectQuery := "SELECT * FROM subject WHERE school_id = $1"
	countQuery := "SELECT COUNT(*) FROM subject WHERE school_id = $1"

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

	err := r.db.SelectContext(ctx, &subjects, selectQuery+filterQuery+orderQuery, selectParams...)
	if err != nil {
		return nil, 0, err
	}

	var total int
	err = r.db.GetContext(ctx, &total, countQuery+filterQuery, filterParams...)
	if err != nil {
		return nil, 0, err
	}

	return subjects, total, nil
}

func (r *repository) UpdateSubject(ctx context.Context, subject Subject) error {
	query := `UPDATE subject SET name = :name, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, subject)
	return err
}

func (r *repository) DeleteSubject(ctx context.Context, subjectID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM subject WHERE id = $1", subjectID)
	return err
}

func (r *repository) AssignTeachersToSubject(ctx context.Context, subjectID uuid.UUID, teacherIDs []uuid.UUID) error {

	jwtClaims, err := jwt.ExtractContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to extract JWT claims: %w", err)
	}
	createdBy := jwtClaims.User.ID

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	committed := false
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				// Note: using fmt.Printf since log might not be imported in repository
				fmt.Printf("error rolling back transaction: %v\n", err)
			}
		}
	}()

	now := time.Now().UnixMilli()
	var teacherSubjects []TeacherSubject

	for _, teacherID := range teacherIDs {
		teacherSubjects = append(teacherSubjects, TeacherSubject{
			ID:        uuid.New(),
			TeacherID: teacherID,
			SubjectID: subjectID,
			CreatedAt: now,
			UpdatedAt: 0,
			CreatedBy: createdBy,
		})
	}

	insertQuery := `INSERT INTO teacher_subject (id, teacher_id, subject_id, created_at, created_by, updated_at) 
					VALUES (:id, :teacher_id, :subject_id, :created_at, :created_by, :updated_at) 
					ON CONFLICT (teacher_id, subject_id, is_deleted) DO NOTHING`

	_, err = tx.NamedExecContext(ctx, insertQuery, teacherSubjects)
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

func (r *repository) GetTeachersBySubjectID(ctx context.Context, subjectID uuid.UUID, query request.GetTeachersBySubjectQuery) ([]Teacher, int, error) {
	baseQuery := `
		SELECT u.id, u.name, u.email, u.is_verified, u.created_at, u.updated_at
		FROM users u
		INNER JOIN teacher_subject ts ON u.id = ts.teacher_id
		WHERE ts.subject_id = $1
	`

	countQuery := `
		SELECT COUNT(*)
		FROM users u
		INNER JOIN teacher_subject ts ON u.id = ts.teacher_id
		WHERE ts.subject_id = $1
	`

	var teachers []Teacher
	limitOrderQuery := fmt.Sprintf(" ORDER BY %s %s LIMIT $2 OFFSET $3", query.OrderBy, query.Order)

	err := r.db.SelectContext(ctx, &teachers, baseQuery+limitOrderQuery, subjectID, query.PageSize, query.GetOffset())
	if err != nil {
		return nil, 0, err
	}

	var total int
	err = r.db.GetContext(ctx, &total, countQuery, subjectID)
	if err != nil {
		return nil, 0, err
	}

	return teachers, total, nil
}
