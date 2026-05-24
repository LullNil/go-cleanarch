package entity1

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/LullNil/go-cleanarch/domain"
	"github.com/LullNil/go-cleanarch/domain/entity1"
	"github.com/LullNil/go-cleanarch/internal/lib/httputil"
	entity1service "github.com/LullNil/go-cleanarch/internal/service/entity1"

	"github.com/go-chi/chi"
)

type Service interface {
	CreateEntity1(ctx context.Context, req *entity1service.CreateRequest) (int64, error)
	UpdateEntity1(ctx context.Context, req *entity1service.UpdateRequest) error
	GetEntity1Details(ctx context.Context, id int64) (*entity1.Entity1, error)
	DeleteEntity1(ctx context.Context, id int64) error
}

type Handler struct {
	service Service
	log     *slog.Logger
}

// NewHandler creates a new entity1 HTTP handler.
func NewHandler(service Service, log *slog.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

// CreateEntity1 creates a new entity1.
func (h *Handler) CreateEntity1(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, h.log, "invalid request body", err, http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateEntity1(r.Context(), req.toServiceRequest())
	if err != nil {
		httputil.WriteError(w, h.log, messageFromError(err), err, statusFromError(err))
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, &createResponse{ID: id})
}

// UpdateEntity1 updates an existing entity1.
func (h *Handler) UpdateEntity1(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, h.log)
	if !ok {
		return
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, h.log, "invalid request body", err, http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateEntity1(r.Context(), req.toServiceRequest(id)); err != nil {
		httputil.WriteError(w, h.log, messageFromError(err), err, statusFromError(err))
		return
	}

	httputil.WriteJSON(w, http.StatusOK, &updateResponse{Status: "updated"})
}

// GetEntity1Details gets entity1 details by ID.
func (h *Handler) GetEntity1Details(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, h.log)
	if !ok {
		return
	}

	entity, err := h.service.GetEntity1Details(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, h.log, messageFromError(err), err, statusFromError(err))
		return
	}

	httputil.WriteJSON(w, http.StatusOK, toGetResponse(entity))
}

// DeleteEntity1 deletes an entity1 by ID.
func (h *Handler) DeleteEntity1(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, h.log)
	if !ok {
		return
	}

	if err := h.service.DeleteEntity1(r.Context(), id); err != nil {
		httputil.WriteError(w, h.log, messageFromError(err), err, statusFromError(err))
		return
	}

	httputil.WriteJSON(w, http.StatusOK, &deleteResponse{Status: "deleted"})
}

func parseID(w http.ResponseWriter, r *http.Request, log *slog.Logger) (int64, bool) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		httputil.WriteError(w, log, "invalid path parameter", domain.ErrInvalidInput, http.StatusBadRequest)
		return 0, false
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httputil.WriteError(w, log, "invalid path parameter", err, http.StatusBadRequest)
		return 0, false
	}

	return id, true
}

func statusFromError(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, domain.ErrInvalidInput):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func messageFromError(err error) string {
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
