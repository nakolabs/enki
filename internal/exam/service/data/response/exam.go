package response

import "github.com/google/uuid"

type ExamResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	SchoolID    uuid.UUID `json:"school_id"`
	SubjectID   uuid.UUID `json:"subject_id"`
	SubjectName string    `json:"subject_name"`
	CreatedAt   int64     `json:"created_at"`
	UpdatedAt   int64     `json:"updated_at"`
}

type GetListExamResponse []ExamResponse

type DetailExamResponse struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	SchoolID    uuid.UUID              `json:"school_id"`
	SubjectID   uuid.UUID              `json:"subject_id"`
	SubjectName string                 `json:"subject_name"`
	Questions   []ExamQuestionResponse `json:"questions"`
	CreatedAt   int64                  `json:"created_at"`
	UpdatedAt   int64                  `json:"updated_at"`
}

type ExamQuestionResponse struct {
	ID            uuid.UUID                `json:"id"`
	Question      string                   `json:"question"`
	QuestionType  string                   `json:"question_type"`
	Options       []QuestionOptionResponse `json:"options,omitempty"`
	CorrectAnswer *string                  `json:"correct_answer,omitempty"` // Only for teachers
}

type QuestionOptionResponse struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type ExamStudentResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Grade     *float64  `json:"grade"`
	IsGraded  bool      `json:"is_graded"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

type GetExamStudentsResponse []ExamStudentResponse

type StudentExamResponse struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	SchoolID    uuid.UUID              `json:"school_id"`
	SubjectID   uuid.UUID              `json:"subject_id"`
	SubjectName string                 `json:"subject_name"`
	Questions   []ExamQuestionResponse `json:"questions"`
	Grade       *float64               `json:"grade"`
	IsSubmitted bool                   `json:"is_submitted"`
	IsGraded    bool                   `json:"is_graded"`
	CreatedAt   int64                  `json:"created_at"`
	UpdatedAt   int64                  `json:"updated_at"`
}

type GetStudentExamsResponse []StudentExamResponse

type StudentExamDetailResponse struct {
	ID          uuid.UUID                   `json:"id"`
	Name        string                      `json:"name"`
	SchoolID    uuid.UUID                   `json:"school_id"`
	SubjectID   uuid.UUID                   `json:"subject_id"`
	SubjectName string                      `json:"subject_name"`
	Questions   []ExamQuestionResponse      `json:"questions"`
	Answers     []StudentExamAnswerResponse `json:"answers"`
	Grade       *float64                    `json:"grade"`
	IsSubmitted bool                        `json:"is_submitted"`
	IsGraded    bool                        `json:"is_graded"`
	CreatedAt   int64                       `json:"created_at"`
	UpdatedAt   int64                       `json:"updated_at"`
}

type StudentExamAnswerResponse struct {
	QuestionID     uuid.UUID `json:"question_id"`
	Answer         string    `json:"answer"`
	SelectedOption *string   `json:"selected_option,omitempty"`
}

type AutoGradeResultResponse struct {
	TotalQuestions      int     `json:"total_questions"`
	CorrectAnswers      int     `json:"correct_answers"`
	Score               float64 `json:"score"`
	MaxScore            float64 `json:"max_score"`
	GradePercentage     float64 `json:"grade_percentage"`
	MultipleChoiceScore float64 `json:"multiple_choice_score"`
	EssayPending        bool    `json:"essay_pending"`
}
