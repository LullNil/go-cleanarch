package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LullNil/go-cleanarch/config"
	httpserver "github.com/LullNil/go-cleanarch/internal/delivery/http"
	"github.com/LullNil/go-cleanarch/internal/lib/logger"
)

type App struct {
	log      *slog.Logger
	modules  *Modules
	services *Services
	server   *httpserver.Server
}

// Run is the entrypoint of the application.
func Run(cfg *config.Config) error {
	app, err := newApp(cfg)
	if err != nil {
		return err
	}
	return app.run()
}

// newApp adjusts all dependencies
func newApp(cfg *config.Config) (*App, error) {
	log := setupLogger(cfg.Env)
	log.Info("initializing application...")

	// External modules
	modules, err := initModules(cfg, log)
	if err != nil {
		return nil, err
	}

	// Services
	services := initServices(modules)

	// HTTP server
	server := httpserver.NewServer(cfg.HTTPServer, log, services.Entity1)

	return &App{
		log:      log,
		modules:  modules,
		services: services,
		server:   server,
	}, nil
}

// run starts all processes and performs graceful shutdown
func (a *App) run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)

	// Start the HTTP server
	go func() {
		if err := a.server.Run(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		a.log.Info("shutdown signal received")
	case err := <-errCh:
		a.log.Error("server stopped with error", slog.String("error", err.Error()))
	}

	// Stop the HTTP server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.log.Error("failed to gracefully shutdown http server", slog.String("error", err.Error()))
	} else {
		a.log.Info("http server stopped gracefully")
	}

	// Close external modules
	a.modules.Close(a.log)

	return nil
}

// setupLogger configures the logger based on the environment
func setupLogger(env string) *slog.Logger {
	switch env {
	case "local":
		return slog.New(logger.NewPrettyHandler(os.Stdout, slog.LevelDebug))
	case "prod":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
}
