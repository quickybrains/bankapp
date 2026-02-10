package accounts

import (
	"encoding/json"
	"io"
	nethttp "net/http"

	"github.com/quickybrains/bankapp/internal/domain"
	"github.com/quickybrains/bankapp/internal/ports/http"
	httpcommon "github.com/quickybrains/bankapp/internal/ports/http/common"
	"github.com/quickybrains/bankapp/internal/service/accounts"

	"github.com/quickybrains/bankapp/internal/ports/log"
)

// Ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ ServerInterface = (*HTTPHandler)(nil)

type HTTPHandler struct {
	service *accounts.Service
	conf    *http.Config

	log log.Logger
}

func NewHandler(srv *accounts.Service, conf *http.Config, log log.Logger) *HTTPHandler {
	return &HTTPHandler{
		service: srv,
		conf:    conf,
		log:     log,
	}
}

func (handler *HTTPHandler) PostAccounts(w nethttp.ResponseWriter, r *nethttp.Request) {
	if r.Body == nil {
		handler.log.Error().Print("request body is nil")

		respondWithStatus(w, nethttp.StatusBadRequest, nil)

		return
	}

	data, err := io.ReadAll(io.LimitReader(r.Body, handler.conf.ReqBodyLimit))
	if err != nil {
		handler.log.Error().Print("io limit read")

		respondWithStatus(w, nethttp.StatusInternalServerError, nil)

		return
	}

	accountReq := &Req{}
	err = json.Unmarshal(data, accountReq)
	if err != nil {
		handler.log.Err(err).Print("json unmarshal account request")

		respondWithStatus(w, nethttp.StatusBadRequest, nil)

		return
	}

	account, err := fromAccountReqToDomainAccount(accountReq)
	if err != nil {
		handleError(err, w, nil)

		return
	}

	err = handler.service.CreateAccount(account)
	if err != nil {
		handler.log.
			Err(err).
			String("account_id", account.ID.String()).
			Print("service create account")

		handleError(err, w, nil)

		return
	}

	res, err := json.Marshal(account)
	if err != nil {
		handler.log.
			Err(err).
			String("account_id", account.ID.String()).
			Print("json marshal")

		respondWithStatus(w, nethttp.StatusInternalServerError, nil)
	}

	respondWithStatus(w, nethttp.StatusOK, res)
}

func (handler *HTTPHandler) GetAccounts(w nethttp.ResponseWriter, r *nethttp.Request, params GetAccountsParams) {
	accounts, err := handler.service.GetAccountsByUserID(params.UserId)
	if err != nil {
		handler.log.
			Err(err).
			String("user_id", params.UserId.String()).
			Print("service get accounts by user id")

		handleError(err, w, nil)

		return
	}

	data, err := json.Marshal(accounts)
	if err != nil {
		handler.log.
			Err(err).
			String("user_id", params.UserId.String()).
			Print("json marshal")

		respondWithStatus(w, nethttp.StatusInternalServerError, nil)

		return
	}

	respondWithStatus(w, nethttp.StatusOK, data)
}

func (handler *HTTPHandler) GetAccountsId(w nethttp.ResponseWriter, r *nethttp.Request, id httpcommon.IdParam) {
	accountID := domain.AccountID(id)
	account, err := handler.service.GetAccountByID(accountID)
	if err != nil {
		handler.log.
			Err(err).
			String("account_id", accountID.String()).
			Print("service get account by id")

		handleError(err, w, nil)

		return
	}

	data, err := json.Marshal(account)
	if err != nil {
		handler.log.
			Err(err).
			String("account_id", accountID.String()).
			Print("json marshal")

		respondWithStatus(w, nethttp.StatusInternalServerError, nil)

		return
	}

	respondWithStatus(w, nethttp.StatusOK, data)
}

func (handler *HTTPHandler) PostAccountsDebit(w nethttp.ResponseWriter, r *nethttp.Request) {
	if r.Body == nil {
		handler.log.Error().Print("request body is nil")

		respondWithStatus(w, nethttp.StatusBadRequest, nil)

		return
	}

	// TODO: config limit
	data, err := io.ReadAll(io.LimitReader(r.Body, 5126))
	if err != nil {
		handler.log.Error().Print("io limit read")

		respondWithStatus(w, nethttp.StatusInternalServerError, nil)

		return
	}

	debitReq := &Debit{}
	err = json.Unmarshal(data, debitReq)
	if err != nil {
		handler.log.Err(err).Print("json unmarshal debit request")

		respondWithStatus(w, nethttp.StatusBadRequest, nil)

		return
	}

	err = handler.service.DebitOnAccount(debitReq.AccountId, debitReq.Amount)
	if err != nil {
		handler.log.Err(err).Print("service debit on account")

		handleError(err, w, nil)

		return
	}

	respondWithStatus(w, nethttp.StatusOK, nil)
}

func (handler *HTTPHandler) PostAccountsWithdraw(w nethttp.ResponseWriter, r *nethttp.Request) {
	if r.Body == nil {
		handler.log.Error().Print("request body is nil")

		respondWithStatus(w, nethttp.StatusBadRequest, nil)

		return
	}

	// TODO: config limit
	data, err := io.ReadAll(io.LimitReader(r.Body, 5126))
	if err != nil {
		handler.log.Error().Print("io limit read")

		respondWithStatus(w, nethttp.StatusInternalServerError, nil)

		return
	}

	withdrawReq := &Withdrawal{}
	err = json.Unmarshal(data, withdrawReq)
	if err != nil {
		handler.log.Err(err).Print("json unmarshal withdraw request")

		respondWithStatus(w, nethttp.StatusBadRequest, nil)

		return
	}

	err = handler.service.WithdrawFromAccount(withdrawReq.AccountId, withdrawReq.Amount)
	if err != nil {
		handler.log.Err(err).Print("service withdraw from account")

		handleError(err, w, nil)

		return
	}

	respondWithStatus(w, nethttp.StatusOK, nil)
}
