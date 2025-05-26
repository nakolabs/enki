package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"enuma-elish/internal/exam/service/data/request"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Exam struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	SchoolID  uuid.UUID `db:"school_id"`
	SubjectID uuid.UUID `db:"subject_id"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

type ExamClass struct {
	ID        uuid.UUID `db:"id"`
	ExamID    uuid.UUID `db:"exam_id"`
	ClassID   uuid.UUID `db:"class_id"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

type ExamGrade struct {
	ID        uuid.UUID `db:"id"`
	ExamID    uuid.UUID `db:"exam_id"`
	StudentID uuid.UUID `db:"student_id"`
	Grade     *float64  `db:"grade"`
	Answers   *string   `db:"answers"` // JSON string of answers
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

type ExamQuestion struct {
	ID         uuid.UUID `db:"id"`
	ExamID     uuid.UUID `db:"exam_id"`
	QuestionID uuid.UUID `db:"question_id"`
	CreatedAt  int64     `db:"created_at"`
	UpdatedAt  int64     `db:"updated_at"`
}

type Question struct {
	ID            uuid.UUID `db:"id"`
	Question      string    `db:"question"`
	QuestionType  string    `db:"question_type"`
	Options       *string   `db:"options"`        // JSON string for multiple choice options
	CorrectAnswer *string   `db:"correct_answer"` // Correct option ID for multiple choice
}

type ExamWithSubject struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	SchoolID    uuid.UUID `db:"school_id"`
	SubjectID   uuid.UUID `db:"subject_id"`
	SubjectName string    `db:"subject_name"`
	CreatedAt   int64     `db:"created_at"`
	UpdatedAt   int64     `db:"updated_at"`
}

type StudentWithGrade struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Grade     *float64  `db:"grade"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

type StudentExamWithAnswers struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	SchoolID    uuid.UUID `db:"school_id"`
	SubjectID   uuid.UUID `db:"subject_id"`
	SubjectName string    `db:"subject_name"`
	Grade       *float64  `db:"grade"`
	Answers     *string   `db:"answers"`
	CreatedAt   int64     `db:"created_at"`
	UpdatedAt   int64     `db:"updated_at"`
}

type Repository interface {
	CreateExam(ctx context.Context, exam Exam, questionIDs []uuid.UUID) error
	GetExamByID(ctx context.Context, examID uuid.UUID) (*ExamWithSubject, error)
	GetListExams(ctx context.Context, query request.GetListExamQuery) ([]ExamWithSubject, int, error)
	UpdateExam(ctx context.Context, examID uuid.UUID, exam Exam) error
	DeleteExam(ctx context.Context, examID uuid.UUID) error

	AssignExamToClass(ctx context.Context, examID, classID uuid.UUID) error
	GetExamQuestions(ctx context.Context, examID uuid.UUID) ([]Question, error)

	GradeExam(ctx context.Context, examID, studentID uuid.UUID, grade float64) error
	AutoGradeExam(ctx context.Context, examID, studentID uuid.UUID, totalScore, maxScore float64) error
	GetExamStudents(ctx context.Context, examID uuid.UUID, query request.GetExamStudentsQuery) ([]StudentWithGrade, int, error)

	// Student exam operations
	SubmitExamAnswers(ctx context.Context, examID, studentID uuid.UUID, answers []request.ExamAnswer) error
	GetStudentExams(ctx context.Context, studentID uuid.UUID, query request.GetStudentExamsQuery) ([]StudentExamWithAnswers, int, error)
	GetStudentExamDetail(ctx context.Context, examID, studentID uuid.UUID) (*StudentExamWithAnswers, error)

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

func (r *repository) CreateExam(ctx context.Context, exam Exam, questionIDs []uuid.UUID) error {
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

	// Insert exam
	insertExamQuery := `INSERT INTO exam (id, name, school_id, subject_id, created_at, updated_at) 
					    VALUES (:id, :name, :school_id, :subject_id, :created_at, :updated_at)`
	_, err = tx.NamedExecContext(ctx, insertExamQuery, exam)
	if err != nil {
		return err
	}

	// Insert exam questions
	now := time.Now().UnixMilli()
	var examQuestions []ExamQuestion
	for _, questionID := range questionIDs {
		examQuestions = append(examQuestions, ExamQuestion{
			ID:         uuid.New(),
			ExamID:     exam.ID,
			QuestionID: questionID,
			CreatedAt:  now,
			UpdatedAt:  0,
		})
	}

	insertQuestionsQuery := `INSERT INTO exam_question (id, exam_id, question_id, created_at, updated_at) 
							 VALUES (:id, :exam_id, :question_id, :created_at, :updated_at)`
	_, err = tx.NamedExecContext(ctx, insertQuestionsQuery, examQuestions)
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

func (r *repository) GetExamByID(ctx context.Context, examID uuid.UUID) (*ExamWithSubject, error) {
	query := `SELECT e.id, e.name, e.school_id, e.subject_id, s.name as subject_name, e.created_at, e.updated_at
			  FROM exam e
			  JOIN subject s ON e.subject_id = s.id
			  WHERE e.id = $1`

	var exam ExamWithSubject
	err := r.db.GetContext(ctx, &exam, query, examID)
	if err != nil {
		return nil, err
	}
	return &exam, nil
}

func (r *repository) GetListExams(ctx context.Context, query request.GetListExamQuery) ([]ExamWithSubject, int, error) {
	baseQuery := `SELECT e.id, e.name, e.school_id, e.subject_id, s.name as subject_name, e.created_at, e.updated_at
				  FROM exam e
				  JOIN subject s ON e.subject_id = s.id
				  WHERE e.school_id = $1`

	countQuery := `SELECT COUNT(*)
				   FROM exam e
				   WHERE e.school_id = $1`

	params := []interface{}{query.SchoolID}
	paramCount := 1

	if query.SubjectID != "" {
		paramCount++
		baseQuery += fmt.Sprintf(" AND e.subject_id = $%d", paramCount)
		countQuery += fmt.Sprintf(" AND e.subject_id = $%d", paramCount)
		params = append(params, query.SubjectID)
	}

	var exams []ExamWithSubject
	limitOrderQuery := fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", query.OrderBy, query.Order, paramCount+1, paramCount+2)
	params = append(params, query.PageSize, query.GetOffset())

	err := r.db.SelectContext(ctx, &exams, baseQuery+limitOrderQuery, params...)
	if err != nil {
		return nil, 0, err
	}

	var total int
	err = r.db.GetContext(ctx, &total, countQuery, params[:paramCount]...)
	if err != nil {
		return nil, 0, err
	}

	return exams, total, nil
}

func (r *repository) UpdateExam(ctx context.Context, examID uuid.UUID, exam Exam) error {
	updateQuery := `UPDATE exam SET name = $1, subject_id = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, updateQuery, exam.Name, exam.SubjectID, exam.UpdatedAt, examID)
	return err
}

