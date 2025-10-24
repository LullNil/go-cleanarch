package entity1

import "context"

type Service interface {
	CreateEntity1(ctx context.Context, req CreateEntity1Request) (int64, error)
	UpdateEntity1(ctx context.Context, req UpdateEntity1Request) error
	GetEntity1Details(ctx context.Context, id int64) (*Entity1, error)
	DeleteEntity1(ctx context.Context, id int64) error
}

// Examples of DTO queries (use case input structs)

type CreateEntity1Request struct {
	Field1 bool   `json:"field1"`
	Field2 int64  `json:"field2"`
	Field3 string `json:"field3"`
}

type UpdateEntity1Request struct {
	ID     int64  `json:"id"`
	Field3 string `json:"field3"`
}
