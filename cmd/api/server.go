package api

import (
	"github.com/dubass83/go-micro-logger/data"
	"github.com/dubass83/go-micro-logger/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	Router     *chi.Mux
	LogStorage data.LogStorage
	Config     util.Config
}

func CreateNewServer(config util.Config, logStorage data.LogStorage) *Server {

	s := &Server{
		Router:     chi.NewRouter(),
		LogStorage: logStorage,
		Config:     config,
	}
	return s
}

func (s *Server) ConfigureCORS() {
	// Configure CORS
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
}

func (s *Server) AddMiddleware() {
	// Mount all Middleware here
	s.Router.Use(middleware.Heartbeat("/ping"))
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
}
