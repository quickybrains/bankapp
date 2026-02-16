package accounts

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/quickybrains/bankapp/internal/domain"
)

func fromAccountReqToDomainAccount(req *AccountReq) (*domain.Account, error) {
	accountID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("uuid new random, err=%w", err)
	}

	var initAmount float64
	if req.InitialDebitAmount != nil {
		initAmount = *req.InitialDebitAmount
	}

	if initAmount < 0 {
		return nil, domain.ErrNegativeAmount(initAmount)
	}

	return &domain.Account{
		ID:       accountID,
		UserID:   req.UserId,
		Fullname: req.Fullname,
		Currency: req.Currency,
		Balance:  initAmount,
	}, nil
}
