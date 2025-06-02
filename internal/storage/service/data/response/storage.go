package response

type StorageResponse struct {
	PublicID  string `json:"public_id"`
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
	Format    string `json:"format"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Bytes     int    `json:"bytes"`
	FileType  string `json:"file_type"`
	LogID     string `json:"log_id"`
}

type DeleteResponse struct {
	Success  bool   `json:"success"`
	PublicID string `json:"public_id"`
	Message  string `json:"message"`
}

type GetFileResponse struct {
	PublicID    string `json:"public_id"`
	Content     []byte `json:"-"` // Don't include in JSON
	ContentType string `json:"content_type,omitempty"`
	Filename    string `json:"filename,omitempty"`
}

type StorageLogResponse struct {
	ID               string  `json:"id"`
	UserID           string  `json:"user_id"`
	PublicID         string  `json:"public_id"`
	OriginalFilename string  `json:"original_filename"`
	FileType         string  `json:"file_type"`
	FileSize         int64   `json:"file_size"`
	MimeType         string  `json:"mime_type"`
	URL              string  `json:"url"`
	SecureURL        string  `json:"secure_url"`
	Folder           *string `json:"folder"`
	Width            *int    `json:"width"`
	Height           *int    `json:"height"`
	Format           *string `json:"format"`
	CreatedAt        int64   `json:"created_at"`
	UpdatedAt        int64   `json:"updated_at"`
}

type StorageHistoryResponse struct {
	Logs []*StorageLogResponse `json:"logs"`
}
