package service

import (
	"context"
	"enuma-elish/internal/auth/repository"
	"enuma-elish/internal/auth/service/data/request"
	commonError "enuma-elish/pkg/error"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s *service) ForgotPassword(ctx context.Context, data request.ForgotPasswordRequest) error {
	u, err := s.repository.GetUserByEmail(ctx, data.Email)
	if err != nil {
		log.Error().Err(err).Str("email", data.Email).Msg("user not found")
		return err
	}

	token := repository.UserForgotPasswordToken{
		Email: u.Email,
		Token: uuid.New().String(),
	}

	err = s.repository.CreateForgotPasswordToken(ctx, &token)
	if err != nil {
		log.Error().Err(err).Str("email", data.Email).Str("token", token.Token).Msg("failed to create token")
		return err
	}

	go func() {
		verifyUrl := fmt.Sprintf("%s/forgot-password?token=%s&email=%s", s.config.Http.FrontendHost, token.Token, data.Email)
		err = s.sendEmail(u.Email, verifyUrl, "verify email")
		if err != nil {
			log.Err(err).Msg("Failed to send email")
		}
	}()

	return nil
}

func (s *service) ForgotPasswordVerify(ctx context.Context, data request.ForgotPasswordVerifyRequest) error {
	_, err := s.repository.GetUserByEmail(ctx, data.Email)
	if err != nil {
		log.Error().Err(err).Str("email", data.Email).Msg("user not found")
		return err
	}

	token, err := s.repository.VerifyForgotPasswordToken(ctx, data.Email)
	if err != nil {
		log.Error().Err(err).Str("token", token.Token).Msg("failed to verify token")
		return err
	}

	if token.Token != data.Token {
		log.Error().Str("token", token.Token).Str("email", data.Email).Msg("invalid token")
		return commonError.ErrInvalidToken
	}

	pass, err := hashPassword(data.NewPassword)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")
		return err
	}

	err = s.repository.UpdatePassword(ctx, data.Email, pass)
	if err != nil {
		log.Error().Err(err).Str("email", data.Email).Msg("failed to update password")
		return err
	}

	s.repository.Redis().Del(ctx, repository.ForgotPasswordTokenKey+":"+data.Email)

	return nil
}
