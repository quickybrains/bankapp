package transactions

import (
	"encoding/json"
	"io"
	nethttp "net/http"

	"github.com/quickybrains/bankapp/internal/domain"
	"github.com/quickybrains/bankapp/internal/ports/http"
	httpcommon "github.com/quickybrains/bankapp/internal/ports/http/common"
	"github.com/quickybrains/bankapp/internal/ports/log"
	"github.com/quickybrains/bankapp/internal/service/transactions"
)

// Ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ ServerInterface = (*HTTPHandler)(nil)

type HTTPHandler struct {
	service *transactions.Service
	conf    *http.Config

	log log.Logger
}

func NewHandler(srv *transactions.Service, conf *http.Config, log log.Logger) *HTTPHandler {
	return &HTTPHandler{
		service: srv,
		conf:    conf,
		log:     log,
	}
}

func (handler *HTTPHandler) PostTransactions(w nethttp.ResponseWriter, r *nethttp.Request) {
	if r.Body == nil {
		handler.log.Error().Print("request body is nil")

		respondWithStatus(w, nethttp.StatusBadRequest, nil)

		return
	}

	data, err := io.ReadAll(io.LimitReader(r.Body, handler.conf.HTTPReqBodyLimit()))
	if err != nil {
		handler.log.Error().Print("io limit read")

		respondWithStatus(w, nethttp.StatusInternalServerError, nil)

		return
	}

	transactionReq := &Req{}
	err = json.Unmarshal(data, transactionReq)
	if err != nil {
		handler.log.Err(err).Print("json unmarshal withdraw request")

		respondWithStatus(w, nethttp.StatusBadRequest, nil)

		return
	}

	transaction, err := fromTransactionReqToDomainTransaction(transactionReq)
	if err != nil {
		handler.log.Err(err).Print("from transaction req to domain transaction")

		handleError(err, w, nil)

		return
	}

	err = handler.service.CreateTransaction(transaction)
	if err != nil {
		handler.log.Err(err).Print("service withdraw from account")

		handleError(err, w, nil)

		return
	}

	data, err = json.Marshal(transaction)
	if err != nil {
		handler.log.Err(err).String("transaction_id", transaction.ID.String()).Print("json marshal")

		handleError(err, w, nil)

		return
	}

	respondWithStatus(w, nethttp.StatusOK, data)
}

func (handler *HTTPHandler) GetTransactionsId(w nethttp.ResponseWriter, r *nethttp.Request, id httpcommon.IdParam) {
	transaction, err := handler.service.GetTransactionByID(domain.TransactionID(id))
	if err != nil {
		handler.log.Err(err).String("transaction_id", id.String()).Print("service get transaction by id")

		handleError(err, w, nil)

		return
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		handler.log.Err(err).String("transaction_id", id.String()).Print("json marshal")

		handleError(err, w, nil)

		return
	}

	respondWithStatus(w, nethttp.StatusOK, data)
}
