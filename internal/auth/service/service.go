package service

import (
	"context"
	"database/sql"
	"enuma-elish/config"
	"enuma-elish/internal/auth/repository"
	"enuma-elish/internal/auth/service/data/request"
	"enuma-elish/internal/auth/service/data/response"
	commonError "enuma-elish/pkg/error"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, data request.Register) error
	Login(ctx context.Context, data request.LoginRequest) (*response.LoginResponse, error)
	VerifyEmail(ctx context.Context, data request.VerifyEmailRequest) error
	Me(ctx context.Context) (*response.UserResponse, error)
	ForgotPassword(ctx context.Context, data request.ForgotPasswordRequest) error
	ForgotPasswordVerify(ctx context.Context, data request.ForgotPasswordVerifyRequest) error
	RefreshToken(ctx context.Context, data request.RefreshTokenRequest) (*response.LoginResponse, error)
	UpdateUser(ctx context.Context, data request.UpdateUserRequest) (*response.UserResponse, error)
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

func (s *service) UpdateUser(ctx context.Context, data request.UpdateUserRequest) (*response.UserResponse, error) {
	// Get user ID from middleware context
	userIDFromCtx := ctx.Value("user_id")
	if userIDFromCtx == nil {
		return nil, fmt.Errorf("user ID not found in context")
	}

	userID, ok := userIDFromCtx.(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("invalid user ID format")
	}

	// Try to get existing profile, if not found create a new one
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, commonError.ErrUserNotFound
		}
		return nil, err
	}

	if data.Email != "" {
		user.Email = data.Email
	}
	if data.Name != "" {
		user.Name = data.Name
	}
	if data.Phone != "" {
		user.Phone = data.Phone
	}
	if data.DateOfBirth != "" {
		user.DateOfBirth = data.DateOfBirth
	}
	if data.Gender != "" {
		user.Gender = data.Gender
	}
	if data.Address != "" {
		user.Address = data.Address
	}
	if data.City != "" {
		user.City = data.City
	}
	if data.Country != "" {
		user.Country = data.Country
	}
	if data.Avatar != "" {
		user.Avatar = data.Avatar
	}
	if data.Bio != "" {
		user.Bio = data.Bio
	}
	if data.ParentName != "" {
		user.ParentName = data.ParentName
	}
	if data.ParentPhone != "" {
		user.ParentPhone = data.ParentPhone
	}
	if data.ParentEmail != "" {
		user.ParentEmail = data.ParentEmail
	}
	user.UpdatedAt = time.Now().UnixMilli()

	// Update/Create profile in database
	updatedProfile, err := s.repository.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return &response.UserResponse{
		ID:          updatedProfile.ID,
		Email:       updatedProfile.Email,
		Name:        updatedProfile.Name,
		Phone:       updatedProfile.Phone,
		DateOfBirth: updatedProfile.DateOfBirth,
		Gender:      updatedProfile.Gender,
		Address:     updatedProfile.Address,
		City:        updatedProfile.City,
		Country:     updatedProfile.Country,
		Avatar:      updatedProfile.Avatar,
		Bio:         updatedProfile.Bio,
		ParentName:  updatedProfile.ParentName,
		ParentPhone: updatedProfile.ParentPhone,
		ParentEmail: updatedProfile.ParentEmail,
		CreatedAt:   updatedProfile.CreatedAt,
		UpdatedAt:   updatedProfile.UpdatedAt,
		DeletedAt:   updatedProfile.DeletedAt,
		DeletedBy:   updatedProfile.DeletedBy.String,
	}, nil
}

func (s *service) Me(ctx context.Context) (*response.UserResponse, error) {
	// Get user ID from middleware context
	userIDFromCtx := ctx.Value("user_id")
	if userIDFromCtx == nil {
		return nil, fmt.Errorf("user ID not found in context")
	}

	userID, ok := userIDFromCtx.(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("invalid user ID format")
	}

	// Get user data
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &response.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Phone:       user.Phone,
		DateOfBirth: user.DateOfBirth,
		Gender:      user.Gender,
		Address:     user.Address,
		City:        user.City,
		Country:     user.Country,
		Avatar:      user.Avatar,
		Bio:         user.Bio,
		ParentName:  user.ParentName,
		ParentPhone: user.ParentPhone,
		ParentEmail: user.ParentEmail,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		DeletedAt:   user.DeletedAt,
		DeletedBy:   user.DeletedBy.String,
	}, nil

}
