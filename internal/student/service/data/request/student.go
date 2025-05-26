package request

import (
	commonHttp "enuma-elish/pkg/http"
	"github.com/google/uuid"
)

type InviteStudentRequest struct {
	SchoolID uuid.UUID `json:"school_id"`
	Emails   []string  `json:"emails"`
}

type UpdateStudentAfterVerifyEmailRequest struct {
	Token    string `json:"token"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type VerifyStudentEmailRequest struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

type GetListStudentQuery struct {
	SchoolID string `form:"school_id" binding:"uuid"`
	commonHttp.Query
}

func (q GetListStudentQuery) Get() (commonHttp.Query, map[string]interface{}) {
	f := map[string]interface{}{
		"school_id": q.SchoolID,
	}
	return q.Query, f
}
