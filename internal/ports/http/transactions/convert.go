package transactions

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/quickybrains/bankapp/internal/domain"
)

func fromTransactionReqToDomainTransaction(reqTx *Req) (*domain.Transaction, error) {
	txID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("uuid new random, err=%w", err)
	}

	if reqTx.Amount < 0 {
		return nil, domain.ErrNegativeAmount(reqTx.Amount)
	}

	return &domain.Transaction{
		ID:            txID,
		FromAccountID: reqTx.FromAccountId,
		ToAccountID:   reqTx.ToAccountId,
		Amount:        reqTx.Amount,
		CreatedAt:     time.Now().UTC(),
	}, nil
}
