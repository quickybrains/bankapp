package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionID = uuid.UUID

type Transaction struct {
	// ID Transaction ID.
	ID uuid.UUID `json:"id"`

	// FromAccountID An ID of an account money being transferred from. The UUID format is used.
	FromAccountID AccountID `json:"fromAccountId"`

	// ToAccountID An ID of an account money being transferred to. The UUID format is used.
	ToAccountID AccountID `json:"toAccountId"`

	// Amount Amount model used for debiting/withdrawal/transferring money. Negative amount is prohibited.
	Amount float64 `json:"amount"`

	Description string `json:"description"`

	// CreatedAt Unix timestamp (seconds or milliseconds since the Unix epoch)
	CreatedAt time.Time `json:"createdAt"`

	// ProcessedAt Unix timestamp (seconds or milliseconds since the Unix epoch)
	ProcessedAt time.Time `json:"processedAt"`
}
