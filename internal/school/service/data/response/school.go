package response

import (
	"github.com/google/uuid"
)

type DetailSchool struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Level     string    `json:"level"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

type ListSchool []DetailSchool