func (r *repository) DeleteExam(ctx context.Context, examID uuid.UUID) error {
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

	// Delete related records first
	_, err = tx.ExecContext(ctx, "DELETE FROM exam_question WHERE exam_id = $1", examID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM exam_class WHERE exam_id = $1", examID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM exam_grade WHERE exam_id = $1", examID)
	if err != nil {
		return err
	}

	// Delete exam
	_, err = tx.ExecContext(ctx, "DELETE FROM exam WHERE id = $1", examID)
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

func (r *repository) AssignExamToClass(ctx context.Context, examID, classID uuid.UUID) error {
	now := time.Now().UnixMilli()
	examClass := ExamClass{
		ID:        uuid.New(),
		ExamID:    examID,
		ClassID:   classID,
		CreatedAt: now,
		UpdatedAt: 0,
	}

	insertQuery := `INSERT INTO exam_class (id, exam_id, class_id, created_at, updated_at) 
					VALUES (:id, :exam_id, :class_id, :created_at, :updated_at)
					ON CONFLICT (exam_id, class_id) DO NOTHING`
	_, err := r.db.NamedExecContext(ctx, insertQuery, examClass)
	return err
}

func (r *repository) GetExamQuestions(ctx context.Context, examID uuid.UUID) ([]Question, error) {
	query := `SELECT q.id, q.question, q.question_type, q.options, q.correct_answer
			  FROM question q
			  JOIN exam_question eq ON q.id = eq.question_id
			  WHERE eq.exam_id = $1`

	var questions []Question
	err := r.db.SelectContext(ctx, &questions, query, examID)
	return questions, err
}

func (r *repository) GradeExam(ctx context.Context, examID, studentID uuid.UUID, grade float64) error {
	now := time.Now().UnixMilli()

	// Upsert grade
	upsertQuery := `INSERT INTO exam_grade (id, exam_id, student_id, grade, created_at, updated_at) 
					VALUES ($1, $2, $3, $4, $5, $6)
					ON CONFLICT (exam_id, student_id) 
					DO UPDATE SET grade = $4, updated_at = $6`

	_, err := r.db.ExecContext(ctx, upsertQuery, uuid.New(), examID, studentID, grade, now, now)
	return err
}

func (r *repository) AutoGradeExam(ctx context.Context, examID, studentID uuid.UUID, totalScore, maxScore float64) error {
	now := time.Now().UnixMilli()

	// Calculate percentage score
	gradePercentage := (totalScore / maxScore) * 100

	// Upsert grade
	upsertQuery := `INSERT INTO exam_grade (id, exam_id, student_id, grade, created_at, updated_at) 
					VALUES ($1, $2, $3, $4, $5, $6)
					ON CONFLICT (exam_id, student_id) 
					DO UPDATE SET grade = $4, updated_at = $6`

	_, err := r.db.ExecContext(ctx, upsertQuery, uuid.New(), examID, studentID, gradePercentage, now, now)
	return err
}

func (r *repository) GetExamStudents(ctx context.Context, examID uuid.UUID, query request.GetExamStudentsQuery) ([]StudentWithGrade, int, error) {
	baseQuery := `SELECT u.id, u.name, u.email, eg.grade, u.created_at, u.updated_at
				  FROM users u
				  JOIN class_student cs ON u.id = cs.student_id
				  JOIN exam_class ec ON cs.class_id = ec.class_id
				  LEFT JOIN exam_grade eg ON u.id = eg.student_id AND eg.exam_id = $1
				  WHERE ec.exam_id = $1`

	countQuery := `SELECT COUNT(*)
				   FROM users u
				   JOIN class_student cs ON u.id = cs.student_id
				   JOIN exam_class ec ON cs.class_id = ec.class_id
				   WHERE ec.exam_id = $1`

	var students []StudentWithGrade
	limitOrderQuery := fmt.Sprintf(" ORDER BY %s %s LIMIT $2 OFFSET $3", query.OrderBy, query.Order)

	err := r.db.SelectContext(ctx, &students, baseQuery+limitOrderQuery, examID, query.PageSize, query.GetOffset())
	if err != nil {
		return nil, 0, err
	}

	var total int
	err = r.db.GetContext(ctx, &total, countQuery, examID)
	if err != nil {
		return nil, 0, err
	}

	return students, total, nil
}

func (r *repository) SubmitExamAnswers(ctx context.Context, examID, studentID uuid.UUID, answers []request.ExamAnswer) error {
	now := time.Now().UnixMilli()

	// Convert answers to JSON
	answersJSON, err := json.Marshal(answers)
	if err != nil {
		return err
	}

	// Upsert exam submission
	upsertQuery := `INSERT INTO exam_grade (id, exam_id, student_id, answers, created_at, updated_at) 
					VALUES ($1, $2, $3, $4, $5, $6)
					ON CONFLICT (exam_id, student_id) 
					DO UPDATE SET answers = $4, updated_at = $6`

	_, err = r.db.ExecContext(ctx, upsertQuery, uuid.New(), examID, studentID, string(answersJSON), now, now)
	return err
}

func (r *repository) GetStudentExams(ctx context.Context, studentID uuid.UUID, query request.GetStudentExamsQuery) ([]StudentExamWithAnswers, int, error) {
	baseQuery := `SELECT e.id, e.name, e.school_id, e.subject_id, s.name as subject_name, 
				  eg.grade, eg.answers, e.created_at, e.updated_at
				  FROM exam e
				  JOIN subject s ON e.subject_id = s.id
				  JOIN exam_class ec ON e.id = ec.exam_id
				  JOIN class_student cs ON ec.class_id = cs.class_id
				  LEFT JOIN exam_grade eg ON e.id = eg.exam_id AND eg.student_id = $1
				  WHERE cs.student_id = $1 AND e.school_id = $2`

	countQuery := `SELECT COUNT(*)
				   FROM exam e
				   JOIN exam_class ec ON e.id = ec.exam_id
				   JOIN class_student cs ON ec.class_id = cs.class_id
				   WHERE cs.student_id = $1 AND e.school_id = $2`

	params := []interface{}{studentID, query.SchoolID}
	paramCount := 2

	if query.SubjectID != "" {
		paramCount++
		baseQuery += fmt.Sprintf(" AND e.subject_id = $%d", paramCount)
		countQuery += fmt.Sprintf(" AND e.subject_id = $%d", paramCount)
		params = append(params, query.SubjectID)
	}

	var exams []StudentExamWithAnswers
	limitOrderQuery := fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", query.OrderBy, query.Order, paramCount+1, paramCount+2)
	params = append(params, query.PageSize, query.GetOffset())

	err := r.db.SelectContext(ctx, &exams, baseQuery+limitOrderQuery, params...)
	if err != nil {
		return nil, 0, err
	}

	var total int
	err = r.db.GetContext(ctx, &total, countQuery, params[:paramCount]...)
	if err != nil {
		return nil, 0, err
	}

	return exams, total, nil
}

func (r *repository) GetStudentExamDetail(ctx context.Context, examID, studentID uuid.UUID) (*StudentExamWithAnswers, error) {
	query := `SELECT e.id, e.name, e.school_id, e.subject_id, s.name as subject_name, 
			  eg.grade, eg.answers, e.created_at, e.updated_at
			  FROM exam e
			  JOIN subject s ON e.subject_id = s.id
			  LEFT JOIN exam_grade eg ON e.id = eg.exam_id AND eg.student_id = $2
			  WHERE e.id = $1`

	var exam StudentExamWithAnswers
	err := r.db.GetContext(ctx, &exam, query, examID, studentID)
	if err != nil {
		return nil, err
	}
	return &exam, nil
}

func (r *repository) Redis() *redis.Client {
	return r.rdb
}

func (r *repository) Tx(ctx context.Context, options *sql.TxOptions) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, options)
}
