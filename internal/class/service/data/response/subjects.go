package response

import "github.com/google/uuid"

type SubjectInClass struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	SchoolID  uuid.UUID `json:"school_id"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

type GetSubjectsByClassResponse []SubjectInClass
