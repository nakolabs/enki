package request

import (
	commonHttp "enuma-elish/pkg/http"

	"github.com/google/uuid"
)

type CreateSubjectRequest struct {
	SchoolID uuid.UUID `json:"school_id" validate:"required"`
	Name     string    `json:"name" validate:"required"`
}

type UpdateSubjectRequest struct {
	Name string `json:"name" validate:"required"`
}

type GetListSubjectQuery struct {
	SchoolID string `form:"school_id" binding:"required,uuid"`
	commonHttp.Query
}

func (q GetListSubjectQuery) Get() (commonHttp.Query, map[string]interface{}) {
	f := map[string]interface{}{
		"school_id": q.SchoolID,
	}
	return q.Query, f
}
