package entity1

import "context"

// AuthClient defines the external authorization port used by entity1 use cases.
type AuthClient interface {
	CanAccessEntity1(ctx context.Context, subjectID string, entity1ID int64) (bool, error)
}
