package response

import (
	"github.com/google/uuid"
)

type PPDBResponse struct {
	ID        uuid.UUID `json:"id"`
	SchoolID  uuid.UUID `json:"school_id"`
	StartAt   int64     `json:"start_at"`
	EndAt     int64     `json:"end_at"`
	Status    string    `json:"status"` // active, inactive
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

type GetListPPDBResponse []PPDBResponse

type PPDBStudentResponse struct {
	ID        uuid.UUID `json:"id"`
	PPDBID    uuid.UUID `json:"ppdb_id"`
	StudentID uuid.UUID `json:"student_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"` // registered, accepted, rejected
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

type GetPPDBRegistrantsResponse []PPDBStudentResponse

type PPDBRegistrationResponse struct {
	ID      uuid.UUID `json:"id"`
	PPDBID  uuid.UUID `json:"ppdb_id"`
	Message string    `json:"message"`
}
