package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// NewValidator sets the validator to a new instance of the validator. New interface
func NewValidator() {
	validate = validator.New()
}

// Validate validates a struct. If validation fails it returns an error describing the validation errors
// @param i The struct to be validated
// @return The error or nil if everything is valid or no
func Validate(i interface{}) error {
	err := validate.Struct(i)
	// err is a validator. ValidationErrors object.
	if err != nil {
		errs := err.(validator.ValidationErrors)
		out := make([]string, len(errs))
		// errors. As err errs.
		if errors.As(err, &errs) {
			// Add a ErrorMessage to the out array.
			for i, fe := range errs {
				out[i] = ErrorMessage(fe.Field(), fe.Tag())
			}
		}

		return errors.New(strings.Join(out, ", "))
	}

	return nil
}

// ErrorMessage return error message from validator error based on translator
func ErrorMessage(field, tag string) string {
	return fmt.Sprintf("Invalid %s value for field %s", tag, field)
}
