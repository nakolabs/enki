package repository

import (
	"context"
	"database/sql"
	"enuma-elish/internal/question/service/data/request"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Question struct {
	ID              uuid.UUID `db:"id"`
	Question        string    `db:"question"`
	QuestionType    string    `db:"question_type"`
	Options         *string   `db:"options"`
	CorrectAnswer   *string   `db:"correct_answer"`
	SchoolID        uuid.UUID `db:"school_id"`
	SubjectID       uuid.UUID `db:"subject_id"`
	DifficultyLevel string    `db:"difficulty_level"`
	Points          int       `db:"points"`
	CreatedAt       int64     `db:"created_at"`
	UpdatedAt       int64     `db:"updated_at"`
}

type QuestionWithSubject struct {
	ID              uuid.UUID `db:"id"`
	Question        string    `db:"question"`
	QuestionType    string    `db:"question_type"`
	Options         *string   `db:"options"`
	CorrectAnswer   *string   `db:"correct_answer"`
	SchoolID        uuid.UUID `db:"school_id"`
	SubjectID       uuid.UUID `db:"subject_id"`
	SubjectName     string    `db:"subject_name"`
	DifficultyLevel string    `db:"difficulty_level"`
	Points          int       `db:"points"`
	CreatedAt       int64     `db:"created_at"`
	UpdatedAt       int64     `db:"updated_at"`
}

type Repository interface {
	CreateQuestion(ctx context.Context, question Question) error
	GetQuestionByID(ctx context.Context, questionID uuid.UUID) (*QuestionWithSubject, error)
	GetListQuestions(ctx context.Context, query request.GetListQuestionQuery) ([]QuestionWithSubject, int, error)
	UpdateQuestion(ctx context.Context, questionID uuid.UUID, question Question) error
	DeleteQuestion(ctx context.Context, questionID uuid.UUID) error
	GetQuestionsByType(ctx context.Context, schoolID, subjectID uuid.UUID, questionType string) ([]QuestionWithSubject, error)

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

func (r *repository) CreateQuestion(ctx context.Context, question Question) error {
	insertQuery := `INSERT INTO question (id, question, question_type, options, correct_answer, school_id, subject_id, difficulty_level, points, created_at, updated_at) 
					VALUES (:id, :question, :question_type, :options, :correct_answer, :school_id, :subject_id, :difficulty_level, :points, :created_at, :updated_at)`

	_, err := r.db.NamedExecContext(ctx, insertQuery, question)
	return err
}

func (r *repository) GetQuestionByID(ctx context.Context, questionID uuid.UUID) (*QuestionWithSubject, error) {
	query := `SELECT q.id, q.question, q.question_type, q.options, q.correct_answer, q.school_id, q.subject_id, s.name as subject_name, 
			  q.difficulty_level, q.points, q.created_at, q.updated_at
			  FROM question q
			  JOIN subject s ON q.subject_id = s.id
			  WHERE q.id = $1`

	var question QuestionWithSubject
	err := r.db.GetContext(ctx, &question, query, questionID)
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *repository) GetListQuestions(ctx context.Context, query request.GetListQuestionQuery) ([]QuestionWithSubject, int, error) {
	baseQuery := `SELECT q.id, q.question, q.question_type, q.options, q.correct_answer, q.school_id, q.subject_id, s.name as subject_name, 
				  q.difficulty_level, q.points, q.created_at, q.updated_at
				  FROM question q
				  JOIN subject s ON q.subject_id = s.id
				  WHERE q.school_id = $1`

	countQuery := `SELECT COUNT(*)
				   FROM question q
				   WHERE q.school_id = $1`

	params := []interface{}{query.SchoolID}
	paramCount := 1

	if query.SubjectID != "" {
		paramCount++
		baseQuery += fmt.Sprintf(" AND q.subject_id = $%d", paramCount)
		countQuery += fmt.Sprintf(" AND q.subject_id = $%d", paramCount)
		params = append(params, query.SubjectID)
	}

	if query.QuestionType != "" {
		paramCount++
		baseQuery += fmt.Sprintf(" AND q.question_type = $%d", paramCount)
		countQuery += fmt.Sprintf(" AND q.question_type = $%d", paramCount)
		params = append(params, query.QuestionType)
	}

	if query.DifficultyLevel != "" {
		paramCount++
		baseQuery += fmt.Sprintf(" AND q.difficulty_level = $%d", paramCount)
		countQuery += fmt.Sprintf(" AND q.difficulty_level = $%d", paramCount)
		params = append(params, query.DifficultyLevel)
	}

	if query.Search != "" {
		paramCount++
		baseQuery += fmt.Sprintf(" AND q.question ILIKE $%d", paramCount)
		countQuery += fmt.Sprintf(" AND q.question ILIKE $%d", paramCount)
		params = append(params, "%"+query.Search+"%")
	}

	var questions []QuestionWithSubject
	limitOrderQuery := fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", query.OrderBy, query.Order, paramCount+1, paramCount+2)
	params = append(params, query.PageSize, query.GetOffset())

	err := r.db.SelectContext(ctx, &questions, baseQuery+limitOrderQuery, params...)
	if err != nil {
		return nil, 0, err
	}

	var total int
	err = r.db.GetContext(ctx, &total, countQuery, params[:paramCount]...)
	if err != nil {
		return nil, 0, err
	}

	return questions, total, nil
}

func (r *repository) UpdateQuestion(ctx context.Context, questionID uuid.UUID, question Question) error {
	updateQuery := `UPDATE question SET question = $1, question_type = $2, options = $3, correct_answer = $4, 
					subject_id = $5, difficulty_level = $6, points = $7, updated_at = $8 WHERE id = $9`

	_, err := r.db.ExecContext(ctx, updateQuery, question.Question, question.QuestionType, question.Options,
		question.CorrectAnswer, question.SubjectID, question.DifficultyLevel, question.Points, question.UpdatedAt, questionID)
	return err
}

func (r *repository) DeleteQuestion(ctx context.Context, questionID uuid.UUID) error {
	// Check if question is used in any exam
	checkQuery := `SELECT COUNT(*) FROM exam_question WHERE question_id = $1`
	var count int
	err := r.db.GetContext(ctx, &count, checkQuery, questionID)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("cannot delete question: it is being used in %d exam(s)", count)
	}

	deleteQuery := `DELETE FROM question WHERE id = $1`
	_, err = r.db.ExecContext(ctx, deleteQuery, questionID)
	return err
}

func (r *repository) GetQuestionsByType(ctx context.Context, schoolID, subjectID uuid.UUID, questionType string) ([]QuestionWithSubject, error) {
	query := `SELECT q.id, q.question, q.question_type, q.options, q.correct_answer, q.school_id, q.subject_id, s.name as subject_name, 
			  q.difficulty_level, q.points, q.created_at, q.updated_at
			  FROM question q
			  JOIN subject s ON q.subject_id = s.id
			  WHERE q.school_id = $1 AND q.subject_id = $2 AND q.question_type = $3
			  ORDER BY q.created_at DESC`

	var questions []QuestionWithSubject
	err := r.db.SelectContext(ctx, &questions, query, schoolID, subjectID, questionType)
	return questions, err
}

func (r *repository) Redis() *redis.Client {
	return r.rdb
}

func (r *repository) Tx(ctx context.Context, options *sql.TxOptions) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, options)
}
