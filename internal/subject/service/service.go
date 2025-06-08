package service

import (
	"context"
	"enuma-elish/config"
	"enuma-elish/internal/subject/repository"
	"enuma-elish/internal/subject/service/data/request"
	"enuma-elish/internal/subject/service/data/response"
	commonError "enuma-elish/pkg/error"
	commonHttp "enuma-elish/pkg/http"
	"enuma-elish/pkg/jwt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreateSubject(ctx context.Context, data request.CreateSubjectRequest) error
	GetDetailSubject(ctx context.Context, subjectID uuid.UUID) (response.DetailSubject, error)
	GetListSubject(ctx context.Context, httpQuery request.GetListSubjectQuery) (response.ListSubject, *commonHttp.Meta, error)
	UpdateSubject(ctx context.Context, subjectID uuid.UUID, data request.UpdateSubjectRequest) error
	DeleteSubject(ctx context.Context, subjectID uuid.UUID) error
	AssignTeachersToSubject(ctx context.Context, data request.AssignTeacherToSubjectRequest) error
	GetTeachersBySubject(ctx context.Context, subjectID uuid.UUID, query request.GetTeachersBySubjectQuery) (response.GetTeachersBySubjectResponse, *commonHttp.Meta, error)
	UpdateSubjectClass(ctx context.Context, data request.UpdateSubjectClassRequest) error
}

type service struct {
	repository repository.Repository
	config     *config.Config
}

func New(repository repository.Repository, config *config.Config) Service {
	return &service{
		repository: repository,
		config:     config,
	}
}

func (s *service) CreateSubject(ctx context.Context, data request.CreateSubjectRequest) error {
	jwtClaim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return commonError.ErrInvalidToken
	}
	now := time.Now().UnixMilli()
	subject := repository.Subject{
		ID:        uuid.New(),
		SchoolID:  data.SchoolID,
		Name:      data.Name,
		CreatedAt: now,
		CreatedBy: jwtClaim.User.ID,
		UpdatedAt: 0,
	}

	err = s.repository.CreateSubject(ctx, subject)
	if err != nil {
		log.Err(err).Msg("Failed to create subject")
		return commonError.ErrInternal
	}

	return nil
}

func (s *service) GetDetailSubject(ctx context.Context, subjectID uuid.UUID) (response.DetailSubject, error) {
	subject, err := s.repository.GetSubjectByID(ctx, subjectID)
	if err != nil {
		log.Err(err).Msg("Failed to get subject detail")
		return response.DetailSubject{}, commonError.ErrInternal
	}

	return response.DetailSubject{
		ID:        subject.ID,
		SchoolID:  subject.SchoolID,
		Name:      subject.Name,
		CreatedAt: subject.CreatedAt,
		UpdatedAt: subject.UpdatedAt,
	}, nil
}

func (s *service) GetListSubject(ctx context.Context, httpQuery request.GetListSubjectQuery) (response.ListSubject, *commonHttp.Meta, error) {
	subjects, total, err := s.repository.GetListSubjects(ctx, httpQuery)
	if err != nil {
		log.Err(err).Msg("Failed to get list subjects")
		return nil, nil, commonError.ErrInternal
	}

	meta := commonHttp.NewMetaFromQuery(httpQuery, total)
	res := make(response.ListSubject, len(subjects))
	for i, subject := range subjects {
		res[i] = response.DetailSubject{
			ID:        subject.ID,
			SchoolID:  subject.SchoolID,
			Name:      subject.Name,
			CreatedAt: subject.CreatedAt,
			UpdatedAt: subject.UpdatedAt,
		}
	}

	return res, meta, nil
}

func (s *service) UpdateSubject(ctx context.Context, subjectID uuid.UUID, data request.UpdateSubjectRequest) error {
	subject, err := s.repository.GetSubjectByID(ctx, subjectID)
	if err != nil {
		log.Err(err).Msg("Failed to get subject")
		return commonError.ErrInternal
	}

	subject.Name = data.Name
	subject.UpdatedAt = time.Now().UnixMilli()

	err = s.repository.UpdateSubject(ctx, *subject)
	if err != nil {
		log.Err(err).Msg("Failed to update subject")
		return commonError.ErrInternal
	}

	return nil
}

func (s *service) DeleteSubject(ctx context.Context, subjectID uuid.UUID) error {
	err := s.repository.DeleteSubject(ctx, subjectID)
	if err != nil {
		log.Err(err).Msg("Failed to delete subject")
		return commonError.ErrInternal
	}

	return nil
}

func (s *service) AssignTeachersToSubject(ctx context.Context, data request.AssignTeacherToSubjectRequest) error {
	err := s.repository.AssignTeachersToSubject(ctx, data.SubjectID, data.TeacherIDs)
	if err != nil {
		log.Err(err).Msg("Failed to assign teachers to subject")
		return err
	}
	return nil
}

func (s *service) GetTeachersBySubject(ctx context.Context, subjectID uuid.UUID, query request.GetTeachersBySubjectQuery) (response.GetTeachersBySubjectResponse, *commonHttp.Meta, error) {
	teachers, total, err := s.repository.GetTeachersBySubjectID(ctx, subjectID, query)
	if err != nil {
		log.Err(err).Msg("Failed to get teachers by subject")
		return response.GetTeachersBySubjectResponse{}, nil, err
	}

	var res response.GetTeachersBySubjectResponse
	for _, teacher := range teachers {
		res = append(res, response.TeacherInSubject{
			ID:         teacher.ID,
			Name:       teacher.Name,
			Email:      teacher.Email,
			IsVerified: teacher.IsVerified,
			CreatedAt:  teacher.CreatedAt,
			UpdatedAt:  teacher.UpdatedAt,
		})
	}

	meta := commonHttp.NewMetaFromQuery(query, total)
	return res, meta, nil
}

func (s *service) UpdateSubjectClass(ctx context.Context, data request.UpdateSubjectClassRequest) error {
	err := s.repository.UpdateSubjectClass(ctx, data.SubjectID, data.OldClassID, data.NewClassID)
	if err != nil {
		log.Err(err).Msg("Failed to update subject class assignment")
		return err
	}
	return nil
}
