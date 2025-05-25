package service

import (
	"context"
	"enuma-elish/config"
	"enuma-elish/internal/auth/repository"
	"enuma-elish/internal/auth/service/data/request"
	"enuma-elish/internal/auth/service/data/response"
)

type Service interface {
	Register(ctx context.Context, data request.Register) error
	Login(ctx context.Context, data request.LoginRequest) (*response.LoginResponse, error)
	VerifyEmail(ctx context.Context, data request.VerifyEmailRequest) error
	Me(ctx context.Context) (*response.ProfileResponse, error)
	ForgotPassword(ctx context.Context, data request.ForgotPasswordRequest) error
	ForgotPasswordVerify(ctx context.Context, data request.ForgotPasswordVerifyRequest) error
	RefreshToken(ctx context.Context, data request.RefreshTokenRequest) (*response.LoginResponse, error)
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
