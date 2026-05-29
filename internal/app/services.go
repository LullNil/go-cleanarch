package app

import (
	"github.com/LullNil/go-cleanarch/config"
	entity1repo "github.com/LullNil/go-cleanarch/internal/repository/postgres"
	entity1cache "github.com/LullNil/go-cleanarch/internal/repository/redis"
	entity1service "github.com/LullNil/go-cleanarch/internal/service/entity1"
)

// Services contains application use case services.
type Services struct {
	Entity1 *entity1service.Service
}

// initServices initializes all services
func initServices(cfg *config.Config, m *Modules, i *Integrations) *Services {
	// Init repositories
	entity1Repo := entity1repo.NewEntity1Repository(m.DB)
	entity1Cache := entity1cache.NewEntity1Cache(m.Redis, cfg.Redis.TTL)

	// Init services
	entity1Svc := entity1service.New(entity1Repo, entity1Cache, i.Entity1Auth)

	return &Services{
		Entity1: entity1Svc,
	}
}
