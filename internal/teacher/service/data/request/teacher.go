package request

import (
	commonHttp "enuma-elish/pkg/http"

	"github.com/google/uuid"
)

type InviteTeacherRequest struct {
	SchoolID uuid.UUID        `json:"school_id"`
	Teachers []TeacherRequest `json:"teachers"`
}

type VerifyTeacherEmailRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type TeacherRequest struct {
	Name       string   `json:"name" validate:"required"`
	Email      string   `json:"email" validate:"required,email"`
	SubjectIDs []string `json:"subject_ids" validate:"uuid"`
	ClassIDs   []string `json:"class_ids" validate:"uuid"`
}

type UpdateTeacherAfterVerifyEmailRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Email    string `json:"email"`
}

type GetListTeacherQuery struct {
	SchoolID   string `form:"school_id" binding:"required,uuid"`
	ClassID    string `form:"class_id"`
	SubjectID  string `form:"subject_id"`
	IsVerified string `form:"is_verified" binding:"omitempty,oneof=true false"`
	commonHttp.Query
}

func (q GetListTeacherQuery) Get() (commonHttp.Query, map[string]interface{}) {
	f := map[string]interface{}{
		"school_id": q.SchoolID,
	}
	if q.ClassID != "" {
		f["class_id"] = q.ClassID
	}
	if q.SubjectID != "" {
		f["subject_id"] = q.SubjectID
	}
	if q.IsVerified != "" {
		f["is_verified"] = q.IsVerified == "true"
	}
	return q.Query, f
}
