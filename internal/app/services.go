package app

import (
	entity1repo "github.com/LullNil/go-cleanarch/internal/repository/postgres"
	entity1service "github.com/LullNil/go-cleanarch/internal/service/entity1"
)

// Services contains application use case services.
type Services struct {
	Entity1 *entity1service.Service
}

// initServices initializes all services.
func initServices(m *Modules) *Services {
	// Init repositories
	entity1Repo := entity1repo.NewEntity1Repository(m.DB)

	// Init services
	entity1Svc := entity1service.New(entity1Repo)

	return &Services{
		Entity1: entity1Svc,
	}
}
