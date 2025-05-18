package request

import "github.com/google/uuid"

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
