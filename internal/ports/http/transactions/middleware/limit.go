package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/quickybrains/bankapp/internal/ports/http/transactions"
	"github.com/quickybrains/bankapp/internal/ports/log"
)

func Limit(lg log.Logger, bodyLimit int64) transactions.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return &limitMiddleware{
			in:        h,
			lg:        lg,
			bodyLimit: bodyLimit,
		}
	}
}

type limitMiddleware struct {
	bodyLimit int64
	lg        log.Logger
	in        http.Handler
}

func (lw *limitMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		data, err := io.ReadAll(io.LimitReader(r.Body, lw.bodyLimit))
		if err != nil {
			lw.lg.Err(err).Print("io read limit reader")
			return
		}

		err = r.Body.Close()
		if err != nil {
			lw.lg.Err(err).Print("http request body close")
		}

		r.Body = io.NopCloser(bytes.NewBuffer(data))
		r.ContentLength = int64(len(data))

		lw.in.ServeHTTP(w, r)
	}
}
