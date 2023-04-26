package handlers

import (
	"errors"
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"net/http"
	"nft/imx"
	"nft/models"
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

	info := imx.MintInformation{
		ContractAddress: collection.ContractAddress,
		TokenID:         data.TokenID,
		Blueprint:       data.Blueprint,
	}

	err = h.imx.CreateToken(r.Context(), &info)
	if err != nil {
		log.Error("error creating token", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	token := models.Token{
		ID:           uuid.New(),
		CollectionID: collectionID,
		TokenID:      data.TokenID,
	}

	err = h.db.CreateToken(&token)
	if err != nil {
		log.Error("error saving token", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, NewTokenResponse(token.ID.String()))
	if err != nil {
		log.Error("error rendering response", err)
	}
}

type TokenRequest struct {
	CollectionID string `json:"collection_id"`
	TokenID      string `json:"token_id"`
	Blueprint    string `json:"blueprint"`
}

func (a *TokenRequest) Bind(r *http.Request) error {
	if len(a.CollectionID) == 0 {
		return errors.New("missing required fields")
	}

	if len(a.TokenID) == 0 {
		return errors.New("missing required fields")
	}

	return nil
}

type TokenResponse struct {
	TokenID string `json:"token_id"`
}

func NewTokenResponse(id string) *TokenResponse {
	resp := &TokenResponse{TokenID: id}
	return resp
}

func (rd *TokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
