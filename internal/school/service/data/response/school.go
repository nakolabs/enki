package response

import (
	"github.com/google/uuid"
)

type ListSchool []DetailSchool

type DetailSchool struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Level       string    `json:"level"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	Province    string    `json:"province"`
	PostalCode  string    `json:"postal_code"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Website     string    `json:"website"`
	Logo        string    `json:"logo"`
	Banner      string    `json:"banner"`
	CreatedAt   int64     `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   int64     `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by,omitempty"`
}
