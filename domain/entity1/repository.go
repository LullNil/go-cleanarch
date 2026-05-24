package entity1

import (
	"context"
)

// Saver stores entity1 records.
type Saver interface {
	Save(ctx context.Context, e *Entity1) (int64, error)
}

// Getter reads entity1 records.
type Getter interface {
	GetByID(ctx context.Context, id int64) (*Entity1, error)
}

// Updater updates entity1 records.
type Updater interface {
	Update(ctx context.Context, e *Entity1) error
}

// Deleter deletes entity1 records.
type Deleter interface {
	Delete(ctx context.Context, id int64) error
}

// Repository defines the entity1 persistence port.
type Repository interface {
	Saver
	Getter
	Updater
	Deleter
}
