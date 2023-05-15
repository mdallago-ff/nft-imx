package handlers

import (
	"errors"
	"net/http"
	"nft/imx"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/oauth"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func (h *Handler) CreateWithdrawal(w http.ResponseWriter, r *http.Request) {
	data := &WithdrawalRequest{}
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

	info := imx.CreateWithdrawalInformation{
		AmountWei: data.AmountWei,
		User:      user,
	}

	withdrawalID, err := h.imx.CreateEthWithdrawal(r.Context(), &info)
	if err != nil {
		log.Error("error creating withdrawal", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, NewWithdrawalResponse(withdrawalID))
	if err != nil {
		log.Error("error rendering response", err)
	}
}

type WithdrawalRequest struct {
	AmountWei string `json:"amount_wei"`
}

func (a *WithdrawalRequest) Bind(r *http.Request) error {
	if len(a.AmountWei) == 0 {
		return errors.New("missing required fields")
	}

	return nil
}

type WithdrawalResponse struct {
	WithdrawalID string `json:"withdrawal_id"`
}

func NewWithdrawalResponse(withdrawalID int32) *WithdrawalResponse {
	resp := &WithdrawalResponse{strconv.FormatInt(int64(withdrawalID), 10)}
	return resp
}

func (rd *WithdrawalResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
