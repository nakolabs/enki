package service

import (
	"context"
	"encoding/json"
	"enuma-elish/config"
	"enuma-elish/internal/exam/repository"
	"enuma-elish/internal/exam/service/data/request"
	"enuma-elish/internal/exam/service/data/response"
	commonHttp "enuma-elish/pkg/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreateExam(ctx context.Context, data request.CreateExamRequest) error
	GetDetailExam(ctx context.Context, examID uuid.UUID) (response.DetailExamResponse, error)
	GetListExams(ctx context.Context, query request.GetListExamQuery) (response.GetListExamResponse, *commonHttp.Meta, error)
	UpdateExam(ctx context.Context, examID uuid.UUID, data request.CreateExamRequest) error
	DeleteExam(ctx context.Context, examID uuid.UUID) error

	AssignExamToClass(ctx context.Context, data request.AssignExamToClassRequest) error
	GradeExam(ctx context.Context, data request.GradeExamRequest) error
	GetExamStudents(ctx context.Context, examID uuid.UUID, query request.GetExamStudentsQuery) (response.GetExamStudentsResponse, *commonHttp.Meta, error)

	// Student exam operations
	SubmitExamAnswers(ctx context.Context, studentID uuid.UUID, data request.SubmitExamAnswersRequest) error
	GetStudentExams(ctx context.Context, studentID uuid.UUID, query request.GetStudentExamsQuery) (response.GetStudentExamsResponse, *commonHttp.Meta, error)
	GetStudentExamDetail(ctx context.Context, examID, studentID uuid.UUID) (response.StudentExamDetailResponse, error)
}

type service struct {
	config     *config.Config
	repository repository.Repository
}

func New(config *config.Config, repository repository.Repository) Service {
	return &service{
		config:     config,
		repository: repository,
	}
}

func (s *service) CreateExam(ctx context.Context, data request.CreateExamRequest) error {
	// Validate that at least one question type is provided
	if err := data.Validate(); err != nil {
		return err
	}

	now := time.Now().UnixMilli()
	exam := repository.Exam{
		ID:        uuid.New(),
		Name:      data.Name,
		SchoolID:  data.SchoolID,
		SubjectID: data.SubjectID,
		CreatedAt: now,
		UpdatedAt: 0,
	}

	// Combine all question IDs
	var allQuestionIDs []uuid.UUID
	allQuestionIDs = append(allQuestionIDs, data.MultipleChoiceIDs...)
	allQuestionIDs = append(allQuestionIDs, data.EssayQuestionIDs...)

	err := s.repository.CreateExam(ctx, exam, allQuestionIDs)
	if err != nil {
		log.Err(err).Msg("Failed to create exam")
		return err
	}

	// Assign exam to class
	err = s.repository.AssignExamToClass(ctx, exam.ID, data.ClassID)
	if err != nil {
		log.Err(err).Msg("Failed to assign exam to class")
		return err
	}

	return nil
}

func (s *service) GetDetailExam(ctx context.Context, examID uuid.UUID) (response.DetailExamResponse, error) {
	exam, err := s.repository.GetExamByID(ctx, examID)
	if err != nil {
		log.Err(err).Msg("Failed to get exam")
		return response.DetailExamResponse{}, err
	}

	questions, err := s.repository.GetExamQuestions(ctx, examID)
	if err != nil {
		log.Err(err).Msg("Failed to get exam questions")
		return response.DetailExamResponse{}, err
	}

	var questionResponses []response.ExamQuestionResponse
	for _, question := range questions {
		questionResponse := response.ExamQuestionResponse{
			ID:            question.ID,
			Question:      question.Question,
			QuestionType:  question.QuestionType,
			CorrectAnswer: question.CorrectAnswer, // Include for teachers
		}

		// Parse options for multiple choice questions
		if question.QuestionType == "multiple_choice" && question.Options != nil {
			var options []map[string]string
			if err := json.Unmarshal([]byte(*question.Options), &options); err == nil {
				for _, option := range options {
					questionResponse.Options = append(questionResponse.Options, response.QuestionOptionResponse{
						ID:   option["id"],
						Text: option["text"],
					})
				}
			}
		}

		questionResponses = append(questionResponses, questionResponse)
	}

	res := response.DetailExamResponse{
		ID:          exam.ID,
		Name:        exam.Name,
		SchoolID:    exam.SchoolID,
		SubjectID:   exam.SubjectID,
		SubjectName: exam.SubjectName,
		Questions:   questionResponses,
		CreatedAt:   exam.CreatedAt,
		UpdatedAt:   exam.UpdatedAt,
	}

	return res, nil
}

