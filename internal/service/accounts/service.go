package accounts

import (
	"context"
	"fmt"

	"github.com/quickybrains/bankapp/internal/domain"
	"github.com/quickybrains/bankapp/internal/ports/log"
)

type checkBalanceFn = func(account *domain.Account) error
type balanceChangedFn = func(*domain.AccountBalanceChangedMsg)

type accountsRepo interface {
	CreateAccount(account *domain.Account) error
	GetAccountByID(id domain.AccountID) (*domain.Account, error)
	GetAccountsByUserID(id domain.UserID) ([]*domain.Account, error)

	DebitOnAccountWithFn(
		id domain.AccountID,
		money float64,
		checkBalanceFn checkBalanceFn,
		fn balanceChangedFn,
	) (*domain.Account, error)
	WithdrawFromAccountWithFn(
		id domain.AccountID,
		money float64,
		checkBalanceFn checkBalanceFn,
		balanceFn balanceChangedFn,
	) (*domain.Account, error)
}

type balanceChangedNotifier interface {
	NotifyBalanceChanged(ctx context.Context, message *domain.AccountBalanceChangedMsg) error
}

type Service struct {
	repo     accountsRepo
	notifier balanceChangedNotifier
	log      log.Logger
}

func NewService(repo accountsRepo, notifier balanceChangedNotifier, log log.Logger) *Service {
	return &Service{
		repo:     repo,
		notifier: notifier,
		log:      log,
	}
}

func (s *Service) CreateAccount(account *domain.Account) error {
	err := s.repo.CreateAccount(account)
	if err != nil {
		return fmt.Errorf("repo create account, err=%w", err)
	}

	return nil
}

func (s *Service) GetAccountByID(id domain.AccountID) (*domain.Account, error) {
	account, err := s.repo.GetAccountByID(id)
	if err != nil {
		return nil, fmt.Errorf("repo get account by id, id=%v, err=%w", id, err)
	}

	return account, nil
}

func (s *Service) GetAccountsByUserID(id domain.UserID) ([]*domain.Account, error) {
	accounts, err := s.repo.GetAccountsByUserID(id)
	if err != nil {
		return nil, fmt.Errorf("repo get accounts by user id, user_id=%v, err=%w", id, err)
	}

	return accounts, nil
}

func (s *Service) DebitOnAccount(id domain.AccountID, money float64) (*domain.Account, error) {
	account, err := s.repo.DebitOnAccountWithFn(
		id,
		money,
		func(account *domain.Account) error {
			if money < 0 {
				return domain.ErrNegativeAmount(money)
			}

			return nil
		},
		func(msg *domain.AccountBalanceChangedMsg) {
			// TODO: context
			// Notifying on changing balance
			err := s.notifier.NotifyBalanceChanged(context.Background(), msg)
			if err != nil {
				s.log.Err(err).Print("notify balance changed")
			}
		},
	)
	if err != nil {
		return nil, fmt.Errorf("repo debit on account, id=%v, err=%w", id, err)
	}

	return account, nil
}

func (s *Service) WithdrawFromAccount(id domain.AccountID, money float64) (*domain.Account, error) {
	account, err := s.repo.WithdrawFromAccountWithFn(
		id,
		money,
		func(account *domain.Account) error {
			if money < 0 {
				return domain.ErrNegativeAmount(money)
			}

			// Checking that account has sufficient funds
			if account.Balance < money {
				return domain.ErrInsufficientFunds()
			}

			return nil
		},
		func(msg *domain.AccountBalanceChangedMsg) {
			// TODO: context
			// Notifying on changing balance
			err := s.notifier.NotifyBalanceChanged(context.Background(), msg)
			if err != nil {
				s.log.Err(err).Print("notify balance changed")
			}
		})
	if err != nil {
		return nil, fmt.Errorf("repo withdraw from account, id=%v, err=%w", id, err)
	}

	return account, nil
}
