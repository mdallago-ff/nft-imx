package handlers

import (
	"nft/db"
	"nft/imx"
)

type Handler struct {
	db  *db.DB
	imx imx.Client
}

func NewHandler(db *db.DB, imx imx.Client) *Handler {
	return &Handler{db, imx}
}
