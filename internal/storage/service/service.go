package service

import (
	"context"
	"enuma-elish/config"
	"enuma-elish/internal/storage/repository"
	"enuma-elish/internal/storage/service/data/request"
	"enuma-elish/internal/storage/service/data/response"
	"enuma-elish/pkg/cloudinary"
	commonError "enuma-elish/pkg/error"
	commonHttp "enuma-elish/pkg/http"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	StoreImage(ctx context.Context, data request.StoreImageRequest) (*response.StorageResponse, error)
	StoreVideo(ctx context.Context, data request.StoreVideoRequest) (*response.StorageResponse, error)
	StoreDocument(ctx context.Context, data request.StoreDocumentRequest) (*response.StorageResponse, error)
	DeleteFile(ctx context.Context, data request.DeleteFileRequest) (*response.DeleteResponse, error)
	GetFile(ctx context.Context, publicID string) (*response.GetFileResponse, error)
	GetStorageHistory(ctx context.Context, httpQuery request.GetStorageHistoryQuery) (*response.StorageHistoryResponse, *commonHttp.Meta, error)
	GetStorageHistoryByType(ctx context.Context, fileType string, httpQuery request.GetStorageHistoryQuery) (*response.StorageHistoryResponse, *commonHttp.Meta, error)
}

type service struct {
	cloudinaryService *cloudinary.Service
	repository        repository.Repository
	config            *config.Config
}

func New(cs *cloudinary.Service, repo repository.Repository, config *config.Config) Service {
	return &service{
		cloudinaryService: cs,
		repository:        repo,
		config:            config,
	}
}

func (s *service) StoreImage(ctx context.Context, data request.StoreImageRequest) (*response.StorageResponse, error) {
	// Get user ID from context
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, commonError.ErrUnauthorized
	}

	// Validate file type
	contentType := data.Header.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		return nil, fmt.Errorf("invalid file type: only JPEG, JPG, PNG, GIF, and WebP are allowed")
	}

	// Validate file size (max 10MB for images)
	if data.Header.Size > 10*1024*1024 {
		return nil, fmt.Errorf("file too large: maximum size is 10MB")
	}

	result, err := s.cloudinaryService.UploadImage(ctx, data.File, data.Header)
	if err != nil {
		log.Err(err).Msg("Failed to store image")
		return nil, commonError.ErrInternal
	}

	// Log the storage operation
	storageLog := &repository.StorageLog{
		UserID:           userID,
		PublicID:         result.PublicID,
		OriginalFilename: data.Header.Filename,
		FileType:         "image",
		FileSize:         data.Header.Size,
		MimeType:         contentType,
		URL:              result.URL,
		SecureURL:        result.SecureURL,
		Folder:           &data.Folder,
		Width:            &result.Width,
		Height:           &result.Height,
		Format:           &result.Format,
	}

	if data.Folder == "" {
		storageLog.Folder = nil
	}

	logResult, err := s.repository.CreateStorageLog(ctx, storageLog)
	if err != nil {
		log.Err(err).Msg("Failed to log storage operation")
	}

	return &response.StorageResponse{
		PublicID:  result.PublicID,
		URL:       result.URL,
		SecureURL: result.SecureURL,
		Format:    result.Format,
		Width:     result.Width,
		Height:    result.Height,
		Bytes:     result.Bytes,
		FileType:  "image",
		LogID:     logResult.ID.String(),
	}, nil
}

