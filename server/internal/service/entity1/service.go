package entity1

import (
	"context"
	"fmt"

	"github.com/LullNil/go-cleanarch/domain/entity1"
)

type service struct {
	repo entity1.Repository
}

func NewService(repo entity1.Repository) entity1.Service {
	return &service{repo: repo}
}

func (s *service) CreateEntity1(ctx context.Context, req entity1.CreateEntity1Request) (int64, error) {
	const op = "service.entity1.CreateEntity1"

	e := &entity1.Entity1{
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

func (s *service) UpdateEntity1(ctx context.Context, req entity1.UpdateEntity1Request) error {
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

func (s *service) GetEntity1Details(ctx context.Context, id int64) (*entity1.Entity1, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) DeleteEntity1(ctx context.Context, id int64) error {
	return fmt.Errorf("not implemented")
}
