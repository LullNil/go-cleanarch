package app

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/LullNil/go-cleanarch/config"
	"github.com/LullNil/go-cleanarch/internal/repository/postgres"
	redisrepo "github.com/LullNil/go-cleanarch/internal/repository/redis"

	goredis "github.com/redis/go-redis/v9"
)

// Modules contains external resources used by the application.
type Modules struct {
	DB    *sql.DB
	Redis *goredis.Client
	// Kafka *kafka.Writer
	// Metrics *prometheus.Registry
}

// initModules initializes all external resources.
func initModules(cfg *config.Config, log *slog.Logger) (*Modules, error) {
	log.Info("connecting external modules...")

	db, err := postgres.ConnectWithRetries(context.Background(), cfg.Postgres, log)
	if err != nil {
		return nil, err
	}

	redisClient, err := redisrepo.ConnectWithRetries(context.Background(), cfg.Redis, log)
	if err != nil {
		return nil, err
	}

	// kafkaWriter := kafka.NewWriter(...)
	// metrics := prometheus.NewRegistry()

	return &Modules{
		DB:    db,
		Redis: redisClient,
		// Kafka: kafkaWriter,
		// Metrics: metrics,
	}, nil
}

// Close closes all external resources.
func (m *Modules) Close(log *slog.Logger) {
	if m.DB != nil {
		log.Debug("closing postgres connection...")
		_ = m.DB.Close()
	}
	if m.Redis != nil {
		log.Debug("closing redis connection...")
		_ = m.Redis.Close()
	}
	// if m.Kafka != nil { m.Kafka.Close() }
}
