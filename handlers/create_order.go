package handlers

import (
	"errors"
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/render"
	"net/http"
	"nft/imx"
	"strconv"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	data := &OrderRequest{}
	if err := render.Bind(r, data); err != nil {
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	amount, err := strconv.ParseUint(data.Amount, 10, 64)
	if err != nil {
		log.Error("error creating order", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	info := imx.OrderInformation{
		ContractAddress: data.ContractAddress,
		TokenID:         data.TokenID,
		Amount:          amount,
	}

	err = h.imx.CreateOrder(r.Context(), &info)
	if err != nil {
		log.Error("error creating order", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, NewOrderResponse(data.TokenID))
	if err != nil {
		log.Error("error rendering response", err)
	}
}

type OrderRequest struct {
	ContractAddress string `json:"contract_address"`
	TokenID         string `json:"token_id"`
	Amount          string `json:"amount"`
}

func (a *OrderRequest) Bind(r *http.Request) error {
	if len(a.ContractAddress) == 0 {
		return errors.New("missing required fields")
	}

	if len(a.TokenID) == 0 {
		return errors.New("missing required fields")
	}

	return nil
}

type OrderResponse struct {
	TokenID string `json:"token_id"`
}

func NewOrderResponse(id string) *OrderResponse {
	resp := &OrderResponse{TokenID: id}
	return resp
}

func (rd *OrderResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