func (s *service) GetListExams(ctx context.Context, query request.GetListExamQuery) (response.GetListExamResponse, *commonHttp.Meta, error) {
	exams, total, err := s.repository.GetListExams(ctx, query)
	if err != nil {
		log.Err(err).Msg("Failed to get exam list")
		return response.GetListExamResponse{}, nil, err
	}

	var res response.GetListExamResponse
	for _, exam := range exams {
		res = append(res, response.ExamResponse{
			ID:          exam.ID,
			Name:        exam.Name,
			SchoolID:    exam.SchoolID,
			SubjectID:   exam.SubjectID,
			SubjectName: exam.SubjectName,
			CreatedAt:   exam.CreatedAt,
			UpdatedAt:   exam.UpdatedAt,
		})
	}

	meta := commonHttp.NewMetaFromQuery(query, total)
	return res, meta, nil
}

func (s *service) UpdateExam(ctx context.Context, examID uuid.UUID, data request.CreateExamRequest) error {
	now := time.Now().UnixMilli()
	exam := repository.Exam{
		Name:      data.Name,
		SubjectID: data.SubjectID,
		UpdatedAt: now,
	}

	err := s.repository.UpdateExam(ctx, examID, exam)
	if err != nil {
		log.Err(err).Msg("Failed to update exam")
		return err
	}

	return nil
}

func (s *service) DeleteExam(ctx context.Context, examID uuid.UUID) error {
	err := s.repository.DeleteExam(ctx, examID)
	if err != nil {
		log.Err(err).Msg("Failed to delete exam")
		return err
	}
	return nil
}

func (s *service) AssignExamToClass(ctx context.Context, data request.AssignExamToClassRequest) error {
	err := s.repository.AssignExamToClass(ctx, data.ExamID, data.ClassID)
	if err != nil {
		log.Err(err).Msg("Failed to assign exam to class")
		return err
	}
	return nil
}

func (s *service) GradeExam(ctx context.Context, data request.GradeExamRequest) error {
	err := s.repository.GradeExam(ctx, data.ExamID, data.StudentID, data.Grade)
	if err != nil {
		log.Err(err).Msg("Failed to grade exam")
		return err
	}
	return nil
}

func (s *service) GetExamStudents(ctx context.Context, examID uuid.UUID, query request.GetExamStudentsQuery) (response.GetExamStudentsResponse, *commonHttp.Meta, error) {
	students, total, err := s.repository.GetExamStudents(ctx, examID, query)
	if err != nil {
		log.Err(err).Msg("Failed to get exam students")
		return response.GetExamStudentsResponse{}, nil, err
	}

	var res response.GetExamStudentsResponse
	for _, student := range students {
		res = append(res, response.ExamStudentResponse{
			ID:        student.ID,
			Name:      student.Name,
			Email:     student.Email,
			Grade:     student.Grade,
			IsGraded:  student.Grade != nil,
			CreatedAt: student.CreatedAt,
			UpdatedAt: student.UpdatedAt,
		})
	}

	meta := commonHttp.NewMetaFromQuery(query, total)
	return res, meta, nil
}

func (s *service) SubmitExamAnswers(ctx context.Context, studentID uuid.UUID, data request.SubmitExamAnswersRequest) error {
	// Submit answers first
	err := s.repository.SubmitExamAnswers(ctx, data.ExamID, studentID, data.Answers)
	if err != nil {
		log.Err(err).Msg("Failed to submit exam answers")
		return err
	}

	// Auto-grade multiple choice questions
	gradeResult, err := s.autoGradeMultipleChoice(ctx, data.ExamID, studentID, data.Answers)
	if err != nil {
		log.Err(err).Msg("Failed to auto-grade multiple choice questions")
		// Don't return error here, as submission was successful
	} else {
		// Store the auto-grade result if no essays need manual grading
		if !gradeResult.EssayPending {
			err = s.repository.AutoGradeExam(ctx, data.ExamID, studentID, gradeResult.Score, gradeResult.MaxScore)
			if err != nil {
				log.Err(err).Msg("Failed to store auto-grade result")
			}
		}
	}

	return nil
}

