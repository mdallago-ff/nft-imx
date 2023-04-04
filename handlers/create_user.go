package handlers

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"net/http"
	"nft/models"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	data := &UserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	//h.imx.CreateUser()

	user := data.User
	user.ID = uuid.NewString()
	user.ApiKey = uuid.NewString()
	//h.db.CreateUser()

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewUserResponse(user))
}

type UserRequest struct {
	*models.User
}

func (a *UserRequest) Bind(r *http.Request) error {
	// a.Article is nil if no Article fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.User == nil {
		return errors.New("missing required fields")
	}

	// a.User is nil if no Userpayload fields are sent in the request. In this app
	// this won't cause a panic, but checks in this Bind method may be required if
	// a.User or futher nested fields like a.User.Name are accessed elsewhere.

	// just a post-process after a decode..
	//a.ProtectedID = ""                                 // unset the protected ID
	//a.Article.Title = strings.ToLower(a.Article.Title) // as an example, we down-case
	return nil
}

type UserResponse struct {
	*models.User

	//User *UserPayload `json:"user,omitempty"`

	// We add an additional field to the response here.. such as this
	// elapsed computed property
	//Elapsed int64 `json:"elapsed"`
}

func NewUserResponse(user *models.User) *UserResponse {
	resp := &UserResponse{User: user}
	return resp
}

func (rd *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}