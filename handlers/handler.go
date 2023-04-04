package handlers

import (
	"nft/db"
	"nft/imx"
)

type Handler struct {
	db  *db.DB
	imx *imx.IMX
}

func NewHandler(db *db.DB, imx *imx.IMX) *Handler {
	return &Handler{db, imx}
}
