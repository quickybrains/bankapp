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
	balanceFn balanceChangedFn,
) error {
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
	if err != nil {
		return fmt.Errorf("repo db begin tx, err=%w", err)
	}

	// First check if accounts exist
	_, err = repo.getAccountWithTx(tx, ctx, transaction.FromAccountID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			repo.log.Err(rollbackErr).Print("tx rollback")
		}

		return domain.ErrNoAccountWithID(transaction.FromAccountID)
	}

	_, err = repo.getAccountWithTx(tx, ctx, transaction.ToAccountID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			repo.log.Err(rollbackErr).Print("tx rollback")
		}

		return domain.ErrNoAccountWithID(transaction.ToAccountID)
	}

	// and we want to check rows affected after query completes
	withdrawQuery := `
	update accounts
	set
		balance = balance - $1
	where id = $2 and balance >= $1
	returning
		balance + $1 as old_balance,
		balance as new_balance
	`

	row := tx.QueryRowContext(ctx, withdrawQuery, transaction.Amount, transaction.FromAccountID)

	var fromOldBalance, fromNewBalance float64
	err = row.Scan(&fromOldBalance, &fromNewBalance)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			repo.log.Err(rollbackErr).Print("tx rollback")
		}

		// Since we checked for accounts here it means
		// that the balance was insufficient
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrInsufficientFunds()
		}

		return fmt.Errorf("repo db row scan updated from account balance, err=%w", err)
	}

	debitQuery := `
	update accounts
	set
		balance = balance + $1
	where
		id = $2
	returning
		balance - $1 as old_balance,
		balance as new_balance
	`

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

	// Once we've dealt with accounts' balances we can insert the transaction itself
	txQuery := `
	insert into
  		transactions (
			id,
			from_account_id,
			to_account_id,
			amount,
			descr,
			created_at,
			processed_at
  		)
	values
  		($1, $2, $3, $4, $5, $6, $7)
	`

	transaction.ProcessedAt = time.Now().UTC()

	_, err = tx.ExecContext(
		ctx,
		txQuery,
		transaction.ID,
		transaction.FromAccountID,
		transaction.ToAccountID,
		transaction.Amount,
		transaction.Description,
		transaction.CreatedAt,
		transaction.ProcessedAt,
	)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			repo.log.Err(rollbackErr).Print("tx rollback")
		}

		return fmt.Errorf("repo db tx exec create transaction query, err=%w", err)
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
		PreviousBalance: toOldBalance,
		CurrentBalance:  toNewBalance,
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoTransactionWithID(id)
		}

		return nil, fmt.Errorf("repo db row scan, err=%w", err)
	}

	return &tx, nil
}

func (repo *transactionsRepo) getAccountWithTx(tx *sql.Tx, ctx context.Context, id domain.AccountID) (*domain.Account, error) {
	getAccountQuery := `
	select
		id,
		fullname,
		currency,
		balance
	from
		accounts
	where
		id = $1
	`

	row := tx.QueryRowContext(ctx, getAccountQuery, id)

	var account domain.Account
	err := row.Scan(&account.ID, &account.Fullname, &account.Currency, &account.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoAccountWithID(id)
		}

		return nil, fmt.Errorf("repo db row scan, err=%w", err)
	}

	return &account, nil
}
