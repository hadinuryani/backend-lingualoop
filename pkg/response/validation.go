package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidationError(c *gin.Context, err error) {
	var errs []FieldError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errs = append(errs, FieldError{
				Field:   e.Field(),
				Message: getValidationMessage(e),
			})
		}

		Error(c, http.StatusBadRequest, MsgBadRequest, errs)
		return
	}

	Error(c, http.StatusBadRequest, MsgBadRequest, err.Error())
}

func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("Field %s wajib diisi", e.Field())
	case "email":
		return "Format email tidak valid"
	case "min":
		return fmt.Sprintf("Field %s minimal %s karakter", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("Field %s maksimal %s karakter", e.Field(), e.Param())
	default:
		return fmt.Sprintf("Field %s tidak valid", e.Field())
	}
}
