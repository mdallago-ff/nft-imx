package server

import (
	"context"
	"nft/imx"
	"nft/models"
)

type ImxDummy struct {
}

func (i ImxDummy) Close() {}

func (i ImxDummy) CreateUser(ctx context.Context, user *models.User) (string, error) {
	return "", nil
}

func (i ImxDummy) CreateCollection(ctx context.Context, info *imx.CollectionInformation) error {
	return nil
}

func (i ImxDummy) CreateMetadata(ctx context.Context, info *imx.MetadataInformation) error {
	return nil
}

func (i ImxDummy) CreateToken(ctx context.Context, info *imx.MintInformation) error {
	return nil
}
