package service

import (
	"context"
	"database/sql"
	"enuma-elish/internal/auth/service/data/request"
	"enuma-elish/internal/auth/service/data/response"
	commonError "enuma-elish/pkg/error"
	"enuma-elish/pkg/jwt"
	"errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s *service) Login(ctx context.Context, data request.LoginRequest) (*response.LoginResponse, error) {
	user, err := s.repository.GetUserByEmail(ctx, data.Email)
	if err != nil {
		log.Err(err).Str("email", data.Email).Msg("User not found")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, commonError.ErrUserNotFound
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		log.Err(err).Str("email", data.Email).Msg("password incorrect")
		return nil, commonError.ErrInvalidPassword
	}

	jwtPayload := map[string]interface{}{
		"email": user.Email,
		"id":    user.ID,
	}

	accessToken, err := jwt.GenerateToken(time.Duration(s.config.JWT.Duration)*time.Hour, jwtPayload, s.config.JWT.Secret)
	if err != nil {
		log.Err(err).Str("email", data.Email).Msg("Error generating token")
		return nil, err
	}

	refreshToken, err := jwt.GenerateToken(time.Duration(s.config.JWT.Duration)*time.Hour, jwtPayload, s.config.JWT.Secret)
	if err != nil {
		log.Err(err).Str("email", data.Email).Msg("Error generating refresh token")
		return nil, err
	}

	res := &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}
