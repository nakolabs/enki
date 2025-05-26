package request

import (
	commonHttp "enuma-elish/pkg/http"

	"github.com/google/uuid"
)

type CreatePPDBRequest struct {
	SchoolID uuid.UUID `json:"school_id" validate:"required"`
	StartAt  int64     `json:"start_at" validate:"required"`
	EndAt    int64     `json:"end_at" validate:"required"`
}

type UpdatePPDBRequest struct {
	ID      uuid.UUID `json:"id" validate:"required"`
	StartAt int64     `json:"start_at" validate:"required"`
	EndAt   int64     `json:"end_at" validate:"required"`
}

type RegisterPPDBRequest struct {
	PPDBID uuid.UUID `json:"ppdb_id" validate:"required"`
	Name   string    `json:"name" validate:"required"`
}

type PPDBSelectionRequest struct {
	PPDBID           uuid.UUID   `json:"ppdb_id" validate:"required"`
	AcceptedStudents []uuid.UUID `json:"accepted_students" validate:"required,min=1"`
}

type GetListPPDBQuery struct {
	SchoolID string `form:"school_id" binding:"uuid"`
	Status   string `form:"status"` // active, inactive, all
	commonHttp.Query
}

func (q GetListPPDBQuery) Get() (commonHttp.Query, map[string]interface{}) {
	f := map[string]interface{}{}

	if q.SchoolID != "" {
		schoolID, err := uuid.Parse(q.SchoolID)
		if err == nil {
			f["school_id"] = schoolID
		}
	}

	if q.Status != "" && q.Status != "all" {
		f["status"] = q.Status
	}

	return q.Query, f
}

type GetPPDBRegistrantsQuery struct {
	PPDBID string `form:"ppdb_id" binding:"uuid"`
	Status string `form:"status"` // registered, accepted, rejected, all
	commonHttp.Query
}

func (q GetPPDBRegistrantsQuery) Get() (commonHttp.Query, map[string]interface{}) {
	f := map[string]interface{}{}

	if q.PPDBID != "" {
		ppdbID, err := uuid.Parse(q.PPDBID)
		if err == nil {
			f["ppdb_id"] = ppdbID
		}
	}

	if q.Status != "" && q.Status != "all" {
		f["status"] = q.Status
	}

	return q.Query, f
}
