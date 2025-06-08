package service

import (
	"context"
	"enuma-elish/config"
	"enuma-elish/internal/student/repository"
	"enuma-elish/internal/student/service/data/request"
	"enuma-elish/internal/student/service/data/response"
	commonError "enuma-elish/pkg/error"
	commonHttp "enuma-elish/pkg/http"
	"enuma-elish/pkg/jwt"
	"fmt"
	"net/smtp"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	InviteStudent(ctx context.Context, data request.InviteStudentRequest) error
	VerifyStudentEmail(ctx context.Context, data request.VerifyStudentEmailRequest) error
	UpdateStudentAfterVerifyEmail(ctx context.Context, data request.UpdateStudentAfterVerifyEmailRequest) error
	GetListStudent(ctx context.Context, httpQuery request.GetListStudentQuery) (response.GetListStudentResponse, *commonHttp.Meta, error)
	GetDetailStudent(ctx context.Context, studentID uuid.UUID) (response.GetDetailStudentResponse, error)
	DeleteStudent(ctx context.Context, studentID uuid.UUID, schoolID uuid.UUID) error
	UpdateStudentClass(ctx context.Context, data request.UpdateStudentClassRequest) error
}

type service struct {
	repository repository.Repository
	config     *config.Config
}

func New(repository repository.Repository, config *config.Config) Service {
	return &service{repository, config}
}

func (s *service) InviteStudent(ctx context.Context, data request.InviteStudentRequest) error {

	claim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to extract claims")
		return err
	}

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
			CreatedBy:  claim.User.ID,
		})
	}

	err = s.repository.CreateStudent(ctx, data.SchoolID, students)
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

			link := fmt.Sprintf("%s/email-verification?email=%s&token=%s&type=invite_student", s.config.Http.FrontendHost, v.Email, token)
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

func (s *service) UpdateStudentAfterVerifyEmail(ctx context.Context, data request.UpdateStudentAfterVerifyEmailRequest) error {

	token, err := s.repository.VerifyEmailToken(ctx, data.Email)
	if err != nil {
		log.Err(err).Msg("Failed to verify email token")
		return commonError.ErrInvalidToken
	}

	if token != data.Token {
		log.Err(commonError.ErrInvalidToken).Msg("Failed to verify email token")
		return commonError.ErrInvalidToken
	}

	student, err := s.repository.GetStudentByEmail(ctx, data.Email)
	if err != nil {
		log.Err(err).Msg("Failed to get student")
		return commonError.ErrUserNotFound
	}

	hashPass, err := hashPassword(data.Password)
	if err != nil {
		log.Err(err).Msg("Failed to hash password")
		return err
	}

	student = &repository.User{
		ID:         student.ID,
		Name:       data.Name,
		Email:      data.Email,
		Password:   hashPass,
		IsVerified: true,
		UpdatedAt:  time.Now().UnixMilli(),
	}

	err = s.repository.UpdateStudent(ctx, *student)
	if err != nil {
		log.Err(err).Msg("Failed to update student")
		return commonError.ErrInternal
	}

	s.repository.Redis().Del(ctx, repository.StudentVerifyEmailTokenKey+":"+data.Email)
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *service) VerifyStudentEmail(ctx context.Context, data request.VerifyStudentEmailRequest) error {
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

func (s *service) GetListStudent(ctx context.Context, httpQuery request.GetListStudentQuery) (response.GetListStudentResponse, *commonHttp.Meta, error) {
	data, total, err := s.repository.GetListStudent(ctx, httpQuery)
	if err != nil {
		log.Err(err).Msg("Failed to get students")
		return response.GetListStudentResponse{}, nil, nil
	}

	res := response.GetListStudentResponse{}
	for _, student := range data {
		res = append(res, response.GetStudentResponse{
			ID:         student.ID,
			Name:       student.Name,
			Email:      student.Email,
			IsVerified: student.IsVerified,
			CreateAt:   student.CreatedAt,
			UpdateAt:   student.UpdatedAt,
		})
	}

	meta := commonHttp.NewMetaFromQuery(httpQuery, total)

	return res, meta, nil
}

func (s *service) DeleteStudent(ctx context.Context, studentID uuid.UUID, schoolID uuid.UUID) error {
	err := s.repository.DeleteStudent(ctx, studentID, schoolID)
	if err != nil {
		log.Err(err).Msg("Failed to delete student")
		return err
	}
	return nil
}

func (s *service) GetDetailStudent(ctx context.Context, studentID uuid.UUID) (response.GetDetailStudentResponse, error) {
	student, err := s.repository.GetStudentByID(ctx, studentID)
	if err != nil {
		log.Err(err).Msg("Failed to get student")
		return response.GetDetailStudentResponse{}, err
	}

	res := response.GetDetailStudentResponse{
		ID:         student.ID,
		Name:       student.Name,
		Email:      student.Email,
		IsVerified: student.IsVerified,
		CreateAt:   student.CreatedAt,
		UpdateAt:   student.UpdatedAt,
	}

	return res, nil
}

func (s *service) UpdateStudentClass(ctx context.Context, data request.UpdateStudentClassRequest) error {
	err := s.repository.UpdateStudentClass(ctx, data.StudentID, data.OldClassID, data.NewClassID)
	if err != nil {
		log.Err(err).Msg("Failed to update student class assignment")
		return err
	}
	return nil
}
