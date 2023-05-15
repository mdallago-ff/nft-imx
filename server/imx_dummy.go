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

func (i ImxDummy) CreateTrade(ctx context.Context, info *imx.CreateTradeInformation) (int32, error) {
	return 0, nil
}

func (i ImxDummy) CreateEthWithdrawal(ctx context.Context, info *imx.CreateWithdrawalInformation) (int32, error) {
	return 1, nil
}

func (i ImxDummy) CompleteEthWithdrawal(ctx context.Context, info *imx.CompleteWithdrawalInformation) error {
	return nil
}
