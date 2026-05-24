package httperror

import (
	"net/http"

	"github.com/LullNil/go-cleanarch/internal/apperr"
)

// StatusCode maps an application error to an HTTP status code.
func StatusCode(err error) int {
	switch apperr.CodeOf(err) {
	case apperr.CodeInvalidArgument:
		return http.StatusBadRequest
	case apperr.CodeNotFound:
		return http.StatusNotFound
	case apperr.CodeAlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// Message maps an application error to a public response message.
func Message(err error) string {
	return apperr.PublicMessage(err)
}
