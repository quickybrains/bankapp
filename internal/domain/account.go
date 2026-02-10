package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type AccountID = uuid.UUID
type UserID = uuid.UUID

type Account struct {
	// Id An ID of an account. The UUID format is used.
	ID uuid.UUID `json:"id"`

	// UserId An ID of a user. The UUID format is used.
	UserID uuid.UUID `json:"userId"`

	// Fullname Name of the account holder.
	Fullname string `json:"fullname"`

	// Currency Account currency.
	Currency string `json:"currency"`

	// Balance Account balance in real units.
	Balance float64 `json:"balance"`
}

type AccountBalanceChangedMsg struct {
	AccountID       AccountID `json:"accountId"`
	PreviousBalance float64   `json:"previousBalance"`
	CurrentBalance  float64   `json:"currentBalance"`
}

func (msg AccountBalanceChangedMsg) String() string {
	return fmt.Sprintf(
		"account: %v | previous balance: %v | current balance: %v",
		msg.AccountID,
		msg.PreviousBalance,
		msg.CurrentBalance,
	)
}
