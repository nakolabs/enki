package request

type CreateSchoolRequest struct {
	Name  string `json:"name" validate:"required"`
	Level string `json:"level" validate:"required,oneof='preschool' 'kindergarten' 'elementary' 'junior' 'senior' 'college"`
}
