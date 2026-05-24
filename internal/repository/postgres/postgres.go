package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/LullNil/go-cleanarch/config"

	_ "github.com/lib/pq"
)

// ConnectToDBWithRetries tries to connect to the database with retries.
func ConnectWithRetries(ctx context.Context, pgConfig config.Postgres, log *slog.Logger) (*sql.DB, error) {
	const op = "postgres.ConnectToDBWithRetries"

	ctx, cancel := context.WithTimeout(ctx, pgConfig.ConnectTimeout)
	defer cancel()

	var db *sql.DB
	var err error

	for i := 0; i < pgConfig.MaxRetries; i++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("%s: context cancelled or timed out before successful connection: %w", op, ctx.Err())
		default:
			log.Debug("Attempting to connect to PostgreSQL...", slog.Int("attempt", i+1), slog.Int("max_attempts", pgConfig.MaxRetries))

			db, err = sql.Open("postgres", pgConfig.DSN)
			if err != nil {
				log.Warn("Failed to open database connection immediately, retrying...", slog.String("err", err.Error()))
				time.Sleep(pgConfig.RetryInterval)
				continue
			}

			// Check database connection
			pingCtx, pingCancel := context.WithTimeout(ctx, pgConfig.RetryInterval)
			err = db.PingContext(pingCtx)
			pingCancel()

			if err == nil {
				log.Info("Successfully connected to PostgreSQL")
				return db, nil
			}

			log.Warn("Failed to ping PostgreSQL, retrying...", slog.String("err", err.Error()))
			if closeErr := db.Close(); closeErr != nil {
				log.Error("Failed to close dangling db connection during retry", slog.String("err", closeErr.Error()))
			}
			time.Sleep(pgConfig.RetryInterval)
		}
	}

	return nil, fmt.Errorf("%s: failed to connect to PostgreSQL after %d retries and %s timeout: %w", op, pgConfig.MaxRetries, pgConfig.ConnectTimeout, err)
}
