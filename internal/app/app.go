package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/LullNil/go-cleanarch/config"
	grpcserver "github.com/LullNil/go-cleanarch/internal/delivery/grpc"
	httpserver "github.com/LullNil/go-cleanarch/internal/delivery/http"
	"github.com/LullNil/go-cleanarch/internal/lib/logger"
)

// App contains application dependencies and lifecycle.
type App struct {
	cfg          *config.Config
	log          *slog.Logger
	modules      *Modules
	integrations *Integrations
	services     *Services
	http         *httpserver.Server
	grpc         *grpcserver.Server
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

	// External integrations
	integrations, err := initIntegrations(cfg, modules, log)
	if err != nil {
		modules.Close(log)
		return nil, err
	}

	// Services
	services := initServices(cfg, modules, integrations)

	httpServer := httpserver.NewServer(cfg.HTTPServer, log, services.Entity1)
	grpcServer := grpcserver.NewServer(cfg.GRPCServer, log, services.Entity1)

	return &App{
		cfg:          cfg,
		log:          log,
		modules:      modules,
		integrations: integrations,
		services:     services,
		http:         httpServer,
		grpc:         grpcServer,
	}, nil
}

// run starts all processes and performs graceful shutdown
func (a *App) run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 2)

	// Start the HTTP server
	go func() {
		if err := a.http.Run(); err != nil {
			errCh <- err
		}
	}()

	// Start the gRPC server
	go func() {
		if err := a.grpc.Run(); err != nil {
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
	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
	defer cancel()

	if err := a.http.Shutdown(shutdownCtx); err != nil {
		a.log.Error("failed to gracefully shutdown http server", slog.String("error", err.Error()))
	} else {
		a.log.Info("http server stopped gracefully")
	}

	if err := a.grpc.Shutdown(shutdownCtx); err != nil {
		a.log.Error("failed to gracefully shutdown grpc server", slog.String("error", err.Error()))
	} else {
		a.log.Info("grpc server stopped gracefully")
	}

	// Close external integrations and modules
	a.integrations.Close(a.log)
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
