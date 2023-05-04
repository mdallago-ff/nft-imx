package server

import (
	"nft/auth"
	"nft/config"
	"nft/db"
	"nft/handlers"
	"nft/imx"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/oauth"
	"github.com/go-chi/render"
)

type Server struct {
	Router *chi.Mux
	config *config.Settings
	db     *db.DB
	imx    imx.Client
}

func NewServer(config *config.Settings, db *db.DB, imx imx.Client) *Server {
	return &Server{chi.NewRouter(), config, db, imx}
}

func (s *Server) Configure() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.URLFormat)
	s.Router.Use(render.SetContentType(render.ContentTypeJSON))

	newHandler := handlers.NewHandler(s.db, s.imx)

	bearerServer := oauth.NewBearerServer(
		s.config.AuthSecret,
		time.Second*time.Duration(s.config.TokenDurationSeconds),
		auth.NewUserVerifier(s.db),
		nil)
	s.Router.Post("/auth", bearerServer.ClientCredentials)

	s.Router.Route("/users", func(r chi.Router) {
		r.Post("/", newHandler.CreateUser)
	})

	s.Router.Group(func(r chi.Router) {
		if !s.config.DebugMode {
			r.Use(oauth.Authorize(s.config.AuthSecret, nil))
		}

		r.Route("/collections", func(r chi.Router) {
			r.Post("/", newHandler.CreateCollection)
		})

		r.Route("/tokens", func(r chi.Router) {
			r.Post("/", newHandler.CreateToken)
		})

		r.Route("/transfers", func(r chi.Router) {
			r.Post("/", newHandler.TransferToken)
		})

		r.Route("/orders", func(r chi.Router) {
			r.Post("/", newHandler.CreateOrder)
		})
	})
}
