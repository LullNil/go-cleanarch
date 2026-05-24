package entity1

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/LullNil/go-cleanarch/domain"
	domainentity1 "github.com/LullNil/go-cleanarch/domain/entity1"
)

// Service implements entity1 use cases.
type Service struct {
	repo  domainentity1.Repository
	cache Cache
}

// New creates a new entity1 service.
func New(repo domainentity1.Repository, cache Cache) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

// CreateEntity1 creates a new entity1.
func (s *Service) CreateEntity1(ctx context.Context, cmd *CreateCommand) (int64, error) {
	const op = "service.entity1.CreateEntity1"

	if cmd == nil || strings.TrimSpace(cmd.Field3) == "" {
		return 0, fmt.Errorf("%s: %w", op, domain.ErrInvalidInput)
	}

	e := &domainentity1.Entity1{
		Field1: cmd.Field1,
		Field2: cmd.Field2,
		Field3: cmd.Field3,
	}

	id, err := s.repo.Save(ctx, e)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to create entity1: %w", op, err)
	}

	e.ID = id
	s.setCache(ctx, e)

	return id, nil
}

// UpdateEntity1 updates an existing entity1.
func (s *Service) UpdateEntity1(ctx context.Context, cmd *UpdateCommand) error {
	const op = "service.entity1.UpdateEntity1"

	if cmd == nil || cmd.ID <= 0 || strings.TrimSpace(cmd.Field3) == "" {
		return fmt.Errorf("%s: %w", op, domain.ErrInvalidInput)
	}

	e, err := s.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("%s: failed to get entity1: %w", op, err)
	}

	e.Field3 = cmd.Field3
	if err := s.repo.Update(ctx, e); err != nil {
		return fmt.Errorf("%s: failed to update entity1: %w", op, err)
	}

	s.setCache(ctx, e)

	return nil
}

// GetEntity1Details returns the details of an entity1.
func (s *Service) GetEntity1Details(ctx context.Context, id int64) (*domainentity1.Entity1, error) {
	const op = "service.entity1.GetEntity1Details"

	if id <= 0 {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrInvalidInput)
	}

	if s.cache != nil {
		e, err := s.cache.Get(ctx, id)
		if err == nil {
			return e, nil
		}
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			slog.WarnContext(ctx, "failed to get entity1 from cache", "error", err, "id", id)
		}
	}

	e, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.setCache(ctx, e)

	return e, nil
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

	if s.cache != nil {
		if err := s.cache.Delete(ctx, id); err != nil {
			slog.WarnContext(ctx, "failed to delete entity1 from cache", "error", err, "id", id)
		}
	}

	return nil
}

func (s *Service) setCache(ctx context.Context, e *domainentity1.Entity1) {
	if s.cache == nil {
		return
	}

	if err := s.cache.Set(ctx, e); err != nil {
		slog.WarnContext(ctx, "failed to set entity1 cache", "error", err, "id", e.ID)
	}
}
