package response

import (
	"github.com/google/uuid"
)

type ListTeacherResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	CreateAt   int64     `json:"create_at"`
	UpdateAt   int64     `json:"update_at"`
}
