package handlers

import (
	"errors"
	"net/http"
	"nft/imx"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/oauth"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func (h *Handler) CreateDeposit(w http.ResponseWriter, r *http.Request) {
	data := &DepositRequest{}
	if err := render.Bind(r, data); err != nil {
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	userID, err := uuid.Parse(r.Context().Value(oauth.CredentialContext).(string))
	if err != nil {
		log.Error("error parsing user", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	user, err := h.db.GetUser(userID)
	if err != nil {
		log.Error("error getting user", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	info := imx.CreateDepositInformation{
		DepositAmountWei: data.AmountWei,
		User:             user,
	}

	hash, err := h.imx.CreateEthDeposit(r.Context(), &info)
	if err != nil {
		log.Error("error creating deposit", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, NewDepositResponse(hash))
	if err != nil {
		log.Error("error rendering response", err)
	}
}

type DepositRequest struct {
	AmountWei string `json:"amount_wei"`
}

func (a *DepositRequest) Bind(r *http.Request) error {
	if len(a.AmountWei) == 0 {
		return errors.New("missing required fields")
	}

	return nil
}

type DepositResponse struct {
	TxHash string `json:"tx_hash"`
}

func NewDepositResponse(txHash string) *DepositResponse {
	resp := &DepositResponse{TxHash: txHash}
	return resp
}

func (rd *DepositResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
