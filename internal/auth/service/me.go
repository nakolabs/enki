package service

import (
	"context"
	"enuma-elish/internal/auth/service/data/response"
	"enuma-elish/pkg/jwt"
	"github.com/rs/zerolog/log"
)

func (s *service) Me(ctx context.Context) (*response.ProfileResponse, error) {
	claims, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to extract claims")
		return nil, err
	}

	u, err := s.repository.GetUserByID(ctx, claims.User.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user")
		return nil, err
	}

	result := response.ProfileResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	return &result, nil
}
