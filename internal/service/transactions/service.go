package transactions

import (
	"context"
	"fmt"

	"github.com/quickybrains/bankapp/internal/domain"
	"github.com/quickybrains/bankapp/internal/ports/log"
)

type checkBalanceFn = func(account *domain.Account) error
type balanceChangedFn = func(*domain.AccountBalanceChangedMsg)

type transactionsRepo interface {
	CreateTransactionWithFn(
		ctx context.Context,
		transaction *domain.Transaction,
		checkFn checkBalanceFn,
		balanceFn balanceChangedFn,
	) error
	GetTransactionByID(ctx context.Context, transaction domain.TransactionID) (*domain.Transaction, error)
}

type balanceChangedNotifier interface {
	NotifyBalanceChanged(ctx context.Context, message *domain.AccountBalanceChangedMsg) error
}

type Service struct {
	repo     transactionsRepo
	notifier balanceChangedNotifier

	log log.Logger
}

func NewService(repo transactionsRepo, notifier balanceChangedNotifier, log log.Logger) *Service {
	return &Service{
		repo:     repo,
		notifier: notifier,
		log:      log,
	}
}

func (s *Service) CreateTransaction(transaction *domain.Transaction) error {
	// TODO: context
	err := s.repo.CreateTransactionWithFn(
		context.Background(),
		transaction,
		func(account *domain.Account) error {
			if account.Balance < transaction.Amount {
				return domain.ErrInsufficientFunds
			}

			return nil
		},
		func(msg *domain.AccountBalanceChangedMsg) {
			err := s.notifier.NotifyBalanceChanged(context.Background(), msg)
			if err != nil {
				s.log.Err(err).Print("notify balance changed")
			}
		})
	if err != nil {
		return fmt.Errorf("repo create transaction withh fn, err=%w", err)
	}

	return nil
}

func (s *Service) GetTransactionByID(id domain.TransactionID) (*domain.Transaction, error) {
	// TODO: context
	tx, err := s.repo.GetTransactionByID(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("repo get transaction by id, transaction id=%v, err=%w", id, err)
	}

	return tx, nil
}
