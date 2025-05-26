package request

import (
	commonHttp "enuma-elish/pkg/http"

	"github.com/google/uuid"
)

type InviteTeacherRequest struct {
	SchoolID uuid.UUID `json:"school_id"`
	Emails   []string  `json:"emails"`
}

type VerifyTeacherEmailRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type UpdateTeacherAfterVerifyEmailRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Email    string `json:"email"`
}

type GetListTeacherQuery struct {
	SchoolID string `form:"school_id" binding:"required,uuid"`
	commonHttp.Query
}

func (q GetListTeacherQuery) Get() (commonHttp.Query, map[string]interface{}) {
	f := map[string]interface{}{
		"school_id": q.SchoolID,
	}
	return q.Query, f
}
