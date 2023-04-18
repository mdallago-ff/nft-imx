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

	imx := imx.NewIMX(config.AlchemyAPIKey, config.L1SignerPrivateKey, config.StarkPrivateKey, config.ProjectID)
	defer imx.Close()

	server := server.NewServer(config, db, imx)
	server.Configure()

	err = http.ListenAndServe(":"+config.Port, server.Router)
	if err != nil {
		log.Fatal("error stopping web server", err)
	}
}
