package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/oauth"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"nft/auth"
	"nft/config"
	"nft/db"
	"nft/handlers"
	"nft/imx"
	"time"
)

func main() {
	log.Println("Starting NFT Marketplace")

	config := config.GetConfig()

	log.Printf("Port: %s", config.Port)
	log.Printf("DebugMode: %t", config.DebugMode)

	migrations := db.NewMigrations(config.DSN)
	err := migrations.Up(context.TODO())
	if err != nil {
		log.Fatal("Error applying migrations")
	}

	db, err := db.NewDB(config.DSN)
	if err != nil {
		log.Fatal("error configuring DB", err)
	}

	imx := imx.NewIMX(config.AlchemyAPIKey, config.L1SignerPrivateKey, config.StarkPrivateKey)
	defer imx.Close()

	handlers := handlers.NewHandler(db, imx)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	s := oauth.NewBearerServer(
		config.AuthSecret,
		time.Second*120,
		auth.NewUserVerifier(db),
		nil)
	r.Post("/auth", s.ClientCredentials)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", handlers.CreateUser)
	})

	r.Group(func(r chi.Router) {
		if !config.DebugMode {
			r.Use(oauth.Authorize(config.AuthSecret, nil))
		}

		r.Route("/collections", func(r chi.Router) {
			r.Post("/", handlers.CreateCollection)
		})

		r.Route("/tokens", func(r chi.Router) {
			r.Post("/", handlers.CreateToken)
		})
	})

	err = http.ListenAndServe(":"+config.Port, r)
	if err != nil {
		log.Fatal("error stopping web server", err)
	}
}
