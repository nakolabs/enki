package service

import (
	"context"
	"enuma-elish/config"
	"enuma-elish/internal/school/repository"
	"enuma-elish/internal/school/service/data/request"
	"enuma-elish/internal/school/service/data/response"

	"github.com/google/uuid"
)

type Service interface {
	CreatSchool(ctx context.Context, data request.CreateSchoolRequest) error
	GetDetailSchool(ctx context.Context, schoolID uuid.UUID) (response.DetailSchool, error)
	GetListSchool(ctx context.Context) (response.ListSchool, error)
	SwitchSchool(ctx context.Context, schoolID uuid.UUID) (string, error)
	DeleteSchool(ctx context.Context, schoolID uuid.UUID) error
	UpdateSchoolProfile(ctx context.Context, schoolID uuid.UUID, data request.UpdateSchoolProfileRequest) (response.DetailSchool, error)
}

type service struct {
	repository repository.Repository
	config     *config.Config
}

func New(r repository.Repository, c *config.Config) Service {
	return &service{
		repository: r,
		config:     c,
	}
}
