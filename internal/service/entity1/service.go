package entity1

import (
	"context"
	"fmt"
	"strings"

	"github.com/LullNil/go-cleanarch/domain"
	domainentity1 "github.com/LullNil/go-cleanarch/domain/entity1"
)

type Service struct {
	repo domainentity1.Repository
}

// New creates a new entity1 service.
func New(repo domainentity1.Repository) *Service {
	return &Service{repo: repo}
}

// CreateEntity1 creates a new entity1.
func (s *Service) CreateEntity1(ctx context.Context, req *CreateRequest) (int64, error) {
	const op = "service.entity1.CreateEntity1"

	if req == nil || strings.TrimSpace(req.Field3) == "" {
		return 0, fmt.Errorf("%s: %w", op, domain.ErrInvalidInput)
	}

	e := &domainentity1.Entity1{
		Field1: req.Field1,
		Field2: req.Field2,
		Field3: req.Field3,
	}

	id, err := s.repo.Save(ctx, e)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to create entity1: %w", op, err)
	}

	return id, nil
}

// UpdateEntity1 updates an existing entity1.
func (s *Service) UpdateEntity1(ctx context.Context, req *UpdateRequest) error {
	const op = "service.entity1.UpdateEntity1"

	if req == nil || req.ID <= 0 || strings.TrimSpace(req.Field3) == "" {
		return fmt.Errorf("%s: %w", op, domain.ErrInvalidInput)
	}

	e, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("%s: failed to get entity1: %w", op, err)
	}

	e.Field3 = req.Field3
	if err := s.repo.Update(ctx, e); err != nil {
		return fmt.Errorf("%s: failed to update entity1: %w", op, err)
	}

	return nil
}

// GetEntity1Details returns the details of an entity1.
func (s *Service) GetEntity1Details(ctx context.Context, id int64) (*domainentity1.Entity1, error) {
	const op = "service.entity1.GetEntity1Details"

	if id <= 0 {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

// DeleteEntity1 deletes an entity1.
func (s *Service) DeleteEntity1(ctx context.Context, id int64) error {
	const op = "service.entity1.DeleteEntity1"

	if id <= 0 {
		return fmt.Errorf("%s: %w", op, domain.ErrInvalidInput)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: failed to delete entity1: %w", op, err)
	}

	return nil
}
