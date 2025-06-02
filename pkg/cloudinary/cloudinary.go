package cloudinary

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type Service struct {
	cld    *cloudinary.Cloudinary
	folder string
}

type UploadResult struct {
	PublicID  string `json:"public_id"`
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
	Format    string `json:"format"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Bytes     int    `json:"bytes"`
}

func New(cloudName, apiKey, apiSecret, folder string) (*Service, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}

	return &Service{
		cld:    cld,
		folder: folder,
	}, nil
}

func (s *Service) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, resourceType string) (*UploadResult, error) {
	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	publicID := fmt.Sprintf("%s/%s/%s", s.folder, resourceType, strings.TrimSuffix(filename, ext))

	// Upload to Cloudinary with private access
	result, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:     publicID,
		ResourceType: resourceType,
		Folder:       s.folder,
		Type:         "private", // Make files private
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload to cloudinary: %w", err)
	}

	return &UploadResult{
		PublicID:  result.PublicID,
		URL:       result.URL,
		SecureURL: result.SecureURL,
		Format:    result.Format,
		Width:     result.Width,
		Height:    result.Height,
		Bytes:     result.Bytes,
	}, nil
}

func (s *Service) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*UploadResult, error) {
	return s.UploadFile(ctx, file, header, "image")
}

func (s *Service) UploadVideo(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*UploadResult, error) {
	return s.UploadFile(ctx, file, header, "video")
}

func (s *Service) UploadDocument(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*UploadResult, error) {
	return s.UploadFile(ctx, file, header, "raw")
}

func (s *Service) DeleteFile(ctx context.Context, publicID string) error {
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from cloudinary: %w", err)
	}
	return nil
}

func (s *Service) GetFileContent(publicID string) ([]byte, string, error) {
	// Determine resource type from public ID
	resourceType := "image"
	if strings.Contains(publicID, "/video/") {
		resourceType = "video"
	} else if strings.Contains(publicID, "/raw/") {
		resourceType = "raw"
	}

	// Build direct access URL using API credentials
	url := fmt.Sprintf("https://%s:%s@res.cloudinary.com/%s/%s/private/%s",
		s.cld.Config.Cloud.APIKey,
		s.cld.Config.Cloud.APISecret,
		s.cld.Config.Cloud.CloudName,
		resourceType,
		publicID)

	// Fetch the file content
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch file from cloudinary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("cloudinary returned status: %d", resp.StatusCode)
	}

	// Read the file content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file content: %w", err)
	}

	// Get content type from response headers
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return content, contentType, nil
}

// Remove or simplify the GetSignedURL method since we don't need it anymore
func (s *Service) Get(publicID string) (string, error) {
	// This method can now return a simple message or be removed
	return fmt.Sprintf("File access through backend: /api/v1/storage/serve/%s", publicID), nil
}
