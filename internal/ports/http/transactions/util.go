package transactions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/quickybrains/bankapp/internal/domain"
	"github.com/quickybrains/bankapp/internal/ports/http/common"
)

func respondWithStatus(w http.ResponseWriter, status int, data []byte) error {
	w.WriteHeader(status)

	if data == nil {
		return nil
	}

	_, err := w.Write(data)
	if err != nil {
		return fmt.Errorf("http response writer write, err=%w", err)
	}

	return nil
}

func respondWithStatusAndError(w http.ResponseWriter, status int, resErr error) error {
	w.WriteHeader(status)

	errResp := &common.ErrorResp{
		Message: resErr.Error(),
	}

	data, err := json.Marshal(errResp)
	if err != nil {
		return fmt.Errorf("json marshal, err=%w", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("http response writer write, err=%w", err)
	}

	return nil
}

func handleError(err error, w http.ResponseWriter, data []byte) error {
	var domainErr *domain.PublicError
	if errors.As(err, &domainErr) {
		return respondWithStatusAndError(w, http.StatusBadRequest, domainErr)
	}

	return respondWithStatus(w, http.StatusInternalServerError, data)
}
