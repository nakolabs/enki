package service

import (
	"context"
	"enuma-elish/config"
	"enuma-elish/internal/teacher/repository"
	"enuma-elish/internal/teacher/service/data/request"
	"enuma-elish/internal/teacher/service/data/response"
	commonHttp "enuma-elish/pkg/http"

	"github.com/google/uuid"
)

type Service interface {
	InviteTeacher(ctx context.Context, data request.InviteTeacherRequest) error
	VerifyTeacherEmail(ctx context.Context, data request.VerifyTeacherEmailRequest) error
	UpdateTeacherAfterVerifyEmail(ctx context.Context, data request.UpdateTeacherAfterVerifyEmailRequest) error
	ListTeachers(ctx context.Context, httpQuery request.GetListTeacherQuery) (response.GetListTeacherResponse, *commonHttp.Meta, error)
	DeleteTeacher(ctx context.Context, teacherID uuid.UUID, schoolID uuid.UUID) error
	GetDetailTeacher(ctx context.Context, teacherID uuid.UUID) (response.GetDetailTeacherResponse, error)
	UpdateTeacherClass(ctx context.Context, data request.UpdateTeacherClassRequest) error
	GetTeacherSubjects(ctx context.Context, teacherID uuid.UUID) ([]response.Subject, error)
}

type service struct {
	config     *config.Config
	repository repository.Repository
}

func New(c *config.Config, r repository.Repository) Service {
	return &service{
		config:     c,
		repository: r,
	}
}