// Helper method to auto-grade multiple choice questions
func (s *service) autoGradeMultipleChoice(ctx context.Context, examID, studentID uuid.UUID, answers []request.ExamAnswer) (*response.AutoGradeResultResponse, error) {
	// Get all questions for the exam
	questions, err := s.repository.GetExamQuestions(ctx, examID)
	if err != nil {
		return nil, err
	}

	// Create maps for quick lookup
	questionMap := make(map[uuid.UUID]repository.Question)
	answerMap := make(map[uuid.UUID]request.ExamAnswer)

	for _, question := range questions {
		questionMap[question.ID] = question
	}

	for _, answer := range answers {
		answerMap[answer.QuestionID] = answer
	}

	var totalQuestions int
	var multipleChoiceQuestions int
	var correctAnswers int
	var hasEssayQuestions bool

	// Grade each question
	for _, question := range questions {
		totalQuestions++

		if question.QuestionType == "multiple_choice" {
			multipleChoiceQuestions++

			// Check if student answered this question
			if studentAnswer, exists := answerMap[question.ID]; exists {
				// For multiple choice, check selected option
				if studentAnswer.SelectedOption != nil && question.CorrectAnswer != nil {
					if *studentAnswer.SelectedOption == *question.CorrectAnswer {
						correctAnswers++
					}
				}
			}
		} else if question.QuestionType == "essay" {
			hasEssayQuestions = true
		}
	}

	// Calculate score (assuming each question has equal weight)
	var score float64
	var maxScore float64 = float64(totalQuestions)

	if multipleChoiceQuestions > 0 {
		// Only count multiple choice for now
		multipleChoiceScore := float64(correctAnswers)
		if hasEssayQuestions {
			// If there are essay questions, only grade multiple choice part
			score = multipleChoiceScore
			maxScore = float64(multipleChoiceQuestions)
		} else {
			// If only multiple choice, this is the final score
			score = multipleChoiceScore
		}
	}

	gradePercentage := (score / maxScore) * 100

	return &response.AutoGradeResultResponse{
		TotalQuestions:      totalQuestions,
		CorrectAnswers:      correctAnswers,
		Score:               score,
		MaxScore:            maxScore,
		GradePercentage:     gradePercentage,
		MultipleChoiceScore: float64(correctAnswers),
		EssayPending:        hasEssayQuestions,
	}, nil
}

func (s *service) GetStudentExams(ctx context.Context, studentID uuid.UUID, query request.GetStudentExamsQuery) (response.GetStudentExamsResponse, *commonHttp.Meta, error) {
	exams, total, err := s.repository.GetStudentExams(ctx, studentID, query)
	if err != nil {
		log.Err(err).Msg("Failed to get student exams")
		return response.GetStudentExamsResponse{}, nil, err
	}

	var res response.GetStudentExamsResponse
	for _, exam := range exams {
		res = append(res, response.StudentExamResponse{
			ID:          exam.ID,
			Name:        exam.Name,
			SchoolID:    exam.SchoolID,
			SubjectID:   exam.SubjectID,
			SubjectName: exam.SubjectName,
			Grade:       exam.Grade,
			IsSubmitted: exam.Answers != nil,
			IsGraded:    exam.Grade != nil,
			CreatedAt:   exam.CreatedAt,
			UpdatedAt:   exam.UpdatedAt,
		})
	}

	meta := commonHttp.NewMetaFromQuery(query, total)
	return res, meta, nil
}

func (s *service) GetStudentExamDetail(ctx context.Context, examID, studentID uuid.UUID) (response.StudentExamDetailResponse, error) {
	exam, err := s.repository.GetStudentExamDetail(ctx, examID, studentID)
	if err != nil {
		log.Err(err).Msg("Failed to get student exam detail")
		return response.StudentExamDetailResponse{}, err
	}

	questions, err := s.repository.GetExamQuestions(ctx, examID)
	if err != nil {
		log.Err(err).Msg("Failed to get exam questions")
		return response.StudentExamDetailResponse{}, err
	}

	var questionResponses []response.ExamQuestionResponse
	for _, question := range questions {
		questionResponse := response.ExamQuestionResponse{
			ID:           question.ID,
			Question:     question.Question,
			QuestionType: question.QuestionType,
			// Don't include correct answer for students
		}

		// Parse options for multiple choice questions
		if question.QuestionType == "multiple_choice" && question.Options != nil {
			var options []map[string]string
			if err := json.Unmarshal([]byte(*question.Options), &options); err == nil {
				for _, option := range options {
					questionResponse.Options = append(questionResponse.Options, response.QuestionOptionResponse{
						ID:   option["id"],
						Text: option["text"],
					})
				}
			}
		}

		questionResponses = append(questionResponses, questionResponse)
	}

	var answerResponses []response.StudentExamAnswerResponse
	if exam.Answers != nil {
		var answers []request.ExamAnswer
		err := json.Unmarshal([]byte(*exam.Answers), &answers)
		if err == nil {
			for _, answer := range answers {
				answerResponses = append(answerResponses, response.StudentExamAnswerResponse{
					QuestionID:     answer.QuestionID,
					Answer:         answer.Answer,
					SelectedOption: answer.SelectedOption,
				})
			}
		}
	}

	res := response.StudentExamDetailResponse{
		ID:          exam.ID,
		Name:        exam.Name,
		SchoolID:    exam.SchoolID,
		SubjectID:   exam.SubjectID,
		SubjectName: exam.SubjectName,
		Questions:   questionResponses,
		Answers:     answerResponses,
		Grade:       exam.Grade,
		IsSubmitted: exam.Answers != nil,
		IsGraded:    exam.Grade != nil,
		CreatedAt:   exam.CreatedAt,
		UpdatedAt:   exam.UpdatedAt,
	}

	return res, nil
}
