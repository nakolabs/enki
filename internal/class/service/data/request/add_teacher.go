package request

import "github.com/google/uuid"

type AddTeacherToClassRequest struct {
	ClassID    uuid.UUID   `json:"class_id" validate:"required"`
	TeacherIDs []uuid.UUID `json:"teacher_ids" validate:"required,min=1"`
}
