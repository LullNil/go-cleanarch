package entity1

import (
	"context"
)

type Saver interface {
	Save(ctx context.Context, e *Entity1) (int64, error)
}

type Getter interface {
	GetByID(ctx context.Context, id int64) (*Entity1, error)
}

type Updater interface {
	Update(ctx context.Context, e *Entity1) error
}

type Repository interface {
	Saver
	Getter
	Updater
}
