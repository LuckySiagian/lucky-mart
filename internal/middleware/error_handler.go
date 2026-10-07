package middleware

import (
	"errors"
	"fmt"

	"backend-sembako/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FieldValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				out := make([]FieldValidationError, len(ve))
				for i, fe := range ve {
					out[i] = FieldValidationError{
						Field:   fe.Field(),
						Message: getValidationMsg(fe),
					}
				}
				helpers.SendValidationError(c, out)
				return
			}

			helpers.SendInternalServerError(c, err.Error())
		}
	}
}

func getValidationMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is a required field", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}
