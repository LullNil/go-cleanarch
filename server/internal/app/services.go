package app

import (
	"github.com/LullNil/go-cleanarch/domain/entity1"
	entity1repo "github.com/LullNil/go-cleanarch/internal/repository/postgres"
	entity1service "github.com/LullNil/go-cleanarch/internal/service/entity1"
)

type Services struct {
	Entity1 entity1.Service
}

// initServices initializes all services.
func initServices(m *Modules) *Services {
	// Init repositories
	entity1Repo := entity1repo.NewEntity1Repository(m.DB)

	// Init services
	entity1Svc := entity1service.NewService(entity1Repo)

	return &Services{
		Entity1: entity1Svc,
	}
}
