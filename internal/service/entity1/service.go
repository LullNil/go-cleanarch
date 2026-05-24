package entity1

import (
	"context"
	"fmt"

	domainentity1 "github.com/LullNil/go-cleanarch/domain/entity1"
)

type Service struct {
	repo domainentity1.Repository
}

func New(repo domainentity1.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateEntity1(ctx context.Context, req CreateRequest) (int64, error) {
	const op = "service.entity1.CreateEntity1"

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

func (s *Service) UpdateEntity1(ctx context.Context, req UpdateRequest) error {
	const op = "service.entity1.UpdateEntity1"

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

func (s *Service) GetEntity1Details(ctx context.Context, id int64) (*domainentity1.Entity1, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) DeleteEntity1(ctx context.Context, id int64) error {
	const op = "service.entity1.DeleteEntity1"

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: failed to delete entity1: %w", op, err)
	}

	return nil
}
