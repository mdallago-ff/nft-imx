package handlers

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/render"
	"net/http"
	"nft/imx"
)

func (h *Handler) CreateToken(w http.ResponseWriter, r *http.Request) {
	data := &TokenRequest{}
	if err := render.Bind(r, data); err != nil {
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	info := imx.MintInformation{
		ContractAddress: data.ContractAddress,
		TokenID:         data.TokenID,
		Blueprint:       data.Blueprint,
	}
	err := h.imx.CreateToken(r.Context(), &info)
	if err != nil {
		log.Error("error creating collection", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, NewTokenResponse())
	if err != nil {
		log.Error("error rendering response", err)
	}
}

type TokenRequest struct {
	ContractAddress string
	TokenID         string
	Blueprint       string
}

func (a *TokenRequest) Bind(r *http.Request) error {

	return nil
}

type TokenResponse struct {
}

func NewTokenResponse() *TokenResponse {
	resp := &TokenResponse{}
	return resp
}

func (rd *TokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
