package errors

import (
	"net/http"
	"reflect"
)

type Errors struct {
	Type    string
	Code    int
	Message string
}

const (
	NotFoundType            = "HTTPStatusNotFound"
	InternalServerErrorType = "HTTPStatusInternalServerError"
	BadRequestType          = "HTTPStatusBadRequest"
)

func (e *Errors) Error() string {
	return e.Message
}

func NewWithCode(code int, message, errType string) error {
	errors := &Errors{
		Type:    errType,
		Code:    code,
		Message: message,
	}

	return errors
}

func NotFound(entity string) error {
	return NewWithCode(http.StatusNotFound, entity+" not found", NotFoundType)
}

func InternalServerError(message string) error {
	return NewWithCode(http.StatusInternalServerError, message, InternalServerErrorType)
}

func GetType(err error) string {
	if err == nil {
		return "HTTPStatusOK"
	}

	if reflect.TypeOf(err).String() == "*errors.Errors" {
		return err.(*Errors).Type
	}

	return InternalServerErrorType
}

func GetCode(err error) int {
	if err == nil {
		return 200
	}

	if reflect.TypeOf(err).String() == "*errors.Errors" {
		return err.(*Errors).Code
	}

	return 500
}

func GetMessage(err error) string {
	if err == nil {
		return "OK"
	}

	if reflect.TypeOf(err).String() == "*errors.Errors" {
		return err.(*Errors).Message
	}

	return err.Error()
}