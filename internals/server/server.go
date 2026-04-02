// Package server - has the router and router configuration
package server

import (
	"fmt"
	"net/http"

	"github.com/DevitoDbug/portfolio/internals/api"
	"github.com/DevitoDbug/portfolio/internals/config"
	"github.com/DevitoDbug/portfolio/internals/utils"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	Port   string
	router *chi.Mux
	Api    *api.Api
}

func NewServer(port string) *Server {
	router := chi.NewRouter()
	api := api.NewAPI()

	return &Server{
		Port:   port,
		router: router,
		Api:    api,
	}
}

func (s *Server) StartServer() error {
	environment, err := config.GetEnvironmentConfig()
	if err != nil {
		return err
	}

	allowedOrigins := utils.GetAllowedOrigins(environment.RunningEnvironment)
	if len(allowedOrigins) == 0 {
		return fmt.Errorf("allowed origins not found")
	}

	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	fileServer := http.FileServer(http.Dir("./internals/web/static"))
	s.router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	s.router.Group(func(r chi.Router) {
		r.Get("/", s.Api.IndexHandler)
		r.Get("/about", s.Api.AboutHandler)
	})

	return http.ListenAndServe(s.Port, s.router)
}
