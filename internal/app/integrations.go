package app

import (
	"log/slog"
	"strings"

	"github.com/LullNil/go-cleanarch/config"
	authintegration "github.com/LullNil/go-cleanarch/internal/integration/auth"
	entity1service "github.com/LullNil/go-cleanarch/internal/service/entity1"
)

// Integrations contains external service clients used by application services.
type Integrations struct {
	Entity1Auth entity1service.AuthClient

	authClient *authintegration.Client
}

// initIntegrations initializes external service clients
func initIntegrations(cfg *config.Config, _ *Modules, log *slog.Logger) (*Integrations, error) {
	log.Info("connecting external integrations...")

	integrations := &Integrations{}

	if strings.TrimSpace(cfg.Integrations.Auth.GRPCTarget) != "" {
		authClient, err := authintegration.NewGRPCClient(cfg.Integrations.Auth)
		if err != nil {
			return nil, err
		}

		integrations.Entity1Auth = authClient
		integrations.authClient = authClient
	}

	return integrations, nil
}

// Close closes all external service clients.
func (i *Integrations) Close(log *slog.Logger) {
	if i == nil {
		return
	}

	if i.authClient != nil {
		log.Debug("closing auth integration...")
		_ = i.authClient.Close()
	}
}
