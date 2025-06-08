package service

import (
	"context"
	"database/sql"
	"enuma-elish/internal/auth/repository"
	"enuma-elish/internal/auth/service/data/request"
	"enuma-elish/internal/auth/service/data/response"
	commonError "enuma-elish/pkg/error"
	"enuma-elish/pkg/jwt"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
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

	userSchoolRole, err := s.repository.GetFirstUserSchoolRoleByUserID(ctx, user.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Err(err).Str("email", data.Email).Msg("failed to get user first user")
			return nil, err
		}

		if user.Role != repository.UserRoleAdmin {
			return nil, commonError.ErrUserNotFound
		}

		userSchoolRole = &repository.UserSchoolRole{}
	}

	now := time.Now()
	exp := now.Add(time.Hour * 2).Unix()
	nbf := now.Unix()
	iat := now.Unix()
	payload := jwt.Payload{
		Exp: exp,
		Iat: iat,
		Nbf: nbf,
		Iss: "genesis",
		Sub: user.ID.String(),
		Aud: "genesis",
		User: jwt.User{
			ID:         user.ID,
			Email:      user.Email,
			SchoolID:   userSchoolRole.SchoolID,
			SchoolRole: userSchoolRole.RoleID,
			UserRole:   user.Role,
		},
	}

	accessToken, err := jwt.GenerateToken(payload, s.config.JWT.Secret)
	if err != nil {
		log.Err(err).Str("email", data.Email).Msg("Error generating token")
		return nil, err
	}

	refreshPayload := payload
	refreshPayload.Exp = now.Add(time.Hour * 240).Unix()
	refreshToken, err := jwt.GenerateToken(refreshPayload, s.config.JWT.Secret)
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

func (s *service) RefreshToken(ctx context.Context, data request.RefreshTokenRequest) (*response.LoginResponse, error) {
	token, err := jwt.Verify(data.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		log.Err(err).Str("refresh_token", data.RefreshToken).Msg("Refresh token invalid")
		return nil, err
	}

	payload, err := jwt.ExtractToken(token)
	if err != nil {
		log.Err(err).Str("refresh_token", data.RefreshToken).Msg("Refresh token invalid")
		return nil, err
	}

	now := time.Now()
	exp := now.Add(time.Hour * 24).Unix()
	nbf := now.Unix()
	iat := now.Unix()

	payload.Exp = exp
	payload.Iat = iat
	payload.Nbf = nbf

	accessToken, err := jwt.GenerateToken(*payload, s.config.JWT.Secret)
	if err != nil {
		log.Err(err).Msg("Error generating token")
		return nil, err
	}

	refreshPayload := payload
	refreshPayload.Exp = now.Add(time.Hour * 240).Unix()
	refreshToken, err := jwt.GenerateToken(*refreshPayload, s.config.JWT.Secret)
	if err != nil {
		log.Err(err).Msg("Error generating refresh token")
		return nil, err
	}

	res := &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}
