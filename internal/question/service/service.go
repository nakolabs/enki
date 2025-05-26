package service

import (
	"context"
	"encoding/json"
	"enuma-elish/config"
	"enuma-elish/internal/question/repository"
	"enuma-elish/internal/question/service/data/request"
	"enuma-elish/internal/question/service/data/response"
	commonHttp "enuma-elish/pkg/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreateQuestion(ctx context.Context, data request.CreateQuestionRequest) error
	GetDetailQuestion(ctx context.Context, questionID uuid.UUID) (response.DetailQuestionResponse, error)
	GetListQuestions(ctx context.Context, query request.GetListQuestionQuery) (response.GetListQuestionResponse, *commonHttp.Meta, error)
	UpdateQuestion(ctx context.Context, questionID uuid.UUID, data request.UpdateQuestionRequest) error
	DeleteQuestion(ctx context.Context, questionID uuid.UUID) error
	GetQuestionsByType(ctx context.Context, query request.GetQuestionsByTypeQuery) (response.QuestionsByTypeResponse, error)
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

func (s *service) CreateQuestion(ctx context.Context, data request.CreateQuestionRequest) error {
	if err := data.Validate(); err != nil {
		return err
	}

	now := time.Now().UnixMilli()

	var optionsJSON *string
	if data.QuestionType == "multiple_choice" && len(data.Options) > 0 {
		optionsBytes, err := json.Marshal(data.Options)
		if err != nil {
			return err
		}
		optionsStr := string(optionsBytes)
		optionsJSON = &optionsStr
	}

	question := repository.Question{
		ID:              uuid.New(),
		Question:        data.Question,
		QuestionType:    data.QuestionType,
		Options:         optionsJSON,
		CorrectAnswer:   data.CorrectAnswer,
		SchoolID:        data.SchoolID,
		SubjectID:       data.SubjectID,
		DifficultyLevel: data.DifficultyLevel,
		Points:          data.Points,
		CreatedAt:       now,
		UpdatedAt:       0,
	}

	err := s.repository.CreateQuestion(ctx, question)
	if err != nil {
		log.Err(err).Msg("Failed to create question")
		return err
	}

	return nil
}

func (s *service) GetDetailQuestion(ctx context.Context, questionID uuid.UUID) (response.DetailQuestionResponse, error) {
	question, err := s.repository.GetQuestionByID(ctx, questionID)
	if err != nil {
		log.Err(err).Msg("Failed to get question")
		return response.DetailQuestionResponse{}, err
	}

	var options []response.QuestionOptionResponse
	if question.QuestionType == "multiple_choice" && question.Options != nil {
		var optionRequests []request.QuestionOptionRequest
		if err := json.Unmarshal([]byte(*question.Options), &optionRequests); err == nil {
			for _, option := range optionRequests {
				options = append(options, response.QuestionOptionResponse{
					ID:   option.ID,
					Text: option.Text,
				})
			}
		}
	}

	res := response.DetailQuestionResponse{
		ID:              question.ID,
		Question:        question.Question,
		QuestionType:    question.QuestionType,
		Options:         options,
		CorrectAnswer:   question.CorrectAnswer,
		SchoolID:        question.SchoolID,
		SubjectID:       question.SubjectID,
		SubjectName:     question.SubjectName,
		DifficultyLevel: question.DifficultyLevel,
		Points:          question.Points,
		CreatedAt:       question.CreatedAt,
		UpdatedAt:       question.UpdatedAt,
	}

	return res, nil
}

func (s *service) GetListQuestions(ctx context.Context, query request.GetListQuestionQuery) (response.GetListQuestionResponse, *commonHttp.Meta, error) {
	questions, total, err := s.repository.GetListQuestions(ctx, query)
	if err != nil {
		log.Err(err).Msg("Failed to get question list")
		return response.GetListQuestionResponse{}, nil, err
	}

	var res response.GetListQuestionResponse
	for _, question := range questions {
		var options []response.QuestionOptionResponse
		if question.QuestionType == "multiple_choice" && question.Options != nil {
			var optionRequests []request.QuestionOptionRequest
			if err := json.Unmarshal([]byte(*question.Options), &optionRequests); err == nil {
				for _, option := range optionRequests {
					options = append(options, response.QuestionOptionResponse{
						ID:   option.ID,
						Text: option.Text,
					})
				}
			}
		}

		res = append(res, response.QuestionResponse{
			ID:              question.ID,
			Question:        question.Question,
			QuestionType:    question.QuestionType,
			Options:         options,
			CorrectAnswer:   question.CorrectAnswer,
			SchoolID:        question.SchoolID,
			SubjectID:       question.SubjectID,
			SubjectName:     question.SubjectName,
			DifficultyLevel: question.DifficultyLevel,
			Points:          question.Points,
			CreatedAt:       question.CreatedAt,
			UpdatedAt:       question.UpdatedAt,
		})
	}

	meta := commonHttp.NewMetaFromQuery(query, total)
	return res, meta, nil
}

func (s *service) UpdateQuestion(ctx context.Context, questionID uuid.UUID, data request.UpdateQuestionRequest) error {
	if err := data.Validate(); err != nil {
		return err
	}

	now := time.Now().UnixMilli()

	var optionsJSON *string
	if data.QuestionType == "multiple_choice" && len(data.Options) > 0 {
		optionsBytes, err := json.Marshal(data.Options)
		if err != nil {
			return err
		}
		optionsStr := string(optionsBytes)
		optionsJSON = &optionsStr
	}

	question := repository.Question{
		Question:        data.Question,
		QuestionType:    data.QuestionType,
		Options:         optionsJSON,
		CorrectAnswer:   data.CorrectAnswer,
		SubjectID:       data.SubjectID,
		DifficultyLevel: data.DifficultyLevel,
		Points:          data.Points,
		UpdatedAt:       now,
	}

	err := s.repository.UpdateQuestion(ctx, questionID, question)
	if err != nil {
		log.Err(err).Msg("Failed to update question")
		return err
	}

	return nil
}

func (s *service) DeleteQuestion(ctx context.Context, questionID uuid.UUID) error {
	err := s.repository.DeleteQuestion(ctx, questionID)
	if err != nil {
		log.Err(err).Msg("Failed to delete question")
		return err
	}
	return nil
}

func (s *service) GetQuestionsByType(ctx context.Context, query request.GetQuestionsByTypeQuery) (response.QuestionsByTypeResponse, error) {
	questions, err := s.repository.GetQuestionsByType(ctx, query.SchoolID, query.SubjectID, query.QuestionType)
	if err != nil {
		log.Err(err).Msg("Failed to get questions by type")
		return response.QuestionsByTypeResponse{}, err
	}

	var res response.QuestionsByTypeResponse
	for _, question := range questions {
		var options []response.QuestionOptionResponse
		if question.QuestionType == "multiple_choice" && question.Options != nil {
			var optionRequests []request.QuestionOptionRequest
			if err := json.Unmarshal([]byte(*question.Options), &optionRequests); err == nil {
				for _, option := range optionRequests {
					options = append(options, response.QuestionOptionResponse{
						ID:   option.ID,
						Text: option.Text,
					})
				}
			}
		}

		res = append(res, response.QuestionResponse{
			ID:              question.ID,
			Question:        question.Question,
			QuestionType:    question.QuestionType,
			Options:         options,
			CorrectAnswer:   question.CorrectAnswer,
			SchoolID:        question.SchoolID,
			SubjectID:       question.SubjectID,
			SubjectName:     question.SubjectName,
			DifficultyLevel: question.DifficultyLevel,
			Points:          question.Points,
			CreatedAt:       question.CreatedAt,
			UpdatedAt:       question.UpdatedAt,
		})
	}

	return res, nil
}
