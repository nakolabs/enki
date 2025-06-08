package request

import (
	commonHttp "enuma-elish/pkg/http"
)

type CreateSchoolRequest struct {
	Name  string `json:"name" validate:"required"`
	Level string `json:"level" validate:"required,oneof='preschool' 'kindergarten' 'elementary' 'junior' 'senior' 'college"`
}

type UpdateSchoolProfileRequest struct {
	Name        string `json:"name" validate:"required"`
	Level       string `json:"level" validate:"required"`
	Description string `json:"description"`
	Address     string `json:"address"`
	City        string `json:"city"`
	Province    string `json:"province"`
	PostalCode  string `json:"postal_code"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Website     string `json:"website"`
	Logo        string `json:"logo"`
	Banner      string `json:"banner"`
}

type GetListSchoolQuery struct {
	commonHttp.Query
	Level string `form:"level"`
}

func (q GetListSchoolQuery) Get() (commonHttp.Query, map[string]interface{}) {
	f := map[string]interface{}{
		"level": q.Level,
	}
	return q.Query, f
}
