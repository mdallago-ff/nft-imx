package handlers

import (
	"errors"
	"net/http"
	"nft/imx"
	"nft/models"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/oauth"
	"github.com/go-chi/render"
	"github.com/google/uuid"
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

	info := imx.CollectionInformation{
		ContractAddress: data.ContractAddress,
		CollectionName:  data.CollectionName,
		MetadataUrl:     data.MetadataUrl,
	}

	err := h.imx.CreateCollection(r.Context(), &info)
	if err != nil {
		log.Error("error creating collection", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	metadataInfo := imx.MetadataInformation{
		ContractAddress: data.ContractAddress,
	}

	for _, f := range data.Fields {
		field := imx.MetadataFieldInformation{Name: f.Name, Type: f.Type}
		metadataInfo.Fields = append(metadataInfo.Fields, field)
	}

	err = h.imx.CreateMetadata(r.Context(), &metadataInfo)
	if err != nil {
		log.Error("error creating metadata", err)
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

	collection := models.Collection{
		ID:              uuid.New(),
		UserID:          userID,
		ContractAddress: data.ContractAddress,
	}
	err = h.db.CreateCollection(&collection)
	if err != nil {
		log.Error("error saving collection", err)
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, NewCollectionResponse(collection.ID.String(), data.CollectionName, data.ContractAddress))
	if err != nil {
		log.Error("error rendering response", err)
	}
}

type CollectionRequest struct {
	ContractAddress string                   `json:"contract_address"`
	CollectionName  string                   `json:"collection_name"`
	MetadataUrl     string                   `json:"metadata_url"`
	Fields          []CollectionFieldRequest `json:"fields"`
}

type CollectionFieldRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (a *CollectionRequest) Bind(r *http.Request) error {
	if len(a.ContractAddress) == 0 {
		return errors.New("missing required fields")
	}

	if len(a.CollectionName) == 0 {
		return errors.New("missing required fields")
	}

	return nil
}

type CollectionResponse struct {
	CollectionID    string `json:"collection_id"`
	ContractAddress string `json:"contract_address"`
	CollectionName  string `json:"collection_name"`
}

func NewCollectionResponse(id string, name string, address string) *CollectionResponse {
	resp := &CollectionResponse{CollectionID: id, CollectionName: name, ContractAddress: address}
	return resp
}

func (rd *CollectionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
