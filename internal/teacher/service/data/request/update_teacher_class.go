package request

import "github.com/google/uuid"

type UpdateTeacherClassRequest struct {
	TeacherID  uuid.UUID `json:"teacher_id" validate:"required"`
	OldClassID uuid.UUID `json:"old_class_id" validate:"required"`
	NewClassID uuid.UUID `json:"new_class_id" validate:"required"`
}
