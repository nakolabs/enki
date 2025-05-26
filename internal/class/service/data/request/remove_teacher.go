package request

import "github.com/google/uuid"

type RemoveTeacherFromClassRequest struct {
	ClassID    uuid.UUID   `json:"class_id" validate:"required"`
	TeacherIDs []uuid.UUID `json:"teacher_ids" validate:"required,min=1"`
}
