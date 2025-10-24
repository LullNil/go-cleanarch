package entity1

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/LullNil/go-cleanarch/domain/entity1"
)

type Handler struct {
	service entity1.Service
	log     *slog.Logger
}

// NewHandler creates a new entity1 HTTP handler.
func NewHandler(service entity1.Service, log *slog.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

// writeJSON is a helper function for sending a JSON response.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// writeError is a helper function for sending errors.
func writeError(w http.ResponseWriter, log *slog.Logger, msg string, err error, status int) {
	log.Error(msg, "error", err)
	writeJSON(w, status, map[string]any{
		"error":   msg,
		"details": err.Error(),
	})
}

// CreateEntity1 creates a new entity1.
func (h *Handler) CreateEntity1(w http.ResponseWriter, r *http.Request) {
	const op = "http.entity1.CreateEntity1"

	var req entity1.CreateEntity1Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, h.log, op+": failed to decode request", err, http.StatusBadRequest)
		return
	}

	// Simple validation
	if strings.TrimSpace(req.Field3) == "" {
		writeError(w, h.log, op+": invalid request", errors.New("field1 is required"), http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateEntity1(r.Context(), req)
	if err != nil {
		writeError(w, h.log, op+": failed to create entity1", err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"id": id})
}

// UpdateEntity1 updates an existing entity1.
func (h *Handler) UpdateEntity1(w http.ResponseWriter, r *http.Request) {
	const op = "http.entity1.UpdateEntity1"

	var req entity1.UpdateEntity1Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, h.log, op+": failed to decode request", err, http.StatusBadRequest)
		return
	}

	if req.ID == 0 {
		writeError(w, h.log, op+": invalid request", errors.New("id is required"), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateEntity1(r.Context(), req); err != nil {
		writeError(w, h.log, op+": failed to update entity1", err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// GetEntity1Details gets entity1 details by ID.
func (h *Handler) GetEntity1Details(w http.ResponseWriter, r *http.Request) {
	const op = "http.entity1.GetEntity1Details"

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, h.log, op+": missing id", errors.New("id query param is required"), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, h.log, op+": invalid id", err, http.StatusBadRequest)
		return
	}

	entity, err := h.service.GetEntity1Details(r.Context(), id)
	if err != nil {
		writeError(w, h.log, op+": failed to get entity1", err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, entity)
}

// DeleteEntity1 deletes an entity1 by ID.
func (h *Handler) DeleteEntity1(w http.ResponseWriter, r *http.Request) {
	const op = "http.entity1.DeleteEntity1"

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, h.log, op+": missing id", errors.New("id query param is required"), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, h.log, op+": invalid id", err, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteEntity1(r.Context(), id); err != nil {
		writeError(w, h.log, op+": failed to delete entity1", err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
