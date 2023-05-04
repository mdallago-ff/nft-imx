package handlers

import (
	"errors"
	"net/http"
	"nft/keys"
	"nft/models"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	data := &UserRequest{}
	if err := render.Bind(r, data); err != nil {
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	mail := data.Mail

	u, err := h.db.GetUserByMail(mail)
	if err != nil {
		log.Error("error getting user", err)
		err = render.Render(w, r, ErrServer(err))
		if err != nil {
			log.Error("error rendering response", err)
		}
		return
	}

	if u == nil {
		pair, err := keys.CreateKeys()
		if err != nil {
			err = render.Render(w, r, ErrInvalidRequest(errors.New("not allowed")))
			if err != nil {
				log.Error("error rendering response", err)
			}
			return
		}

		user := models.User{}
		user.ID = uuid.New()
		user.ApiKey = uuid.NewString()
		user.Mail = mail
		user.Private = pair.Private //TODO do not store private keys in plain text. We must use a vault or similar.
		user.Public = pair.Public
		user.Address = pair.Address

		err = h.db.CreateUser(&user)
		if err != nil {
			log.Error("error saving user", err)
			err = render.Render(w, r, ErrServer(err))
			if err != nil {
				log.Error("error rendering response", err)
			}
			return
		}
		u = &user
	}

	if len(u.StarkKey) == 0 {
		starkKey, err := h.imx.CreateUser(r.Context(), u)
		if err != nil {
			log.Error("error creating user", err)
			err = render.Render(w, r, ErrServer(err))
			if err != nil {
				log.Error("error rendering response", err)
			}
			return
		}

		u.StarkKey = starkKey
		err = h.db.UpdateUser(u)
		if err != nil {
			log.Error("error saving user", err)
			err = render.Render(w, r, ErrServer(err))
			if err != nil {
				log.Error("error rendering response", err)
			}
			return
		}
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, NewUserResponse(u))
	if err != nil {
		log.Error("error rendering response", err)
	}
}

type UserRequest struct {
	Mail string `json:"mail"`
}

func (a *UserRequest) Bind(r *http.Request) error {
	if len(a.Mail) == 0 {
		return errors.New("missing required fields")
	}

	return nil
}

type UserResponse struct {
	*models.User
}

func NewUserResponse(user *models.User) *UserResponse {
	resp := &UserResponse{User: user}
	return resp
}

func (rd *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
