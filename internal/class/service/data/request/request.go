package request

import (
	commonHttp "enuma-elish/pkg/http"

	"github.com/google/uuid"
)

type CreateClassRequest struct {
	SchoolID uuid.UUID `json:"school_id" validate:"required"`
	Name     string    `json:"name" validate:"required,min=1,max=100"`
}

type UpdateClassRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type GetListClassQuery struct {
	commonHttp.Query
	SchoolID string `form:"school_id" validate:"required", binding:"uuid"`
}
