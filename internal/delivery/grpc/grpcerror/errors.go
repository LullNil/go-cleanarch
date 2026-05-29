package grpcerror

import (
	"github.com/LullNil/go-cleanarch/internal/apperr"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Code maps an application error to a gRPC code.
func Code(err error) codes.Code {
	switch apperr.CodeOf(err) {
	case apperr.CodeInvalidArgument:
		return codes.InvalidArgument
	case apperr.CodeNotFound:
		return codes.NotFound
	case apperr.CodeAlreadyExists:
		return codes.AlreadyExists
	case apperr.CodePermissionDenied:
		return codes.PermissionDenied
	default:
		return codes.Internal
	}
}

// Status maps an application error to a gRPC status error.
func Status(err error) error {
	return status.Error(Code(err), apperr.PublicMessage(err))
}
