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

func (h *Handler) CreateTrade(w http.ResponseWriter, r *http.Request) {
	data := &TradeRequest{}
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

	orderID, err := strconv.ParseInt(data.OrderID, 10, 32)
	if err != nil {
		log.Error("error creating order", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	info := imx.CreateTradeInformation{
		OrderID: int32(orderID),
		User:    user,
	}

	tradeID, err := h.imx.CreateTrade(r.Context(), &info)
	if err != nil {
		log.Error("error creating trade", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, NewTradeResponse(tradeID))
	if err != nil {
		log.Error("error rendering response", err)
	}
}

type TradeRequest struct {
	OrderID string `json:"order_id"`
}

func (a *TradeRequest) Bind(r *http.Request) error {
	if len(a.OrderID) == 0 {
		return errors.New("missing required fields")
	}

	return nil
}

type TradeResponse struct {
	TradeID string `json:"trade_id"`
}

func NewTradeResponse(id int32) *TradeResponse {
	resp := &TradeResponse{TradeID: string(id)}
	return resp
}

func (rd *TradeResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
