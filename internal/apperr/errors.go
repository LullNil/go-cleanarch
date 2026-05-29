package apperr

import (
	"errors"
	"fmt"

	"github.com/LullNil/go-cleanarch/domain"
)

// Code identifies an application error category.
type Code string

const (
	// CodeInvalidArgument indicates invalid input.
	CodeInvalidArgument Code = "invalid_argument"

	// CodeNotFound indicates that a requested resource does not exist.
	CodeNotFound Code = "not_found"

	// CodeAlreadyExists indicates that a resource already exists.
	CodeAlreadyExists Code = "already_exists"

	// CodePermissionDenied indicates that the caller is not allowed to perform an action.
	CodePermissionDenied Code = "permission_denied"

	// CodeInternal indicates an unexpected error.
	CodeInternal Code = "internal"
)

// Error wraps an error with an application error code.
type Error struct {
	code Code
	err  error
}

// New creates an application error.
func New(code Code, message string) *Error {
	return &Error{
		code: code,
		err:  errors.New(message),
	}
}

// Wrap wraps an error with an application error code.
func Wrap(code Code, err error) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		code: code,
		err:  err,
	}
}

// Code returns the application error code.
func (e *Error) Code() Code {
	return e.code
}

// Error returns the error message.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.code, e.err)
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	return e.err
}

// CodeOf returns the application error code for err.
func CodeOf(err error) Code {
	if err == nil {
		return ""
	}

	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Code()
	}

	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return CodeInvalidArgument
	case errors.Is(err, domain.ErrNotFound):
		return CodeNotFound
	case errors.Is(err, domain.ErrAlreadyExists):
		return CodeAlreadyExists
	case errors.Is(err, domain.ErrPermissionDenied):
		return CodePermissionDenied
	default:
		return CodeInternal
	}
}

// PublicMessage returns a safe public message for err.
func PublicMessage(err error) string {
	switch CodeOf(err) {
	case CodeInvalidArgument:
		return "invalid request"
	case CodeNotFound:
		return "resource not found"
	case CodeAlreadyExists:
		return "resource already exists"
	case CodePermissionDenied:
		return "permission denied"
	default:
		return "internal server error"
	}
}
