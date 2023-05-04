package main

import (
	"context"
	"log"
	"net/http"
	"nft/config"
	"nft/db"
	"nft/imx"
	"nft/server"
)

func main() {
	log.Println("Starting NFT Marketplace")

	settings := config.GetConfig()

	log.Printf("Port: %s", settings.Port)
	log.Printf("DebugMode: %t", settings.DebugMode)

	migrations := db.NewMigrations(settings.DSN)
	err := migrations.Up(context.TODO())
	if err != nil {
		log.Fatal("Error applying migrations")
	}

	newDB, err := db.NewDB(settings.DSN)
	if err != nil {
		log.Fatal("error configuring DB", err)
	}

	imxClient := imx.NewIMX(settings.AlchemyAPIKey, settings.L1SignerPrivateKey, settings.StarkPrivateKey, settings.ProjectID)
	defer imxClient.Close()

	newServer := server.NewServer(settings, newDB, imxClient)
	newServer.Configure()

	err = http.ListenAndServe(":"+settings.Port, newServer.Router)
	if err != nil {
		log.Fatal("error stopping web server", err)
	}
}
