package request

import "github.com/google/uuid"

type RemoveStudentFromClassRequest struct {
	ClassID    uuid.UUID   `json:"class_id" validate:"required"`
	StudentIDs []uuid.UUID `json:"student_ids" validate:"required,min=1"`
}
