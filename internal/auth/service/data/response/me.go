package response

import "github.com/google/uuid"

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	DateOfBirth string    `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Avatar      string    `json:"avatar"`
	Bio         string    `json:"bio"`
	ParentName  string    `json:"parent_name"`
	ParentPhone string    `json:"parent_phone"`
	ParentEmail string    `json:"parent_email"`
	CreatedAt   int64     `json:"created_at"`
	UpdatedAt   int64     `json:"updated_at"`
	DeletedAt   int64     `json:"deleted_at"`
	DeletedBy   string    `json:"deleted_by"`
}
