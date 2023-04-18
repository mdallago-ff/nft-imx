package handlers

import (
	"errors"
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/render"
	"net/http"
	"nft/imx"
)

func (h *Handler) TransferToken(w http.ResponseWriter, r *http.Request) {
	data := &TransferTokenRequest{}
	if err := render.Bind(r, data); err != nil {
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	info := imx.TransferInformation{
		ContractAddress: data.ContractAddress,
		TokenID:         data.TokenID,
		ReceiverAddress: data.ReceiverAddress,
	}
	err := h.imx.TransferToken(r.Context(), &info)
	if err != nil {
		log.Error("error transfering token", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, NewTransferTokenResponse(data.TokenID))
	if err != nil {
		log.Error("error rendering response", err)
	}
}

type TransferTokenRequest struct {
	ContractAddress string `json:"contract_address"`
	TokenID         string `json:"token_id"`
	ReceiverAddress string `json:"receiver_address"`
}

func (a *TransferTokenRequest) Bind(r *http.Request) error {
	if len(a.ContractAddress) == 0 {
		return errors.New("missing required fields")
	}

	if len(a.TokenID) == 0 {
		return errors.New("missing required fields")
	}

	if len(a.ReceiverAddress) == 0 {
		return errors.New("missing required fields")
	}

	return nil
}

type TransferTokenResponse struct {
	TokenID string `json:"token_id"`
}

func NewTransferTokenResponse(id string) *TransferTokenResponse {
	resp := &TransferTokenResponse{TokenID: id}
	return resp
}

func (rd *TransferTokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
