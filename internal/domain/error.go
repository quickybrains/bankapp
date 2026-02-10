package domain

import (
	"errors"
	"slices"
)

var (
	ErrNegativeAmount    = errors.New("negative amount is prohibited")
	ErrInsufficientFunds = errors.New("insufficient funds")

	ErrNoAccountWithID     = errors.New("found no account with given id")
	ErrNoTransactionWithID = errors.New("found no transaction with given id")
)

var allDomainErrors = []error{ErrNegativeAmount, ErrInsufficientFunds}

func IsDomainError(inErr error) bool {
	return slices.IndexFunc(allDomainErrors, func(err error) bool {
		return errors.Is(inErr, err)
	}) != -1
}
