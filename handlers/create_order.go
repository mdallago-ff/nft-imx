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

	collectionID, err := uuid.Parse(data.CollectionID)
	if err != nil {
		log.Error("error parsing collection", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	collection, err := h.db.GetCollection(collectionID)
	if err != nil {
		log.Error("error getting collection", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	if collection == nil {
		err = errors.New("collection missing")
		log.Error("error getting collection", err)
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

	if collection.UserID != userID {
		err = errors.New("invalid collection")
		log.Error("invalid collection", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	tokenID, err := uuid.Parse(data.TokenID)
	if err != nil {
		log.Error("error parsing token", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	token, err := h.db.GetToken(tokenID)
	if err != nil {
		log.Error("error getting token", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	if token == nil {
		err = errors.New("token missing")
		log.Error("error getting token", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	if token.CollectionID != collectionID {
		err = errors.New("invalid token")
		log.Error("invalid token", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	info := imx.OrderInformation{
		ContractAddress: collection.ContractAddress,
		TokenID:         data.TokenID,
		Amount:          amount,
	}

	orderID, err := h.imx.CreateOrder(r.Context(), &info)
	if err != nil {
		log.Error("error creating order", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, NewOrderResponse(orderID))
	if err != nil {
		log.Error("error rendering response", err)
	}
}

type OrderRequest struct {
	CollectionID string `json:"collection_id"`
	TokenID      string `json:"token_id"`
	Amount       string `json:"amount"`
}

func (a *OrderRequest) Bind(r *http.Request) error {
	if len(a.CollectionID) == 0 {
		return errors.New("missing required fields")
	}

	if len(a.TokenID) == 0 {
		return errors.New("missing required fields")
	}

	return nil
}

type OrderResponse struct {
	OrderID string `json:"order_id"`
}

func NewOrderResponse(id int32) *OrderResponse {
	resp := &OrderResponse{OrderID: string(id)}
	return resp
}

func (rd *OrderResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
