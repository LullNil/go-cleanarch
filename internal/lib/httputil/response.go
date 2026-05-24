package httputil

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	chimiddleware "github.com/go-chi/chi/middleware"
)

type errorResponse struct {
	Error string `json:"error"`
}

// WriteJSON writes a JSON response.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// WriteError writes an error response.
func WriteError(w http.ResponseWriter, log *slog.Logger, message string, err error, status int) {
	log.Error(message, "error", err)
	WriteJSON(w, status, &errorResponse{Error: message})
}

// WriteRequestError writes an error response and logs request correlation fields.
func WriteRequestError(w http.ResponseWriter, r *http.Request, log *slog.Logger, message string, err error, status int) {
	log.Error(message,
		"error", err,
		"request_id", RequestID(r),
		"trace_id", TraceID(r),
	)
	WriteJSON(w, status, &errorResponse{Error: message})
}

// RequestID returns the request id from chi middleware or request headers.
func RequestID(r *http.Request) string {
	if requestID := chimiddleware.GetReqID(r.Context()); requestID != "" {
		return requestID
	}
	return r.Header.Get("X-Request-Id")
}

// TraceID returns the trace id from request headers.
func TraceID(r *http.Request) string {
	if traceID := r.Header.Get("X-Trace-Id"); traceID != "" {
		return traceID
	}

	traceparent := r.Header.Get("Traceparent")
	parts := strings.Split(traceparent, "-")
	if len(parts) >= 2 {
		return parts[1]
	}

	return ""
}
