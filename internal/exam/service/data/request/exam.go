package request

import (
	commonHttp "enuma-elish/pkg/http"
	"fmt"

	"github.com/google/uuid"
)

type CreateExamRequest struct {
	Name              string      `json:"name" validate:"required"`
	SchoolID          uuid.UUID   `json:"school_id" validate:"required"`
	SubjectID         uuid.UUID   `json:"subject_id" validate:"required"`
	ClassID           uuid.UUID   `json:"class_id" validate:"required"`
	MultipleChoiceIDs []uuid.UUID `json:"multiple_choice_ids"`
	EssayQuestionIDs  []uuid.UUID `json:"essay_question_ids"`
}

// Custom validation to ensure at least one question type is provided
func (r CreateExamRequest) Validate() error {
	if len(r.MultipleChoiceIDs) == 0 && len(r.EssayQuestionIDs) == 0 {
		return fmt.Errorf("at least one question (multiple choice or essay) must be provided")
	}
	return nil
}

type AssignExamToClassRequest struct {
	ExamID  uuid.UUID `json:"exam_id" validate:"required"`
	ClassID uuid.UUID `json:"class_id" validate:"required"`
}

type GradeExamRequest struct {
	ExamID    uuid.UUID `json:"exam_id" validate:"required"`
	StudentID uuid.UUID `json:"student_id" validate:"required"`
	Grade     float64   `json:"grade" validate:"required,min=0,max=100"`
}

type SubmitExamAnswersRequest struct {
	ExamID  uuid.UUID    `json:"exam_id" validate:"required"`
	Answers []ExamAnswer `json:"answers" validate:"required,min=1"`
}

type ExamAnswer struct {
	QuestionID     uuid.UUID `json:"question_id" validate:"required"`
	Answer         string    `json:"answer" validate:"required"`
	SelectedOption *string   `json:"selected_option,omitempty"` // For multiple choice
}

type GetStudentExamsQuery struct {
	SchoolID  string `form:"school_id" binding:"uuid"`
	SubjectID string `form:"subject_id"`
	ClassID   string `form:"class_id"`
	commonHttp.Query
}

func (q GetStudentExamsQuery) Get() (commonHttp.Query, map[string]interface{}) {
	f := map[string]interface{}{
		"school_id": q.SchoolID,
	}
	if q.SubjectID != "" {
		f["subject_id"] = q.SubjectID
	}
	if q.ClassID != "" {
		f["class_id"] = q.ClassID
	}
	return q.Query, f
}

type GetListExamQuery struct {
	SchoolID  string `form:"school_id" binding:"uuid"`
	SubjectID string `form:"subject_id"`
	ClassID   string `form:"class_id"`
	commonHttp.Query
}

func (q GetListExamQuery) Get() (commonHttp.Query, map[string]interface{}) {
	f := map[string]interface{}{
		"school_id": q.SchoolID,
	}
	if q.SubjectID != "" {
		f["subject_id"] = q.SubjectID
	}
	if q.ClassID != "" {
		f["class_id"] = q.ClassID
	}
	return q.Query, f
}

type GetExamStudentsQuery struct {
	commonHttp.Query
}
