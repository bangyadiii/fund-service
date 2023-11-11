package errors

import (
	"fmt"
	"reflect"
)

type Error struct {
	code    ErrorCode
	message string
	orig    error
	errors  ErrorMap
}

type ErrorCode int
type ErrorMap map[string]string

const (
	NotFoundType            ErrorCode = 404
	InternalServerErrorType ErrorCode = 500
	BadRequestType          ErrorCode = 400
	ValidationErrorType     ErrorCode = 400
	AuthorizationErrorType  ErrorCode = 403
	AuthenticationErrorType ErrorCode = 401
	ConflictErrorType       ErrorCode = 409
)

func (err *Error) Errors() ErrorMap {
	if err.orig == nil {
		return nil
	}

	return err.errors
}

func (err *Error) Error() string {
	if err.orig != nil {
		return fmt.Sprintf("%s: %v", err.message, err.orig)
	}
	return err.message
}

func WrapErrorf(orig error, code ErrorCode, errors ErrorMap, format string, a ...interface{}) error {
	return &Error{
		code:    code,
		message: fmt.Sprintf(format, a...),
		orig:    orig,
		errors:  errors,
	}
}

func NewErrorf(code ErrorCode, errors ErrorMap, format string, a ...interface{}) error {
	return WrapErrorf(nil, code, errors, format, a...)
}

// UnWarp returns the wrapped error, if any
func (err Error) UnWarp() error {
	return err.orig
}

func NewWithCode(code ErrorCode, message string) error {
	err := &Error{
		code:    code,
		message: message,
	}

	return err
}

//func NotFound(entity string) error {
//	return NewWithCode(http.StatusNotFound, entity+" not found", NotFoundType)
//}
//
//func InternalServerError(message string) error {
//	return NewWithCode(http.StatusInternalServerError, message, InternalServerErrorType)
//}

//func AuthorizationError(message string) error {
//	return NewWithCode(http.StatusUnauthorized, message, ValidationErrorType)
//}
//
//func AuthenticationError(message string) error {
//	return NewWithCode(401, message, AuthenticationErrorType)
//}
//
//func ConflictError(message string) error {
//	return NewWithCode(409, message, ConflictErrorType)
//}

func (err *Error) GetCode() ErrorCode {
	return err.code
}

func GetMessage(err error) string {
	if err == nil {
		return "OK"
	}

	if reflect.TypeOf(err).String() == "errors.Error" {
		return err.(*Error).message
	}

	return err.Error()
}
