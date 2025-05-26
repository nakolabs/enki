package response

import "github.com/google/uuid"

type QuestionResponse struct {
	ID              uuid.UUID                `json:"id"`
	Question        string                   `json:"question"`
	QuestionType    string                   `json:"question_type"`
	Options         []QuestionOptionResponse `json:"options,omitempty"`
	CorrectAnswer   *string                  `json:"correct_answer,omitempty"`
	SchoolID        uuid.UUID                `json:"school_id"`
	SubjectID       uuid.UUID                `json:"subject_id"`
	SubjectName     string                   `json:"subject_name"`
	DifficultyLevel string                   `json:"difficulty_level"`
	Points          int                      `json:"points"`
	CreatedAt       int64                    `json:"created_at"`
	UpdatedAt       int64                    `json:"updated_at"`
}

type QuestionOptionResponse struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type GetListQuestionResponse []QuestionResponse

type DetailQuestionResponse struct {
	ID              uuid.UUID                `json:"id"`
	Question        string                   `json:"question"`
	QuestionType    string                   `json:"question_type"`
	Options         []QuestionOptionResponse `json:"options,omitempty"`
	CorrectAnswer   *string                  `json:"correct_answer,omitempty"`
	SchoolID        uuid.UUID                `json:"school_id"`
	SubjectID       uuid.UUID                `json:"subject_id"`
	SubjectName     string                   `json:"subject_name"`
	DifficultyLevel string                   `json:"difficulty_level"`
	Points          int                      `json:"points"`
	CreatedAt       int64                    `json:"created_at"`
	UpdatedAt       int64                    `json:"updated_at"`
}

type QuestionsByTypeResponse []QuestionResponse
