package httpserver

import (
	"context"
	"log/slog"
	stdhttp "net/http"
	"time"

	"github.com/LullNil/go-cleanarch/config"
	entity1http "github.com/LullNil/go-cleanarch/internal/delivery/http/entity1"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// Server wraps the HTTP server lifecycle.
type Server struct {
	log    *slog.Logger
	server *stdhttp.Server
}

// NewServer creates a new HTTP server.
func NewServer(cfg config.HTTPServer, log *slog.Logger, entity1Service entity1http.Service) *Server {
	router := newRouter(cfg, log, entity1Service)

	return &Server{
		log: log,
		server: &stdhttp.Server{
			Addr:         cfg.Port,
			Handler:      router,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

// Run starts the HTTP server.
func (s *Server) Run() error {
	s.log.Info("starting http server", slog.String("addr", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && err != stdhttp.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown gracefully shuts down the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// newRouter creates a new HTTP router
func newRouter(cfg config.HTTPServer, log *slog.Logger, entity1Service entity1http.Service) stdhttp.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.New(corsOptions(cfg.CORS)).Handler)

	entity1Handler := entity1http.NewHandler(entity1Service, log)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/entity1", func(r chi.Router) {
			r.Post("/", entity1Handler.CreateEntity1)
			r.Get("/{id}", entity1Handler.GetEntity1Details)
			r.Put("/{id}", entity1Handler.UpdateEntity1)
			r.Delete("/{id}", entity1Handler.DeleteEntity1)
		})
	})

	if cfg.EnableSwagger {
		registerSwaggerRoutes(r)
	}

	return r
}

func corsOptions(cfg config.CORS) cors.Options {
	if len(cfg.AllowedOrigins) == 0 {
		cfg.AllowedOrigins = []string{"*"}
	}
	if len(cfg.AllowedMethods) == 0 {
		cfg.AllowedMethods = []string{
			stdhttp.MethodGet,
			stdhttp.MethodPost,
			stdhttp.MethodPatch,
			stdhttp.MethodPut,
			stdhttp.MethodDelete,
			stdhttp.MethodOptions,
		}
	}
	if len(cfg.AllowedHeaders) == 0 {
		cfg.AllowedHeaders = []string{
			"Content-Type",
			"Authorization",
		}
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 5 * time.Minute
	}

	return cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   cfg.AllowedMethods,
		AllowedHeaders:   cfg.AllowedHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           int(cfg.MaxAge.Seconds()),
	}
}
