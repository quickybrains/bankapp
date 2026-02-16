package middleware

import (
	"io"
	"net/http"

	"github.com/quickybrains/bankapp/internal/ports/http/transactions"
	"github.com/quickybrains/bankapp/internal/ports/log"
)

func Logging(lg log.Logger) transactions.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return &loggingMiddleware{
			in: h,
			lg: lg,
		}
	}
}

type loggingMiddleware struct {
	in http.Handler
	lg log.Logger
}

func (lw *loggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			lw.lg.Err(err).Print("io read all")
			return
		}

		lw.lg.With("request_body", string(data)).Info().Print("Incoming request")

		lw.in.ServeHTTP(w, r)
	}
}
