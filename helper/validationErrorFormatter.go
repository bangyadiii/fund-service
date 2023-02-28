package helper

import "github.com/go-playground/validator/v10"

func FormatErrorValidation(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return errors

}
