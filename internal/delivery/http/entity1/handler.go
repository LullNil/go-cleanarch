package entity1

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/LullNil/go-cleanarch/domain"
	"github.com/LullNil/go-cleanarch/domain/entity1"
	"github.com/LullNil/go-cleanarch/internal/lib/httputil"
	entity1service "github.com/LullNil/go-cleanarch/internal/service/entity1"

	"github.com/go-chi/chi"
)

type service interface {
	CreateEntity1(ctx context.Context, req entity1service.CreateRequest) (int64, error)
	UpdateEntity1(ctx context.Context, req entity1service.UpdateRequest) error
	GetEntity1Details(ctx context.Context, id int64) (*entity1.Entity1, error)
	DeleteEntity1(ctx context.Context, id int64) error
}

type Handler struct {
	service service
	log     *slog.Logger
}

// NewHandler creates a new entity1 HTTP handler.
func NewHandler(service service, log *slog.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

// CreateEntity1 creates a new entity1.
func (h *Handler) CreateEntity1(w http.ResponseWriter, r *http.Request) {
	const op = "http.entity1.CreateEntity1"

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, h.log, op+": failed to decode request", err, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Field3) == "" {
		httputil.WriteError(w, h.log, op+": invalid request", errors.New("field3 is required"), http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateEntity1(r.Context(), entity1service.CreateRequest{
		Field1: req.Field1,
		Field2: req.Field2,
		Field3: req.Field3,
	})
	if err != nil {
		httputil.WriteError(w, h.log, op+": failed to create entity1", err, statusFromError(err))
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, map[string]any{"id": id})
}

// UpdateEntity1 updates an existing entity1.
func (h *Handler) UpdateEntity1(w http.ResponseWriter, r *http.Request) {
	const op = "http.entity1.UpdateEntity1"

	id, ok := parseID(w, r, h.log, op)
	if !ok {
		return
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, h.log, op+": failed to decode request", err, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Field3) == "" {
		httputil.WriteError(w, h.log, op+": invalid request", errors.New("field3 is required"), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateEntity1(r.Context(), entity1service.UpdateRequest{
		ID:     id,
		Field3: req.Field3,
	}); err != nil {
		httputil.WriteError(w, h.log, op+": failed to update entity1", err, statusFromError(err))
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// GetEntity1Details gets entity1 details by ID.
func (h *Handler) GetEntity1Details(w http.ResponseWriter, r *http.Request) {
	const op = "http.entity1.GetEntity1Details"

	id, ok := parseID(w, r, h.log, op)
	if !ok {
		return
	}

	entity, err := h.service.GetEntity1Details(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, h.log, op+": failed to get entity1", err, statusFromError(err))
		return
	}

	httputil.WriteJSON(w, http.StatusOK, entity)
}

// DeleteEntity1 deletes an entity1 by ID.
func (h *Handler) DeleteEntity1(w http.ResponseWriter, r *http.Request) {
	const op = "http.entity1.DeleteEntity1"

	id, ok := parseID(w, r, h.log, op)
	if !ok {
		return
	}

	if err := h.service.DeleteEntity1(r.Context(), id); err != nil {
		httputil.WriteError(w, h.log, op+": failed to delete entity1", err, statusFromError(err))
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func parseID(w http.ResponseWriter, r *http.Request, log *slog.Logger, op string) (int64, bool) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		httputil.WriteError(w, log, op+": missing id", errors.New("id path param is required"), http.StatusBadRequest)
		return 0, false
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httputil.WriteError(w, log, op+": invalid id", err, http.StatusBadRequest)
		return 0, false
	}

	return id, true
}

func statusFromError(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
