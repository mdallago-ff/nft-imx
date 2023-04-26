package test

import (
	"github.com/google/uuid"
	"nft/models"
)

func CreateDummyUser(id uuid.UUID, mail string) *models.User {
	return &models.User{
		ID:      id,
		ApiKey:  uuid.NewString(),
		Mail:    mail,
		Private: "",
		Public:  "",
		Address: "",
	}
}

func CreateDummyCollection(id uuid.UUID, userID uuid.UUID, contractAddress string) *models.Collection {
	return &models.Collection{
		ID:              id,
		UserID:          userID,
		ContractAddress: contractAddress,
	}
}

func CreateDummyToken(id uuid.UUID, collectionID uuid.UUID, tokenID string) *models.Token {
	return &models.Token{
		ID:           id,
		CollectionID: collectionID,
		TokenID:      tokenID,
	}
}
