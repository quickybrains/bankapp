package domain

import (
	"errors"
	"fmt"
)

type PublicError struct {
	err error
}

func (e *PublicError) Error() string {
	return e.err.Error()
}

func ErrNegativeAmount(amount float64) *PublicError {
	return &PublicError{fmt.Errorf("negative amount %v is prohibited", amount)}
}

func ErrInsufficientFunds() *PublicError {
	return &PublicError{errors.New("insufficient funds")}
}

func ErrNoAccountsForUserID(id UserID) *PublicError {
	return &PublicError{fmt.Errorf("found no accounts for given userId %s", id.String())}
}

func ErrNoAccountWithID(id AccountID) *PublicError {
	return &PublicError{fmt.Errorf("found no account with given id %s", id.String())}
}

func ErrNoTransactionWithID(id TransactionID) *PublicError {
	return &PublicError{fmt.Errorf("found no transaction with given id %s", id.String())}
}
