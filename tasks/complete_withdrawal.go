package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
)

const (
	TypeCompleteWithdrawal = "withdrawal:complete"
)

type CompleteWithdrawalPayload struct {
	WithdrawalID int
}

func NewCompleteWithdrawalTask(withdrawalID int) (*asynq.Task, error) {
	payload, err := json.Marshal(CompleteWithdrawalPayload{withdrawalID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeCompleteWithdrawal, payload), nil
}

func HandleCompleteWithdrawalTask(ctx context.Context, t *asynq.Task) error {
	var p CompleteWithdrawalPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("withdrawal_id=%d", p.WithdrawalID)
	return nil
}
