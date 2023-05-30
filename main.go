package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"nft/config"
	"nft/db"
	"nft/imx"
	"nft/server"
	"nft/tasks"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
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

	imxClient, err := imx.NewIMX(settings.AlchemyAPIKey, settings.L1SignerPrivateKey, settings.StarkPrivateKey, settings.ProjectID)
	if err != nil {
		log.Fatal("error configuring imx", err)
	}

	defer imxClient.Close()

	asyncClient := asynq.NewClient(asynq.RedisClientOpt{Addr: settings.RedisUrl})

	newServer := server.NewServer(settings, newDB, imxClient, asyncClient)
	newServer.Configure()

	httpServer := &http.Server{Addr: ":" + settings.Port, Handler: newServer.Router}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	asynqServer := asynq.NewServer(
		asynq.RedisClientOpt{Addr: settings.RedisUrl},
		asynq.Config{
			IsFailure: func(err error) bool {
				// this type of error will not count as an error in asynq and will not affect retry count
				// the task will be scheduled again for execution.
				var notReady imx.WithdrawalNotReadyError
				return !errors.As(err, &notReady)
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.Handle(tasks.TypeCompleteWithdrawal, tasks.NewCompleteWithdrawalProcessor(imxClient, newDB))

	if err := asynqServer.Start(mux); err != nil {
		log.Fatalf("could not run asynq server: %v", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancelShutdown := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		asynqServer.Stop()
		asynqServer.Shutdown()

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}

		serverStopCtx()
		cancelShutdown()
	}()

	err = httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("error stopping web server", err)
	}

	<-serverCtx.Done()
}
