package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"nft/db"
	"nft/imx"

	"github.com/google/uuid"

	"github.com/hibiken/asynq"
)

const (
	TypeCompleteWithdrawal = "withdrawal:complete"
)

type CompleteWithdrawalPayload struct {
	WithdrawalID int32
	UserID       uuid.UUID
}

func NewCompleteWithdrawalTask(withdrawalID int32, userID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(CompleteWithdrawalPayload{withdrawalID, userID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeCompleteWithdrawal, payload), nil
}

type CompleteWithdrawalProcessor struct {
	imx imx.Client
	db  *db.DB
}

func (processor *CompleteWithdrawalProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p CompleteWithdrawalPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("withdrawal_id=%d", p.WithdrawalID)

	user, err := processor.db.GetUser(p.UserID)
	if err != nil {
		return err
	}

	if user == nil {
		return fmt.Errorf("user not exists: %v: %w", p.UserID, asynq.SkipRetry)
	}

	info := imx.CompleteWithdrawalInformation{
		User:         user,
		WithdrawalID: p.WithdrawalID,
	}

	err = processor.imx.CompleteEthWithdrawal(ctx, &info)
	if err != nil {
		return err
	}

	return nil
}

func NewCompleteWithdrawalProcessor(imx imx.Client, db *db.DB) *CompleteWithdrawalProcessor {
	return &CompleteWithdrawalProcessor{imx, db}
}
