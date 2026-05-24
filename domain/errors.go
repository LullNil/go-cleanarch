package domain

import "errors"

var (
	// ErrInternalServerError will throw when an unexpected error occurs
	ErrInternalServerError = errors.New("internal server error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("item not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("item already exists")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given param is not valid")
)
