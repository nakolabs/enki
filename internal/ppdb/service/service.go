package service

import (
	"context"
	"database/sql"
	"enuma-elish/config"
	"enuma-elish/internal/ppdb/repository"
	"enuma-elish/internal/ppdb/service/data/request"
	"enuma-elish/internal/ppdb/service/data/response"
	commonError "enuma-elish/pkg/error"
	commonHttp "enuma-elish/pkg/http"
	"enuma-elish/pkg/jwt"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreatePPDB(ctx context.Context, data request.CreatePPDBRequest) error
	UpdatePPDB(ctx context.Context, data request.UpdatePPDBRequest) error
	DeletePPDB(ctx context.Context, id uuid.UUID) error
	GetPPDBByID(ctx context.Context, id uuid.UUID) (response.PPDBResponse, error)
	GetListPPDB(ctx context.Context, query request.GetListPPDBQuery) (response.GetListPPDBResponse, *commonHttp.Meta, error)

	RegisterPPDB(ctx context.Context, data request.RegisterPPDBRequest) (response.PPDBRegistrationResponse, error)
	GetPPDBRegistrants(ctx context.Context, query request.GetPPDBRegistrantsQuery) (response.GetPPDBRegistrantsResponse, *commonHttp.Meta, error)
	SelectPPDBStudents(ctx context.Context, data request.PPDBSelectionRequest) error
}

type service struct {
	repository repository.Repository
	config     *config.Config
}

func New(repository repository.Repository, config *config.Config) Service {
	return &service{
		repository: repository,
		config:     config,
	}
}

func (s *service) CreatePPDB(ctx context.Context, data request.CreatePPDBRequest) error {
	claim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return err
	}

	if claim.User.SchoolID != data.SchoolID {
		return commonError.ErrUnauthorized
	}

	if data.StartAt >= data.EndAt {
		return errors.New("start date must be before end date")
	}

	now := time.Now().UnixMilli()
	ppdb := &repository.PPDB{
		ID:        uuid.New(),
		SchoolID:  data.SchoolID,
		StartAt:   data.StartAt,
		EndAt:     data.EndAt,
		CreatedAt: now,
		UpdatedAt: 0,
	}

	err = s.repository.CreatePPDB(ctx, ppdb)
	if err != nil {
		log.Err(err).Msg("Failed to create PPDB")
		return err
	}

	return nil
}

func (s *service) UpdatePPDB(ctx context.Context, data request.UpdatePPDBRequest) error {
	claim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return err
	}

	ppdb, err := s.repository.GetPPDBByID(ctx, data.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return commonError.ErrNotFound
		}
		log.Err(err).Msg("Failed to get PPDB")
		return err
	}

	if claim.User.SchoolID != ppdb.SchoolID {
		return commonError.ErrUnauthorized
	}

	if data.StartAt >= data.EndAt {
		return errors.New("start date must be before end date")
	}

	ppdb.StartAt = data.StartAt
	ppdb.EndAt = data.EndAt
	ppdb.UpdatedAt = time.Now().UnixMilli()

	err = s.repository.UpdatePPDB(ctx, ppdb)
	if err != nil {
		log.Err(err).Msg("Failed to update PPDB")
		return err
	}

	return nil
}

func (s *service) DeletePPDB(ctx context.Context, id uuid.UUID) error {
	claim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return err
	}

	ppdb, err := s.repository.GetPPDBByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return commonError.ErrNotFound
		}
		log.Err(err).Msg("Failed to get PPDB")
		return err
	}

	if claim.User.SchoolID != ppdb.SchoolID {
		return commonError.ErrUnauthorized
	}

	err = s.repository.DeletePPDB(ctx, id)
	if err != nil {
		log.Err(err).Msg("Failed to delete PPDB")
		return err
	}

	return nil
}

func (s *service) GetPPDBByID(ctx context.Context, id uuid.UUID) (response.PPDBResponse, error) {
	ppdb, err := s.repository.GetPPDBByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response.PPDBResponse{}, commonError.ErrNotFound
		}
		log.Err(err).Msg("Failed to get PPDB")
		return response.PPDBResponse{}, err
	}

	now := time.Now().UnixMilli()
	status := "inactive"
	if ppdb.StartAt <= now && ppdb.EndAt >= now {
		status = "active"
	}

	return response.PPDBResponse{
		ID:        ppdb.ID,
		SchoolID:  ppdb.SchoolID,
		StartAt:   ppdb.StartAt,
		EndAt:     ppdb.EndAt,
		Status:    status,
		CreatedAt: ppdb.CreatedAt,
		UpdatedAt: ppdb.UpdatedAt,
	}, nil
}

