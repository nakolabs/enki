package request

import "github.com/google/uuid"

type UpdateSubjectClassRequest struct {
	SubjectID  uuid.UUID `json:"subject_id" validate:"required"`
	OldClassID uuid.UUID `json:"old_class_id" validate:"required"`
	NewClassID uuid.UUID `json:"new_class_id" validate:"required"`
}
