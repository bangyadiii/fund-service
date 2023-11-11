package validator

import (
	"backend-crowdfunding/sdk/errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type XValidator struct {
	validator *validator.Validate
}

var validate = validator.New()

func Validate(data interface{}) errors.ErrorMap {
	validationErrors := make(errors.ErrorMap)

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			validationErrors[err.Field()] = fmt.Sprintf("The %s field must implements %s", err.Field(), err.ActualTag())
		}
		return validationErrors
	}
	return nil
}