func (s *service) GetListPPDB(ctx context.Context, query request.GetListPPDBQuery) (response.GetListPPDBResponse, *commonHttp.Meta, error) {
	// Validate school_id parameter
	if query.SchoolID == "" {
		return response.GetListPPDBResponse{}, nil, errors.New("school_id is required")
	}

	_, err := uuid.Parse(query.SchoolID)
	if err != nil {
		return response.GetListPPDBResponse{}, nil, errors.New("invalid school_id format")
	}

	ppdbs, meta, err := s.repository.GetListPPDB(ctx, query)
	if err != nil {
		log.Err(err).Msg("Failed to get PPDB list")
		return response.GetListPPDBResponse{}, nil, err
	}

	return response.GetListPPDBResponse(ppdbs), meta, nil
}

func (s *service) RegisterPPDB(ctx context.Context, data request.RegisterPPDBRequest) (response.PPDBRegistrationResponse, error) {
	// Extract user from JWT token
	claim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return response.PPDBRegistrationResponse{}, commonError.ErrUnauthorized
	}

	// Check if PPDB exists and is active
	ppdb, err := s.repository.GetPPDBByID(ctx, data.PPDBID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response.PPDBRegistrationResponse{}, commonError.ErrNotFound
		}
		log.Err(err).Msg("Failed to get PPDB")
		return response.PPDBRegistrationResponse{}, err
	}

	now := time.Now().UnixMilli()
	if ppdb.StartAt > now || ppdb.EndAt < now {
		return response.PPDBRegistrationResponse{}, errors.New("PPDB registration is not active")
	}

	// Check if already registered using user ID
	existing, err := s.repository.GetPPDBStudentByPPDBIDAndUserID(ctx, data.PPDBID, claim.User.ID)
	if err == nil && existing != nil {
		return response.PPDBRegistrationResponse{}, errors.New("already registered for this PPDB")
	}

	ppdbStudent := &repository.PPDBStudent{
		ID:        uuid.New(),
		PPDBID:    data.PPDBID,
		StudentID: claim.User.ID,    // Use authenticated user ID
		Name:      data.Name,        // Use authenticated user name
		Email:     claim.User.Email, // Use authenticated user email
		Status:    "registered",
		CreatedAt: now,
		UpdatedAt: 0,
	}

	err = s.repository.RegisterPPDB(ctx, ppdbStudent)
	if err != nil {
		log.Err(err).Msg("Failed to register PPDB")
		return response.PPDBRegistrationResponse{}, err
	}

	return response.PPDBRegistrationResponse{
		ID:      ppdbStudent.ID,
		PPDBID:  data.PPDBID,
		Message: "Successfully registered for PPDB",
	}, nil
}

func (s *service) GetPPDBRegistrants(ctx context.Context, query request.GetPPDBRegistrantsQuery) (response.GetPPDBRegistrantsResponse, *commonHttp.Meta, error) {
	// Validate ppdb_id parameter
	if query.PPDBID == "" {
		return response.GetPPDBRegistrantsResponse{}, nil, errors.New("ppdb_id is required")
	}

	_, err := uuid.Parse(query.PPDBID)
	if err != nil {
		return response.GetPPDBRegistrantsResponse{}, nil, errors.New("invalid ppdb_id format")
	}

	registrants, meta, err := s.repository.GetPPDBRegistrants(ctx, query)
	if err != nil {
		log.Err(err).Msg("Failed to get PPDB registrants")
		return response.GetPPDBRegistrantsResponse{}, nil, err
	}

	return response.GetPPDBRegistrantsResponse(registrants), meta, nil
}

func (s *service) SelectPPDBStudents(ctx context.Context, data request.PPDBSelectionRequest) error {
	claim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return err
	}

	ppdb, err := s.repository.GetPPDBByID(ctx, data.PPDBID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return commonError.ErrNotFound
		}
		log.Err(err).Msg("Failed to get PPDB")
		return err
	}

	if claim.User.SchoolID != ppdb.SchoolID {
		return commonError.ErrUnauthorized
	}

	// Update accepted students
	err = s.repository.UpdatePPDBStudentStatus(ctx, data.PPDBID, data.AcceptedStudents, "accepted")
	if err != nil {
		log.Err(err).Msg("Failed to update PPDB student status")
		return err
	}

	return nil
}
