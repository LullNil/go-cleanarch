package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/LullNil/go-cleanarch/config"

	goredis "github.com/redis/go-redis/v9"
)

// ConnectWithRetries tries to connect to Redis with retries.
func ConnectWithRetries(ctx context.Context, cfg config.Redis, log *slog.Logger) (*goredis.Client, error) {
	const op = "redis.ConnectWithRetries"

	ctx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	var err error
	for i := 0; i < cfg.MaxRetries; i++ {
		log.Debug("attempting to connect to Redis", slog.Int("attempt", i+1), slog.Int("max_attempts", cfg.MaxRetries))

		if err = client.Ping(ctx).Err(); err == nil {
			log.Info("successfully connected to Redis")
			return client, nil
		}

		select {
		case <-ctx.Done():
			_ = client.Close()
			return nil, fmt.Errorf("%s: context cancelled or timed out before successful connection: %w", op, ctx.Err())
		case <-time.After(cfg.RetryInterval):
		}
	}

	_ = client.Close()
	return nil, fmt.Errorf("%s: failed to connect to Redis after %d retries: %w", op, cfg.MaxRetries, err)
}
