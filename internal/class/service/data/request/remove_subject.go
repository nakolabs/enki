package request

import "github.com/google/uuid"

type RemoveSubjectFromClassRequest struct {
	ClassID    uuid.UUID   `json:"class_id" validate:"required"`
	SubjectIDs []uuid.UUID `json:"subject_ids" validate:"required,min=1"`
}
