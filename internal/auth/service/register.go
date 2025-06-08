package service

import (
	"context"
	"database/sql"
	"enuma-elish/internal/auth/repository"
	"enuma-elish/internal/auth/service/data/request"
	commonError "enuma-elish/pkg/error"
	"errors"
	"fmt"
	"net/smtp"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Register(ctx context.Context, data request.Register) error {
	u, err := s.repository.GetUserByEmail(ctx, data.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Err(err).Msg("Failed to get user by email")
		return err

	}

	if u != nil {
		log.Err(err).Msg("Email already exists")
		return commonError.ErrEmailAlreadyExists
	}

	hashPass, err := hashPassword(data.Password)
	if err != nil {
		log.Err(err).Msg("Failed to hash password")
		return err
	}

	user := repository.UserVerifyEmailToken{
		Token:    uuid.New().String(),
		Email:    data.Email,
		Password: hashPass,
		Name:     data.Name,
	}

	err = s.repository.CreateVerifyEmailToken(ctx, &user)
	if err != nil {
		log.Err(err).Msg("Failed to create verify email token")
		return err
	}

	go func() {
		verifyUrl := fmt.Sprintf("%s/email-verification?token=%s&email=%s&type=registration", s.config.Http.FrontendHost, user.Token, data.Email)
		err = s.sendEmail(user.Email, verifyUrl, "verify email")
		if err != nil {
			log.Err(err).Msg("Failed to send email")
		}
	}()

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *service) sendEmail(to string, msg, subject string) error {
	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)
	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, msg))

	addr := fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port)
	err := smtp.SendMail(addr, auth, s.config.SMTP.Username, []string{to}, message)
	if err != nil {
		log.Err(err).Msg("Failed to send email")
		return err
	}

	return nil
}

func (s *service) VerifyEmail(ctx context.Context, data request.VerifyEmailRequest) error {
	u, err := s.repository.VerifyEmailToken(ctx, data.Email)
	if err != nil {
		log.Err(err).Msg("Failed to verify email token")
		return commonError.ErrInvalidToken
	}

	if u.Token != data.Token {
		log.Err(commonError.ErrInvalidToken).Msg("Failed to verify email token")
		return commonError.ErrInvalidToken
	}

	user := repository.User{
		Email:      u.Email,
		Password:   u.Password,
		Name:       u.Name,
		IsVerified: true,
	}

	err = s.repository.CreateUser(ctx, &user)
	if err != nil {
		log.Err(err).Msg("Failed to create user")
		return err
	}

	key := repository.VerifyEmailTokenKey + ":" + u.Email
	s.repository.Redis().Del(ctx, key)

	return nil
}
