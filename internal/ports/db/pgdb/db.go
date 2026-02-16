package pgdb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/quickybrains/bankapp/internal/ports/db"
	"github.com/quickybrains/bankapp/internal/ports/log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type pgDB struct {
	db     *sql.DB
	config *db.Config
	lg     log.Logger
}

func NewPGDB(config *db.Config, lg log.Logger) (*pgDB, error) {
	db := &pgDB{
		config: config,
		lg:     lg,
	}

	err := db.init()
	if err != nil {
		return nil, fmt.Errorf("db init, err=%w", err)
	}

	return db, nil
}

func (db *pgDB) PingContext(ctx context.Context) error {
	return db.db.PingContext(ctx)
}

func (db *pgDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return db.db.QueryRowContext(ctx, query, args...)
}

func (db *pgDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args...)
}

func (db *pgDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

func (db *pgDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.db.BeginTx(ctx, opts)
}

func (db *pgDB) init() error {
	dsn := db.config.DSN()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("pgxpool parse config, dsn=%s, err=%w", dsn, err)
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), db.config.PoolConnectTimeoutInSeconds())
	defer cancelFn()

	db.lg.With("dsn", dsn).Info().Print("Connecting to db")

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("pgxpool new with config, dsn=%s, err=%w", dsn, err)
	}

	db.db = stdlib.OpenDBFromPool(pool)

	var testOne int
	if err := db.db.QueryRow("select 1").Scan(&testOne); err != nil {
		return fmt.Errorf("db ping, err=%w", err)
	}

	db.lg.Info().Print("db connected")

	return nil
}
