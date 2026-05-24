package entity1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LullNil/go-cleanarch/domain/entity1"
	"github.com/LullNil/go-cleanarch/internal/delivery/http/httperror"
	"github.com/LullNil/go-cleanarch/internal/lib/httputil"
	entity1service "github.com/LullNil/go-cleanarch/internal/service/entity1"
)

type Service interface {
	CreateEntity1(ctx context.Context, cmd *entity1service.CreateCommand) (int64, error)
	UpdateEntity1(ctx context.Context, cmd *entity1service.UpdateCommand) error
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
		httputil.WriteRequestError(w, r, h.log, "invalid request body", err, http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateEntity1(r.Context(), &entity1service.CreateCommand{
		Field1: req.Field1,
		Field2: req.Field2,
		Field3: req.Field3,
	})
	if err != nil {
		httputil.WriteRequestError(w, r, h.log, httperror.Message(err), err, httperror.StatusCode(err))
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, &createResponse{ID: id})
}

// UpdateEntity1 updates an existing entity1.
func (h *Handler) UpdateEntity1(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.PathInt64(w, r, h.log, "id")
	if !ok {
		return
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteRequestError(w, r, h.log, "invalid request body", err, http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateEntity1(r.Context(), &entity1service.UpdateCommand{
		ID:     id,
		Field3: req.Field3,
	}); err != nil {
		httputil.WriteRequestError(w, r, h.log, httperror.Message(err), err, httperror.StatusCode(err))
		return
	}

	httputil.WriteJSON(w, http.StatusOK, &updateResponse{Status: "updated"})
}

// GetEntity1Details gets entity1 details by ID.
func (h *Handler) GetEntity1Details(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.PathInt64(w, r, h.log, "id")
	if !ok {
		return
	}

	entity, err := h.service.GetEntity1Details(r.Context(), id)
	if err != nil {
		httputil.WriteRequestError(w, r, h.log, httperror.Message(err), err, httperror.StatusCode(err))
		return
	}

	httputil.WriteJSON(w, http.StatusOK, newGetResponse(entity))
}

// DeleteEntity1 deletes an entity1 by ID.
func (h *Handler) DeleteEntity1(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.PathInt64(w, r, h.log, "id")
	if !ok {
		return
	}

	if err := h.service.DeleteEntity1(r.Context(), id); err != nil {
		httputil.WriteRequestError(w, r, h.log, httperror.Message(err), err, httperror.StatusCode(err))
		return
	}

	httputil.WriteJSON(w, http.StatusOK, &deleteResponse{Status: "deleted"})
}
