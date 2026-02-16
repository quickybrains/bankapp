package domain

import (
	"errors"
)

type PublicError struct {
	err error
}

func (e *PublicError) Error() string {
	return e.err.Error()
}

var (
	ErrNegativeAmount    = &PublicError{errors.New("negative amount is prohibited")}
	ErrInsufficientFunds = &PublicError{errors.New("insufficient funds")}

	ErrNoAccountWithID     = &PublicError{errors.New("found no account with given id")}
	ErrNoTransactionWithID = &PublicError{errors.New("found no transaction with given id")}
)
