package entity1

import (
	"context"

	domainentity1 "github.com/LullNil/go-cleanarch/domain/entity1"
)

// Cache defines the entity1 cache port.
type Cache interface {
	Get(ctx context.Context, id int64) (*domainentity1.Entity1, error)
	Set(ctx context.Context, e *domainentity1.Entity1) error
	Delete(ctx context.Context, id int64) error
}
