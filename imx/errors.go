package imx

import "fmt"

type WithdrawalNotReadyError struct {
	CurrentStatus string
}

func (c WithdrawalNotReadyError) Error() string {
	return fmt.Sprintf("Withdrawal needs to be confirmed to complete it. Current status %s", c.CurrentStatus)
}

func NewWithdrawalNotReadyError(currentStatus string) WithdrawalNotReadyError {
	return WithdrawalNotReadyError{currentStatus}
}
