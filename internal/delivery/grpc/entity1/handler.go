package entity1

import (
	"context"

	domainentity1 "github.com/LullNil/go-cleanarch/domain/entity1"
	"github.com/LullNil/go-cleanarch/internal/delivery/grpc/grpcerror"
	"github.com/LullNil/go-cleanarch/internal/delivery/grpc/pb"
	entity1service "github.com/LullNil/go-cleanarch/internal/service/entity1"
)

// Service defines the entity1 use cases required by the gRPC handler.
type Service interface {
	CreateEntity1(ctx context.Context, req *entity1service.CreateRequest) (int64, error)
	UpdateEntity1(ctx context.Context, req *entity1service.UpdateRequest) error
	GetEntity1Details(ctx context.Context, id int64) (*domainentity1.Entity1, error)
	DeleteEntity1(ctx context.Context, id int64) error
}

// Handler handles entity1 gRPC requests.
type Handler struct {
	pb.UnimplementedEntity1ServiceServer

	service Service
}

// NewHandler creates a new entity1 gRPC handler.
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateEntity1 creates a new entity1.
func (h *Handler) CreateEntity1(ctx context.Context, req *pb.CreateEntity1Request) (*pb.CreateEntity1Response, error) {
	id, err := h.service.CreateEntity1(ctx, &entity1service.CreateRequest{
		Field1: req.GetField1(),
		Field2: req.GetField2(),
		Field3: req.GetField3(),
	})
	if err != nil {
		return nil, grpcerror.Status(err)
	}

	return &pb.CreateEntity1Response{Id: id}, nil
}

// UpdateEntity1 updates an existing entity1.
func (h *Handler) UpdateEntity1(ctx context.Context, req *pb.UpdateEntity1Request) (*pb.UpdateEntity1Response, error) {
	if err := h.service.UpdateEntity1(ctx, &entity1service.UpdateRequest{
		ID:     req.GetId(),
		Field3: req.GetField3(),
	}); err != nil {
		return nil, grpcerror.Status(err)
	}

	return &pb.UpdateEntity1Response{Status: "updated"}, nil
}

// GetEntity1Details gets entity1 details by ID.
func (h *Handler) GetEntity1Details(ctx context.Context, req *pb.GetEntity1DetailsRequest) (*pb.GetEntity1DetailsResponse, error) {
	entity, err := h.service.GetEntity1Details(ctx, req.GetId())
	if err != nil {
		return nil, grpcerror.Status(err)
	}

	return &pb.GetEntity1DetailsResponse{
		Id:     entity.ID,
		Field1: entity.Field1,
		Field2: entity.Field2,
		Field3: entity.Field3,
	}, nil
}

// DeleteEntity1 deletes entity1 by ID.
func (h *Handler) DeleteEntity1(ctx context.Context, req *pb.DeleteEntity1Request) (*pb.DeleteEntity1Response, error) {
	if err := h.service.DeleteEntity1(ctx, req.GetId()); err != nil {
		return nil, grpcerror.Status(err)
	}

	return &pb.DeleteEntity1Response{Status: "deleted"}, nil
}
