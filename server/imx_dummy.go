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

func (i ImxDummy) TransferToken(ctx context.Context, info *imx.TransferInformation) error {
	return nil
}

func (i ImxDummy) CreateOrder(ctx context.Context, info *imx.OrderInformation) (int32, error) {
	return 0, nil
}

func (i ImxDummy) CreateEthDeposit(ctx context.Context, info *imx.CreateDepositInformation) (string, error) {
	return "hash", nil
}
