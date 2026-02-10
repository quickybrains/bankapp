package domain

import "github.com/google/uuid"

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

	// CreatedAt Unix timestamp (seconds or milliseconds since the Unix epoch)
	CreatedAt int64 `json:"createdAt"`

	// ProcessedAt Unix timestamp (seconds or milliseconds since the Unix epoch)
	ProcessedAt int64 `json:"processedAt"`
}
