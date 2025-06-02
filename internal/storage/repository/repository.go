package repository

import (
	"context"
	"enuma-elish/internal/storage/service/data/request"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StorageLog struct {
	ID               uuid.UUID `db:"id" json:"id"`
	UserID           uuid.UUID `db:"user_id" json:"user_id"`
	PublicID         string    `db:"public_id" json:"public_id"`
	OriginalFilename string    `db:"original_filename" json:"original_filename"`
	FileType         string    `db:"file_type" json:"file_type"`
	FileSize         int64     `db:"file_size" json:"file_size"`
	MimeType         string    `db:"mime_type" json:"mime_type"`
	URL              string    `db:"url" json:"url"`
	SecureURL        string    `db:"secure_url" json:"secure_url"`
	Folder           *string   `db:"folder" json:"folder"`
	Width            *int      `db:"width" json:"width"`
	Height           *int      `db:"height" json:"height"`
	Format           *string   `db:"format" json:"format"`
	CreatedAt        int64     `db:"created_at" json:"created_at"`
	UpdatedAt        int64     `db:"updated_at" json:"updated_at"`
}

type Repository interface {
	CreateStorageLog(ctx context.Context, log *StorageLog) (*StorageLog, error)
	GetStorageLogsByUserID(ctx context.Context, userID uuid.UUID, query request.GetStorageHistoryQuery) ([]*StorageLog, int, error)
	GetStorageLogByPublicID(ctx context.Context, publicID string) (*StorageLog, error)
	DeleteStorageLog(ctx context.Context, publicID string) error
	GetStorageLogsByFileType(ctx context.Context, userID uuid.UUID, fileType string, query request.GetStorageHistoryQuery) ([]*StorageLog, int, error)
}

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateStorageLog(ctx context.Context, log *StorageLog) (*StorageLog, error) {
	query := `
		INSERT INTO storage_log (
			user_id, public_id, original_filename, file_type, file_size, 
			mime_type, url, secure_url, folder, width, height, format
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		) RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx, query,
		log.UserID, log.PublicID, log.OriginalFilename, log.FileType, log.FileSize,
		log.MimeType, log.URL, log.SecureURL, log.Folder, log.Width, log.Height, log.Format,
	).Scan(&log.ID, &log.CreatedAt, &log.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return log, nil
}

func (r *repository) GetStorageLogsByUserID(ctx context.Context, userID uuid.UUID, query request.GetStorageHistoryQuery) ([]*StorageLog, int, error) {
	// Count total records
	countQuery := `SELECT COUNT(*) FROM storage_log WHERE user_id = $1`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated records using standard pagination
	dataQuery := `
		SELECT id, user_id, public_id, original_filename, file_type, file_size,
			   mime_type, url, secure_url, folder, width, height, format,
			   created_at, updated_at
		FROM storage_log 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, dataQuery, userID, query.PageSize, query.GetOffset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*StorageLog
	for rows.Next() {
		log := &StorageLog{}
		err := rows.Scan(
			&log.ID, &log.UserID, &log.PublicID, &log.OriginalFilename, &log.FileType,
			&log.FileSize, &log.MimeType, &log.URL, &log.SecureURL, &log.Folder,
			&log.Width, &log.Height, &log.Format, &log.CreatedAt, &log.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}

func (r *repository) GetStorageLogsByFileType(ctx context.Context, userID uuid.UUID, fileType string, query request.GetStorageHistoryQuery) ([]*StorageLog, int, error) {
	// Count total records
	countQuery := `SELECT COUNT(*) FROM storage_log WHERE user_id = $1 AND file_type = $2`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID, fileType).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated records using standard pagination
	dataQuery := `
		SELECT id, user_id, public_id, original_filename, file_type, file_size,
			   mime_type, url, secure_url, folder, width, height, format,
			   created_at, updated_at
		FROM storage_log 
		WHERE user_id = $1 AND file_type = $2
		ORDER BY created_at DESC 
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, dataQuery, userID, fileType, query.PageSize, query.GetOffset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*StorageLog
	for rows.Next() {
		log := &StorageLog{}
		err := rows.Scan(
			&log.ID, &log.UserID, &log.PublicID, &log.OriginalFilename, &log.FileType,
			&log.FileSize, &log.MimeType, &log.URL, &log.SecureURL, &log.Folder,
			&log.Width, &log.Height, &log.Format, &log.CreatedAt, &log.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}

func (r *repository) GetStorageLogByPublicID(ctx context.Context, publicID string) (*StorageLog, error) {
	query := `
		SELECT id, user_id, public_id, original_filename, file_type, file_size,
			   mime_type, url, secure_url, folder, width, height, format,
			   created_at, updated_at
		FROM storage_log 
		WHERE public_id = $1`

	var storageLog StorageLog
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(
		&storageLog.ID, &storageLog.UserID, &storageLog.PublicID, &storageLog.OriginalFilename,
		&storageLog.FileType, &storageLog.FileSize, &storageLog.MimeType, &storageLog.URL,
		&storageLog.SecureURL, &storageLog.Folder, &storageLog.Width, &storageLog.Height,
		&storageLog.Format, &storageLog.CreatedAt, &storageLog.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &storageLog, nil
}

func (r *repository) DeleteStorageLog(ctx context.Context, publicID string) error {
	query := `DELETE FROM storage_log WHERE public_id = $1`
	_, err := r.db.ExecContext(ctx, query, publicID)
	return err
}

func (r *repository) CountStorageLogsByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM storage_log WHERE user_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
