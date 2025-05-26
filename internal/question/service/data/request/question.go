package request

import (
	commonHttp "enuma-elish/pkg/http"
	"fmt"

	"github.com/google/uuid"
)

type CreateQuestionRequest struct {
	Question        string                  `json:"question" validate:"required"`
	QuestionType    string                  `json:"question_type" validate:"required,oneof=multiple_choice essay"`
	Options         []QuestionOptionRequest `json:"options,omitempty"`
	CorrectAnswer   *string                 `json:"correct_answer,omitempty"`
	SchoolID        uuid.UUID               `json:"school_id" validate:"required"`
	SubjectID       uuid.UUID               `json:"subject_id" validate:"required"`
	DifficultyLevel string                  `json:"difficulty_level" validate:"required,oneof=easy medium hard"`
	Points          int                     `json:"points" validate:"required,min=1"`
}

type QuestionOptionRequest struct {
	ID   string `json:"id" validate:"required"`
	Text string `json:"text" validate:"required"`
}

func (r CreateQuestionRequest) Validate() error {
	if r.QuestionType == "multiple_choice" {
		if len(r.Options) < 2 {
			return fmt.Errorf("multiple choice questions must have at least 2 options")
		}
		if r.CorrectAnswer == nil {
			return fmt.Errorf("multiple choice questions must have a correct answer")
		}

		// Validate that correct answer exists in options
		correctAnswerExists := false
		for _, option := range r.Options {
			if option.ID == *r.CorrectAnswer {
				correctAnswerExists = true
				break
			}
		}
		if !correctAnswerExists {
			return fmt.Errorf("correct answer must match one of the option IDs")
		}
	} else if r.QuestionType == "essay" {
		if len(r.Options) > 0 {
			return fmt.Errorf("essay questions should not have options")
		}
		if r.CorrectAnswer != nil {
			return fmt.Errorf("essay questions should not have a correct answer")
		}
	}

	return nil
}

type UpdateQuestionRequest struct {
	Question        string                  `json:"question" validate:"required"`
	QuestionType    string                  `json:"question_type" validate:"required,oneof=multiple_choice essay"`
	Options         []QuestionOptionRequest `json:"options,omitempty"`
	CorrectAnswer   *string                 `json:"correct_answer,omitempty"`
	SubjectID       uuid.UUID               `json:"subject_id" validate:"required"`
	DifficultyLevel string                  `json:"difficulty_level" validate:"required,oneof=easy medium hard"`
	Points          int                     `json:"points" validate:"required,min=1"`
}

func (r UpdateQuestionRequest) Validate() error {
	if r.QuestionType == "multiple_choice" {
		if len(r.Options) < 2 {
			return fmt.Errorf("multiple choice questions must have at least 2 options")
		}
		if r.CorrectAnswer == nil {
			return fmt.Errorf("multiple choice questions must have a correct answer")
		}

		// Validate that correct answer exists in options
		correctAnswerExists := false
		for _, option := range r.Options {
			if option.ID == *r.CorrectAnswer {
				correctAnswerExists = true
				break
			}
		}
		if !correctAnswerExists {
			return fmt.Errorf("correct answer must match one of the option IDs")
		}
	} else if r.QuestionType == "essay" {
		if len(r.Options) > 0 {
			return fmt.Errorf("essay questions should not have options")
		}
		if r.CorrectAnswer != nil {
			return fmt.Errorf("essay questions should not have a correct answer")
		}
	}

	return nil
}

type GetListQuestionQuery struct {
	SchoolID        string `form:"school_id" binding:"required,uuid"`
	SubjectID       string `form:"subject_id"`
	QuestionType    string `form:"question_type"`
	DifficultyLevel string `form:"difficulty_level"`
	commonHttp.Query
}

func (q GetListQuestionQuery) Get() (commonHttp.Query, map[string]interface{}) {
	f := map[string]interface{}{
		"school_id": q.SchoolID,
	}
	if q.SubjectID != "" {
		f["subject_id"] = q.SubjectID
	}
	if q.QuestionType != "" {
		f["question_type"] = q.QuestionType
	}
	if q.DifficultyLevel != "" {
		f["difficulty_level"] = q.DifficultyLevel
	}
	return q.Query, f
}

type GetQuestionsByTypeQuery struct {
	SchoolID     uuid.UUID `form:"school_id" binding:"required,uuid"`
	SubjectID    uuid.UUID `form:"subject_id" binding:"required,uuid"`
	QuestionType string    `form:"question_type" binding:"required,oneof=multiple_choice essay"`
}
