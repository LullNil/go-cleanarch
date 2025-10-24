package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LullNil/go-cleanarch/config"
	"github.com/LullNil/go-cleanarch/internal/lib/logger"
)

type App struct {
	log      *slog.Logger
	modules  *Modules
	services *Services
	server   *http.Server
}

// Run - entrypoint of the application.
func Run(cfg *config.Config) error {
	app, err := newApp(cfg)
	if err != nil {
		return err
	}
	return app.run()
}

// newApp adjusts all dependencies.
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
	router := initRouter(log, services)
	server := &http.Server{
		Addr:         cfg.HTTPServer.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
	}

	return &App{
		log:      log,
		modules:  modules,
		services: services,
		server:   server,
	}, nil
}

// run starts all processes and performs graceful shutdown.
func (a *App) run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)

	go func() {
		a.log.Info("starting http server", slog.String("addr", a.server.Addr))
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

// setupLogger configures the logger based on the environment.
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