func (s *service) StoreVideo(ctx context.Context, data request.StoreVideoRequest) (*response.StorageResponse, error) {
	// Get user ID from context
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, commonError.ErrUnauthorized
	}

	// Validate file type
	contentType := data.Header.Header.Get("Content-Type")
	if !isValidVideoType(contentType) {
		return nil, fmt.Errorf("invalid file type: only MP4, AVI, MOV, and WebM are allowed")
	}

	// Validate file size (max 100MB for videos)
	if data.Header.Size > 100*1024*1024 {
		return nil, fmt.Errorf("file too large: maximum size is 100MB")
	}

	result, err := s.cloudinaryService.UploadVideo(ctx, data.File, data.Header)
	if err != nil {
		log.Err(err).Msg("Failed to store video")
		return nil, commonError.ErrInternal
	}

	// Log the storage operation
	storageLog := &repository.StorageLog{
		UserID:           userID,
		PublicID:         result.PublicID,
		OriginalFilename: data.Header.Filename,
		FileType:         "video",
		FileSize:         data.Header.Size,
		MimeType:         contentType,
		URL:              result.URL,
		SecureURL:        result.SecureURL,
		Folder:           &data.Folder,
		Width:            &result.Width,
		Height:           &result.Height,
		Format:           &result.Format,
	}

	if data.Folder == "" {
		storageLog.Folder = nil
	}

	logResult, err := s.repository.CreateStorageLog(ctx, storageLog)
	if err != nil {
		log.Err(err).Msg("Failed to log storage operation")
	}

	return &response.StorageResponse{
		PublicID:  result.PublicID,
		URL:       result.URL,
		SecureURL: result.SecureURL,
		Format:    result.Format,
		Width:     result.Width,
		Height:    result.Height,
		Bytes:     result.Bytes,
		FileType:  "video",
		LogID:     logResult.ID.String(),
	}, nil
}

func (s *service) StoreDocument(ctx context.Context, data request.StoreDocumentRequest) (*response.StorageResponse, error) {
	// Get user ID from context
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, commonError.ErrUnauthorized
	}

	// Validate file type
	contentType := data.Header.Header.Get("Content-Type")
	if !isValidDocumentType(contentType) {
		return nil, fmt.Errorf("invalid file type: only PDF, DOC, DOCX, XLS, XLSX, PPT, and PPTX are allowed")
	}

	// Validate file size (max 50MB for documents)
	if data.Header.Size > 50*1024*1024 {
		return nil, fmt.Errorf("file too large: maximum size is 50MB")
	}

	result, err := s.cloudinaryService.UploadDocument(ctx, data.File, data.Header)
	if err != nil {
		log.Err(err).Msg("Failed to store document")
		return nil, commonError.ErrInternal
	}

	// Log the storage operation
	storageLog := &repository.StorageLog{
		UserID:           userID,
		PublicID:         result.PublicID,
		OriginalFilename: data.Header.Filename,
		FileType:         "document",
		FileSize:         data.Header.Size,
		MimeType:         contentType,
		URL:              result.URL,
		SecureURL:        result.SecureURL,
		Folder:           &data.Folder,
		Format:           &result.Format,
	}

	if data.Folder == "" {
		storageLog.Folder = nil
	}

	logResult, err := s.repository.CreateStorageLog(ctx, storageLog)
	if err != nil {
		log.Err(err).Msg("Failed to log storage operation")
	}

	return &response.StorageResponse{
		PublicID:  result.PublicID,
		URL:       result.URL,
		SecureURL: result.SecureURL,
		Format:    result.Format,
		Width:     result.Width,
		Height:    result.Height,
		Bytes:     result.Bytes,
		FileType:  "document",
		LogID:     logResult.ID.String(),
	}, nil
}

func (s *service) DeleteFile(ctx context.Context, data request.DeleteFileRequest) (*response.DeleteResponse, error) {
	// Delete from Cloudinary
	err := s.cloudinaryService.DeleteFile(ctx, data.PublicID)
	if err != nil {
		log.Err(err).Msg("Failed to delete file from cloudinary")
		return &response.DeleteResponse{
			Success:  false,
			PublicID: data.PublicID,
			Message:  "failed to delete file",
		}, commonError.ErrInternal
	}

	// Delete from storage log
	err = s.repository.DeleteStorageLog(ctx, data.PublicID)
	if err != nil {
		log.Err(err).Msg("Failed to delete storage log")
	}

	return &response.DeleteResponse{
		Success:  true,
		PublicID: data.PublicID,
		Message:  "file deleted successfully",
	}, nil
}

