package request

import (
	commonHttp "enuma-elish/pkg/http"
	"mime/multipart"
)

type StoreFileRequest struct {
	File   multipart.File        `json:"-"`
	Header *multipart.FileHeader `json:"-"`
	Folder string                `json:"folder,omitempty"`
}

type StoreImageRequest struct {
	File   multipart.File        `json:"-"`
	Header *multipart.FileHeader `json:"-"`
	Folder string                `json:"folder,omitempty"`
}

type StoreVideoRequest struct {
	File   multipart.File        `json:"-"`
	Header *multipart.FileHeader `json:"-"`
	Folder string                `json:"folder,omitempty"`
}

type StoreDocumentRequest struct {
	File   multipart.File        `json:"-"`
	Header *multipart.FileHeader `json:"-"`
	Folder string                `json:"folder,omitempty"`
}

type DeleteFileRequest struct {
	PublicID string `json:"public_id" binding:"required"`
}

type GetStorageHistoryQuery struct {
	commonHttp.Query
	FileType string `query:"file_type"`
}
