package service

import (
	"context"
	"database/sql"
	"enuma-elish/config"
	"enuma-elish/internal/class/repository"
	"enuma-elish/internal/class/service/data/request"
	"enuma-elish/internal/class/service/data/response"
	commonError "enuma-elish/pkg/error"
	commonHttp "enuma-elish/pkg/http"
	"enuma-elish/pkg/jwt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreateClass(ctx context.Context, data request.CreateClassRequest) error
	GetDetailClass(ctx context.Context, classID uuid.UUID) (response.DetailClass, error)
	GetListClass(ctx context.Context, httpQuery request.GetListClassQuery) (response.ListClass, *commonHttp.Meta, error)
	UpdateClass(ctx context.Context, classID uuid.UUID, data request.UpdateClassRequest) error
	DeleteClass(ctx context.Context, classID uuid.UUID) error
	AddStudentsToClass(ctx context.Context, data request.AddStudentToClassRequest) error
	GetStudentsByClass(ctx context.Context, classID uuid.UUID, query request.GetStudentsByClassQuery) (response.GetStudentsByClassResponse, *commonHttp.Meta, error)
	AddTeachersToClass(ctx context.Context, data request.AddTeacherToClassRequest) error
	GetTeachersByClass(ctx context.Context, classID uuid.UUID, query request.GetTeachersByClassQuery) (response.GetTeachersByClassResponse, *commonHttp.Meta, error)
	AddSubjectsToClass(ctx context.Context, data request.AddSubjectToClassRequest) error
	GetSubjectsByClass(ctx context.Context, classID uuid.UUID, query request.GetSubjectsByClassQuery) (response.GetSubjectsByClassResponse, *commonHttp.Meta, error)
	RemoveTeachersFromClass(ctx context.Context, data request.RemoveTeacherFromClassRequest) error
	RemoveStudentsFromClass(ctx context.Context, data request.RemoveStudentFromClassRequest) error
	RemoveSubjectsFromClass(ctx context.Context, data request.RemoveSubjectFromClassRequest) error
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

func (s *service) CreateClass(ctx context.Context, data request.CreateClassRequest) error {

	jwtClaim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return commonError.ErrInvalidToken
	}

	now := time.Now().UnixMilli()
	class := repository.Class{
		ID:        uuid.New(),
		SchoolID:  data.SchoolID,
		Name:      data.Name,
		CreatedAt: now,
		CreatedBy: jwtClaim.User.ID,
		UpdatedAt: 0,
	}

	err = s.repository.CreateClass(ctx, class)
	if err != nil {
		log.Err(err).Msg("Failed to create class")
		return commonError.ErrInternal
	}

	return nil
}

func (s *service) GetDetailClass(ctx context.Context, classID uuid.UUID) (response.DetailClass, error) {
	class, err := s.repository.GetClassByID(ctx, classID)
	if err != nil {
		log.Err(err).Msg("Failed to get class detail")
		return response.DetailClass{}, commonError.ErrInternal
	}

	return response.DetailClass{
		ID:        class.ID,
		SchoolID:  class.SchoolID,
		Name:      class.Name,
		CreatedAt: class.CreatedAt,
		UpdatedAt: class.UpdatedAt,
	}, nil
}

func (s *service) GetListClass(ctx context.Context, httpQuery request.GetListClassQuery) (response.ListClass, *commonHttp.Meta, error) {
	classes, total, err := s.repository.GetListClasses(ctx, httpQuery)
	if err != nil {
		log.Err(err).Msg("Failed to get list classes")
		return nil, nil, commonError.ErrInternal
	}

	meta := commonHttp.NewMetaFromQuery(httpQuery, total)
	res := make(response.ListClass, len(classes))
	for i, class := range classes {
		res[i] = response.DetailClass{
			ID:        class.ID,
			SchoolID:  class.SchoolID,
			Name:      class.Name,
			CreatedAt: class.CreatedAt,
			UpdatedAt: class.UpdatedAt,
		}
	}

	return res, meta, nil
}

func (s *service) UpdateClass(ctx context.Context, classID uuid.UUID, data request.UpdateClassRequest) error {

	jwtClaim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return commonError.ErrInvalidToken
	}

	class, err := s.repository.GetClassByID(ctx, classID)
	if err != nil {
		log.Err(err).Msg("Failed to get class")
		return commonError.ErrInternal
	}

	class.Name = data.Name
	class.UpdatedAt = time.Now().UnixMilli()
	class.UpdatedBy = sql.NullString{String: jwtClaim.User.ID.String(), Valid: true}

	err = s.repository.UpdateClass(ctx, *class)
	if err != nil {
		log.Err(err).Msg("Failed to update class")
		return commonError.ErrInternal
	}

	return nil
}

