package service

import (
	"context"
	"enuma-elish/config"
	"enuma-elish/internal/student/repository"
	"enuma-elish/internal/student/service/data/request"
	commonError "enuma-elish/pkg/error"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"net/smtp"
	"time"
)

type Service interface {
	InviteStudent(ctx context.Context, data request.InviteStudentRequest) error
	VerifyTeacherEmail(ctx context.Context, data request.VerifyStudentEmailRequest) error
	UpdateTeacherAfterVerifyEmail(ctx context.Context, data request.UpdateStudentAfterVerifyEmailRequest) error
}

type service struct {
	repository repository.Repository
	config     *config.Config
}

func New(repository repository.Repository, config *config.Config) Service {
	return &service{repository, config}
}

func (s *service) InviteStudent(ctx context.Context, data request.InviteStudentRequest) error {

	var students []repository.User
	now := time.Now().UnixMilli()

	for _, email := range data.Emails {
		students = append(students, repository.User{
			ID:         uuid.New(),
			Email:      email,
			Password:   "",
			IsVerified: false,
			CreatedAt:  now,
			UpdatedAt:  0,
		})
	}

	err := s.repository.CreateStudent(ctx, data.SchoolID, students)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create student")
		return err
	}

	go func([]repository.User) {
		for _, v := range students {
			token, err := s.repository.CreateStudentVerifyEmailToken(context.Background(), v.Email)
			if err != nil {
				log.Err(err).Str("email", v.Email).Msg("create student verify token")
			}

			link := fmt.Sprintf("%s/email-verification?email=%s&token=%s", s.config.Http.FrontendHost, v.Email, token)
			err = s.sendEmail(v.Email, link, "Email Verification")
			if err != nil {
				log.Err(err).Str("email", v.Email).Msg("send email")
			}
		}
	}(students)

	return nil
}

func (s *service) sendEmail(to string, msg, subject string) error {
	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)
	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, msg))

	addr := fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port)
	err := smtp.SendMail(addr, auth, s.config.SMTP.Username, []string{to}, message)
	if err != nil {
		log.Err(err).Msg("failed to send email")
		return err
	}

	return nil
}

func (s *service) UpdateTeacherAfterVerifyEmail(ctx context.Context, data request.UpdateStudentAfterVerifyEmailRequest) error {

	token, err := s.repository.VerifyEmailToken(ctx, data.Email)
	if err != nil {
		log.Err(err).Msg("Failed to verify email token")
		return commonError.ErrInvalidToken
	}

	if token != data.Token {
		log.Err(commonError.ErrInvalidToken).Msg("Failed to verify email token")
		return commonError.ErrInvalidToken
	}

	teacher, err := s.repository.GetStudentByEmail(ctx, data.Email)
	if err != nil {
		log.Err(err).Msg("Failed to get teacher")
		return commonError.ErrUserNotFound
	}

	hashPass, err := hashPassword(data.Password)
	if err != nil {
		log.Err(err).Msg("Failed to hash password")
		return err
	}

	teacher = &repository.User{
		ID:         teacher.ID,
		Name:       data.Name,
		Email:      data.Email,
		Password:   hashPass,
		IsVerified: true,
		UpdatedAt:  time.Now().UnixMilli(),
	}

	err = s.repository.UpdateTeacher(ctx, *teacher)
	if err != nil {
		log.Err(err).Msg("Failed to update teacher")
		return commonError.ErrInternal
	}

	s.repository.Redis().Del(ctx, repository.StudentVerifyEmailTokenKey+":"+data.Email)
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *service) VerifyTeacherEmail(ctx context.Context, data request.VerifyStudentEmailRequest) error {
	token, err := s.repository.VerifyEmailToken(ctx, data.Email)
	if err != nil {
		log.Err(err).Msg("Failed to verify student email token")
		return commonError.ErrInvalidToken
	}

	if token != data.Token {
		log.Err(commonError.ErrInvalidToken).Msg("Failed to verify student email token")
		return commonError.ErrInvalidToken
	}

	return nil
}
