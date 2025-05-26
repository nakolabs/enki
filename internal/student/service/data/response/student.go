package response

import "github.com/google/uuid"

type GetStudentResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	CreateAt   int64     `json:"created_at"`
	UpdateAt   int64     `json:"updated_at"`
}

type GetListStudentResponse []GetStudentResponse

type GetDetailStudentResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	CreateAt   int64     `json:"created_at"`
	UpdateAt   int64     `json:"updated_at"`
}
