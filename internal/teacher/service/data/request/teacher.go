package request

import (
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
