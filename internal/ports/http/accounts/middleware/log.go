package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/quickybrains/bankapp/internal/ports/http/accounts"
	"github.com/quickybrains/bankapp/internal/ports/log"
)

func Logging(lg log.Logger) accounts.MiddlewareFunc {
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

		lw.lg.Info().Print("Request headers params:")

		headers := r.Header.Clone()
		for k, v := range headers {
			lw.lg.With(k, strings.Join(v, ", ")).Info().Print("")
		}

		lw.lg.Info().Print("Request query params:")

		params := r.URL.Query()
		for k, v := range params {
			lw.lg.With(k, strings.Join(v, ", ")).Info().Print("")
		}

		lw.lg.With("request_body", string(data)).Info().Print("")

		if err := r.Body.Close(); err != nil {
			lw.lg.Err(err).Print("http request body close")

		}

		r.Body = io.NopCloser(bytes.NewBuffer(data))
		r.ContentLength = int64(len(data))

		lw.in.ServeHTTP(w, r)
	}
}
