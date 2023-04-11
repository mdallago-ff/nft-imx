package handlers

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/render"
	"net/http"
	"nft/models"
)

func (h *Handler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	data := &CollectionRequest{}
	if err := render.Bind(r, data); err != nil {
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	render.Status(r, http.StatusCreated)
	/*err := render.Render(w, r, NewCollectionResponse(u))
	if err != nil {
		log.Error("error rendering response", err)
	}*/
}

type CollectionRequest struct {
}

func (a *CollectionRequest) Bind(r *http.Request) error {
	/*if len(a.Mail) == 0 {
		return errors.New("missing required fields")
	}*/

	return nil
}

type CollectionResponse struct {
	*models.User
}

func NewCollectionResponse(user *models.User) *CollectionResponse {
	resp := &CollectionResponse{User: user}
	return resp
}

func (rd *CollectionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
