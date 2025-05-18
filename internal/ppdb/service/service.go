package service

import (
	"context"
	"enuma-elish/config"
	"enuma-elish/internal/ppdb/repository"
	"enuma-elish/internal/ppdb/service/data/request"
)

type Service interface {
	CreatSchool(ctx context.Context, data request.CreateSchoolRequest) error
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
