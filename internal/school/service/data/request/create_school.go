package request

type CreateSchoolRequest struct {
	Name  string `json:"name" validate:"required"`
	Level string `json:"level" validate:"required,oneof='preschool' 'kindergarten' 'elementary' 'junior' 'senior' 'college"`
}

type UpdateSchoolProfileRequest struct {
	Name        string  `json:"name" validate:"required"`
	Level       string  `json:"level" validate:"required"`
	Description *string `json:"description"`
	Address     *string `json:"address"`
	City        *string `json:"city"`
	Province    *string `json:"province"`
	PostalCode  *string `json:"postal_code"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email"`
	Website     *string `json:"website"`
	Logo        *string `json:"logo"`
	Banner      *string `json:"banner"`
}
