package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/quickybrains/bankapp/internal/adapters"
	"github.com/quickybrains/bankapp/internal/ports/db/pgdb"
	accountshttp "github.com/quickybrains/bankapp/internal/ports/http/accounts"
	accounts_middleware "github.com/quickybrains/bankapp/internal/ports/http/accounts/middleware"
	transactionshttp "github.com/quickybrains/bankapp/internal/ports/http/transactions"
	transactions_middleware "github.com/quickybrains/bankapp/internal/ports/http/transactions/middleware"
	"github.com/quickybrains/bankapp/internal/ports/log"
	"github.com/quickybrains/bankapp/internal/ports/log/zerolog"
	"github.com/quickybrains/bankapp/internal/ports/notifier/rabbitmq"
	"github.com/quickybrains/bankapp/internal/service/accounts"
	"github.com/quickybrains/bankapp/internal/service/transactions"
)

type Node struct {
	config *Config
	srv    *http.Server
	logger log.Logger

	errCh chan error
}

func StartNewNode(config *Config) (*Node, error) {
	var logger log.Logger = zerolog.New(config.Log).With("component", "app_node")

	db, err := pgdb.NewPGDB(config.DB, logger.With("component", "db"))
	if err != nil {
		panic(err)
	}

	accountsRepo := adapters.NewAccountsRepo(db, config.DB, logger.With("component", "accounts_repo"))
	transactionsRepo := adapters.NewTransactionsRepo(
		db,
		config.DB,
		logger.With("component", "transactions_repo"),
	)

	notifier, err := rabbitmq.NewNotifier(config.Notifier, logger.With("component", "notifier"))
	if err != nil {
		panic(err)
	}

	accountsService := accounts.NewService(
		accountsRepo,
		notifier,
		logger.With("component", "accounts_service"),
	)
	transactionsService := transactions.NewService(
		transactionsRepo,
		notifier,
		logger.With("component", "transactions_service"),
	)

	accountsHandler := accountshttp.NewHandler(
		accountsService,
		config.HTTP,
		logger.With("component", "accounts_http_handler"),
	)

	transactionsHandler := transactionshttp.NewHandler(
		transactionsService,
		config.HTTP,
		logger.With("component", "transactions_http_handler"),
	)

	mx := http.NewServeMux()

	accountshttp.HandlerWithOptions(accountsHandler, accountshttp.StdHTTPServerOptions{
		BaseRouter: mx,
		Middlewares: []accountshttp.MiddlewareFunc{
			accounts_middleware.Logging(logger.With("component", "logging_middleware")),
			accounts_middleware.Limit(logger.With("component", "limit_middleware"), config.HTTP.HTTPReqBodyLimit()),
		},
	})
	transactionshttp.HandlerWithOptions(transactionsHandler, transactionshttp.StdHTTPServerOptions{
		BaseRouter: mx,
		Middlewares: []transactionshttp.MiddlewareFunc{
			transactions_middleware.Logging(logger.With("component", "logging_middleware")),
			transactions_middleware.Limit(logger.With("component", "limit_middleware"), config.HTTP.HTTPReqBodyLimit()),
		},
	})

	srv := &http.Server{
		Addr:         config.HTTP.ServerURL(),
		ReadTimeout:  config.HTTP.ReadTimeout(),
		WriteTimeout: config.HTTP.WriteTimeout(),

		Handler: mx,
	}

	node := &Node{
		config: config,
		srv:    srv,
		logger: logger,

		errCh: make(chan error, 1),
	}

	go func() {
		node.errCh <- node.srv.ListenAndServe()
	}()

	node.logger.Info().Print("Node started")

	return nil, nil
}

func (n *Node) Stop() {
	ctx, cancelFn := context.WithTimeout(context.Background(), n.config.ShutdownTimeoutInSeconds())
	defer cancelFn()

	err := n.srv.Shutdown(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		n.logger.Err(err).Print("Server shutdown failed")
	}

	// After server shutdown if there is any error it's already in errCh
	select {
	case err := <-n.errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			n.logger.Err(err).Print("Server shutdown failed")
		}
	default:
	}
}
