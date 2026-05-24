package httputil

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, log *slog.Logger, message string, err error, status int) {
	log.Error(message, "error", err)
	WriteJSON(w, status, map[string]any{
		"error": message,
	})
}