func (s *service) DeleteClass(ctx context.Context, classID uuid.UUID) error {
	err := s.repository.DeleteClass(ctx, classID)
	if err != nil {
		log.Err(err).Msg("Failed to delete class")
		return commonError.ErrInternal
	}

	return nil
}

func (s *service) AddStudentsToClass(ctx context.Context, data request.AddStudentToClassRequest) error {
	err := s.repository.AddStudentsToClass(ctx, data.ClassID, data.StudentIDs)
	if err != nil {
		log.Err(err).Msg("Failed to add students to class")
		return err
	}
	return nil
}

func (s *service) GetStudentsByClass(ctx context.Context, classID uuid.UUID, query request.GetStudentsByClassQuery) (response.GetStudentsByClassResponse, *commonHttp.Meta, error) {
	students, total, err := s.repository.GetStudentsByClassID(ctx, classID, query)
	if err != nil {
		log.Err(err).Msg("Failed to get students by class")
		return response.GetStudentsByClassResponse{}, nil, err
	}

	var res response.GetStudentsByClassResponse
	for _, student := range students {
		res = append(res, response.StudentInClass{
			ID:         student.ID,
			Name:       student.Name,
			Email:      student.Email,
			IsVerified: student.IsVerified,
			CreatedAt:  student.CreatedAt,
			UpdatedAt:  student.UpdatedAt,
		})
	}

	meta := commonHttp.NewMetaFromQuery(query, total)
	return res, meta, nil
}

func (s *service) AddTeachersToClass(ctx context.Context, data request.AddTeacherToClassRequest) error {
	err := s.repository.AddTeachersToClass(ctx, data.ClassID, data.TeacherIDs)
	if err != nil {
		log.Err(err).Msg("Failed to add teachers to class")
		return err
	}
	return nil
}

func (s *service) GetTeachersByClass(ctx context.Context, classID uuid.UUID, query request.GetTeachersByClassQuery) (response.GetTeachersByClassResponse, *commonHttp.Meta, error) {
	teachers, total, err := s.repository.GetTeachersByClassID(ctx, classID, query)
	if err != nil {
		log.Err(err).Msg("Failed to get teachers by class")
		return response.GetTeachersByClassResponse{}, nil, err
	}

	var res response.GetTeachersByClassResponse
	for _, teacher := range teachers {
		res = append(res, response.TeacherInClass{
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

func (s *service) AddSubjectsToClass(ctx context.Context, data request.AddSubjectToClassRequest) error {
	err := s.repository.AddSubjectsToClass(ctx, data.ClassID, data.SubjectIDs)
	if err != nil {
		log.Err(err).Msg("Failed to add subjects to class")
		return err
	}
	return nil
}

func (s *service) GetSubjectsByClass(ctx context.Context, classID uuid.UUID, query request.GetSubjectsByClassQuery) (response.GetSubjectsByClassResponse, *commonHttp.Meta, error) {
	subjects, total, err := s.repository.GetSubjectsByClassID(ctx, classID, query)
	if err != nil {
		log.Err(err).Msg("Failed to get subjects by class")
		return response.GetSubjectsByClassResponse{}, nil, err
	}

	var res response.GetSubjectsByClassResponse
	for _, subject := range subjects {
		res = append(res, response.SubjectInClass{
			ID:        subject.ID,
			Name:      subject.Name,
			SchoolID:  subject.SchoolID,
			CreatedAt: subject.CreatedAt,
			UpdatedAt: subject.UpdatedAt,
		})
	}

	meta := commonHttp.NewMetaFromQuery(query, total)
	return res, meta, nil
}

func (s *service) RemoveTeachersFromClass(ctx context.Context, data request.RemoveTeacherFromClassRequest) error {
	err := s.repository.RemoveTeachersFromClass(ctx, data.ClassID, data.TeacherIDs)
	if err != nil {
		log.Err(err).Msg("Failed to remove teachers from class")
		return err
	}
	return nil
}

func (s *service) RemoveStudentsFromClass(ctx context.Context, data request.RemoveStudentFromClassRequest) error {
	err := s.repository.RemoveStudentsFromClass(ctx, data.ClassID, data.StudentIDs)
	if err != nil {
		log.Err(err).Msg("Failed to remove students from class")
		return err
	}
	return nil
}

func (s *service) RemoveSubjectsFromClass(ctx context.Context, data request.RemoveSubjectFromClassRequest) error {
	err := s.repository.RemoveSubjectsFromClass(ctx, data.ClassID, data.SubjectIDs)
	if err != nil {
		log.Err(err).Msg("Failed to remove subjects from class")
		return err
	}
	return nil
}
