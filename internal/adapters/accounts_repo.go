package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/quickybrains/bankapp/internal/domain"
	"github.com/quickybrains/bankapp/internal/ports/db"
	"github.com/quickybrains/bankapp/internal/ports/log"
)

type accountsRepo struct {
	db   db.CommonDB
	conf *db.Config
	log  log.Logger
}

func NewAccountsRepo(db db.CommonDB, conf *db.Config, log log.Logger) *accountsRepo {
	return &accountsRepo{
		db:   db,
		conf: conf,
		log:  log,
	}
}

func (repo *accountsRepo) CreateAccount(ctx context.Context, account *domain.Account) error {
	query := `
	insert into
		accounts (id, user_id, fullname, currency, balance, created_at)
	values
		($1, $2, $3, $4, $5, $6)
	`

	ctx, cancelFn := context.WithTimeout(ctx, repo.conf.QueryTimeoutInSeconds())
	defer cancelFn()

	_, err := repo.db.ExecContext(
		ctx,
		query,
		account.ID,
		account.UserID,
		account.Fullname,
		account.Currency,
		account.Balance,
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("repo db exec context, err=%w", err)
	}

	return nil
}

func (repo *accountsRepo) GetAccountByID(ctx context.Context, id domain.AccountID) (*domain.Account, error) {
	query := `
	select
		id,
		user_id,
		fullname,
		balance
	from
		accounts
	where
		id = $1
	`

	var account domain.Account

	ctx, cancelFn := context.WithTimeout(ctx, repo.conf.QueryTimeoutInSeconds())
	defer cancelFn()

	row := repo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&account.ID, &account.UserID, &account.Fullname, &account.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoAccountWithID(id)
		}

		return nil, fmt.Errorf("repo db row scan, err=%w", err)
	}

	return &account, nil
}

func (repo *accountsRepo) GetAccountsByUserID(ctx context.Context, user_id domain.UserID) ([]*domain.Account, error) {
	query := `
	select id, user_id, fullname, balance 
	from accounts 
	where user_id = $1
	`

	ctx, cancelFn := context.WithTimeout(ctx, repo.conf.QueryTimeoutInSeconds())
	defer cancelFn()

	rows, err := repo.db.QueryContext(ctx, query, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoAccountsForUserID(user_id)
		}

		return nil, fmt.Errorf("repo query context, err=%w", err)
	}

	accounts := make([]*domain.Account, 0, 10)
	for rows.Next() {
		var account domain.Account
		err = rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Fullname,
			&account.Balance,
		)
		if err != nil {
			return nil, fmt.Errorf("db rows scan, err=%w", err)
		}

		accounts = append(accounts, &account)
	}

	if len(accounts) == 0 {
		return nil, domain.ErrNoAccountsForUserID(user_id)
	}

	return accounts, nil
}

type checkBalanceFn = func(account *domain.Account) error
type balanceChangedFn = func(*domain.AccountBalanceChangedMsg)

func (repo *accountsRepo) DebitOnAccountWithFn(
	ctx context.Context,
	id domain.AccountID,
	money float64,
	checkFn checkBalanceFn,
	onChangeFn balanceChangedFn,
) (*domain.Account, error) {
	getAccountCtx, cancelFn := context.WithTimeout(ctx, repo.conf.QueryTimeoutInSeconds())
	defer cancelFn()

	account, err := repo.GetAccountByID(getAccountCtx, id)
	if err != nil {
		return nil, fmt.Errorf("get account by id, err=%w", err)
	}

	err = checkFn(account)
	if err != nil {
		return nil, fmt.Errorf("check fn, err=%w", err)
	}

	query := `
	update accounts
	set
		balance = balance + $1
	where
		id = $2
	returning
		-- returning full account info on debit
		id,
		user_id,
		fullname,
		currency,
		-- returning old and new balance for notification
		balance - $1 as old_balance,
		balance as new_balance
	`

	updateBalanceCtx, cancelFn := context.WithTimeout(ctx, repo.conf.QueryTimeoutInSeconds())
	defer cancelFn()

	row := repo.db.QueryRowContext(updateBalanceCtx, query, money, id)

	var oldBalance, newBalance float64
	err = row.Scan(
		&account.ID,
		&account.UserID,
		&account.Fullname,
		&account.Currency,
		&oldBalance,
		&newBalance,
	)
	if err != nil {
		return nil, fmt.Errorf("repo db row scan account balance after update, err=%w", err)
	}

	balanceChanged := &domain.AccountBalanceChangedMsg{
		AccountID:       id,
		PreviousBalance: oldBalance,
		CurrentBalance:  newBalance,
	}

	onChangeFn(balanceChanged)

	account.Balance = newBalance

	return account, nil
}

func (repo *accountsRepo) WithdrawFromAccountWithFn(
	ctx context.Context,
	id domain.AccountID,
	money float64,
	checkFn checkBalanceFn,
	onChangeFn balanceChangedFn,
) (*domain.Account, error) {
	getAccountCtx, cancelFn := context.WithTimeout(ctx, repo.conf.QueryTimeoutInSeconds())
	defer cancelFn()

	account, err := repo.GetAccountByID(getAccountCtx, id)
	if err != nil {
		return nil, fmt.Errorf("repo get account by id, err=%w", err)
	}

	query := `
	update accounts
	set
		balance = balance - $1
	where id = $2 and balance >= $1
	returning
		balance + $1 as old_balance,
		balance as new_balance
	`

	withdrawCtx, cancelFn := context.WithTimeout(ctx, repo.conf.QueryTimeoutInSeconds())
	defer cancelFn()

	row := repo.db.QueryRowContext(withdrawCtx, query, money, id)

	var oldBalance, newBalance float64
	err = row.Scan(&oldBalance, &newBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrInsufficientFunds()
		}

		return nil, fmt.Errorf("repo db row scan account balance after update, err=%w", err)
	}

	// TODO
	onChangeFn(&domain.AccountBalanceChangedMsg{
		AccountID:       id,
		PreviousBalance: oldBalance,
		CurrentBalance:  newBalance,
	})

	account.Balance = newBalance

	return account, nil
}
