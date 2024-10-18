package api

import (
	"github.com/rs/zerolog/log"

	"github.com/dubass83/go-micro-logger/data"
	"github.com/dubass83/go-micro-logger/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Server struct {
	Router  *chi.Mux
	MongoDB *mongo.Client
	Config  util.Config
}

func CreateNewServer(config util.Config) *Server {
	client, err := data.New(config)
	if err != nil {
		log.Fatal().Err(err).Msg("mongodb connection failed")
	}
	s := &Server{
		Router:  chi.NewRouter(),
		MongoDB: client,
		Config:  config,
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
