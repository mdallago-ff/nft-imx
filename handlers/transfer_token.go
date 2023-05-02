package handlers

import (
	"errors"
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/oauth"
	"github.com/go-chi/render"
	"github.com/google/uuid"
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

	info := imx.TransferInformation{
		ContractAddress: collection.ContractAddress,
		TokenID:         data.TokenID,
		ReceiverAddress: data.ReceiverAddress,
	}

	err = h.imx.TransferToken(r.Context(), &info)
	if err != nil {
		log.Error("error transferring token", err)
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
	CollectionID    string `json:"collection_id"`
	TokenID         string `json:"token_id"`
	ReceiverAddress string `json:"receiver_address"`
}

func (a *TransferTokenRequest) Bind(r *http.Request) error {
	if len(a.CollectionID) == 0 {
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
