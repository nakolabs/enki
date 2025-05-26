package request

import "github.com/google/uuid"

type AssignTeacherToSubjectRequest struct {
	SubjectID  uuid.UUID   `json:"subject_id" validate:"required"`
	TeacherIDs []uuid.UUID `json:"teacher_ids" validate:"required,min=1"`
}
