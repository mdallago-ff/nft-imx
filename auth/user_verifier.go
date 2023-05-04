package auth

import (
	"crypto/subtle"
	"errors"
	"net/http"
	"nft/db"

	"github.com/go-chi/oauth"
	"github.com/google/uuid"
)

// UserVerifier provides user credentials verifier for testing.
type UserVerifier struct {
	db *db.DB
}

func NewUserVerifier(db *db.DB) *UserVerifier {
	return &UserVerifier{db}
}

// ValidateUser validates username and password returning an error if the user credentials are wrong
func (*UserVerifier) ValidateUser(username, password, scope string, r *http.Request) error {
	return nil
}

// ValidateClient validates clientID and secret returning an error if the client credentials are wrong
func (u *UserVerifier) ValidateClient(clientID, clientSecret, scope string, r *http.Request) error {
	id, err := uuid.Parse(clientID)
	if err != nil {
		return errors.New("wrong client")
	}

	user, err := u.db.GetUser(id)
	if err != nil {
		return errors.New("wrong client")
	}

	if user == nil {
		return errors.New("wrong client")
	}

	if subtle.ConstantTimeCompare([]byte(user.ApiKey), []byte(clientSecret)) == 0 {
		return errors.New("wrong client")
	}

	return nil
}

// ValidateCode validates token ID
func (*UserVerifier) ValidateCode(clientID, clientSecret, code, redirectURI string, r *http.Request) (string, error) {
	return "", nil
}

// AddClaims provides additional claims to the token
func (*UserVerifier) AddClaims(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	claims := make(map[string]string)
	return claims, nil
}

// AddProperties provides additional information to the token response
func (*UserVerifier) AddProperties(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	props := make(map[string]string)
	return props, nil
}

// ValidateTokenID validates token ID
func (*UserVerifier) ValidateTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}

// StoreTokenID saves the token id generated for the user
func (*UserVerifier) StoreTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}
