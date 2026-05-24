package app

import (
	"log/slog"
	"net/http"

	entity1http "github.com/LullNil/go-cleanarch/internal/delivery/http/entity1"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// initRouter initializes HTTP routes.
func initRouter(log *slog.Logger, services *Services) http.Handler {
	// Router
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods: []string{
			"GET",
			"POST",
			"PATCH",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
	}).Handler)

	// HTTP handlers
	entity1Handler := entity1http.NewHandler(services.Entity1, log)

	// Routes
	r.Route("/v1", func(r chi.Router) {
		r.Route("/entity1", func(r chi.Router) {
			r.Post("/", entity1Handler.CreateEntity1)
			r.Get("/{id}", entity1Handler.GetEntity1Details)
			r.Put("/{id}", entity1Handler.UpdateEntity1)
			r.Delete("/{id}", entity1Handler.DeleteEntity1)
		})
	})

	return r
}
