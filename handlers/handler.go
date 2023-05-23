package handlers

import (
	"nft/db"
	"nft/imx"

	"github.com/hibiken/asynq"
)

type Handler struct {
	db          *db.DB
	imx         imx.Client
	asynqClient *asynq.Client
}

func NewHandler(db *db.DB, imx imx.Client, asynqClient *asynq.Client) *Handler {
	return &Handler{db, imx, asynqClient}
}
