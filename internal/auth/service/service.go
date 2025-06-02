package service

import (
	"context"
	"enuma-elish/config"
	"enuma-elish/internal/auth/repository"
	"enuma-elish/internal/auth/service/data/request"
	"enuma-elish/internal/auth/service/data/response"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, data request.Register) error
	Login(ctx context.Context, data request.LoginRequest) (*response.LoginResponse, error)
	VerifyEmail(ctx context.Context, data request.VerifyEmailRequest) error
	Me(ctx context.Context) (*response.UserResponse, error)
	Profile(ctx context.Context) (*response.UserProfileResponse, error)
	ForgotPassword(ctx context.Context, data request.ForgotPasswordRequest) error
	ForgotPasswordVerify(ctx context.Context, data request.ForgotPasswordVerifyRequest) error
	RefreshToken(ctx context.Context, data request.RefreshTokenRequest) (*response.LoginResponse, error)
	UpdateProfile(ctx context.Context, data request.UpdateProfileRequest) (*response.ProfileResponse, error)
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

func (s *service) UpdateProfile(ctx context.Context, data request.UpdateProfileRequest) (*response.ProfileResponse, error) {
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
	profile, err := s.repository.GetProfileByUserID(ctx, userID)
	if err != nil {
		// Profile doesn't exist, create new one with provided data
		profile = &repository.Profile{
			UserID:      userID,
			FirstName:   data.FirstName,
			LastName:    data.LastName,
			Phone:       data.Phone,
			DateOfBirth: data.DateOfBirth,
			Gender:      data.Gender,
			Address:     data.Address,
			City:        data.City,
			Country:     data.Country,
			Avatar:      data.Avatar,
			Bio:         data.Bio,
			ParentName:  data.ParentName,
			ParentPhone: data.ParentPhone,
			ParentEmail: data.ParentEmail,
		}
	} else {
		// Profile exists, update only provided fields
		if data.FirstName != nil {
			profile.FirstName = data.FirstName
		}
		if data.LastName != nil {
			profile.LastName = data.LastName
		}
		if data.Phone != nil {
			profile.Phone = data.Phone
		}
		if data.DateOfBirth != nil {
			profile.DateOfBirth = data.DateOfBirth
		}
		if data.Gender != nil {
			profile.Gender = data.Gender
		}
		if data.Address != nil {
			profile.Address = data.Address
		}
		if data.City != nil {
			profile.City = data.City
		}
		if data.Country != nil {
			profile.Country = data.Country
		}
		if data.Avatar != nil {
			profile.Avatar = data.Avatar
		}
		if data.Bio != nil {
			profile.Bio = data.Bio
		}
		if data.ParentName != nil {
			profile.ParentName = data.ParentName
		}
		if data.ParentPhone != nil {
			profile.ParentPhone = data.ParentPhone
		}
		if data.ParentEmail != nil {
			profile.ParentEmail = data.ParentEmail
		}
	}

	// Update/Create profile in database
	updatedProfile, err := s.repository.UpdateProfile(ctx, userID, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return &response.ProfileResponse{
		ID:          updatedProfile.ID.String(),
		UserID:      updatedProfile.UserID.String(),
		FirstName:   updatedProfile.FirstName,
		LastName:    updatedProfile.LastName,
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
	}, nil
}

func (s *service) Profile(ctx context.Context) (*response.UserProfileResponse, error) {
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

	// Get profile data (might not exist)
	profile, err := s.repository.GetProfileByUserID(ctx, userID)
	var profileData *response.ProfileResponse
	if err == nil {
		// Profile exists, map it to response
		profileData = &response.ProfileResponse{
			ID:          profile.ID.String(),
			UserID:      profile.UserID.String(),
			FirstName:   profile.FirstName,
			LastName:    profile.LastName,
			Phone:       profile.Phone,
			DateOfBirth: profile.DateOfBirth,
			Gender:      profile.Gender,
			Address:     profile.Address,
			City:        profile.City,
			Country:     profile.Country,
			Avatar:      profile.Avatar,
			Bio:         profile.Bio,
			ParentName:  profile.ParentName,
			ParentPhone: profile.ParentPhone,
			ParentEmail: profile.ParentEmail,
			CreatedAt:   profile.CreatedAt,
			UpdatedAt:   profile.UpdatedAt,
		}
	}
	// If profile doesn't exist, profileData will be nil

	return &response.UserProfileResponse{
		ID:         user.ID.String(),
		Name:       user.Name,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		Profile:    profileData,
	}, nil
}