func (s *service) GetFile(ctx context.Context, publicID string) (*response.GetFileResponse, error) {
	// Check if file exists in our database
	storageLog, err := s.repository.GetStorageLogByPublicID(ctx, publicID)
	if err != nil {
		return nil, commonError.ErrNotFound
	}

	// Get file content from Cloudinary
	content, contentType, err := s.cloudinaryService.GetFileContent(publicID)
	if err != nil {
		log.Err(err).Msg("Failed to get file content")
		return nil, commonError.ErrInternal
	}

	return &response.GetFileResponse{
		PublicID:    publicID,
		Content:     content,
		ContentType: contentType,
		Filename:    storageLog.OriginalFilename,
	}, nil
}

func (s *service) GetStorageHistory(ctx context.Context, httpQuery request.GetStorageHistoryQuery) (*response.StorageHistoryResponse, *commonHttp.Meta, error) {
	// Get user ID from context
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, nil, commonError.ErrUnauthorized
	}

	// Get logs
	logs, total, err := s.repository.GetStorageLogsByUserID(ctx, userID, httpQuery)
	if err != nil {
		log.Err(err).Msg("Failed to get storage history")
		return nil, nil, commonError.ErrInternal
	}

	meta := commonHttp.NewMetaFromQuery(httpQuery, total)

	// Convert to response format
	var responseLogs []*response.StorageLogResponse
	for _, storageLog := range logs {
		responseLogs = append(responseLogs, &response.StorageLogResponse{
			ID:               storageLog.ID.String(),
			UserID:           storageLog.UserID.String(),
			PublicID:         storageLog.PublicID,
			OriginalFilename: storageLog.OriginalFilename,
			FileType:         storageLog.FileType,
			FileSize:         storageLog.FileSize,
			MimeType:         storageLog.MimeType,
			URL:              storageLog.URL,
			SecureURL:        storageLog.SecureURL,
			Folder:           storageLog.Folder,
			Width:            storageLog.Width,
			Height:           storageLog.Height,
			Format:           storageLog.Format,
			CreatedAt:        storageLog.CreatedAt,
			UpdatedAt:        storageLog.UpdatedAt,
		})
	}

	result := &response.StorageHistoryResponse{
		Logs: responseLogs,
	}

	return result, meta, nil
}

func (s *service) GetStorageHistoryByType(ctx context.Context, fileType string, httpQuery request.GetStorageHistoryQuery) (*response.StorageHistoryResponse, *commonHttp.Meta, error) {
	// Get user ID from context
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, nil, commonError.ErrUnauthorized
	}

	// Get logs by file type
	logs, total, err := s.repository.GetStorageLogsByFileType(ctx, userID, fileType, httpQuery)
	if err != nil {
		log.Err(err).Msg("Failed to get storage history by type")
		return nil, nil, commonError.ErrInternal
	}

	meta := commonHttp.NewMetaFromQuery(httpQuery, total)

	// Convert to response format
	var responseLogs []*response.StorageLogResponse
	for _, storageLog := range logs {
		responseLogs = append(responseLogs, &response.StorageLogResponse{
			ID:               storageLog.ID.String(),
			UserID:           storageLog.UserID.String(),
			PublicID:         storageLog.PublicID,
			OriginalFilename: storageLog.OriginalFilename,
			FileType:         storageLog.FileType,
			FileSize:         storageLog.FileSize,
			MimeType:         storageLog.MimeType,
			URL:              storageLog.URL,
			SecureURL:        storageLog.SecureURL,
			Folder:           storageLog.Folder,
			Width:            storageLog.Width,
			Height:           storageLog.Height,
			Format:           storageLog.Format,
			CreatedAt:        storageLog.CreatedAt,
			UpdatedAt:        storageLog.UpdatedAt,
		})
	}

	result := &response.StorageHistoryResponse{
		Logs: responseLogs,
	}

	return result, meta, nil
}

func (s *service) getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userIDFromCtx := ctx.Value("user_id")
	if userIDFromCtx == nil {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}

	userID, ok := userIDFromCtx.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user ID format")
	}

	return userID, nil
}

// Helper functions for file type validation
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}
	return contains(validTypes, contentType)
}

func isValidVideoType(contentType string) bool {
	validTypes := []string{
		"video/mp4",
		"video/avi",
		"video/quicktime",
		"video/webm",
	}
	return contains(validTypes, contentType)
}

func isValidDocumentType(contentType string) bool {
	validTypes := []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
	}
	return contains(validTypes, contentType)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}
