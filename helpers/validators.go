package helpers

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validator = validator.New()

func CreateValidationErrorMessages(err error) string {
	// چک کردن نوع ارور
	if err != nil {
		var errMessages []string
		e, ok := err.(validator.ValidationErrors)
		// تبدیل ارور به ValidationErrors
		if ok {
			for _, e := range e {
				errMessages = append(errMessages, fmt.Sprintf("Field '%s' is required", e.Field()))
			}

			return strings.Join(errMessages, ", ")
		} else {
			return err.Error()
		}
	}

	return "hi"
}
