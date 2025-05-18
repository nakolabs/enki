package middleware

import (
	commonError "enuma-elish/pkg/error"
	"enuma-elish/pkg/http"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"strings"
)

func ErrorParser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if err := c.Errors.Last(); err != nil {
			var apiErr commonError.Error
			if errors.As(err.Err, &apiErr) {
				response := http.NewResponse().
					SetCode(apiErr.Code).
					SetMessage(apiErr.Error())

				c.JSON(apiErr.Code, response)
				return
			}

			var validationError validator.ValidationErrors
			if errors.As(err.Err, &validationError) {
				log.Err(err).Msg("validation error 123")
				errorResponse := map[string][]string{}
				for _, v := range validationError {
					key := strings.ToLower(v.Field())
					errorResponse[key] = append(errorResponse[key], validationReadableError(v))
				}

				response := http.NewResponse().
					SetCode(422).
					SetErrors(errorResponse).
					SetMessage("validation error")

				c.JSON(422, response)
				return
			}

			if err.Type == gin.ErrorTypeBind {
				response := http.NewResponse().
					SetCode(400).
					SetErrors(gin.H{"_general": bindingReadableError(err.Err)}).
					SetMessage("invalid request body")

				c.JSON(400, response)
				return
			}

			response := http.NewResponse().
				SetCode(500).
				SetMessage("internal server error")
			c.JSON(500, response)
		}
	}
}

func validationReadableError(e validator.FieldError) string {
	field := e.Field()
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "required_with":
		return fmt.Sprintf("%s is required when %s is present", field, e.Param())
	case "required_without":
		return fmt.Sprintf("%s is required when %s is absent", field, e.Param())
	case "required_if":
		return fmt.Sprintf("%s is required if %s", field, e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, e.Param())
	case "eq":
		return fmt.Sprintf("%s must be equal to %s", field, e.Param())
	case "ne":
		return fmt.Sprintf("%s must not be equal to %s", field, e.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", field, e.Param())
	case "contains":
		return fmt.Sprintf("%s must contain %s", field, e.Param())
	case "excludes":
		return fmt.Sprintf("%s must not contain %s", field, e.Param())
	case "startswith":
		return fmt.Sprintf("%s must start with %s", field, e.Param())
	case "endswith":
		return fmt.Sprintf("%s must end with %s", field, e.Param())
	case "numeric":
		return fmt.Sprintf("%s must be a numeric value", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "boolean":
		return fmt.Sprintf("%s must be a boolean value", field)
	case "datetime":
		return fmt.Sprintf("%s must be a valid datetime", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func bindingReadableError(err error) string {
	if err == nil {
		return ""
	}

	errStr := err.Error()

	if strings.Contains(errStr, "cannot unmarshal") {
		parts := strings.Split(errStr, "into Go struct field")
		if len(parts) == 2 {
			right := strings.TrimSpace(parts[1])

			fieldParts := strings.Split(right, " ")
			if len(fieldParts) >= 1 {
				fieldName := strings.Split(fieldParts[0], ".")
				field := fieldName[len(fieldName)-1]

				return fmt.Sprintf("Field '%s' has invalid type", field)
			}
		}
	}

	return "Invalid request format"
}
