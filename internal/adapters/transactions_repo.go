package adapters

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/quickybrains/bankapp/internal/domain"
	"github.com/quickybrains/bankapp/internal/ports/db"
	"github.com/quickybrains/bankapp/internal/ports/log"
)

type transactionsRepo struct {
	db   db.CommonDB
	conf *db.Config
	log  log.Logger
}

func NewTransactionsRepo(db db.CommonDB, conf *db.Config, log log.Logger) *transactionsRepo {
	return &transactionsRepo{
		db:   db,
		conf: conf,
		log:  log,
	}
}

func (repo *transactionsRepo) CreateTransactionWithFn(
	ctx context.Context,
	transaction *domain.Transaction,
	checkFn checkBalanceFn,
	balanceFn balanceChangedFn,
) error {
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
	if err != nil {
		return fmt.Errorf("repo db begin tx, err=%w", err)
	}

	// First get an account
	account, err := repo.getAccountWithTx(tx, ctx, transaction.FromAccountID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			repo.log.Err(rollbackErr).Print("tx rollback")
		}

		return fmt.Errorf("repo get account with tx, err=%w", err)
	}

	// Check account balance by calling provided func
	err = checkFn(account)
	if err != nil {
		return fmt.Errorf("check balance fn, err=%w", err)
	}

	withdrawQuery := `
	update account
	set
		balance = balance - $1
	where
		id = $2
	returning
		old.balance as old_balance,
		new.balance as new_balance
	`

	debitQuery := `
	update account
	set
		balance = balance + $1
	where
		id = $2
	returning
		old.balance as old_balance,
		new.balance as new_balance
	`

	row := tx.QueryRowContext(ctx, withdrawQuery, transaction.Amount, transaction.FromAccountID)

	var fromOldBalance, fromNewBalance float64
	err = row.Scan(&fromOldBalance, &fromNewBalance)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			repo.log.Err(rollbackErr).Print("tx rollback")
		}

		return fmt.Errorf("repo db row scan updated from account balance, err=%w", err)
	}

	var toOldBalance, toNewBalance float64
	row = tx.QueryRowContext(ctx, debitQuery, transaction.Amount, transaction.ToAccountID)
	err = row.Scan(&toOldBalance, &toNewBalance)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			repo.log.Err(rollbackErr).Print("tx rollback")
		}

		return fmt.Errorf("repo db row scan updated to account balance, err=%w", err)
	}

	err = tx.Commit()
	if err != nil {
		repo.log.Err(err).Print("tx commit")

		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			repo.log.Err(rollbackErr).Print("tx rollback")
		}
	}

	// Calling balanceFn callbacks with actions on
	// accounts' balances being changed
	// only after successful tx commit
	balanceFn(&domain.AccountBalanceChangedMsg{
		AccountID:       transaction.FromAccountID,
		PreviousBalance: fromOldBalance,
		CurrentBalance:  fromNewBalance,
	})

	balanceFn(&domain.AccountBalanceChangedMsg{
		AccountID:       transaction.ToAccountID,
		PreviousBalance: fromOldBalance,
		CurrentBalance:  fromNewBalance,
	})

	return nil
}

func (repo *transactionsRepo) GetTransactionByID(
	ctx context.Context,
	id domain.TransactionID,
) (*domain.Transaction, error) {
	query := `
	select
		id,
		from_account_id,
		to_account_id,
		amount,
		created_at,
		processed_at
	from
		transactions
	where
		id = $1
	`

	row := repo.db.QueryRowContext(ctx, query, id)

	var tx domain.Transaction
	err := row.Scan(
		&tx.ID,
		&tx.FromAccountID,
		&tx.ToAccountID,
		&tx.Amount,
		&tx.CreatedAt,
		&tx.ProcessedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("repo db row scan, err=%w", err)
	}

	return &tx, nil
}

func (repo *transactionsRepo) getAccountWithTx(tx *sql.Tx, ctx context.Context, id domain.AccountID) (*domain.Account, error) {
	getAccountQuery := `
	select
		id,
		fullname,
		balance
	from
		accounts
	where
		id = $1
	`

	row := tx.QueryRowContext(ctx, getAccountQuery, id)

	var account domain.Account
	err := row.Scan(&account.ID, &account.Fullname, &account.Balance)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			repo.log.Err(rollbackErr).Print("tx rollback")
		}

		return nil, fmt.Errorf("repo db row scan, err=%w", err)
	}

	return &account, nil
}
