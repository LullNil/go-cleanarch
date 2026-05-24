package httperror

import (
	"errors"
	"net/http"

	"github.com/LullNil/go-cleanarch/domain"
)

// StatusCode maps an application error to an HTTP status code.
func StatusCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrAlreadyExists):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// Message maps an application error to a public response message.
func Message(err error) string {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return "invalid request"
	case errors.Is(err, domain.ErrNotFound):
		return "resource not found"
	case errors.Is(err, domain.ErrAlreadyExists):
		return "resource already exists"
	default:
		return "internal server error"
	}
}
